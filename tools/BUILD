VERSION = "1.10.2"

hashes = {
    "darwin_amd64": "083a425eae3104dd2fb2cef02d8c118df8828418568fe8c2844ca050aacb0066",
    "darwin_arm64": "ea8d443e5ee94ff1057196920d7b5bbc3e8b044b19fd3c4474247bf48ad3be00",
    "freebsd_amd64": "76b2422cad8e32f7552309f548ad99c9e9fe6b2f29964d5e5f49afb750e1d1c4",
    "linux_amd64": "bfcf5ff496ffe8633657b7ed5a37dc2d47526db398bd1d04394caea8170f85ad",
    "linux_arm64": "7992465c31e93cf7ac3299aeea983b4865a92dfa5a957feec2ba226f9d01d83e",
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
