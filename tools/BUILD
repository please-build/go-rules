VERSION = "1.24.1"

hashes = {
    "darwin_amd64": "0b00408859b67e48aabdf632bc7debf9fb1016896725bc7fe3e8b527ce89720a",
    "darwin_arm64": "192076228d266b1ec42b56b8d0da50171263b8f684d4e838e7f53c32d8c9a256",
    "freebsd_amd64": "3a3a9e7e80c9f037d7f184065adb0b79d1747fc4d3f0f0d88a3e1aea8de27f37",
    "linux_amd64": "f448e1f23438821e750054240cfc897b1a8baffdc2a53adda10f879d3bdd0dc9",
    "linux_arm64": "9b7d5c4efe577fcf2020f7a500f3508660b313ecf0b7473768a2a131169268a2",
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
