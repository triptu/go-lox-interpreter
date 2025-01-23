#!/usr/bin/env python3


# source - https://github.com/benhoyt/loxlox/blob/master/test.py
# NOTE: this was originally copied from the github.com/munificent/craftinginterpreters
# whic is now using Dart instead of python
# also refer to - https://github.com/munificent/craftinginterpreters/blob/benhoyt-patch-4/util/test.py

from os import listdir
from os.path import abspath, dirname, isdir, join, realpath, relpath, splitext
import re
from subprocess import Popen, PIPE
import sys
from typing import Literal

# Runs the tests.
REPO_DIR = dirname(realpath(__file__))
print("running tests from " + REPO_DIR)

expectedOutputPattern = re.compile(r"// expect: ?(.*)")
expectedErrorPattern = re.compile(r"// (Error.*)")
# for when errors from both interpreters are different
errorLinePattern = re.compile(r"// \[((java|c) )?line (\d+)\] (Error.*)")
expectedRuntimeErrorPattern = re.compile(r"// expect runtime error: (.+)")
# matches optional column number and ignores it, not making part of the capture group
syntaxErrorPattern = re.compile(r"\[.*line (\d+)(?::\d+)?\] (Error.+)")
stackTracePattern = re.compile(r"\[line (\d+)\]")
nonTestPattern = re.compile(r"// nontest")


class TestSuite:
    def __init__(
        self,
        suite_name: str,
        language: str,
        args: list[str],
        tests: dict[str, Literal["pass", "skip"]],
    ):
        self.name = suite_name
        self.language = language
        self.args = args
        self.tests = tests  # {test_path: "pass" | "skip"} more specific path override less specific

    def __str__(self):
        s = "---------- Test Suite ----------\n"
        s += "name: " + self.name + "\n"
        s += "language: " + self.language + "\n"
        s += "args: " + str(self.args) + "\n"
        s += "tests: " + str(self.tests) + "\n"
        s += "--------------------------------"
        return s


class TestRunner:
    def __init__(self):
        self.reset()

    def reset(self, suite_name: TestSuite = None, filter_path: str = ""):
        self.passed = 0
        self.failed = 0
        self.num_skipped = 0
        self.expectations = 0
        self.curr_suite: TestSuite = suite_name
        self.filter_path = filter_path

    def should_run(self, test_path: str):
        if not self.filter_path:
            return True
        rel_path = relpath(test_path, join(REPO_DIR, "test"))
        return rel_path.startswith(self.filter_path)


Runner = TestRunner()
TEST_SUITES: dict[str, TestSuite] = {}
C_SUITE_NAMES: list[str] = []
GO_SUITE_NAMES: list[str] = []


def populate_go_tests():
    def add_to_go_suite(test_name, tests_meta):
        command = "run"
        if test_name == "chap04_scanning":
            command = "tokenize"
        elif test_name == "chap06_parsing":
            command = "parse"
        elif test_name == "chap07_evaluating":
            command = "evaluate"
        args = ["./build/golox", command]
        TEST_SUITES[test_name] = TestSuite(test_name, "go", args, tests_meta)
        GO_SUITE_NAMES.append(test_name)

    earlyChapters = {
        "test/scanning": "skip",
        "test/expressions": "skip",
    }

    # not implemented the concept of nan yet
    noNaNEquality = {
        "test/number/nan_equality.lox": "skip",
    }

    # limit tests are for clox
    noLanguageLimits = {
        "test/limit/loop_too_large.lox": "skip",
        "test/limit/no_reuse_constants.lox": "skip",
        "test/limit/too_many_constants.lox": "skip",
        "test/limit/too_many_locals.lox": "skip",
        "test/limit/too_many_upvalues.lox": "skip",
        # Rely on implementing language for stack overflow checking.
        "test/limit/stack_overflow.lox": "skip",
    }

    noClasses = {
        "test/assignment/to_this.lox": "skip",
        "test/call/object.lox": "skip",
        "test/class": "skip",
        "test/closure/close_over_method_parameter.lox": "skip",
        "test/constructor": "skip",
        "test/field": "skip",
        "test/inheritance": "skip",
        "test/method": "skip",
        "test/number/decimal_point_at_eof.lox": "skip",
        "test/number/trailing_dot.lox": "skip",
        "test/operator/equals_class.lox": "skip",
        "test/operator/equals_method.lox": "skip",
        "test/operator/not_class.lox": "skip",
        "test/regression/394.lox": "skip",
        "test/super": "skip",
        "test/this": "skip",
        "test/return/in_method.lox": "skip",
        "test/variable/local_from_method.lox": "skip",
    }

    noFunctions = {
        "test/call": "skip",
        "test/closure": "skip",
        "test/for/closure_in_body.lox": "skip",
        "test/for/return_closure.lox": "skip",
        "test/for/return_inside.lox": "skip",
        "test/for/syntax.lox": "skip",
        "test/function": "skip",
        "test/operator/not.lox": "skip",
        "test/regression/40.lox": "skip",
        "test/return": "skip",
        "test/unexpected_character.lox": "skip",
        "test/while/closure_in_body.lox": "skip",
        "test/while/return_closure.lox": "skip",
        "test/while/return_inside.lox": "skip",
    }

    noResolution = {
        "test/closure/assign_to_shadowed_later.lox": "skip",
        "test/function/local_mutual_recursion.lox": "skip",
        "test/variable/collide_with_parameter.lox": "skip",
        "test/variable/duplicate_local.lox": "skip",
        "test/variable/duplicate_parameter.lox": "skip",
        "test/variable/early_bound.lox": "skip",
        # Broken because we haven"t fixed it yet by detecting the error.
        # "test/return/at_top_level.lox": "skip",
        "test/variable/use_local_in_initializer.lox": "skip",
    }

    add_to_go_suite(
        "golox",
        {
            "test": "pass",
            # These are just for earlier chapters.
            **earlyChapters,
            **noNaNEquality,
            **noLanguageLimits,
        },
    )

    add_to_go_suite(
        "chap04_scanning",
        {
            # No real interpreter yet till this chapter
            "test": "skip",
            "test/scanning": "pass",
        },
    )

    # No test for chapter 5. It just has a hardcoded main() in AstPrinter.

    add_to_go_suite(
        "chap06_parsing",
        {
            # No real interpreter yet till this chapter
            "test": "skip",
            "test/expressions/parse.lox": "pass",
        },
    )

    add_to_go_suite(
        "chap07_evaluating",
        {
            # No real interpreter yet till this chapter
            "test": "skip",
            "test/expressions/evaluate.lox": "pass",
        },
    )

    add_to_go_suite(
        "chap08_statements",
        {
            "test": "pass",
            **earlyChapters,
            **noNaNEquality,
            **noLanguageLimits,
            **noFunctions,
            **noResolution,
            **noClasses,
            # No control flow.
            "test/block/empty.lox": "skip",
            "test/for": "skip",
            "test/if": "skip",
            "test/logical_operator": "skip",
            "test/while": "skip",
            "test/variable/unreached_undefined.lox": "skip",
        },
    )

    add_to_go_suite(
        "chap09_control",
        {
            "test": "pass",
            **earlyChapters,
            **noNaNEquality,
            **noLanguageLimits,
            **noFunctions,
            **noResolution,
            **noClasses,
        },
    )

    add_to_go_suite(
        "chap10_functions",
        {
            "test": "pass",
            **earlyChapters,
            **noNaNEquality,
            **noLanguageLimits,
            **noResolution,
            **noClasses,
        },
    )

    add_to_go_suite(
        "chap11_resolving",
        {
            "test": "pass",
            **earlyChapters,
            **noNaNEquality,
            **noLanguageLimits,
            **noClasses,
        },
    )

    add_to_go_suite(
        "chap12_classes",
        {
            "test": "pass",
            **earlyChapters,
            **noLanguageLimits,
            **noNaNEquality,
            # No inheritance.
            "test/class/local_inherit_other.lox": "skip",
            "test/class/local_inherit_self.lox": "skip",
            "test/class/inherit_self.lox": "skip",
            "test/class/inherited_method.lox": "skip",
            "test/inheritance": "skip",
            "test/regression/394.lox": "skip",
            "test/super": "skip",
        },
    )

    add_to_go_suite(
        "chap13_inheritance",
        {
            "test": "pass",
            **earlyChapters,
            **noNaNEquality,
            **noLanguageLimits,
        },
    )


def populate_clox_tests():
    def add_to_c_suite(name, tests):
        if name == "clox":
            path = "build/cloxd"
        else:
            path = "build/" + name

        TEST_SUITES[name] = TestSuite(name, "c", [path], tests)
        C_SUITE_NAMES.append(name)

    add_to_c_suite(
        "clox",
        {
            "test": "pass",
            # These are just for earlier chapters.
            "test/scanning": "skip",
            "test/expressions": "skip",
        },
    )

    add_to_c_suite(
        "chap17_compiling",
        {
            # No real interpreter yet till this chapter
            "test": "skip",
            "test/expressions/evaluate.lox": "pass",
        },
    )

    add_to_c_suite(
        "chap18_types",
        {
            # No real interpreter yet till this chapter
            "test": "skip",
            "test/expressions/evaluate.lox": "pass",
        },
    )

    add_to_c_suite(
        "chap19_strings",
        {
            # No real interpreter yet till this chapter
            "test": "skip",
            "test/expressions/evaluate.lox": "pass",
        },
    )

    add_to_c_suite(
        "chap20_hash",
        {
            # No real interpreter yet till this chapter
            "test": "skip",
            "test/expressions/evaluate.lox": "pass",
        },
    )

    add_to_c_suite(
        "chap21_global",
        {
            "test": "pass",
            # These are just for earlier chapters.
            "test/scanning": "skip",
            "test/expressions": "skip",
            # No control flow.
            "test/block/empty.lox": "skip",
            "test/for": "skip",
            "test/if": "skip",
            "test/limit/loop_too_large.lox": "skip",
            "test/logical_operator": "skip",
            "test/variable/unreached_undefined.lox": "skip",
            "test/while": "skip",
            # No blocks.
            "test/assignment/local.lox": "skip",
            "test/variable/in_middle_of_block.lox": "skip",
            "test/variable/in_nested_block.lox": "skip",
            "test/variable/scope_reuse_in_different_blocks.lox": "skip",
            "test/variable/shadow_and_local.lox": "skip",
            "test/variable/undefined_local.lox": "skip",
            # No local variables.
            "test/block/scope.lox": "skip",
            "test/variable/duplicate_local.lox": "skip",
            "test/variable/shadow_global.lox": "skip",
            "test/variable/shadow_local.lox": "skip",
            "test/variable/use_local_in_initializer.lox": "skip",
            # No functions.
            "test/call": "skip",
            "test/closure": "skip",
            "test/function": "skip",
            "test/limit/no_reuse_constants.lox": "skip",
            "test/limit/stack_overflow.lox": "skip",
            "test/limit/too_many_constants.lox": "skip",
            "test/limit/too_many_locals.lox": "skip",
            "test/limit/too_many_upvalues.lox": "skip",
            "test/regression/40.lox": "skip",
            "test/return": "skip",
            "test/unexpected_character.lox": "skip",
            "test/variable/collide_with_parameter.lox": "skip",
            "test/variable/duplicate_parameter.lox": "skip",
            "test/variable/early_bound.lox": "skip",
            # No classes.
            "test/assignment/to_this.lox": "skip",
            "test/class": "skip",
            "test/constructor": "skip",
            "test/field": "skip",
            "test/inheritance": "skip",
            "test/method": "skip",
            "test/number/decimal_point_at_eof.lox": "skip",
            "test/number/trailing_dot.lox": "skip",
            "test/operator/equals_class.lox": "skip",
            "test/operator/equals_method.lox": "skip",
            "test/operator/not.lox": "skip",
            "test/operator/not_class.lox": "skip",
            "test/super": "skip",
            "test/this": "skip",
            "test/variable/local_from_method.lox": "skip",
        },
    )

    add_to_c_suite(
        "chap22_local",
        {
            "test": "pass",
            # These are just for earlier chapters.
            "test/scanning": "skip",
            "test/expressions": "skip",
            # No control flow.
            "test/block/empty.lox": "skip",
            "test/for": "skip",
            "test/if": "skip",
            "test/limit/loop_too_large.lox": "skip",
            "test/logical_operator": "skip",
            "test/variable/unreached_undefined.lox": "skip",
            "test/while": "skip",
            # No functions.
            "test/call": "skip",
            "test/closure": "skip",
            "test/function": "skip",
            "test/limit/no_reuse_constants.lox": "skip",
            "test/limit/stack_overflow.lox": "skip",
            "test/limit/too_many_constants.lox": "skip",
            "test/limit/too_many_locals.lox": "skip",
            "test/limit/too_many_upvalues.lox": "skip",
            "test/regression/40.lox": "skip",
            "test/return": "skip",
            "test/unexpected_character.lox": "skip",
            "test/variable/collide_with_parameter.lox": "skip",
            "test/variable/duplicate_parameter.lox": "skip",
            "test/variable/early_bound.lox": "skip",
            # No classes.
            "test/assignment/to_this.lox": "skip",
            "test/class": "skip",
            "test/constructor": "skip",
            "test/field": "skip",
            "test/inheritance": "skip",
            "test/method": "skip",
            "test/number/decimal_point_at_eof.lox": "skip",
            "test/number/trailing_dot.lox": "skip",
            "test/operator/equals_class.lox": "skip",
            "test/operator/equals_method.lox": "skip",
            "test/operator/not.lox": "skip",
            "test/operator/not_class.lox": "skip",
            "test/super": "skip",
            "test/this": "skip",
            "test/variable/local_from_method.lox": "skip",
        },
    )

    add_to_c_suite(
        "chap23_jumping",
        {
            "test": "pass",
            # These are just for earlier chapters.
            "test/scanning": "skip",
            "test/expressions": "skip",
            # No functions.
            "test/call": "skip",
            "test/closure": "skip",
            "test/for/closure_in_body.lox": "skip",
            "test/for/return_closure.lox": "skip",
            "test/for/return_inside.lox": "skip",
            "test/for/syntax.lox": "skip",
            "test/function": "skip",
            "test/limit/no_reuse_constants.lox": "skip",
            "test/limit/stack_overflow.lox": "skip",
            "test/limit/too_many_constants.lox": "skip",
            "test/limit/too_many_locals.lox": "skip",
            "test/limit/too_many_upvalues.lox": "skip",
            "test/regression/40.lox": "skip",
            "test/return": "skip",
            "test/unexpected_character.lox": "skip",
            "test/variable/collide_with_parameter.lox": "skip",
            "test/variable/duplicate_parameter.lox": "skip",
            "test/variable/early_bound.lox": "skip",
            "test/while/closure_in_body.lox": "skip",
            "test/while/return_closure.lox": "skip",
            "test/while/return_inside.lox": "skip",
            # No classes.
            "test/assignment/to_this.lox": "skip",
            "test/class": "skip",
            "test/constructor": "skip",
            "test/field": "skip",
            "test/inheritance": "skip",
            "test/method": "skip",
            "test/number/decimal_point_at_eof.lox": "skip",
            "test/number/trailing_dot.lox": "skip",
            "test/operator/equals_class.lox": "skip",
            "test/operator/equals_method.lox": "skip",
            "test/operator/not.lox": "skip",
            "test/operator/not_class.lox": "skip",
            "test/super": "skip",
            "test/this": "skip",
            "test/variable/local_from_method.lox": "skip",
        },
    )

    add_to_c_suite(
        "chap24_calls",
        {
            "test": "pass",
            # These are just for earlier chapters.
            "test/scanning": "skip",
            "test/expressions": "skip",
            # No closures.
            "test/closure": "skip",
            "test/for/closure_in_body.lox": "skip",
            "test/for/return_closure.lox": "skip",
            "test/function/local_recursion.lox": "skip",
            "test/limit/too_many_upvalues.lox": "skip",
            "test/regression/40.lox": "skip",
            "test/while/closure_in_body.lox": "skip",
            "test/while/return_closure.lox": "skip",
            # No classes.
            "test/assignment/to_this.lox": "skip",
            "test/call/object.lox": "skip",
            "test/class": "skip",
            "test/constructor": "skip",
            "test/field": "skip",
            "test/inheritance": "skip",
            "test/method": "skip",
            "test/number/decimal_point_at_eof.lox": "skip",
            "test/number/trailing_dot.lox": "skip",
            "test/operator/equals_class.lox": "skip",
            "test/operator/equals_method.lox": "skip",
            "test/operator/not.lox": "skip",
            "test/operator/not_class.lox": "skip",
            "test/return/in_method.lox": "skip",
            "test/super": "skip",
            "test/this": "skip",
            "test/variable/local_from_method.lox": "skip",
        },
    )

    add_to_c_suite(
        "chap25_closures",
        {
            "test": "pass",
            # These are just for earlier chapters.
            "test/scanning": "skip",
            "test/expressions": "skip",
            # No classes.
            "test/assignment/to_this.lox": "skip",
            "test/call/object.lox": "skip",
            "test/class": "skip",
            "test/closure/close_over_method_parameter.lox": "skip",
            "test/constructor": "skip",
            "test/field": "skip",
            "test/inheritance": "skip",
            "test/method": "skip",
            "test/number/decimal_point_at_eof.lox": "skip",
            "test/number/trailing_dot.lox": "skip",
            "test/operator/equals_class.lox": "skip",
            "test/operator/equals_method.lox": "skip",
            "test/operator/not.lox": "skip",
            "test/operator/not_class.lox": "skip",
            "test/return/in_method.lox": "skip",
            "test/super": "skip",
            "test/this": "skip",
            "test/variable/local_from_method.lox": "skip",
        },
    )

    add_to_c_suite(
        "chap26_garbage",
        {
            "test": "pass",
            # These are just for earlier chapters.
            "test/scanning": "skip",
            "test/expressions": "skip",
            # No classes.
            "test/assignment/to_this.lox": "skip",
            "test/call/object.lox": "skip",
            "test/class": "skip",
            "test/closure/close_over_method_parameter.lox": "skip",
            "test/constructor": "skip",
            "test/field": "skip",
            "test/inheritance": "skip",
            "test/method": "skip",
            "test/number/decimal_point_at_eof.lox": "skip",
            "test/number/trailing_dot.lox": "skip",
            "test/operator/equals_class.lox": "skip",
            "test/operator/equals_method.lox": "skip",
            "test/operator/not.lox": "skip",
            "test/operator/not_class.lox": "skip",
            "test/return/in_method.lox": "skip",
            "test/super": "skip",
            "test/this": "skip",
            "test/variable/local_from_method.lox": "skip",
        },
    )

    add_to_c_suite(
        "chap27_classes",
        {
            "test": "pass",
            # These are just for earlier chapters.
            "test/scanning": "skip",
            "test/expressions": "skip",
            # No inheritance.
            "test/class/local_inherit_self.lox": "skip",
            "test/class/inherit_self.lox": "skip",
            "test/class/inherited_method.lox": "skip",
            "test/inheritance": "skip",
            "test/super": "skip",
            # No methods.
            "test/assignment/to_this.lox": "skip",
            "test/class/local_reference_self.lox": "skip",
            "test/class/reference_self.lox": "skip",
            "test/closure/close_over_method_parameter.lox": "skip",
            "test/constructor": "skip",
            "test/field/get_and_set_method.lox": "skip",
            "test/field/method.lox": "skip",
            "test/field/method_binds_this.lox": "skip",
            "test/method": "skip",
            "test/operator/equals_class.lox": "skip",
            "test/operator/equals_method.lox": "skip",
            "test/return/in_method.lox": "skip",
            "test/this": "skip",
            "test/variable/local_from_method.lox": "skip",
        },
    )

    add_to_c_suite(
        "chap28_methods",
        {
            "test": "pass",
            # These are just for earlier chapters.
            "test/scanning": "skip",
            "test/expressions": "skip",
            # No inheritance.
            "test/class/local_inherit_self.lox": "skip",
            "test/class/inherit_self.lox": "skip",
            "test/class/inherited_method.lox": "skip",
            "test/inheritance": "skip",
            "test/super": "skip",
        },
    )

    add_to_c_suite(
        "chap29_superclasses",
        {
            "test": "pass",
            # These are just for earlier chapters.
            "test/scanning": "skip",
            "test/expressions": "skip",
        },
    )

    add_to_c_suite(
        "chap30_optimization",
        {
            "test": "pass",
            # These are just for earlier chapters.
            "test/scanning": "skip",
            "test/expressions": "skip",
        },
    )


class Test:
    def __init__(self, path):
        self.path = path  # e.g. tests/assignment/associativity.lox
        self.output = []
        self.compile_errors = set()
        self.runtime_error_line = 0
        self.runtime_error_message = None
        self.exit_code = 0
        self.failures = []

    def parse(self):
        # Get the path components.
        parts = self.path.split("/")
        subpath = ""
        pass_or_skip = None

        # Figure out whether the test is skipped. We don't break out of this loop because
        # we want lines for more specific paths to override more general ones.
        for part in parts:
            if subpath:
                subpath += "/"
            subpath += part

            if subpath in Runner.curr_suite.tests:
                pass_or_skip = Runner.curr_suite.tests[subpath]

        if not pass_or_skip or pass_or_skip == "skip":
            if not pass_or_skip:
                print(
                    'Unknown test state(whether to run or skip) for "{}, skipping it".'.format(
                        self.path
                    )
                )
            Runner.num_skipped += 1
            return

        line_num = 1
        with open(self.path, "r") as file:
            for line in file:
                match = expectedOutputPattern.search(line)
                if match:
                    self.output.append((match.group(1), line_num))
                    Runner.expectations += 1

                match = expectedErrorPattern.search(line)
                if match:
                    compile_err_with_line = "[line {0}] {1}".format(
                        line_num, match.group(1)
                    )
                    self.compile_errors.add(compile_err_with_line)

                    # If we expect a compile error, it should exit with EX_DATAERR.
                    self.exit_code = 65
                    Runner.expectations += 1

                match = errorLinePattern.search(line)
                if match:
                    # The two interpreters are slightly different in terms of which
                    # cascaded errors may appear after an initial compile error because
                    # their panic mode recovery is a little different. To handle that,
                    # the tests can indicate if an error line should only appear for a
                    # certain interpreter.
                    language = match.group(2)
                    if (
                        not language
                        or language == Runner.curr_suite.language
                        or (language == "java" and Runner.curr_suite.language == "go")
                    ):
                        self.compile_errors.add(f"[line {match[3]}] {match[4]}")

                        # If we expect a compile error, it should exit with EX_DATAERR.
                        self.exit_code = 65
                        Runner.expectations += 1

                match = expectedRuntimeErrorPattern.search(line)
                if match:
                    self.runtime_error_line = line_num
                    self.runtime_error_message = match.group(1)
                    # If we expect a runtime error, it should exit with EX_SOFTWARE.
                    self.exit_code = 70
                    Runner.expectations += 1

                match = nonTestPattern.search(line)
                if match:
                    # Not a test file at all, so ignore it.
                    return False

                line_num += 1

        # If we got here, it's a valid test.
        return True

    def run(self):
        # Invoke the test suite and run the test.
        args = Runner.curr_suite.args[:]
        args.append(self.path)
        proc = Popen(args, stdin=PIPE, stdout=PIPE, stderr=PIPE)

        # print("running test", self.path)

        out, err = proc.communicate()
        self.validate(proc.returncode, out, err)

    def validate(self, exit_code, out, err):
        if self.compile_errors and self.runtime_error_message:
            self.fail("Test error: Cannot expect both compile and runtime errors.")
            return

        try:
            out = out.decode("utf-8").replace("\r\n", "\n")
            err = err.decode("utf-8").replace("\r\n", "\n")
        except Exception as e:
            self.fail("Error decoding output.", e)

        error_lines = err.split("\n")

        # Validate that an expected runtime error occurred.
        if self.runtime_error_message:
            self.validate_runtime_error(error_lines)
        else:
            self.validate_compile_errors(error_lines)

        self.validate_exit_code(exit_code, error_lines)
        self.validate_output(out)

    def validate_runtime_error(self, error_lines):
        if len(error_lines) < 2:
            self.fail(
                'Expected runtime error "{0}" and got none.', self.runtime_error_message
            )
            return

        # Skip any compile errors. This can happen if there is a compile error in
        # a module loaded by the module being tested.
        line = 0
        while syntaxErrorPattern.search(error_lines[line]):
            line += 1

        if error_lines[line] != self.runtime_error_message:
            self.fail(
                'Expected runtime error "{0}" and got:', self.runtime_error_message
            )
            self.fail(error_lines[line])

        # Make sure the stack trace has the right line. Skip over any lines that
        # come from builtin libraries.
        match = False
        stack_lines = error_lines[line + 1 :]
        for stack_line in stack_lines:
            match = stackTracePattern.search(stack_line)
            if match:
                break

        if not match:
            self.fail("Expected stack trace with line numbers and got:")
            for stack_line in stack_lines:
                self.fail(stack_line)
        else:
            stack_line = int(match.group(1))
            if stack_line != self.runtime_error_line:
                self.fail(
                    "Expected runtime error on line {0} but was on line {1}.",
                    self.runtime_error_line,
                    stack_line,
                )

    def validate_compile_errors(self, error_lines):
        # Validate that every compile error was expected.
        found_errors = set()
        num_unexpected = 0
        for line in error_lines:
            match = syntaxErrorPattern.search(line)
            if match:
                error = f"[line {match.group(1)}] {match.group(2)}"
                if error in self.compile_errors:
                    found_errors.add(error)
                else:
                    if num_unexpected < 10:
                        self.fail("Unexpected compile error:")
                        self.fail(line)
                    num_unexpected += 1
            elif line != "":
                if num_unexpected < 10:
                    self.fail("Unexpected output on stderr:")
                    self.fail(line)
                num_unexpected += 1

        if num_unexpected > 10:
            self.fail("(truncated " + str(num_unexpected - 10) + " more..)")

        # Validate that every expected error occurred.
        for error in self.compile_errors - found_errors:
            self.fail("Missing expected compile error: {0}", error)

    def validate_exit_code(self, exit_code, error_lines):
        if exit_code == self.exit_code:
            return

        if len(error_lines) > 10:
            error_lines = error_lines[0:10]
            error_lines.append("(truncated..)")
        self.fail(
            "Expected exit code {0} and got {1}. Stderr:", self.exit_code, exit_code
        )
        self.failures += error_lines

    def validate_output(self, out):
        # Remove the trailing last empty line.
        out_lines = out.split("\n")
        if out_lines[-1] == "":
            del out_lines[-1]

        index = 0
        for line in out_lines:
            if index >= len(self.output):
                self.fail('Got output "{0}" when none was expected.', line)
            elif self.output[index][0] != line:
                self.fail(
                    'Expected output "{0}" on line {1} and got "{2}".',
                    self.output[index][0],
                    self.output[index][1],
                    line,
                )
            index += 1

        while index < len(self.output):
            self.fail(
                'Missing expected output "{0}" on line {1}.',
                self.output[index][0],
                self.output[index][1],
            )
            index += 1

    def fail(self, message, *args):
        if args:
            message = message.format(*args)
        self.failures.append(message)

    def __str__(self):
        s = "---- Test {0} expectations ----\n".format(self.path)
        # print all the parsed information. this will be called after
        # parse is called and before run is called
        for i in range(len(self.output)):
            s += "{0}: {1}\n".format(self.output[i][0], self.output[i][1])

        if self.failures:
            s += "--- Failures ----\n"
            for f in self.failures:
                s += f + "\n"

        return s


## --------------- Print utils start ---------------


def supports_ansi():
    return sys.platform != "win32" and sys.stdout.isatty()


def color_text(text, color):
    """Converts text to a string and wraps it in the ANSI escape sequence for
    color, if supported."""

    if not supports_ansi():
        return str(text)

    return color + str(text) + "\033[0m"


def green(text):
    return color_text(text, "\033[32m")


def pink(text):
    return color_text(text, "\033[91m")


def red(text):
    return color_text(text, "\033[31m")


def yellow(text):
    return color_text(text, "\033[33m")


def gray(text):
    return color_text(text, "\033[1;30m")


def print_line(line=None):
    if supports_ansi():
        # Erase the line.
        print("\033[2K", end="")
        # Move the cursor to the beginning.
        print("\r", end="")
    else:
        print()
    if line:
        print(line, end="")
        sys.stdout.flush()


## --------------- Print utils end ---------------


# runs the callback on all files in the directory at any depth
def run_recursively(dir, callback):
    dir = abspath(dir)
    for file in sorted(listdir(dir)):
        nfile = join(dir, file)
        if isdir(nfile):
            run_recursively(nfile, callback)
        else:
            callback(nfile)


# this runs for every single test file
def run_script(test_path):
    if (
        "benchmark" in test_path
        or splitext(test_path)[1] != ".lox"
        or not Runner.should_run(test_path)
    ):
        return

    # Make a nice short path relative to the working directory and normalize it to use "/" since
    test_path = relpath(test_path).replace("\\", "/")

    # Read the test and parse out the expectations.
    test = Test(test_path)
    if not test.parse():
        return  # It's a skipped or non-test file.

    test.run()

    # Display the results.
    if len(test.failures) == 0:
        Runner.passed += 1
        print_line(green("PASS") + ": " + test_path)
    else:
        Runner.failed += 1
        print_line(red("FAIL") + ": " + test_path)
        print("")
        for failure in test.failures:
            print("      " + pink(failure))


def run_suite(name, filter_path: str):
    Runner.reset(TEST_SUITES[name], filter_path)

    run_recursively(join(REPO_DIR, "test"), run_script)
    print_line()

    if Runner.failed == 0:
        print(
            "All "
            + green(Runner.passed)
            + " tests passed ("
            + str(Runner.expectations)
            + " expectations)."
        )
    else:
        print(
            green(Runner.passed)
            + " tests passed. "
            + red(Runner.failed)
            + " tests failed."
        )

    return Runner.failed == 0


def make_go_build():
    command = "go build -o ./build/golox ./cmd/cli/main.go"
    proc = Popen(command, shell=True)
    proc.wait()
    if proc.returncode != 0:
        print("Error building golox")
        sys.exit(1)


def main(argv):
    populate_go_tests()
    # populate_clox_tests()

    if len(argv) < 2 or len(argv) > 3:
        print("Usage: test.py <test_suite> [filter]")
        print("<test_suite> should be one of:\n- " + "\n- ".join(TEST_SUITES.keys()))
        sys.exit(1)

    filter_path = ""  # only run tests in matching path
    if len(argv) == 3:
        filter_path = argv[2]

    suite_name = sys.argv[1]
    if suite_name not in TEST_SUITES:
        print('Unknown test suite "{}"'.format(argv[1]))
        sys.exit(1)

    make_go_build()

    if not run_suite(suite_name, filter_path):
        sys.exit(1)


if __name__ == "__main__":
    main(sys.argv)
