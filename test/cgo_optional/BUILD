subinclude("//build_defs:cgo")

cgo_library(
    name = "cgo_optional",
    srcs = [],
    go_srcs = ["go.go"],
)

go_test(
    name = "cgo_optional_test",
    srcs = ["go_test.go"],
    deps = [
        ":cgo_optional",
        "///third_party/go/github.com_stretchr_testify//assert",
    ],
)
