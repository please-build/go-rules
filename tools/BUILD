VERSION = "1.9.1"

hashes = {
    "darwin_amd64": "27fe56ac049f62010d11cacff459e422f2f7c4934f239ed83fd56f1e56aeaacb",
    "darwin_arm64": "06fb2913b681cb47c033371a98762f059083c43708d200d0849a6f42cf7761dc",
    "freebsd_amd64": "0445f79d3d74f4d3ac636c4656ac84551418902b8d23097a375bcfbb250c5862",
    "linux_amd64": "ae0033ce76915151d2dd5a4a7f162c8e0cdbe3d8a053e1948852975205b2f67a",
    "linux_arm64": "b1a317782e59f90330f9c80de23cf7dc4664c63d93a91b9d94ca614161ace494",
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
