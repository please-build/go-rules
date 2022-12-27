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
        "//tools/please_go/covervars",
        "//tools/please_go/embed",
        "//tools/please_go/filter",
        "//tools/please_go/install",
        "//tools/please_go/packageinfo",
        "//tools/please_go/test",
        "//tools/please_go/generate",
        "//tools/please_go/goget",
    ],
)

genrule(
    name = "bootstrap",
    srcs = [
        "go.mod",
        "go.sum",
        "please_go.go",
    ],
    outs = ["please_go.bootstrap"],
    binary = True,
    sandbox = False,
    cmd = "(cd $PKG_DIR && $TOOL build please_go.go && mv please_go $OUT && $TOOL clean -modcache) || ($TOOL clean -modcache; exit 1)",
    requires = ["go_src"],
    tools = [CONFIG.GO.GO_TOOL],
    visibility = ["PUBLIC"],
    deps = [
        "//tools/please_go/covervars:srcs",
        "//tools/please_go/embed:srcs",
        "//tools/please_go/filter:srcs",
        "//tools/please_go/install:srcs",
        "//tools/please_go/install/exec:srcs",
        "//tools/please_go/install/toolchain:srcs",
        "//tools/please_go/test:srcs",
        "//tools/please_go/packageinfo:srcs",
        "//tools/please_go/generate:srcs",
        "//tools/please_go/goget:srcs",
    ],
)

genrule(
    name = "version",
    srcs = ["VERSION"],
    outs = ["version.build_defs"],
    cmd = "echo PLEASE_GO_VERSION = \\\"$(cat $SRCS)\\\" > $OUT",
    visibility = [
        "//package:all",
        "//tools:all",
    ],
)
