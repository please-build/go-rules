subinclude("//build_defs:go")

go_binary(
    name = "plz-gopackagesdriver",
    srcs = ["main.go"],
    deps = [
        "//third_party/go:flags",
        "//tools/driver/packages",
    ],
    visibility = ["PUBLIC"],
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
