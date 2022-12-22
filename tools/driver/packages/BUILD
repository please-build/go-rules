subinclude("//build_defs:go")

go_library(
    name = "packages",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    deps = [
        "//third_party/go:flags",
        "//third_party/go:xsync",
        "//third_party/go:xtools",
    ],
    visibility = ["//tools/driver/..."],
)
