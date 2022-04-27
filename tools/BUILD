VERSION = "0.1.2"

remote_file(
    name = "please_go",
    binary = True,
    hashes = [
        "c6148a0cb89ffeb975a68588035bb3c050ddecae7833b464545a973e20ab469a",  # linux_amd64
        "4a248ae9b5836b5f10a1095a4c8c557c50f8b8becaa038104ecd7051ef778627",  # linux_arm64
        "e2a4cfabbd5c34d5434d0c47a117fdb9880401d8c114bc883743e2d03a60385c",  # darwin_amd64
        "eed909e27e433458a0df94f2d5b2a842ddd2f66c6a668c5d853e82d23eb424c4",  # darwin_arm64
        "e2339f5e96d5c222b7b7e27add933049043d1ee01fe7704e21b14d29e913059d",  # freebsd_amd64
    ],
    labels = ["manual"],  # Needed to avoid picking this up for coverage
    url = f"https://github.com/please-build/go-rules/releases/download/please-go-v{VERSION}/please_go-{VERSION}-{CONFIG.OS}_{CONFIG.ARCH}",
    visibility = ["PUBLIC"],
)
