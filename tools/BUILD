VERSION = "1.8.0"

hashes = {
    "darwin_amd64": "d45b876a50daa784a5921e4c46e3a6b18dc276cdd1cbca8c8fda9a4508ef26e1",
    "darwin_arm64": "497f358983f1515c577b4e8de07a1bef123bd61e64daef528c12a672bfba335e",
    "freebsd_amd64": "10c25f8029494fcf200a8fce0b39d0fc2195a467604582045fc024f9bd97d306",
    "linux_amd64": "cdf31ab729e61b3f95aa1a350437935c8bf0bef3420067424d417e26b5ea9019",
    "linux_arm64": "d8faabfa8a03d21c12678fc947a277455ae9c17bd0fd00e48937b6db2c5392ec",
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
