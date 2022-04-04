subinclude("//:version")

remote_file(
    name = "please_go",
    url = f"https://github.com/please-build/go-rules/releases/download/v{VERSION}/please_go-{VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
    visibility = ["PUBLIC"],
    labels = ["manual"], # Needed to avoid picking this up for coverage
)
