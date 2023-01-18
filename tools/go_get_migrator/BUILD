subinclude("//build_defs:go")

go_binary(
    name = "go_get_migrator",
    srcs = glob("*.go"),
    deps = [
        "//tools/please_go/generate",
        "//third_party/go:buildtools",
        "//third_party/go:mod",
        "//third_party/go:flags",
    ],
)