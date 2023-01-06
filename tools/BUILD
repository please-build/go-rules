VERSION = "1.1.2"

remote_file(
    name = "please_go",
    binary = True,
    hashes = [
        "319c93ed2dfbc023dc8aa9c8b6205dac0f6a83130ed4049d6f01e122034468b4",  # darwin_amd64
        "4b48329fb43cfca3fb86754927a0b7cbbdc96035f722af18872c78d76797b20f",  # darwin_arm64
        "62084453432a8f3f4a3df5cd12c7fb4fd50e0014b2551b7afae4e5f37fddb402",  # freebsd_amd64
        "f227d4497f1129bcc1bb2de276e32e7d92a526eefe001fa2cd49ca13f8b96caa",  # linux_amd64
        "ab705179b2a46cb6fb23dd8381ece0985a3ae0c0fcbf8d4886ccddd05cd94f7d",  # linux_arm64
    ],
    labels = ["manual"],  # Needed to avoid picking this up for coverage
    url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
    visibility = ["PUBLIC"],
)

gentest(
    name = "hashes_test",
    data = [
        "///darwin_amd64//tools:please_go",
        "///darwin_arm64//tools:please_go",
        "///freebsd_amd64//tools:please_go",
        "///linux_amd64//tools:please_go",
        "///linux_arm64//tools:please_go",
    ],
    test_cmd = "exit 0", # If we manage to build the data targets, then the hashes were good.
    no_test_output = True,
)