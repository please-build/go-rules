VERSION = "1.23.0"

hashes = {
    "darwin_amd64": "19e514e2525f7eab008b97a87713e3d84f3c2b0991d015ca91fb38f41e9bfc76",
    "darwin_arm64": "78c8e5420f86b3d4491b59c071c55007847b6e841b1821fb203cc243e2d54f3f",
    "freebsd_amd64": "8e80494192ea5b8777002d6d2351555a01a0d0c0c5885e59ad836d733c2b3d8e",
    "linux_amd64": "bad4e86e9b961465d95c17f3f6a2011ee58d6c8ebf3e4dd746589a37ddd3aed7",
    "linux_arm64": "216bf7d2a9adfeffdd9d3eb2fba755a5f22608078e8f94e0ab5dd09bb2dc20dc",
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
