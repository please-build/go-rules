VERSION = "1.1.0"

remote_file(
    name = "please_go",
    binary = True,
    hashes = [
        "86a9160cc7f6d4922234f4eee42a473ff93851d5743dd3cfb656560b53d746f3",  # darwin_amd64
        "8daefbfd90c3dcc8b043963199375c70b343efcc1adedb77e3ee2bda884b8ec5",  # darwin_arm64
        "2874143e089b4139693927142c488f4e80b235f1fcbbf031a3cefcd5774467de",  # freebsd_amd64
        "da7a21e1bde921095b788ee468877548063d60ab25979e0d2a0583a283bf943b",  # linux_amd64
        "199e3342772402d06e820fee75cffff94381145abc2f26937ef5431a2fc12ec1",  # linux_arm64
    ],
    labels = ["manual"],  # Needed to avoid picking this up for coverage
    url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
    visibility = ["PUBLIC"],
)
