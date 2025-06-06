subinclude("//build_defs:go")

go_binary(
    name = "please_go",
    srcs = [
        "please_go.go",
    ],
    labels = ["link:plz-out/please"],
    visibility = ["PUBLIC"],
    deps = [
        "///third_party/go/github.com_peterebden_go-cli-init_v5//flags",
        "//tools/please_go/cover",
        "//tools/please_go/embed",
        "//tools/please_go/filter",
        "//tools/please_go/generate",
        "//tools/please_go/install",
        "//tools/please_go/modinfo",
        "//tools/please_go/packageinfo",
        "//tools/please_go/test",
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
    cmd = "(cd $PKG_DIR && $TOOL build please_go.go && mv please_go $OUT && $TOOL clean -modcache) || ($TOOL clean -modcache; exit 1)",
    requires = ["go_src"],
    sandbox = False,
    tools = [CONFIG.GO.GO_TOOL],
    visibility = ["PUBLIC"],
    deps = [
        "//tools/please_go/cover:srcs",
        "//tools/please_go/embed:srcs",
        "//tools/please_go/filter:srcs",
        "//tools/please_go/generate:srcs",
        "//tools/please_go/generate/gomoddeps:srcs",
        "//tools/please_go/install:srcs",
        "//tools/please_go/install/exec:srcs",
        "//tools/please_go/install/toolchain:srcs",
        "//tools/please_go/modinfo:srcs",
        "//tools/please_go/packageinfo:srcs",
        "//tools/please_go/test:srcs",
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
