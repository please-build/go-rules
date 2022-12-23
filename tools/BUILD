VERSION = "1.0.0"

remote_file(
    name = "please_go",
    binary = True,
    hashes = [
        "bbede3a0a606d07cb7a894ed0e1bbf4bc85a2ae1400e14899ae559115689f0b2",  # darwin_amd64
        "d9f3ac04eccae6006b0daca42a167d13e1ef8f1fa3dbe7fa9ab879f22d634211",  # darwin_arm64
        "22486b809c5c9a0f399c064ac7f37eac69d85f5b96f5498a67ca2c9e6b11733b",  # freebsd_amd64
        "934fd993e23c0f07a9179b76f27a4e4d861b8477d230cb18e4b980a33d52d53c",  # linux_amd64
        "5fe68fc73b49f7d8c25f4f4e83b4f2abd538106e29551290623ed920191fff6a",  # linux_arm64
    ],
    labels = ["manual"],  # Needed to avoid picking this up for coverage
    url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
    visibility = ["PUBLIC"],
)
