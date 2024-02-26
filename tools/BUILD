VERSION = "1.10.4"

hashes = {
    "darwin_amd64": "080a0b71a863411758fa2e54fae1940eaf97ce803ddea4b4839c6471079bfda1",
    "darwin_arm64": "56d77d3d1708c2178d2d3b775bf0b801f0ffe2fa8842245d7de2d7d090834b7a",
    "freebsd_amd64": "7c301c9a6bbf12e5bc16cccd76ca65aa8bd0747f148cf385b830788e684e4c91",
    "linux_amd64": "c7f4467f10bf7cd1a7e92e808a0fbb8492a4d770dbc030300610498658f517b4",
    "linux_arm64": "bb3556d3f12bdebd874d71b12c45ad99f34b372d465d54604cd682ede485ae55",
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
