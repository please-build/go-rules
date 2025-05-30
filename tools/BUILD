VERSION = "1.18.0"

hashes = {
    "darwin_amd64": "d03eee83000ed6cece1f061eef2a25e7b3e0906508b6ad372dc55602459cc774",
    "darwin_arm64": "4fdfb922546589081565b7e603a75a492b5e639ff7a82d58fa48acb1226a33e2",
    "freebsd_amd64": "b07476b8bf7412f4b3d0ba276553bc9abf73f71a66ed225860698aa839f44050",
    "linux_amd64": "a5b8293fbb05c7f08191ca7dd026a0e711955e9fb2b526e220556f390c0f5e03",
    "linux_arm64": "b8f39c0c1fdd5f551b020786be5e8030bdb7b09f1e327585b890767ce1bf1135",
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
