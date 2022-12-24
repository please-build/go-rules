VERSION = "1.1.0"

remote_file(
    name = "please_go",
    binary = True,
    hashes = [
        "3b430242c3351346e3bb83e6beba88b43bd575f7581677f8fbd127071d310191",  # darwin_amd64
        "3bd07425420ed604f5bf59323f770cdb442a80873bdde90792bfdf0dfa43a500",  # darwin_arm64
        "8785b22b12f8141461f3a7c5930a485576a20969b223208bbe2b50ee7e79ac8f",  # freebsd_amd64
        "13a8d7c857c76928a578c5853f0cfda5ab2f605398f99f3b6f0e36fa0e31b5c5",  # linux_amd64
        "4000f0213fe9081483053bcb95ce01238c1d1f6387b36747fb1b8447019c6ca7",  # linux_arm64
    ],
    labels = ["manual"],  # Needed to avoid picking this up for coverage
    url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
    visibility = ["PUBLIC"],
)
