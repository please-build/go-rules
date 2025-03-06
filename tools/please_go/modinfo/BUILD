subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs = ["modinfo.go"],
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "modinfo",
    srcs = ["modinfo.go"],
    visibility = [
        "//tools/please_go/...",
    ],
    deps = [
        "///third_party/go/golang.org_x_mod//module",
    ],
)

go_test(
    name = "modinfo_test",
    srcs = ["modinfo_test.go"],
    external = True,
    deps = [
        "///third_party/go/github.com_stretchr_testify//assert",
        "///third_party/go/github.com_stretchr_testify//require",
    ],
)
