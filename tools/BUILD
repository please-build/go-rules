VERSION = "1.25.0"

hashes = {
    "darwin_amd64": "d18b242735aef9fc48a4d1e41ef14d57c2de0a185661954851b0cf09a4beac15",
    "darwin_arm64": "33a29068894bc10e245281faa55957c785e20844bf0362f4eebdba3f375ee464",
    "freebsd_amd64": "ea35c5d1d605ba680ba82d467e718be9d8d6b5ce1fa2a144a3824b2d2bd323cc",
    "linux_amd64": "86b40c9450e5ee15ab08e7aa3ec8e73068b7626f787088f91e96c5e35c0b7cc7",
    "linux_arm64": "69e6ae65cdaadb3bcc1e413ef38b147b97cdb9542e86f0698958394487142eb3",
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
