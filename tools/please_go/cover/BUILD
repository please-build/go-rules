subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(
        include = ["*.go"],
        exclude = ["*_test.go"],
    ),
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "cover",
    srcs = ["cover.go"],
    visibility = ["//tools/please_go/..."],
    deps = ["//tools/please_go/install/toolchain"],
)

go_test(
    name = "cover_test",
    srcs = ["cover_test.go"],
    deps = [
        ":cover",
        "///third_party/go/github.com_stretchr_testify//assert",
        "///third_party/go/github.com_stretchr_testify//require",
    ],
)
