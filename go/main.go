// RiveScript Test Suite: Go Test Runner
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strings"

	yaml "gopkg.in/yaml.v2"

	rivescript "github.com/aichaos/rivescript-go"
	rss "github.com/aichaos/rivescript-go/src"
)

// TestCase wraps each RiveScript test.
type TestCase struct {
	file     string
	name     string
	username string
	rs       *rivescript.RiveScript
	steps    []TestStep
}

// RootSchema is the root of the YAML structure.
type RootSchema map[string]TestSchema

// TestSchema describes the YAML test files.
type TestSchema struct {
	Username string `yaml:"username"`
	UTF8     bool   `yaml:"utf8"`
	Debug    bool   `yaml:"debug"`
	Tests    []TestStep
}

// TestStep describes the YAML structure for the actual tests.
type TestStep struct {
	Source string            `yaml:"source"`
	Input  string            `yaml:"input"`
	Reply  interface{}       `yaml:"reply"`
	Assert map[string]string `yaml:"assert"`
	Set    map[string]string `yaml:"set"`
}

// NewTestCase initializes a new test.
func NewTestCase(file, name string, opts TestSchema) *TestCase {
	username := opts.Username
	if username == "" {
		username = "localuser"
	}

	return &TestCase{
		file:     file,
		name:     name,
		username: username,
		rs: rivescript.New(&rivescript.Config{
			Debug: opts.Debug,
			UTF8:  opts.UTF8,
		}),
		steps: opts.Tests,
	}
}

// Run steps through the test cases and runs them.
func (t *TestCase) Run() {
	var hasErrors bool
	for _, step := range t.steps {
		var err error

		if step.Source != "" {
			t.source(step)
		} else if step.Input != "" {
			err = t.input(step)
		} else if len(step.Set) > 0 {
			t.set(step)
		} else if len(step.Assert) > 0 {
			err = t.get(step)
		} else {
			log.Printf("Unsupported test step")
		}

		if err != nil {
			t.fail(err)
			hasErrors = true
			break
		}
	}

	var sym string
	if hasErrors {
		sym = `❌`
	} else {
		sym = `✓`
	}
	fmt.Printf("%s %s#%s\n", sym, t.file, t.name)
}

// source handles a `source` step, which parses RiveScript code.
func (t *TestCase) source(step TestStep) {
	t.rs.Stream(step.Source)
	t.rs.SortReplies()
}

// input handles an `input` step, which tests the brain for a reply.
func (t *TestCase) input(step TestStep) error {
	reply, err := t.rs.Reply(t.username, step.Input)
	if err != nil {
		return t.expectedError(step, reply, err)
	}

	// Random replies?
	if expect, ok := step.Reply.([]interface{}); ok {
		pass := false
		for _, candidate := range expect {
			cmp, ok := candidate.(string)
			if !ok {
				return fmt.Errorf(
					"Error",
				)
			}
			if cmp == reply {
				pass = true
				break
			}
		}

		if !pass {
			return fmt.Errorf(
				"Did not get expected reply for input: %s\n"+
					"Expected one of: %v\n"+
					"            Got: %s",
				step.Input,
				expect,
				reply,
			)
		}
	} else if expect, ok := step.Reply.(string); ok {
		if reply != strings.TrimSpace(expect) {
			return fmt.Errorf(
				"Did not get expected reply for input: %s\n"+
					"Expected: %s\n"+
					"     Got: %s",
				step.Input,
				expect,
				reply,
			)
		}
	} else {
		return fmt.Errorf(
			"YAML error: `reply` was neither a `string` nor a `[]string` "+
				"at %s test %s (input %s); reply was: '%v' (type %s)",
			t.file,
			t.name,
			step.Input,
			step.Reply,
			reflect.TypeOf(step.Reply),
		)
	}

	return nil
}

// expectedError inspects a Reply() error to see if it was expected by the test.
func (t *TestCase) expectedError(step TestStep, reply string, err error) error {
	// Map of expected errors to their string counterpart from the test file.
	goodErrors := map[string]error{
		"ERR: No Reply Matched": rss.ErrNoTriggerMatched,
	}

	if expect, ok := goodErrors[step.Reply.(string)]; ok {
		if err == expect {
			return nil
		}
	}

	return fmt.Errorf(
		"Got unexpected error from Reply (input step: %s; expected: %v): %s",
		step.Input,
		step.Reply,
		err,
	)
}

// set handles a `set` step, which sets user variables.
func (t *TestCase) set(step TestStep) {
	for key, value := range step.Set {
		t.rs.SetUservar(t.username, key, value)
	}
}

// get handles an `assert` step, which tests user variables.
func (t *TestCase) get(step TestStep) error {
	for key, expect := range step.Assert {
		value, err := t.rs.GetUservar(t.username, key)
		if err != nil {
			return err
		}
		if value != expect {
			return fmt.Errorf(
				"Did not get expected user variable: %s\n"+
					"Expected: %s\n"+
					"     Got: %s",
				key,
				expect,
				value,
			)
		}
	}

	return nil
}

// fail handles a failed test.
func (t *TestCase) fail(err error) {
	banner := fmt.Sprintf("Failed: %s#%s", t.file, t.name)
	fmt.Printf("%s\n%s\n%s\n\n",
		banner,
		strings.Repeat("=", len(banner)),
		err,
	)
}

func main() {
	tests, err := filepath.Glob("../tests/*.yml")
	if err != nil {
		panic(err)
	}

	for _, filename := range tests {
		yamlSource, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		data := RootSchema{}
		yaml.Unmarshal(yamlSource, &data)

		for name, opts := range data {
			test := NewTestCase(filename, name, opts)
			test.Run()
		}
	}
}
