subinclude("//build_defs:go")

go_library(
    name = "filter",
    srcs = ["filter.go"],
    visibility = ["//tools/please_go"],
)

filegroup(
    name = "srcs",
    srcs = glob(["*.go"]),
    visibility = ["//tools/please_go:bootstrap"],
)
