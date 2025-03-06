subinclude("///go//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "generate",
    srcs = [
        "generate.go",
        "rules.go",
    ],
    visibility = ["//tools/..."],
    deps = [
        "///third_party/go/github.com_bazelbuild_buildtools//build",
        "///third_party/go/github.com_bazelbuild_buildtools//edit",
        "//tools/please_go/generate/gomoddeps",
    ],
)

go_test(
    name = "generate_test",
    srcs = glob(["*_test.go"]),
    deps = [
        ":generate",
        "///third_party/go/github.com_stretchr_testify//assert",
    ],
)
