subinclude(
    "//build_defs:go",
)

go_library(
    name = "packages",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    visibility = ["//tools/driver/..."],
    deps = [
        "///third_party/go/github.com_peterebden_go-cli-init_v5//logging",
        "///third_party/go/golang.org_x_sync//errgroup",
        "///third_party/go/golang.org_x_term//:term",
        "///third_party/go/golang.org_x_tools//go/packages",
        "//tools/please_go/packageinfo",
    ],
)
