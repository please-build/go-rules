VERSION = "1.21.0"

hashes = {
    "darwin_amd64": "ed105de84226521ee70756558f028c63a380da433eede6e01bb2abe76fb63b90",
    "darwin_arm64": "2553c598cc7eb9915072d080da17f9a8e5b398ef8c548c29ec46203d7da96b83",
    "freebsd_amd64": "e46e349363fccd0865e7744f19c2de684d5ee21f4e7611bdf802c599c09bcdea",
    "linux_amd64": "1e2407960085e4735523e3d27023a1483037ac62dfd4e7473c6fd4f5a58f2471",
    "linux_arm64": "e2aea3747fe88b386d0bc9244804665bf03ea491f30a36ddeff116a317d2dc48",
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
