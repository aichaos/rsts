#!/usr/bin/env python3
# -*- coding: utf-8 -*-
from __future__ import unicode_literals

# RiveScript Test Suite: Python Test Runner

import codecs
from glob import glob
from os.path import basename
import six
from rivescript import RiveScript
import yaml

class TestCase:
    def __init__(self, file, name, opts):
        self.file = file
        self.name = name
        self.rs = RiveScript(
            debug=opts.get("debug", False),
            utf8=opts.get("utf8", False),
        )
        self.username = opts.get("username", "localuser")
        self.steps    = opts["tests"]

    def run(self):
        errors = False
        for step in self.steps:
            try:
                if "source" in step:
                    self.source(step)
                elif "input" in step:
                    self.input(step)
                elif "set" in step:
                    self.set(step)
                elif "assert" in step:
                    self.get(step)
                else:
                    self.warn("Unsupported test step")
            except Exception as e:
                self.fail(e)
                errors = True
                break

        sym = "❌" if errors else "✓"
        print(sym + " " + self.file + "#" + self.name)

    def source(self, step):
        self.rs.stream(step["source"])
        self.rs.sort_replies()

    def set(self, step):
        self.rs.set_uservars(self.username, step["set"])

    def get(self, step):
        for key, expect in step["assert"].items():
            cmp = self.rs.get_uservar(self.username, key)
            if cmp != expect:
                raise AssertionError("Did not get expected user variable: {}\n"
                    "Expected: {}\n"
                    "     Got: {}".format(key, expect, cmp)
                )

    def input(self, step):
        reply = self.rs.reply(self.username, step["input"])
        if type(step["reply"]) is list:
            ok = False
            for candidate in step["reply"]:
                if candidate == reply:
                    ok = True
                    break

            if not ok:
                raise AssertionError("Did not get expected reply for input: {}\n"
                    "Expected one of: {}\n"
                    "            Got: {}".format(
                        step["input"],
                        repr(step["reply"]),
                        repr(reply),
                    )
                )
        else:
            if reply != step["reply"].strip():
                raise AssertionError("Did not get expected reply for input: {}\n"
                    "Expected: {}\n"
                    "     Got: {}".format(
                        step["input"],
                        repr(step["reply"]),
                        repr(reply),
                    ))

    def fail(self, e):
        banner = "Failed: {}#{}".format(self.file, self.name)
        banner += "\n" + "=" * len(banner) + "\n"
        print(banner + str(e) + "\n\n")

    def warn(self, message):
        print(message)

if __name__ == "__main__":
    if six.PY2:
        print("This test was written for Python 3, and you're using Python 2.")
        quit()

    # Get all the test files.
    tests = glob("../tests/*.yml")

    for filename in sorted(tests):
        # Parse the YAML file.
        with codecs.open(filename, "r", "utf-8") as fh:
            data = yaml.load(fh.read())

        # from pprint import pprint
        # pprint(data)

        for name, opts in data.items():
            test = TestCase(basename(filename), name, opts)
            test.run()
