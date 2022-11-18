subinclude("//build_defs:go", "///e2e//build_defs:e2e")

go_binary(
    name = "main",
    srcs = ["main.go"],
    definitions = {"github.com/please-build/go-rules/test/definitions/version.Version": "v1.0.0"},
    deps = ["//test/definitions/version"],
)

plz_e2e_test(
    name = "definitions_test",
    cmd = "plz run //test/definitions:main",
    expect_output_contains = "v1.0.0",
)
