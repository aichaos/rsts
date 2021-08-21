#!/usr/bin/env perl

use 5.14.0;
use lib ".";

use strict;
use warnings;
use RiveScript;
use YAML::Tiny;
use Data::Dumper;

sub main {
    say "RiveScript ${RiveScript::VERSION}";

    my @tests = sort <../tests/*.yml>;
    # sort @tests;

    my $failCount = 0;
    my $testCount = 0;
    foreach my $filename (@tests) {
        open(my $fh, '<', $filename) or die $!;
        my @yaml = <$fh>;
        close($fh);

        my ($hashref, $arrayref, $string) = Load(join("\n", @yaml));
        
        foreach my $test_name (keys %{$hashref}) {
            $testCount++;
            my $hasErrors = run_test($filename, $test_name, $hashref->{$test_name});
            if ($hasErrors) {
                $failCount++;
            }
        }
    }

    my $passes = $testCount - $failCount;
    say "Passes $passes tests out of $testCount ($failCount failed)";
}

sub run_test {
    my ($filename, $test_name, $opts) = @_;

    my $username = $opts->{username} || 'localuser';
    my $rs = new RiveScript(
        utf8 => $opts->{utf8},
        debug => $opts->{debug},
    );
    my $hasErrors;

    foreach my $step (@{$opts->{tests}}) {

        if (defined $step->{source} && $step->{source} ne "") {
            # - source: stream RiveScript code
            $rs->stream($step->{source});
            $rs->sortReplies();
        } elsif (defined $step->{input} && $step->{input} ne "") {
            # - input: check an input and reply
            my $reply = $rs->reply($username, $step->{input});

            # Expecting random replies?
            if (ref($step->{reply}) eq "ARRAY") {
                my $pass;
                foreach my $expect (@{$step->{reply}}) {
                    chomp $expect;
                    if ($reply eq $expect) {
                        $pass = 1;
                        last;
                    } else {
                    }
                }

                if (!$pass) {
                    my $expected = join ", ", @{$step->{reply}};
                    say "Did not get expected reply for input: $step->{input}\n"
                        . "Expected one of: $expected\n"
                        . "            Got: $reply";
                    $hasErrors = 1;
                }
            } elsif ($reply ne $step->{reply}) {
                chomp $step->{reply};
                if ($reply ne $step->{reply}) {
                    say "Did not get expected reply for input: $step->{input}\n"
                        . "Expected: $step->{reply}\n"
                        . "     Got: $reply";
                    $hasErrors = 1;
                }
                
            }
        } elsif (defined $step->{set} && scalar keys %{$step->{set}} > 0) {
            # - set: user variables
            foreach my $key (keys %{$step->{set}}) {
                $rs->setUservar($username, $key, $step->{set}->{$key});
            }
        } elsif (defined $step->{assert} && scalar keys %{$step->{assert}} > 0) {
            # - assert: check user variables
            foreach my $key (keys %{$step->{assert}}) {
                my $expect = $step->{assert}->{$key};
                my $actual = $rs->getUservar($username, $key);
                if ($expect ne $actual) {
                    say "Did not get expected user variable: $key\n"
                        . "Expected: $expect\n"
                        . "     Got: $actual";
                    $hasErrors = 1;
                }
            }
        } else {
            say "Unsupported test step in $test_name";
            $hasErrors = 1;
        }
    }

    my $sym = $hasErrors ? '×' : '✓';
    say "$sym $filename#$test_name";
    return $hasErrors;
}

main();