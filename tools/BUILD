VERSION = "1.19.0"

hashes = {
    "darwin_amd64": "5545418c33bcffacf3fc17d3ae533404b5892d24d9aff1658b3746d620f0046b",
    "darwin_arm64": "22880ff53c278f071a11b49c7394a5ceada03851d23fd0cce57116daccc0447f",
    "freebsd_amd64": "c4293cc3c48fc6ce3cc4ee9db7ab4229476d63cee8f404ea83c6b79677abd33c",
    "linux_amd64": "804af54d8dfd5eace871e7b91e0503311a896e3718f8a765c2c053ae05a345f6",
    "linux_arm64": "554d7a2afbda8f887b31f06d52610b6cc8ff0c1d487fa9ab706cee0c92af3978",
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
