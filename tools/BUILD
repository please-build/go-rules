VERSION = "1.7.0"

hashes = {
    "darwin_amd64": "e2c3c8a113f086a79c9a1f0aecb127ebec508275b8e516f9edbf64ad0d7a34b8",
    "darwin_arm64": "d495f652afc118503c00ac9ea537d8b81082ecbc90ba7f8e001a4c0096ee90e8",
    "freebsd_amd64": "101cabd05f8c3e13563cb239e529294b6185f017fca960616faf261f66371e0f,",
    "linux_amd64": "f5e9c53d718c0ec6cc355644291442577e33239602946091727f0690f3bf2efe",
    "linux_arm64": "442fb650c7caa6c747649060f76a420d2e4aa3f7034ea651b5340e4e29b538a8",
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
