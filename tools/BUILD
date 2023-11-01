VERSION = "1.6.1"

hashes = {
    "darwin_amd64": "5447e997a92cb9374238f55b5ee6541577df578926818c232574724fe1730f06",
    "darwin_arm64": "64b328a11347e581122c9290d85e4426ca9d8185b28deb8ead49c55f90697511",
    "freebsd_amd64": "364671c696692de98d4cd22a1cc165dd2c1dd3167c51c5fef8fcec0782b249ab",
    "linux_amd64": "8f7e1b3b2ec7afbdc056b2373bb0578889ebd6b354096a370527ff3decaff0a8",
    "linux_arm64": "5c2213d4638d97d37cd68a4f7d2998286fcbc533203d54410b6766ac8b046bfa",
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
