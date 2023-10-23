VERSION = "1.4.4"

hashes = {
    "darwin_amd64": "bb36e910341d0d77a90d6468d46db3ae6af6946dd9f3279fe8dec305d3404516",
    "darwin_arm64": "0da19be6e67d1e997f0a1a84f2dfc065993b142c5ef7d4a30cacc4c6cb411b73",
    "freebsd_amd64": "9f63aa4db858c2818b8d6bc9ab76bd643833185aa4c6be1162fb5090be07852e",
    "linux_amd64": "504c4737b48e4c21de546e1f01d21124fc1a73a2b05c69d93503830da00f9c7a",
    "linux_arm64": "f319a8803c7b9d08ecbe955fc7e671c3506319145f5d3ddd96fa8671edd91ebb",
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
