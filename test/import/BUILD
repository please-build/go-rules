subinclude("//build_defs:go")

go_test(
    name = "import_test",
    srcs = ["import_test.go"],
    deps = [
        "//test/import/legacy:foo",
        "//test/message",
        "///third_party/go/github.com_stretchr_testify//assert",
    ],
)
