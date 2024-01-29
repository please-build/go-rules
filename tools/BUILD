VERSION = "1.9.2"

hashes = {
    "darwin_amd64": "a7f14eeda9f8575755a71bc92d5faa6eda35f8e3eae2cca088e6492f67067a5c",
    "darwin_arm64": "658e1bf7c0ced39527e76a5f36ba8934548f055b039745238dd99059535b7037",
    "freebsd_amd64": "4a44a894551b0be6f568369ab28c4eb28106a3dcff8d90a2408f5d41a558b8a3",
    "linux_amd64": "e78d76a8b67736b83f1e4d436608fc8d6ad11af87aa1b9ff62464f71d941ebd5",
    "linux_arm64": "5c5282ed3d26a7631a7cb1e7f006200911480794b7a71cb48241dac470353749",
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
