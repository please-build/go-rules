remote_file(
    name = "gosec_download",
    hashes = ["226bd8825b7aed3d454446d1ec094f817f37859dded4211a5b707d0f36c5fdb7"],
    test_only = True,
    url = ["https://github.com/securego/gosec/releases/download/v2.14.0/gosec_2.14.0_linux_amd64.tar.gz"],
)

genrule(
    name = "gosec",
    srcs = [":gosec_download"],
    outs = ["gosec"],
    binary = True,
    cmd = "$TOOLS x $SRCS -s v2.14.0",
    test_only = True,
    tools = [CONFIG.ARCAT_TOOL],
    visibility = ["//tools/driver/..."],
)

remote_file(
    name = "puku",
    binary = True,
    url = f"https://github.com/please-build/puku/releases/download/v{CONFIG.PUKU_VERSION}/puku-{CONFIG.PUKU_VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
)
