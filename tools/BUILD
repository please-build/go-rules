VERSION = "1.20.0"

hashes = {
    "darwin_amd64": "c78ad796230f92547b66ad6fb2fb1ffd7ed7d6a98f5693b448c0e1ca49d26b18",
    "darwin_arm64": "7a049e7411b569eaef389c2a9e1bc649552de57ca867c33adddb1bea28c78f79",
    "freebsd_amd64": "d2800b1fcf49781d2f4197d9c7842f5b9210b50946c9dce1f14eeb4cde501b96",
    "linux_amd64": "69d6e5e935971f2e1c97b8ee11a6c9f0e501f23a24e1ed517841d3c08511cbfa",
    "linux_arm64": "65be124a799af1e70b4ee89f6f33c986a3a3673505ef081d77a49452bf26ff0a",
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
