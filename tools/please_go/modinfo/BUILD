subinclude("//build_defs:go")

filegroup(
    name = "srcs",
    srcs = ["modinfo.go"],
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "modinfo",
    srcs = ["modinfo.go"],
    visibility = [
        "//tools/please_go/...",
    ],
)

go_test(
    name = "modinfo_test",
    srcs = ["modinfo_test.go"],
    external = True,
    deps = [
        ":modinfo",
        "//third_party/go:testify",
    ],
)
