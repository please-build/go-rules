subinclude("///go//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "gomoddeps",
    srcs = [
        ":srcs",
    ],
    visibility = ["//tools/please_go/generate/..."],
    deps = [
        "//third_party/go:mod",
    ],
)

go_test(
    name = "gomoddeps_test",
    srcs = glob(["*_test.go"]),
    data = [
        "//tools/please_go/generate/gomoddeps/test_data:test_go_mod_files",
    ],
    deps = [
        "//third_party/go:testify",
        ":gomoddeps",
    ],
)