subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(["*.go"]),
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "test",
    srcs = [
        "gotest.go",
        "write_test_main.go",
    ],
    visibility = ["//tools/please_go/..."],
)

go_test(
    name = "write_test_main_test",
    srcs = ["write_test_main_test.go"],
    data = glob(["test_data/**/*.go"]),
    deps = [
        ":test",
        "///third_party/go/github.com_stretchr_testify//assert",
    ],
)
