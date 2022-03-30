subinclude("//build_defs:go")

go_binary(
    name = "please_go",
    srcs = [
        "please_go.go",
    ],
    labels = ["link:plz-out/please"],
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/go:flags",
        "//tools/please_go/install",
        "//tools/please_go/test",
    ],
)


genrule(
    name = "bootstrap",
    outs = ["please_go.bootstrap"],
    srcs = [
        "go.mod",
        "go.sum",
        "please_go.go",
    ],
    deps = [
        "//tools/please_go/test:srcs",
        "//tools/please_go/install:srcs",
        "//tools/please_go/install/exec:srcs",
        "//tools/please_go/install/toolchain:srcs",
        "//tools/please_go/embed:srcs",
    ],
    requires = ["go_src"],
    binary = True,
    cmd = "(cd $PKG_DIR && go build please_go.go && mv please_go $OUT); go clean -modcache",
    visibility = ["PUBLIC"],
)
