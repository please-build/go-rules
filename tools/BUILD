VERSION = "1.3.2"

remote_file(
    name = "please_go",
    binary = True,
    hashes = [
        "b368e6086b2d397b71417d2ab4ed62672dc18ac3f9d590029daa5cd833135fc6",  # darwin_amd64
        "a0f4b6d187874fc82e60b2311d0efc8eda191c5fe6b25ff662e02100fc3a537f",  # darwin_arm64
        "c14addd0318eeead57bd06f1a795798519b2962bacc92574f6f1b461d02ed8b4",  # freebsd_amd64
        "6d565098752d235110a1967d53c2ca24dea76f24172c9a7ba0abdb6df656aad0",  # linux_amd64
        "6bfa7562790bb4bccef7a0b4c4c3b567a30edade090a74010e6945a1eb1efe6d",  # linux_arm64
    ],
    labels = ["manual"],  # Needed to avoid picking this up for coverage
    url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
    visibility = ["PUBLIC"],
)

filegroup(
    name = "download_test",
    labels = ["manual"],
    srcs = [
        "///darwin_amd64//tools:please_go",
        "///darwin_arm64//tools:please_go",
        "///freebsd_amd64//tools:please_go",
        "///linux_amd64//tools:please_go",
        "///linux_arm64//tools:please_go",
    ],
)