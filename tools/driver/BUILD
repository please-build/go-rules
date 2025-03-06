subinclude(
    "//build_defs:go",
)

go_binary(
    name = "plz-gopackagesdriver",
    srcs = ["main.go"],
    visibility = ["PUBLIC"],
    deps = [
        "///third_party/go/github.com_peterebden_go-cli-init_v5//flags",
        "///third_party/go/github.com_peterebden_go-cli-init_v5//logging",
        "///third_party/go/golang.org_x_tools//go/packages",
        "//tools/driver/packages",
    ],
)

genrule(
    name = "version",
    srcs = ["VERSION"],
    outs = ["version.build_defs"],
    cmd = "echo DRIVER_VERSION = \\\"$(cat $SRCS)\\\" > $OUT",
    visibility = [
        "//package:all",
        "//tools:all",
    ],
)
