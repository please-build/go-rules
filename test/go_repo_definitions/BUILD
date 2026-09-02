subinclude("//test/build_defs:e2e")

please_repo_e2e_test(
    name = "go_repo_definitions_test",
    expected_output = {
        "test.out": "go_repo value",
    },
    plz_command = "plz run ///third_party/go/test.please.build_module//cmd > test.out",
    repo = "test_repo",
)
