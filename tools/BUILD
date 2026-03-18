VERSION = "1.22.0"

hashes = {
    "darwin_amd64": "ed105de84226521ee70756558f028c63a380da433eede6e01bb2abe76fb63b90",
    "darwin_arm64": "8e45c9f56cd5d9f8a43e6ed9f5735695aa2745e6822be6eabf83ab7b53eeb4a8",
    "freebsd_amd64": "d8f5b5b8d1f9a0afd6cbb1db5ac9b11b3c1bb7cf4692f1309909542f2d984cf4",
    "linux_amd64": "97edb846709d03c63d0dcc4562a5b6bd97d6f1cfc075d69d816f47b582d751f2",
    "linux_arm64": "3bd706d4332e455e480eaa4ef933094056793e178a74ec1fb9245fe30bb7eaac",
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
