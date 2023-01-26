subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs=glob("*.go"),
)

go_binary(
    name = "go_get_migrator",
    srcs = [":srcs"],
    deps = [
        "//tools/please_go/generate",
        "//third_party/go:buildtools",
        "//third_party/go:mod",
        "//third_party/go:flags",
    ],
)