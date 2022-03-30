subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(["*.go"]),
    visibility = ["//tools/please_go:bootstrap"]
)
go_library(
    name = "exec",
    srcs = [
        "exec.go",
    ],
    visibility = ["//tools/please_go/install/..."],
)
