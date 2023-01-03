subinclude("//build_defs:go")

go_toolchain(
    name = "toolchain",
    architectures = [
        "darwin_amd64",
        "darwin_arm64",
        "freebsd_amd64",
        "linux_amd64",
        "linux_arm64",
        # Installing takes a while so only do it when we need it for release
        # This extra check can be removed once we fix booleans in build config
    ] if CONFIG.RELEASE and CONFIG.RELEASE != "False" else [],
    tags = [
        "netgo",
        "osusergo",
    ],
    version = "1.18",
)

go_module(
    name = "testify",
    install = [
        "assert",
        "require",
    ],
    licences = ["MIT"],
    module = "github.com/stretchr/testify",
    version = "v1.7.0",
    visibility = ["PUBLIC"],
    deps = [
        ":difflib",
        ":spew",
        ":yaml",
    ],
)

go_module(
    name = "difflib",
    install = ["difflib"],
    licences = ["BSD-3-Clause"],
    module = "github.com/pmezard/go-difflib",
    version = "v1.0.0",
    visibility = ["PUBLIC"],
)

go_module(
    name = "spew",
    install = ["spew"],
    licences = ["ISC"],
    module = "github.com/davecgh/go-spew",
    version = "v1.1.1",
    visibility = ["PUBLIC"],
)

go_module(
    name = "yaml",
    licences = ["MIT"],
    module = "gopkg.in/yaml.v3",
    version = "v3.0.0-20210107192922-496545a6307b",
    visibility = ["PUBLIC"],
)

go_module(
    name = "go-flags",
    licences = ["BSD-3-Clause"],
    module = "github.com/thought-machine/go-flags",
    version = "v1.6.0",
    deps = [":xsys"],
)

go_module(
    name = "xsys",
    install = ["..."],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/sys",
    version = "039c03cc5b867cd7b06a19ff375be5c945c80b10",
)

go_module(
    name = "humanize",
    licences = ["MIT"],
    module = "github.com/dustin/go-humanize",
    version = "v1.0.0",
)

go_module(
    name = "flags",
    install = [
        "flags",
        "logging",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/peterebden/go-cli-init/v5",
    version = "v5.1.0",
    visibility = ["PUBLIC"],
    deps = [
        ":go-flags",
        ":humanize",
        ":xcrypto",
        ":logging",
        ":xsys",
    ],
)

go_module(
    name = "xxhash",
    licences = ["MIT"],
    module = "github.com/cespare/xxhash/v2",
    version = "v2.1.2",
    visibility = ["PUBLIC"],
)

go_module(
    name = "buildtools",
    install = [
        "api_proto",
        "build",
        "build_proto",
        "edit",
        "file",
        "labels",
        "lang",
        "tables",
        "wspace",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/bazelbuild/buildtools",
    version = "v0.0.0-20210920153738-d6daef01a1a2",
    visibility = ["PUBLIC"],
    deps = [
        ":protobuf",
        ":google-protobuf",
    ],
)

go_module(
    name = "google-protobuf",
    install = [
        "encoding/prototext",
        "encoding/protowire",
        "internal/...",
        "proto",
        "reflect/protodesc",
        "reflect/protoreflect",
        "reflect/protoregistry",
        "runtime/protoiface",
        "runtime/protoimpl",
        "types/descriptorpb",
    ],
    licences = ["BSD-3-Clause"],
    module = "google.golang.org/protobuf",
    version = "v1.27.1",
    visibility = ["PUBLIC"],
)

go_module(
    name = "protobuf",
    install = ["proto"],
    licences = ["BSD-3-Clause"],
    module = "github.com/golang/protobuf",
    version = "v1.5.2",
    visibility = ["PUBLIC"],
    deps = [":google-protobuf"],
)

go_module(
    name = "mod",
    install = [
        "modfile",
        "internal/lazyregexp",
        "module",
        "semver",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/mod",
    version = "v0.5.0",
    visibility = ["PUBLIC"],
    deps = [
        ":xerrors",
    ],
)

go_module(
    name = "xerrors",
    install = [
        ".",
        "internal",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/xerrors",
    version = "v0.0.0-20200804184101-5ec99f83aff1",
    visibility = ["PUBLIC"],
)

go_module(
    name = "xtools",
    install = [
        "go/packages",
        "go/gcexportdata",
        "go/internal/packagesdriver",
        "internal/event/...",
        "internal/pkgbits",
        "internal/gcimporter",
        "internal/gocommand",
        "internal/packagesinternal",
        "internal/typeparams",
        "internal/typesinternal",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/tools",
    version = "v0.4.0",
    deps = [
        ":mod",
        ":xsys",
    ],
    visibility = ["PUBLIC"],
)

go_module(
    name = "logging",
    licences = ["BSD-3-Clause"],
    module = "gopkg.in/op/go-logging.v1",
    version = "v1.0.0-20160211212156-b2cb9fa56473",
    deps = [":xcrypto"],
)

go_module(
    name = "xcrypto",
    install = [
        "ssh/terminal",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/crypto",
    version = "v0.0.0-20210920023735-84f357641f63",
    deps = [":xterm"],
)

go_module(
    name = "xterm",
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/term",
    version = "v0.0.0-20210615171337-6886f2dfbf5b",
    deps = [":xsys"],
    visibility = ["PUBLIC"],
)

# This is not really necessary, it's here to test the case of a separate go_mod_download
# for the gopackagesdriver.
go_mod_download(
    name = "xsync_download",
    module = "golang.org/x/sync",
    version = "v0.1.0",
    licences = ["BSD-3-Clause"],
)

go_module(
    name = "xsync",
    install = ["errgroup"],
    module = "golang.org/x/sync",
    download = ":xsync_download",
    visibility = ["PUBLIC"],
)
