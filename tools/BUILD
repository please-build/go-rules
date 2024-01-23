VERSION = "1.8.5"

hashes = {
    "darwin_amd64": "c7f3b8fa9abae9d00800308e200e109adf27095cabc3c838907fcb1603d92c4a",
    "darwin_arm64": "ee755f39be27013dad2040b5cad3db7d98ced21418105bd88bf58548790d0951",
    "freebsd_amd64": "2cfd3fca2cd0c83e7d3bd77591ba2a264bfa95aed1fbe2d4cac40924333613e2",
    "linux_amd64": "2709d3fb06a0c10919e2d724116a42e871f26df68d8b699bb91a874997f763c0,",
    "linux_arm64": "c43cca962cc123b3753aafc97577cdc9456461e6e4cb9a68991c08647fe30265",
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
