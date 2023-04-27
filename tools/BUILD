VERSION = "1.3.0"

remote_file(
    name = "please_go",
    binary = True,
    hashes = [
        "2f6744d095273ad977260d64260db92b316c7acaf12351df90bb7fbdf974ea77",  # darwin_amd64
        "a44fb6c1ac4734ee0af1f218576e751e8811ec73e6422d254af80ab1f73b0807",  # darwin_arm64
        "465e4da37b9157d4da8bc616f8d53ecc4abfd3a017e7cab3cd5237dbc35a561b",  # freebsd_amd64
        "f6c54b2c5d316cd8ac22a40556ef250e08cba42c7fa4dd1dbc4e1284c85f07e4",  # linux_amd64
        "adeb01a1335da76281879c9707a084fb260ad11bea9b6c6aead46ea9c4607161",  # linux_arm64
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