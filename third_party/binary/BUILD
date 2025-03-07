remote_file(
    name = "gosec_download",
    hashes = ["38e596e1ad580666e5efaa159fbaf63bd5873ea1e74beacbdf785b025b52e1db"],
    test_only = True,
    url = ["https://github.com/securego/gosec/releases/download/v2.22.0/gosec_2.22.0_linux_amd64.tar.gz"],
)

genrule(
    name = "gosec",
    srcs = [":gosec_download"],
    outs = ["gosec"],
    binary = True,
    cmd = "$TOOLS x $SRCS -s v2.22.0",
    test_only = True,
    tools = [CONFIG.ARCAT_TOOL],
    visibility = ["//tools/driver/..."],
)

remote_file(
    name = "puku",
    binary = True,
    url = f"https://github.com/please-build/puku/releases/download/v{CONFIG.PUKU_VERSION}/puku-{CONFIG.PUKU_VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
)
