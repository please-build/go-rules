VERSION = "1.24.0"

hashes = {
    "darwin_amd64": "10ba63f78cd622601d83f53f4603a4e5d946ac2f331eba3af2a1627dc59ca908",
    "darwin_arm64": "3b886e6b00c4e6cdaa633752439860d771f6cc93720fdaf899ef01d5991e7473",
    "freebsd_amd64": "becb883a13c7f00eca71219a8539c5cc80eead8e0104e07c2d02ec6c6900ec7d",
    "linux_amd64": "92fbc45047eb825c0bded235447b9c455716f270ce174d107798f5770a4201ab",
    "linux_arm64": "bec438aa687021ad97ab1d626b538644b0a8827784e6458d24939e726f3db9c5",
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
