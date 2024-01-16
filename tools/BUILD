VERSION = "1.8.2"

hashes = {
    "darwin_amd64": "90780ff069f112b656eb0d5b8318068243b5e96e5c64fdafb9cf3ebc66881e17",
    "darwin_arm64": "ecf810b8ccfa8028476aaf2549bd5b4aa8fb3db7622125686f17dd5e72ab0cac",
    "freebsd_amd64": "e3b31e118e73ce0934c1cbb944b53f3fedeef7533ae9e71a3508eefd3d288d06",
    "linux_amd64": "cdf31ab729e61b3f95aa1a350437935c8bf0bef3420067424d417e26b5ea9019",
    "linux_arm64": "12b5317b32d862528b4a2dea65f475acbe9f5269cb46595c70b01a1f4aa46814",
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
