VERSION = "1.4.0"

remote_file(
    name = "please_go",
    binary = True,
    hashes = [
        "52a49477d64b6b66b960b0952a66b312c1002dd1dc1a59b58772c8a279ea2071",  # darwin_amd64
        "bc38888501974890b4559d08b321da31bbf1e3d714ef52c09c538ec608d3d004",  # darwin_arm64
        "d1dddad9917ed9e34855c6cbf7b749bd66109552329e79b0621a077d4f3308e4",  # freebsd_amd64
        "73baa823c91cdfc9db3051a4249e94aff4cbe5e2e3c339e1f03a2801ccc2cddf",  # linux_amd64
        "fa0fd037a4917dcbbcb0dcffcce1c14158a4a51016a210658066be5a02122b8c",  # linux_arm64
    ],
    labels = ["manual"],  # Needed to avoid picking this up for coverage
    url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
    visibility = ["PUBLIC"],
)

filegroup(
    name = "download_test",
    labels = ["manual"],
    srcs = [
        "///darwin_amd64//tools:please_go",
        "///darwin_arm64//tools:please_go",
        "///freebsd_amd64//tools:please_go",
        "///linux_amd64//tools:please_go",
        "///linux_arm64//tools:please_go",
    ],
)