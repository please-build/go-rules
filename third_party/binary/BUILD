remote_file(
    name = "gosec_download",
    url = ["https://github.com/securego/gosec/releases/download/v2.14.0/gosec_2.14.0_linux_amd64.tar.gz"],
    hashes = ["226bd8825b7aed3d454446d1ec094f817f37859dded4211a5b707d0f36c5fdb7"],
    test_only = True,
)

genrule(
    name = "gosec",
    srcs = [":gosec_download"],
    tools = [CONFIG.ARCAT_TOOL],
    outs = ["gosec"],
    cmd = "$TOOLS x $SRCS -s v2.14.0",
    binary = True,
    test_only = True,
    visibility = ["//tools/driver/..."],
)

