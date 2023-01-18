subinclude("///go//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(["*.go"]),
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "generate",
    srcs = ["generate.go", "rules.go", "update.go"],
    visibility = ["//tools/..."],
    deps = [
        "//third_party/go:buildtools",
        "//third_party/go:mod",
    ],
)
