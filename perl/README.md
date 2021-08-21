# Test Runner for RiveScript-Perl

To run this test:

```bash
$ perl test.pl
```

You may need to get the RiveScript module installed. Sometimes this is
available in your package manager, like `apt install librivescript-perl`
on Debian-likes or `sudo cpan RiveScript` to install it manually.
I'm not familiar with modern Perl dependency management, pull
requests welcome to improve this!

## Test Results (Aug 20 2021): 30/32

There were a couple of bugs in the Perl RiveScript module (latest
version was 2.0.3, Aug 26 2016) which were found by the test suite
and are fixed on the latest github branch at
<https://github.com/aichaos/rivescript-perl> and will be released
on CPAN at some point in the near future:

* Nested tags had an infinite loop bug, triggerable by the text
  `<a href="/?q=<star>">` in a reply segment.
* {weight} tags in triggers that had spaces next to them were not
  trimming the spaces.

There are still two other failures with the Perl RiveScript module:

* It doesn't support (@array) syntax in -Response segments, which
  is like a {random} tag but inserts a random value from the array.
* It has a Unicode matching bug with the `_` wildcard, which
  translates to regexp `(\w+)` and does not match words with
  Unicode symbols in them.

```
RiveScript v2.0.3
✓ ../tests/begin.yml#simple_begin_block
✓ ../tests/begin.yml#no_begin_block
✓ ../tests/begin.yml#blocked_begin_block
✓ ../tests/begin.yml#conditional_begin_block
✓ ../tests/bot-variables.yml#global_variables
✓ ../tests/bot-variables.yml#bot_variables
✓ ../tests/math.yml#addition
✓ ../tests/options.yml#test_concat_newline_with_conditionals
✓ ../tests/options.yml#test_concat_none_with_conditionals
✓ ../tests/options.yml#concat
✓ ../tests/options.yml#test_concat_space_with_conditionals
✓ ../tests/replies.yml#continuations
✓ ../tests/replies.yml#embedded_tags
✓ ../tests/replies.yml#questionmark
✓ ../tests/replies.yml#set_uservars
✓ ../tests/replies.yml#redirect_with_undefined_vars
Did not get expected reply for input: test random array
Expected one of: Testing alpha array., Testing beta array., Testing gamma array.
            Got: Testing (@greek) array.
Did not get expected reply for input: test two random arrays
Expected one of: Testing another alpha array., Testing another beta array., Testing another gamma array., Trying another alpha array., Trying another beta array., Trying another gamma array.
            Got: (@Test) another (@greek) array.
Did not get expected reply for input: test more arrays
Expected one of: I'm testing more alpha (@arrays)., I'm testing more beta (@arrays)., I'm testing more gamma (@arrays)., I'm trying more alpha (@arrays)., I'm trying more beta (@arrays)., I'm trying more gamma (@arrays).
            Got: I'm (@test) more (@greek) (@arrays).
Did not get expected reply for input: random format hello world
Expected one of: HELLO WORLD, hello world, Hello World, Hello world
            Got: (@format)
× ../tests/replies.yml#reply_arrays
✓ ../tests/replies.yml#redirects
✓ ../tests/replies.yml#random
✓ ../tests/replies.yml#redirect_with_undefined_input
✓ ../tests/replies.yml#conditions
✓ ../tests/replies.yml#previous
✓ ../tests/substitutions.yml#person_substitutions
✓ ../tests/substitutions.yml#message_substitutions
✓ ../tests/test-spec.yml#test_name
✓ ../tests/triggers.yml#alternatives_and_optionals
✓ ../tests/triggers.yml#trigger_arrays
✓ ../tests/triggers.yml#atomic
✓ ../tests/triggers.yml#weighted_triggers
✓ ../tests/triggers.yml#wildcards
✓ ../tests/unicode.yml#unicode
Did not get expected reply for input: My name is Bảo
Expected: Nice to meet you, bảo.
     Got: No match.
× ../tests/unicode.yml#wildcards
Passes 30 tests out of 32 (2 failed)
```