VERSION = "1.4.3"

hashes = {
    "darwin_amd64": "9209ecdcf63d8c1b55a19c751fb5db3c6e3886088776b4c58b9ec0e3bc253154",
    "darwin_arm64": "ff88728286bd528ce6f798e8ec81bd668248365de8fe73955018498e9a0ba2e9",
    "freebsd_amd64": "cb5a9bad2405c6fae2fb64a8a5bc5aee2567b43c79108c77b37558df25cc4950",
    "linux_amd64": "d4e670a9803ceba54f8e899d3d3fddb5d41cc3106bec047b4fb8e85e87fd593d",
    "linux_arm64": "57c89ebe4da7dfb3401e2be0be7e4a21ffc3d9d866cc198a4479f7c7f734be07",
}

for a, h in hashes.items():
    native = f"{CONFIG.OS}_{CONFIG.ARCH}" == a
    remote_file(
        name = "please_go" if native else f"please_go_{a}",
        binary = True,
        hashes = [h],
        labels = None if native else ["manual"],  # Needed to avoid picking this up for other architectures for coverage
        url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{a}",
        visibility = ["PUBLIC"] if native else None,
    )
