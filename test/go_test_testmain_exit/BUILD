subinclude("//test/build_defs:e2e")

# The test binary must exit non-zero when a test fails even when TestMain doesn't call os.Exit.
# The binary is run directly because plz test derives failures from the test output not the exit code.
please_repo_e2e_test(
    name = "go_test_testmain_exit_test",
    expected_output = {
        "exit_code": "1",
    },
    plz_command = "plz build //foo:foo_test; ./plz-out/bin/foo/foo_test > /dev/null 2>&1; echo $? > exit_code",
    repo = "test_repo",
)
