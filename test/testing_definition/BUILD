subinclude("//build_defs:go")
subinclude("//test/build_defs:e2e")

go_test(
    name = "testing_definition_test",
    srcs = ["testing_definition_test.go"],
    deps = [
        "///third_party/go/github.com_stretchr_testify//assert",
    ],
)

please_repo_e2e_test(
    name = "testing_definition_integration_test",
    plz_command = "plz test --show_all_output //sample:testing_flag_test -- -test.v",
    repo = "test_repo",
    expect_output_contains = {
        "output": "testing.Testing() = true",
    },
)
