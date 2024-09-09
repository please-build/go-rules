VERSION = "1.15.1"

hashes = {
    "darwin_amd64": "e7fd623e923a5f2f26b6510a100db36a4af74a178aaef0989f9ba9991633c82b",
    "darwin_arm64": "bdf6574d75efd26635539143dc6a2434234f34f621f730e29bc017efb173ebd3",
    "freebsd_amd64": "004e30a5ede1af591b0496ebbb5a8f56b507c4e291fa42bc3b22f6ac58c1c521",
    "linux_amd64": "b2f6a01c7c69d8686b5a4eb83411811d848c8697b2a36f6be431fa2a54555026",
    "linux_arm64": "7e9e6f2f471c9ca1063c5a64043431d079a53f9d02050d781fd112d4141786fc",
}

for a, h in hashes.items():
    native = f"{CONFIG.OS}_{CONFIG.ARCH}" == a
    remote_file(
        name = "please_go" if native else f"please_go_{a}",
        binary = True,
        hashes = [h],
        url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{a}",
        visibility = ["PUBLIC"] if native else None,
    )
