subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(["*.go"]),
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "install",
    srcs = [
        "install.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//tools/please_go/embed",
        "//tools/please_go/install/exec",
        "//tools/please_go/install/toolchain",
    ],
)

go_test(
    name = "install_test",
    srcs = ["install_test.go"],
    data = {
        "test_data": ["test_data"],
        "go_tool": ["//third_party/go:toolchain"],
    },
    deps = [
        ":install",
        "///third_party/go/github.com_stretchr_testify//assert",
        "///third_party/go/github.com_stretchr_testify//require",
        "//tools/please_go/install/exec",
    ],
)
