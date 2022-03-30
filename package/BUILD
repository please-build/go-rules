subinclude("//:version")

def cross_compile(arch):
    return build_rule(
        name = f"arcat_{arch}",
        srcs = [f"///{arch}//tools/please_go"],
        outs = [f"arcat-{VERSION}-{arch}"],
        cmd = "mv $SRC $OUT",
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
    labels = ["hlink:plz-out/package"],
)