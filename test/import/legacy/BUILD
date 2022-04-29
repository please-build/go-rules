subinclude("//build_defs:go")

package(go = {"legacy_imports": True})

go_library(
    name = "foo",
    srcs = ["foo.go"],
    visibility = ["//test/import/..."],
)
