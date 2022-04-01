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
        "//tools/please_go/filter",
        "//tools/please_go/embed",
        "//tools/please_go/test",
        "//tools/please_go/covervars",
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
        "//tools/please_go/covervars:srcs",
        "//tools/please_go/filter:srcs",
        "//tools/please_go/install:srcs",
        "//tools/please_go/install/exec:srcs",
        "//tools/please_go/install/toolchain:srcs",
        "//tools/please_go/embed:srcs",
    ],
    tools = [CONFIG.GO.GO_TOOL],
    requires = ["go_src"],
    binary = True,
    cmd = "(cd $PKG_DIR && $TOOL build please_go.go && mv please_go $OUT && $TOOL clean -modcache) || ($TOOL clean -modcache; exit 1)",
    visibility = ["PUBLIC"],
)
