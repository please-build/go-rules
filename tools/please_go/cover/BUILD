subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs = ["cover.go"],
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "cover",
    srcs = ["cover.go"],
    visibility = ["//tools/please_go/..."],
    deps = ["//tools/please_go/install/toolchain"],
)
