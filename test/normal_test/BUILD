subinclude("///go//build_defs:go")

go_library(
    name = "foo",
    srcs = ["foo.go"],
)

go_test(
    name = "foo_test",
    srcs = ["foo_test.go"],
    deps = [":foo"],
)