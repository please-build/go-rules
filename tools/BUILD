VERSION = "1.4.2"

hashes = {
    "darwin_amd64": "0df80f71e6a8e8c289c601bdb1937e39960afebbfaf2db5f2dfe32bd3da6beff",
    "darwin_arm64": "96722929660b865022caf2cf963ea57d98321bf1d07a9ea1cb4daf901b45fedc",
    "freebsd_amd64": "fbd7305a108b2c1abf2eee13189cd4aca23fb97578c909fa7d8614dbc974d6b4",
    "linux_amd64": "2df16b7a15c37ea0bffcb82c85ea49c0336138d63d88c67722966b3646fbd5b3",
    "linux_arm64": "4afd983b8d9ac848496c5d9ca18c5774435eeeef743155fe0804b3b2673ded96",
}

for a, h in hashes.items():
    native = f"{CONFIG.OS}_{CONFIG.ARCH}" == a
    remote_file(
        name = "please_go" if native else f"please_go_{a}",
        binary = True,
        hashes = [h],
        labels = None if native else ["manual"],  # Needed to avoid picking this up for other architectures for coverage
        url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{a}",
        visibility = ["PUBLIC"] if native else None,
    )
