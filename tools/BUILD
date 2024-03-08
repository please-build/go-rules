VERSION = "1.10.5"

hashes = {
    "darwin_amd64": "2d6da3bd774d954ab3694498813f9e97d06e61efe94a3c4a1a751ab73f8944a6",
    "darwin_arm64": "0226b1d816b8456ef6ec81123856832494317cf0d324f0d5dbfe9e5b82a2d0e0",
    "freebsd_amd64": "c64320c701f07b813e91e5a0f23ec73081e134c2b4b216e04b43c5faa2bc3e4c",
    "linux_amd64": "456f368cf32f359a36f8f1464c36c87427df5e5048382c9667aea1913424a384",
    "linux_arm64": "b13adebaee8426bccdcc9a2dae1bd4f3e3605eccce5c79c81fc76cea00c0078d",
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
