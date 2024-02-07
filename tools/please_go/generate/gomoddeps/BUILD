subinclude("///go//build_defs:go")

go_library(
    name = "gomoddeps",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    visibility = ["//tools/please_go/generate/..."],
    deps = [
        "//third_party/go:mod",
    ],
)

go_test(
    name = "gomoddeps_test",
    srcs = glob(["*_test.go"]),
    data = [
        "//tools/please_go/generate/gomoddeps/testdata:test_go_mod_files",
    ],
    deps = [
        "//third_party/go:testify",
        ":gomoddeps",
    ],
)