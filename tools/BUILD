VERSION = "1.10.3"

hashes = {
    "darwin_amd64": "c40189fde3b628528e46c6793d0f685774e9686d8942cc327d0234364a10b89e",
    "darwin_arm64": "3100756285218aa5effec7b8fe5753c1ad085ac8a1393efbd16f9c1990fc89b3",
    "freebsd_amd64": "c3302d480aa9b42272334357666983648bcc0b2e1a2cacb15b1cef2c97eee6b7",
    "linux_amd64": "09f15195b373e3fc7b66055e3b81523e38c48d6b417009d24bf7db2940a04361",
    "linux_arm64": "c5bd21a22a9efe6bab6d5f1a9083ee9e3beb8a779157f03771d868d70631e6de",
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
