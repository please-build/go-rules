subinclude("///shell//build_defs:shell", "//test/build_defs:e2e")

please_repo_e2e_test(
    name = "labels_test",
    plz_command = "plz query alltargets -i external > targets_with_external_label.txt",
    expected_output = {"targets_with_external_label.txt": "//:benchmark_test\n//:labels_test"},
    repo = "test_repo",
)
