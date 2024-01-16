VERSION = "1.8.3"

hashes = {
    "darwin_amd64": "3a15575069f2e45d8e5e781ece6d329f896843a1c32c7100fe458e672676a18b",
    "darwin_arm64": "59dbf5836d4baf2c83d45492adc285b4174fb136bafd7f532d0141e0352da522",
    "freebsd_amd64": "0bbdee4fe8fc08a2c4374be1a3c324035789310ac2c7d48788415c98eede04ee",
    "linux_amd64": "ec9ff49bdac27e2f0a8fa014b7a219d746d3239e379a5634c6485ec096be38ee",
    "linux_arm64": "17ff1b0b4a68fedcfb3a7d12410e1c0e05274f782530538622bcd3a5db44965c",
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
