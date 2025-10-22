subinclude(
    "//build_defs:go",
)

go_toolchain(
    name = "toolchain",
    install_std = False,
    version = "1.25.2",
)

go_stdlib(
    name = "std",
    tags = [
        "netgo",
        "osusergo",
    ],
)

go_repo(
    licences = ["MIT"],
    module = "github.com/stretchr/testify",
    version = "v1.7.1",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "github.com/pmezard/go-difflib",
    version = "v1.0.0",
)

go_repo(
    licences = ["ISC"],
    module = "github.com/davecgh/go-spew",
    version = "v1.1.1",
)

go_repo(
    licences = ["MIT"],
    module = "gopkg.in/yaml.v3",
    version = "v3.0.0-20210107192922-496545a6307b",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "github.com/thought-machine/go-flags",
    version = "v1.6.2",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/sys",
    version = "v0.36.0",
)

go_repo(
    licences = ["MIT"],
    module = "github.com/dustin/go-humanize",
    version = "v1.0.0",
)

go_repo(
    licences = ["Apache-2.0"],
    module = "github.com/peterebden/go-cli-init/v5",
    version = "v5.2.0",
)

go_repo(
    licences = ["MIT"],
    module = "github.com/cespare/xxhash/v2",
    version = "v2.3.0",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/mod",
    version = "v0.29.0",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/xerrors",
    version = "v0.0.0-20200804184101-5ec99f83aff1",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/tools",
    version = "v0.37.0",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "gopkg.in/op/go-logging.v1",
    version = "v1.0.0-20160211212156-b2cb9fa56473",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/crypto",
    version = "v0.26.0",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/term",
    version = "v0.23.0",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/sync",
    version = "v0.17.0",
)

go_repo(
    licences = ["BSD-3"],
    module = "github.com/golang/protobuf",
    version = "v1.4.3",
)

go_repo(
    licences = ["BSD-3"],
    module = "google.golang.org/protobuf",
    version = "v1.25.0",
)

go_repo(
    licences = ["Apache-2.0"],
    module = "github.com/thought-machine/please",
    strip = [
        "test",
        "tools/please_go_test/gotest/test_data",
    ],
    version = "v12.2.0+incompatible",
)

go_repo(
    licences = ["CC0-1.0"],
    module = "github.com/please-build/please-rules",
    version = "v0.0.0-20240510030102-b6acfa733ea1",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "github.com/jessevdk/go-flags",
    version = "v1.5.0",
)

go_repo(
    licences = ["Apache-2.0"],
    module = "github.com/please-build/go-rules/tools/please_go",
    version = "v0.0.0-20230117162622-4a2c3e910628",
)

go_mod_download(
    name = "github.com_please-build_buildtools_dl",
    labels = ["go_replace_directive"],
    licences = ["Apache-2.0"],
    module = "github.com/please-build/buildtools",
    version = "v0.0.0-20221110131218-762712d8ce3f",
)

go_repo(
    download = ":github.com_please-build_buildtools_dl",
    labels = ["go_replace_directive"],
    module = "github.com/bazelbuild/buildtools",
)
