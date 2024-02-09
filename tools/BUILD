VERSION = "1.10.0"

hashes = {
    "darwin_amd64": "2f8c6ce9d9141de444e598475a2ea574d341fab952c8eecf3a840522bc9fb10a",
    "darwin_arm64": "22d77f14edf9800e4a07221768e717ff8119bbef1aba876b1d94fe79d1d50a2a",
    "freebsd_amd64": "26858675da9cdab99685bf49b1963c6dedeadf004f03369e974e8fcfe099b9fa",
    "linux_amd64": "54b5c41156eec49edfa66e3a96e87f2d637322972fb78ca40dc3c741d883444f",
    "linux_arm64": "4811ff93e9e09c3f8f30b5ec9fcbd3ffada1ebb9bdcce4f4fe01f995c2014036",
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
