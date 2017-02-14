package com.rivescript

import org.junit.Test
import org.junit.runner.RunWith
import org.junit.runners.Parameterized
import org.yaml.snakeyaml.Yaml

import static org.hamcrest.Matchers.is
import static org.hamcrest.Matchers.isIn
import static org.junit.Assert.assertThat
import static org.junit.runners.Parameterized.Parameter
import static org.junit.runners.Parameterized.Parameters

/**
 * RiveScript Test Suite: Java Test Runner
 *
 * @author Marcel Overdijk
 */
@RunWith(Parameterized.class)
class RiveScriptTestSuiteTests {

    @Parameters(name = "{0}#{1}")
    public static Collection<Object[]> data() throws Exception {
        def testsDir = System.getProperty("user.dir") + "/../tests"
        def files = new FileNameFinder().getFileNames(testsDir, "*.yml")
        def parameters = []
        files.each {
            def file = new File(it)
            def data = new Yaml().load(new FileInputStream(file))
            data.each { name, opts ->
                 parameters << [file.name, name, opts].toArray()
            }
        }
        return parameters
    }

    @Parameter(value = 0)
    public String filename

    @Parameter(value = 1)
    public String name

    @Parameter(value = 2)
    public Map<String, Object> opts

    @Test
    public void test() {
        def utf8 = opts.utf8 ?: false
        def username = opts.username ?: "localuser"
        def steps = opts.tests

        Config config = utf8 ? Config.utf8() : Config.basic()

        RiveScript rs = new RiveScript(config)

        steps.each { step ->
            if ("source" in step) {
                rs.stream(step.source)
                rs.sortReplies()
            } else if ("input" in step) {
                def reply = rs.reply(username, step.input)
                def matcher
                if (step.reply instanceof List) {
                    matcher = isIn(step.reply)
                } else {
                    matcher = is(step.reply.trim())
                }
                assertThat("Did not get expected reply for input: ${step.input}", reply, matcher)
            } else if ("set" in step) {
                step.set.each { name, value ->
                    rs.setUservar(username, name, value.toString())
                }
            } else if ("assert" in step) {
                step.assert.each { name, expected ->
                    def actual = rs.getUservar(username, name)
                    assertThat("Did not get expected user variable: ${name}", actual, is(expected))
                }
            } else {
                throw new IllegalStateException("Unsupported test step")
            }
        }
    }
}
