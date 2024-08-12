VERSION = "1.13.0"

hashes = {
    "darwin_amd64": "60ea0b01e1568a5320b408d0f1c0cf1b49c1bb617fb47bf508152061f2860eda",
    "darwin_arm64": "b345d04bd4fe504f1c1dd88d9efa39a252fe6c60b2f25538436076c2a61c9411",
    "freebsd_amd64": "37d9a5f937ad765302977f72fdabd0496f9001f1357a661d2bcb14262397939a",
    "linux_amd64": "282ea88e8d7af0786facd6143894eda857bae1fc6bffa8a038f392d8cacf0926",
    "linux_arm64": "e55ee1f08190070f5b730dabc52edbf04ab1cd17f00ee5ee566f052191a3e75c",
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
