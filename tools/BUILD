VERSION = "1.16.0"

hashes = {
    "darwin_amd64": "0fcb68181111c65b9c786dd20fa6091341ad535c6d7a42c940150c7c9ed51868",
    "darwin_arm64": "2c00787d0f6e18d74eb69139ed6c16df252cca8ca3cc915e59d6ca222b3bff55",
    "freebsd_amd64": "65b8f69d044e1a4ddf9a8efc605e8526a987e24e0f825f16333bdd40060605e3",
    "linux_amd64": "8e3f70d0f988c66d59dd22569bc8e111b82cbfd242e52ef2792e6f698165e947",
    "linux_arm64": "f97f690877b99e7e66c60e74db35993cbe60306fcfbe1201c109bc0c61e0dcd6",
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
