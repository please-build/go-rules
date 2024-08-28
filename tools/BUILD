VERSION = "1.15.0"

hashes = {
    "darwin_amd64": "483a4519b746c6b36f9f748b2aa9bb59595bd8e02166e1802e07de0a2ca3177d",
    "darwin_arm64": "a721600ab3b2f1e0639931ec7b33bc1a8faed257ed7b69b86416f4ea8c35c700",
    "freebsd_amd64": "dd1e45fd3095c00c7dedebac9f9f7219f748c25feb7be81b1df68278d2f305aa",
    "linux_amd64": "88bef7b855445e9f1fa0ef23a4ae3f5003eddf5131b02ca80c6b5fd5908876da",
    "linux_arm64": "0f23866a1707fb46bf4956c169f68cb913b1ccccb1618309ca6b7e9dae1cc571",
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
