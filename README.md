# RiveScript Test Suite (RSTS)

This is the unified test suite for the [RiveScript][1] chatbot scripting
language. It serves as a centralized set of RiveScript tests that can verify
the accuracy of any implementation of RiveScript.

If you're developing your own implementation of RiveScript, passing the tests
in this set should give you the confidence that your implementation is at least
as correct as the official implementations.

As of August 20, 2021 the following official implementations of RiveScript
pass the entire test suite:

* [x] Go
* [x] Java
* [x] JavaScript
* [ ] Perl (couple issues, [see report](perl/))
* [x] Python

The YAML test files are in the `tests/` directory, and individual test runner
scripts for each programming language are in their respective directories.

Each programming language directory should be kept minimal: just a test runner
script and optional metadata (for dependency declaration, etc.) should be in
each directory.

# Test Schema

For maximum language interoperability, the test files are written in [YAML][2]
format. There are YAML libraries available for most popular programming
languages.

The test files follow this format:

```yaml
# Each top level key is a unique name for the test.
test_name:
  # Tests can optionally define these global options to their RiveScript instance
  username: "localuser"  # Override the username; this is the default one.
  utf8: true

  # The tests themselves are broken into a series of "Test Actions": it's an
  # array of tasks to be processed from top to bottom, and each action can
  # stream RiveScript source code (additive), test inputs and outputs and
  # manage user variables. See "Test Actions" below for more details.
  tests:
    # Stream in RiveScript source code to test against. Most tests will start with
    # one of these, and some tests may have multiple source statements that should
    # be streamed over top of the existing brain.
    - source: |
        // RiveScript source is included as a YAML
        // literal-style multi-line string.
        + hello bot
        - Hello human!

        + how are you
        - Good, you?
        - Alright, how are you?

        + my name is *
        - <set name=<formal>>Nice to meet you, <get name>.

        + who am i
        - Your name is <get name>, right?

    # A simple input and output
    - input: "Hello bot!"
      reply: "Hello human!"

    # You can test random outputs too, just use an array for the 'reply'
    - input: "How are you?"
      reply:
        - Good, you?
        - Alright, how are you?

    # Testing setting a user variable...
    - input: "My name is Alice"
      reply: "Nice to meet you, Alice."

    # ...and verify it was set
    - input: "Who am I?"
      reply: "Your name is Alice, right?"

    # ...by asserting the user variable is correct.
    - assert:
        name: "Alice"

    # Or you can set the user variable directly:
    - set:
        name: "Bob"

    # And verify:
    - assert:
        name: "Bob"

    # Another way to verify:
    - input: "Who am I?"
      reply: "Your name is Bob, right?"

    # And if you need to stream in additional code:
    - source: |
        + what is your name
        - I'm just a RiveScript bot.

    - input: "What is your name?"
      reply: "I'm just a RiveScript bot."
```

## Test Actions

The types of actions that your test runner must support are as follows:

* Stream RiveScript source code. Each test should begin with a blank brain, and
  each `stream` action should stream its source code on top of the existing
  brain (updating the reply set it already has).

    ```yaml
    - source: |
      + hello bot
      - Hello human!
    ```

* Check an input and reply (or replies):

    ```yaml
    # A one-to-one mapping
    - input: "Hello bot"
      reply: "Hello human"

    # A one-to-many mapping, if random replies can be given
    - input: "Hello bot"
      reply:
        - "Hello human"
        - "Hi there"
    ```

* Set a user variable:

    ```yaml
    # You can set as many variables as you need. Use key/value pairs
    # for mapping a user variable to its value.
    - set:
        name: "Alice"
        age: "5"
    ```

* Assert the values of user variables:

    ```yaml
    # Similar to `set`, provide key/value pairs of variable names and
    # their expected values. Multiple variables can be used.
    - assert:
        name: "Alice"
        age: "5"
    ```

## Writing a Test Runner

To verify that your implementation of the test runner is processing the YAML
syntax correctly, there are a pair of files named `test-spec.yml` and
`test-spec.json` which contain identical data, but one is in YAML format and
the other is in JSON. Use your programming language's deep comparison function
to verify both files are parsed the same.

# See Also

* [RiveScript.com][1]
* [The AiChaos GitHub Organization][3]

[1]: https://www.rivescript.com/
[2]: http://yaml.org/
[3]: https://github.com/aichaos
