subinclude("//build_defs:go")

go_binary(
    name = "main",
    srcs = ["main.go"],
)

PLZ_VERSION = "16.21.3"

remote_file(
    name = "please",
    url = f"https://github.com/thought-machine/please/releases/download/v{PLZ_VERSION}/please_{PLZ_VERSION}_${OS}_${ARCH}",
    out = "please",
    test_only = True,
    binary = True,
    visibility = ["//test/..."],
)
