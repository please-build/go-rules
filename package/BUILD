subinclude("//tools/please_go:version")

def cross_compile(arch):
    return build_rule(
        name = f"please_go_{arch}",
        srcs = [f"///{arch}//tools/please_go"],
        outs = [f"please_go-{VERSION}-{arch}"],
        cmd = "mv $SRC $OUT",
        labels = ["manual"],
    )

filegroup(
    name = "release_files",
    srcs = [cross_compile(arch) for arch in [
        "darwin_amd64",
        "darwin_arm64",
        "freebsd_amd64",
        "linux_amd64",
        "linux_arm64",
    ]],
    labels = ["hlink:plz-out/package", "manual"],
)