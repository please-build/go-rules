subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(["*.go"]),
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "packageinfo",
    srcs = glob(["*.go"]),
    deps = [
        "//third_party/go:xtools",
    ],
    visibility = [
        "//tools/please_go/...",
        "//tools/driver/...",
    ],
)
