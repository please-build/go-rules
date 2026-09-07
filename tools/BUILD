VERSION = "1.26.1"

hashes = {
    "darwin_amd64": "c3a934cff6c3f0dfc2a142db2ac5a06381c8a6d6ec022b3dc9d2c81065c81db3",
    "darwin_arm64": "0f62beb68377a60c6378b1cef82dd4db5ded0cd393d60b1ac3adcf7fc8409476",
    "freebsd_amd64": "87638ae1b81fe6b443a0484d2996e71455b2bd0b33911f4dabad89c337db71b8",
    "linux_amd64": "f21a4ba65e17e8fee7a4064bb26b1a0c455674f93243f84f8cbc92df7ae95e96",
    "linux_arm64": "dccae11d79403e5966801d7cb68e01de865544e293dc939dfea6a103492fef36",
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
