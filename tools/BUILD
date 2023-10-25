VERSION = "1.5.0"

hashes = {
    "darwin_amd64": "dc0893949401375221065aefeb7f7aa83dcef6d5f67942274bc0f8563bda0566",
    "darwin_arm64": "eb3afaca176856eb36770c9da15c9c6aa0e513f96063f7a6d3ac91d127f876a6",
    "freebsd_amd64": "21a6a31f214480801055a581c0b345e5a4ef250b055492f5c6c16fbe5c342555",
    "linux_amd64": "7deadd1d61efcc8931f907cb1953d9c61d89c2d21d8212691a51199061531e60",
    "linux_arm64": "02e8a6d7634af74b72879c1f36680370bc76963379727428d1e55a70cc8e8d44",
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
