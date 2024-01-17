VERSION = "1.8.4"

hashes = {
    "darwin_amd64": "d513d09c16615f4f105af24c39f73e0a79ee81e8cc66ff73803562ff3f9e279f",
    "darwin_arm64": "860e913980839ce9a3609b2c5fe38486bf0246f9dfe5ce409abb0ddfd54f0c82",
    "freebsd_amd64": "81734f5b3ffb13a733c68a403ab4d4d26a57bd12c9d8a82f0b07372308dc12d1",
    "linux_amd64": "77e068f88db3e25c4c6bfab521fba12315b9fb6b102e08d9ee1fa134f7e89a59",
    "linux_arm64": "262732b5565886075c518690ab0d0b9500e5eba89cdbc29f5e4bf46242b82900",
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
