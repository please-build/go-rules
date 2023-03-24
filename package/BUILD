subinclude("//tools/please_go:version", "//tools/driver:version")

tools = {
    "please_go": PLEASE_GO_VERSION,
    "driver": DRIVER_VERSION,
}

target_names = {
    "driver": "plz-gopackagesdriver",
}

architectures = [
    "darwin_amd64",
    "darwin_arm64",
    "freebsd_amd64",
    "linux_amd64",
    "linux_arm64",
]

def cross_compile(tool, version, arch):
    target = target_names.get(tool, tool)
    return build_rule(
        name = f"{tool}_{arch}",
        srcs = [f"///{arch}//tools/{tool}:{target}"],
        outs = [f"{tool}-{version}-{arch}"],
        cmd = "mv $SRC $OUT",
        labels = ["manual"],
    )

for tool, version in tools.items():
    filegroup(
        name = f"{tool}_release_files",
        srcs = [cross_compile(tool, version, arch) for arch in architectures],
        labels = [
            f"hlink:plz-out/package/{tool}",
            "manual",
        ],
    )
