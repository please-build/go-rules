VERSION = "0.1.1"

remote_file(
    name = "please_go",
    url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
    visibility = ["PUBLIC"],
    binary = True,
    labels = ["manual"], # Needed to avoid picking this up for coverage
    hashes = [
        "83fded9411fd56c739c3e19db36612e2893d94231c5869b291de3f1937f4bb43", # linux_amd64
        "18afcada1bc0833391614c0fdad2efc61e0a4176f1747c43de0a0baf82f6712a", # linux_arm64
        "f33254b6f981af7b748ff5aefbb4e9b49c7c11ab1ed71bb8957e5fd486f1ddd8", # darwin_amd64
        "a57fbbb8532b8d1fe63bf74fbf7745224b0bed2aec9633114d7cee8593df242e", # darwin_arm64
        "48bd61402a5a11445c791d3520467770ddfc495fd4c2deb491580b664bd42200", # freebsd_amd64
    ],
)
