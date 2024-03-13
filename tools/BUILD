VERSION = "1.10.6"

hashes = {
    "darwin_amd64": "8c68525f6e25dd0b07c8b3aa55e297505ec6884fa5f4f3bae6804b013d19325c",
    "darwin_arm64": "1edb5c2a189ff34cbeb65f4006858fb9cd20a80394780b573267442b8832750f",
    "freebsd_amd64": "51ba14804fe22456f4892cee66ba69d1ea966a505846c9a627a17f34e14182aa",
    "linux_amd64": "10aa72ccb08b5b4333a492c99a001a1011da3ac6431d474771faf69d9a635eb4",
    "linux_arm64": "5cd9519ceff4ff693ccc8bef36bdf509fbd8f934ad032f091725a9342c177e4f",
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
