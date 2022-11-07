VERSION = "0.1.4"

remote_file(
    name = "please_go",
    binary = True,
    hashes = [
        "e59d445371aba528942b3d9b67ca0ee2b443ea31ac48fff0c91d88037c758806",  # darwin_amd64
        "0e99ad95d8df074763518de069793786752c5a19174396654d452efff9504a62",  # darwin_arm64
        "34c292fd0b88b1387becf7d67639ab400d780f8bb7b8cb914ab0a01a51a00ef3",  # freebsd_amd64
        "2bfb2cedfc92aa0660f48b20b2db6ed98f4f22c21ee0f037d6679721f3244d91",  # linux_amd64
        "34515d6f55976872285f34ab75e8fcc42b51143b67775498b1706b3d68aea29e",  # linux_arm64
    ],
    labels = ["manual"],  # Needed to avoid picking this up for coverage
    url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
    visibility = ["PUBLIC"],
)
