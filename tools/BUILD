VERSION = "1.9.0"

hashes = {
    "darwin_amd64": "2e12719762f30c71c579f65b97b4493ef4ac33b4b2e3cb74a936858a3519ac56",
    "darwin_arm64": "502ef05ad52c061145eddbc4f3b0ae376c1bfc031078450b9e9db960db7a074a",
    "freebsd_amd64": "e96102e0b7deab1e34ccce9a91c3999ae8a82797f4bfff25537bedccbe76c3f1",
    "linux_amd64": "5ab64b67a2868cd645730a1f99eb487e7cfbe05e5efe516ffbc5a147184c415e",
    "linux_arm64": "eedc9ad8048995d2f056143508c60bebc4bf74873937342047c07a12da7e49b6",
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
