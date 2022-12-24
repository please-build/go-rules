subinclude("///go//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(["*.go"]),
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "goget",
    srcs = ["go_get.go"],
    visibility = ["//tools/..."],
    deps = [
        "//third_party/go:mod",
    ],
)
