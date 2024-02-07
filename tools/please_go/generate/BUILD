subinclude("///go//build_defs:go")

filegroup(
    name = "srcs",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    visibility = ["//tools/please_go:bootstrap"],
)

go_library(
    name = "generate",
    srcs = [":srcs"],
    visibility = ["//tools/..."],
    deps = [
        "//third_party/go:buildtools",
        "//tools/please_go/generate/gomoddeps",
    ],
)

go_test(
    name = "generate_test",
    srcs = glob(["*_test.go"]),
    deps = [
        "//third_party/go:testify",
        ":generate",
    ],
)
