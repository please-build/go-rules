subinclude("//build_defs:go")

go_toolchain(
    name = "toolchain",
    install_std = False,
    version = "1.22.0",
)

go_stdlib(
    name = "std",
    tags = [
        "netgo",
        "osusergo",
    ],
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
    version = "v0.13.0",
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
        ":logging",
        ":xcrypto",
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
        "go/gcexportdata",
        "go/internal/packagesdriver",
        "go/packages",
        "internal/event/...",
        "internal/gcimporter",
        "internal/gocommand",
        "internal/packagesinternal",
        "internal/pkgbits",
        "internal/typeparams",
        "internal/typesinternal",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/tools",
    version = "v0.4.0",
    visibility = ["PUBLIC"],
    deps = [
        ":mod",
        ":xsys",
    ],
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
    visibility = ["PUBLIC"],
    deps = [":xsys"],
)

# This is not really necessary, it's here to test the case of a separate go_mod_download
# for the gopackagesdriver.
go_mod_download(
    name = "xsync_download",
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/sync",
    version = "v0.1.0",
)

go_module(
    name = "xsync",
    download = ":xsync_download",
    install = ["errgroup"],
    module = "golang.org/x/sync",
    visibility = ["PUBLIC"],
)

go_mod_download(
    name = "build_tools_dl",
    licences = ["Apache-2.0"],
    module = "github.com/please-build/buildtools",
    version = "v0.0.0-20230725091334-24cce64e067d",
)

go_module(
    name = "buildtools",
    download = ":build_tools_dl",
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
    module = "github.com/bazelbuild/buildtools",
    visibility = ["PUBLIC"],
    deps = [
        ":google.golang.org.protobuf",
        ":protobuf",
    ],
)

go_module(
    name = "protobuf",
    install = [
        "jsonpb",
        "proto",
    ],
    licences = ["BSD-3"],
    module = "github.com/golang/protobuf",
    version = "v1.4.3",
    visibility = ["PUBLIC"],
    deps = [":google.golang.org.protobuf"],
)

go_module(
    name = "google.golang.org.protobuf",
    install = [
        "encoding/protojson",
        "encoding/prototext",
        "encoding/protowire",
        "internal/descfmt",
        "internal/descopts",
        "internal/detrand",
        "internal/encoding/defval",
        "internal/encoding/json",
        "internal/encoding/messageset",
        "internal/encoding/tag",
        "internal/encoding/text",
        "internal/errors",
        "internal/fieldsort",
        "internal/filedesc",
        "internal/filetype",
        "internal/flags",
        "internal/genid",
        "internal/impl",
        "internal/mapsort",
        "internal/pragma",
        "internal/set",
        "internal/strs",
        "internal/version",
        "proto",
        "reflect/protoreflect",
        "reflect/protoregistry",
        "runtime/protoiface",
        "runtime/protoimpl",
    ],
    licences = ["BSD-3"],
    module = "google.golang.org/protobuf",
    version = "v1.25.0",
    visibility = ["PUBLIC"],
)
