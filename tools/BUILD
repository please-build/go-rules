VERSION = "1.6.0"

hashes = {
    "darwin_amd64": "3c0902fd24aff274d0938e5b98493d6106a64b3c938ac0b5f5f6ced1d8e39470",
    "darwin_arm64": "fb068451cbf76c8cf615a1da3494ea9fef4669f6c53900eb40793b185a72dab4",
    "freebsd_amd64": "a4c0c8463fde9cef05e12d88b713d40c78da48a7f6d10bc9d5dffddd5188e938",
    "linux_amd64": "24dadf9e0f18f9e1451ffdd6bf66b02b262546db2ac12b72f087be3c9ad60efe",
    "linux_arm64": "dea795a77b545419e4875486023b80766e5dffe7d687aba10f361a7241f395dc",
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
