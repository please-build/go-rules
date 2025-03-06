subinclude("//build_defs:go")

go_binary(
    name = "require_provide",
    srcs = ["main.go"],
)

go_test(
    name = "require_provide_test",
    srcs = ["require_provide_test.go"],
    deps = [
        ":require_provide",
        "///third_party/go/github.com_stretchr_testify//assert",
    ],
)
