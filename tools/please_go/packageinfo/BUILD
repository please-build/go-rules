subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(["*.go"]),
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "packageinfo",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    visibility = [
        "//tools/driver/...",
        "//tools/please_go/...",
    ],
    deps = ["///third_party/go/golang.org_x_tools//go/packages"],
)

go_test(
    name = "test",
    srcs = glob(["*_test.go"]),
    deps = [
        ":packageinfo",
        "///third_party/go/github.com_stretchr_testify//assert",
        "///third_party/go/github.com_stretchr_testify//require",
        "///third_party/go/golang.org_x_tools//go/packages",
    ],
)
