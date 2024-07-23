VERSION = "1.12.0"

hashes = {
    "darwin_amd64": "01dbdbdea6f6f5827537f0ed0b1be8f705b6200691a01738bd80ab083567cf8e",
    "darwin_arm64": "6d485c9c1c4d0040ccc03d973419f0e728db187cd70ea3ff1482e8d6a54749f3",
    "freebsd_amd64": "ab3784f2ac9fadb051d38c8098514e5fa5c6512c30275453edc92ab7c6f2b864",
    "linux_amd64": "76b2764433ed4b90d496b3e0bda79ece24776e9d4fb66911c4acd4186f20e497",
    "linux_arm64": "01faa09dbd59f1777727852a501e8af644c786d50dbb3be22f467a97f7d0fe3e",
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
