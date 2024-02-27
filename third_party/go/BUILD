subinclude("//build_defs:go")

go_toolchain(
    name = "toolchain",
    install_std = False,
    version = "1.21.1",
)

go_stdlib(
    name = "std",
    tags = [
        "netgo",
        "osusergo",
    ],
)

go_repo(
    module = "github.com/stretchr/testify",
    name = "testify",
    install = [
        "assert",
        "require",
    ],
    licences = ["MIT"],
    version = "v1.7.0",
)

go_repo(
    module = "github.com/pmezard/go-difflib",
    name = "difflib",
    install = ["difflib"],
    licences = ["BSD-3-Clause"],
    version = "v1.0.0",
)

go_repo(
    module = "github.com/davecgh/go-spew",
    name = "spew",
    install = ["spew"],
    licences = ["ISC"],
    version = "v1.1.1",
)

go_repo(
    module = "gopkg.in/yaml.v3",
    name = "yaml",
    install = ["."],
    licences = ["MIT"],
    version = "v3.0.0-20210107192922-496545a6307b",
)

go_repo(
    module = "github.com/thought-machine/go-flags",
    name = "go-flags",
    install = ["."],
    licences = ["BSD-3-Clause"],
    version = "v1.6.0",
)

go_repo(
    module = "golang.org/x/sys",
    name = "xsys",
    install = ["..."],
    licences = ["BSD-3-Clause"],
    version = "v0.13.0",
)

go_repo(
    module = "github.com/dustin/go-humanize",
    name = "humanize",
    install = ["."],
    licences = ["MIT"],
    version = "v1.0.0",
)

go_repo(
    module = "github.com/peterebden/go-cli-init/v5",
    name = "flags",
    install = [
        "flags",
        "logging",
    ],
    licences = ["Apache-2.0"],
    version = "v5.1.0",
)

go_repo(
    module = "github.com/cespare/xxhash/v2",
    name = "xxhash",
    install = ["."],
    licences = ["MIT"],
    version = "v2.1.2",
)

go_repo(
    module = "golang.org/x/mod",
    name = "mod",
    install = [
        "modfile",
        "module",
        "semver",
    ],
    licences = ["BSD-3-Clause"],
    version = "v0.5.0",
)

go_repo(
    module = "golang.org/x/xerrors",
    name = "xerrors",
    install = ["."],
    licences = ["BSD-3-Clause"],
    version = "v0.0.0-20200804184101-5ec99f83aff1",
)

go_repo(
    module = "golang.org/x/tools",
    name = "xtools",
    install = [
        "go/gcexportdata",
        "go/packages",
    ],
    licences = ["BSD-3-Clause"],
    version = "v0.4.0",
)

go_repo(
    module = "gopkg.in/op/go-logging.v1",
    name = "logging",
    install = ["."],
    licences = ["BSD-3-Clause"],
    version = "v1.0.0-20160211212156-b2cb9fa56473",
)

go_repo(
    module = "golang.org/x/crypto",
    name = "xcrypto",
    install = ["ssh/terminal"],
    licences = ["BSD-3-Clause"],
    version = "v0.0.0-20210920023735-84f357641f63",
)

go_repo(
    module = "golang.org/x/term",
    name = "xterm",
    install = ["."],
    licences = ["BSD-3-Clause"],
    version = "v0.0.0-20210615171337-6886f2dfbf5b",
)

go_repo(
    module = "golang.org/x/sync",
    name = "xsync",
    install = ["errgroup"],
    licences = ["BSD-3-Clause"],
    version = "v0.1.0",
)

go_repo(
    module = "github.com/bazelbuild/buildtools",
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
    version = "v0.0.0-20230725091334-24cce64e067d",
)

go_repo(
    module = "github.com/golang/protobuf",
    name = "protobuf",
    install = [
        "jsonpb",
        "proto",
    ],
    licences = ["BSD-3"],
    version = "v1.4.3",
)

go_repo(
    module = "google.golang.org/protobuf",
    name = "google.golang.org.protobuf",
    install = [
        "encoding/protojson",
        "encoding/prototext",
        "encoding/protowire",
        "proto",
        "reflect/protoreflect",
        "reflect/protoregistry",
        "runtime/protoiface",
        "runtime/protoimpl",
    ],
    licences = ["BSD-3"],
    version = "v1.25.0",
)
