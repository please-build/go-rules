# Go Rules

This repo provides [Go](https://go.dev/) build rules for the [Please](https://please.build) build system.

Go is especially suited to writing command line tools and utilities. Binaries can be run with `plz run`, or used
as tools for [custom rules](https://please.build/codelabs/genrule/#0).

# Installation

First, add the plugin to your Please repo. In `plugins/BUILD`:

```python
plugin_repo(
    name = "go",
    revision = "vx.x.x", # A go-rules version
)
```

Next, define a version of Go for the plugin to use. In `third_party/go/BUILD`:

```python
subinclude("///go//build_defs:go")

go_toolchain(
    name = "toolchain",
    version = "1.x.x", # A Go version (see https://go.dev/dl/)
)
```

Finally, configure the plugin in `.plzconfig`:

```ini
[Plugin "go"]
Target = //plugins:go
GoTool = //third_party/go:toolchain|go
ImportPath = github.com/example/repo
```

# Basic usage

```python
subinclude("///go//build_defs:go")

go_library(
    name = "lib",
    srcs = ["lib.go"],
    deps = ["//some:package"],
)

go_test(
    name = "lib_test",
    srcs = ["lib_test.go"],
    deps = [
        ":lib",
        # Third party dependencies are added to a subrepo
        "///third_party/go/github.com_stretchr_testify//assert",
        # If you've passed in the packages you need via the install arg (see below) then you can depend on them like so
        "//third_party/go:testify",
    ],
)
```

Compile binaries with `go_binary`:

```python
subinclude("///go//build_defs:go")

go_binary(
    name = "bin",
    srcs = ["main.go"],
)
```

Introduce dependencies on third-party Go modules with `go_repo`:

```python
subinclude("///go//build_defs:go")

# We can give direct modules a name, and install list so we can reference them nicely as :testify
go_repo(
    name = "testify",
    module = "github.com/stretchr/testify",
    version = "v1.8.2",
    # We add the subset of packages we actually depend on here
    install = [
        "assert",
        "require",
    ]
)

# Indirect modules are referenced internally, so we don't have to name them if we don't want to
go_repo(module = "github.com/davecgh/go-spew", version="v1.1.1")
go_repo(module = "github.com/pmezard/go-difflib", version="v1.0.0")
go_repo(module = "github.com/stretchr/objx", version="v0.5.0")
go_repo(module = "gopkg.in/yaml.v3", version="v3.0.1")
```
