VERSION = "1.26.0"

hashes = {
    "darwin_amd64": "5eff55f42a6a5471feb38fb87ff4f1eeefa2bc222c8b5cfcca41c822423e57ca",
    "darwin_arm64": "a7914480f2a49ebfd6ca713622cde0e9143d53a39ec7e67c502b1d8c7d00db33",
    "freebsd_amd64": "09641cc1191a419f149ad44d8c49b79bbcb523d12fb181c0443d7fb5d6b773a1",
    "linux_amd64": "918ba0cbc644a48eb75645a017ea2a65fbe65494e9adf41162c408888e7ff10f",
    "linux_arm64": "9beefe4d653cedd7014a7e94336c4eccf20fb15770b6290422db3f26841fd535",
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
