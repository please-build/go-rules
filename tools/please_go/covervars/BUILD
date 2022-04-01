subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(["*.go"]),
    visibility = ["//tools/please_go:bootstrap"]
)

go_library(
    name = "covervars",
    srcs = [
        "cover_vars.go",
    ],
    visibility = ["PUBLIC"],
)