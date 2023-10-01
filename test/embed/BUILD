subinclude("//build_defs:go")

go_library(
    name = "embed",
    srcs = ["embed.go"],
    resources = ["hello.txt", "subdir"],
)

go_test(
    name = "embed_test",
    srcs = ["embed_test.go"],
    resources = ["hello.txt", "subdir"],
    deps = [
        ":embed",
        "//third_party/go:testify",
    ],
)
