subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(["*.go"]),
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "test",
    srcs = [
        "find_cover_vars.go",
        "gotest.go",
        "write_test_main.go",
    ],
    visibility = ["//tools/please_go/..."],
    deps = [
        "//tools/please_go/install/toolchain",
    ],
)

go_test(
    name = "find_cover_vars_test",
    srcs = ["find_cover_vars_test.go"],
    data = [
        "test_data",
    ],
    deps = [
        ":test",
        "//third_party/go:testify",
    ],
)

go_test(
    name = "write_test_main_test",
    srcs = ["write_test_main_test.go"],
    data = glob(["test_data/**/*.go"]),
    deps = [
        ":test",
        "//third_party/go:testify",
    ],
)
