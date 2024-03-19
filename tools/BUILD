VERSION = "1.11.0"

hashes = {
    "darwin_amd64": "ce2e58f0aaaff7b5694b84ffb15d767904f71116df35ad01f284a63b3d1ef0c5",
    "darwin_arm64": "6d485c9c1c4d0040ccc03d973419f0e728db187cd70ea3ff1482e8d6a54749f3",
    "freebsd_amd64": "e79aabbb5de73b4c1ee5b8db02079e829aa53e9b7ef8b7a857bd39647c2a8f99",
    "linux_amd64": "191a20991043e9700ac4121f3c92f7bf17f059750659b6a8d97921bd0f58a9e7",
    "linux_arm64": "49384f2ca4c1082b7c27fc89be2635c14ffab1010264cd24663efe986a4ab374",
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
