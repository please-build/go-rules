VERSION = "1.4.1"

remote_file(
    name = "please_go",
    binary = True,
    hashes = [
        "39ba7a7478c16306c7758e81aa53d7ff2908894a8fbc31b1809c2ede7950911e",  # please_go-1.4.1-darwin_amd64
        "342a23306fca55ec7c3c56ccc96af8e61131897c5ba9c245d2d6c4bc293ee1fd",  # please_go-1.4.1-darwin_arm64
        "5273e338aa07f2f2b1a68a08a14acebee52f8a9758f209b0c9deb7cc0a3b90f4",  # please_go-1.4.1-freebsd_amd64
        "beb961aec753e4e9b819e8ab259c30c019229903e27f51b5095322b8c9c0ada7",  # please_go-1.4.1-linux_amd64
        "1644712fa35b076b7d4b65e4e250fe9a2c01516561841aa362568f3a26f4f0d2",  # please_go-1.4.1-linux_arm64
    ],
    labels = ["manual"],  # Needed to avoid picking this up for coverage
    url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
    visibility = ["PUBLIC"],
)

filegroup(
    name = "download_test",
    srcs = [
        "///darwin_amd64//tools:please_go",
        "///darwin_arm64//tools:please_go",
        "///freebsd_amd64//tools:please_go",
        "///linux_amd64//tools:please_go",
        "///linux_arm64//tools:please_go",
    ],
    labels = ["manual"],
)
