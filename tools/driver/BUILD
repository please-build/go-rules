subinclude("//build_defs:go")

go_binary(
    name = "plz-gopackagesdriver",
    srcs = ["main.go"],
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/go:flags",
        "//third_party/go:xtools",
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
