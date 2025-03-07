subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(
        include = ["*.go"],
        exclude = ["*_test.go"],
    ),
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "embed",
    srcs = ["embed.go"],
    visibility = ["PUBLIC"],
)

go_test(
    name = "embed_test",
    srcs = ["embed_test.go"],
    data = glob(["test_data/*"]),
    deps = [
        ":embed",
        "///third_party/go/github.com_stretchr_testify//require",
        "///third_party/go/github.com_stretchr_testify//assert",
    ],
)
