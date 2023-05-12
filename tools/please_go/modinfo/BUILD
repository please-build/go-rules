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
)

    
