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
    name = "gomoddeps",
    srcs = [
        "gomoddeps.go",
    ],
    visibility = ["//tools/please_go/generate/..."],
    deps = [
        "///third_party/go/golang.org_x_mod//modfile",
    ],
)

go_test(
    name = "gomoddeps_test",
    srcs = glob(["*_test.go"]),
    data = [
        "//tools/please_go/generate/gomoddeps/test_data:test_go_mod_files",
    ],
    deps = [
        ":gomoddeps",
        "///third_party/go/github.com_stretchr_testify//assert",
    ],
)
