VERSION = "1.6.2"

hashes = {
    "darwin_amd64": "f0693b790af48da82f4a96107ce5b60ba1e1bf22bab959cf9ff1a71de7e53609",
    "darwin_arm64": "eeccc50be855eb9b7bff851f57f1611f3d003833ebae93abcc8b9f4a964f4d6e",
    "freebsd_amd64": "071f11d5913df02db5922baf0bd2f929b531514cbb86ff794f211db2b191505e",
    "linux_amd64": "47c8065f789cb677f885545ebf273763b0e3185a5e75be1f9b6541b6d3e52d94",
    "linux_arm64": "f17403d926a70afe80326e34305c090efea5dccba5162839ff67641805e8aec6",
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
