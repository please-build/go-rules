VERSION = "0.1.3"

remote_file(
    name = "please_go",
    binary = True,
    hashes = [
        "89514d72757de68b2e3b5a55140b7d899a241fd64c7536d01c0c8116c463fd54",  # darwin_amd64
        "55866317074e633dbb8ed6e6ede82b89e421ad34ac0c9bb1831afbbaad1c44fe",  # darwin_arm64
        "da708d56473d0eb9ce40b64bce286a80255e103fc86533a318d8fd7fefc0474d",  # freebsd_amd64
        "b9b981106b6950acc1c51bc11fb25f91ee51f5a0c5c2da6248b81e2d97c946ce",  # linux_amd64
        "eb55e9d2e16009f00ac61a1fe01c4347412bc6df8dde760d7fdd1f6a2b5b23cc",  # linux_arm64
    ],
    labels = ["manual"],  # Needed to avoid picking this up for coverage
    url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
    visibility = ["PUBLIC"],
)
