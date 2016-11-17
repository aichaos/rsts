#!/usr/bin/env node

// RiveScript Test Suite: JavaScript Test Runner

const fs = require("fs"),
	yaml = require("js-yaml"),
	RiveScript = require("rivescript");

class TestCase {
	constructor(file, name, opts) {
		this.file = file;
		this.name = name;
		this.rs = new RiveScript({
			debug: opts.debug || false,
			utf8: opts.utf8 || false,
		});
		this.username = opts.username || "localuser";
		this.steps    = opts.tests;
	}

	run() {
		var errors = false;
		for (let i in this.steps) {
			let step = this.steps[i];

			try {
				if (step.source) {
					this.source(step);
				}
				else if (step.input) {
					this.input(step);
				}
				else if (step.set) {
					this.set(step);
				}
				else if (step.assert) {
					this.assert(step);
				}
				else {
					this.warn("Unsupported test step");
				}
			}
			catch (e) {
				this.fail(e);
				errors = true;
				break;
			}
		}

		let sym = errors ? "❌" : "✓";
		console.log(sym + " " + this.file + "#" + this.name);
	}

	source(step) {
		this.rs.stream(step.source);
		this.rs.sortReplies();
	}

	set(step) {
		this.rs.setUservars(this.username, step.set);
	}

	assert(step) {
		for (var key in step.assert) {
			let cmp = this.rs.getUservar(this.username, key);
			if (cmp !== step.assert[key]) {
				throw "Did not get expected user variable: " + key + "\n"
					+ "Expected: '" + step.assert[key] + "'\n"
					+ "     Got: '" + cmp + "'";
			}
		}
	}

	input(step) {
		var reply = this.rs.reply(this.username, step.input);
		if (typeof(step.reply) === "string") {
			if (reply !== step.reply.trim()) {
				throw "Did not get expected reply for input: " + step.input + "\n"
					+ "Expected: '" + step.reply + "'\n"
					+ "     Got: '" + reply + "'";
			}
		}
		else {
			let ok = false;
			for (let i in step.reply) {
				if (reply === step.reply[i]) {
					ok = true;
					break;
				}
			}

			if (!ok) {
				throw "Did not get expected reply for input: " + step.input + "\n"
					+ "Expected one of: " + JSON.stringify(step.reply) + "\n"
					+ "            Got: " + reply;
			}
		}
	}

	fail(message) {
		if (message.message !== undefined) message = message.message;
		let banner = "Failed: " + this.file + "#" + this.name;
		banner += "\n" + "=".repeat(banner.length) + "\n";
		console.error(banner + message + "\n\n");
	}

	warn(message) {
		console.warn(message);
	}
}

// Get all the test files.
const tests = fs.readdirSync("../tests");
for (let i in tests) {
	let filename = tests[i];

	// Parse the YAML file.
	if (!filename.endsWith(".yml")) continue;
	let data = yaml.safeLoad(fs.readFileSync("../tests/" + filename, "utf8"));
	for (let name in data) {
		if (!data.hasOwnProperty(name)) {
			continue;
		}

		let test = new TestCase(filename, name, data[name]);
		test.run();
	}
}
