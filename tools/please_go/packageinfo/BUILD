subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(["*.go"]),
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "packageinfo",
    srcs = glob(["*.go"]),
    visibility = [
        "//tools/driver/...",
        "//tools/please_go/...",
    ],
    deps = [
        "///third_party/go/golang.org_x_tools//go/packages",
    ],
)
