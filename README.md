# Go Rules
This repo provides Golang build rules for the [Please](https://please.build) build system. 

# Basic usage

First add the plugin to your project. In `plugins/BUILD`:
```python
plugin_repo(
    name = "go",
    revision="vx.x.x",
)
```

Then add the plugin config:
```
[Plugin "go"]
Target = //plugins:go
ImportPath = github.com/example/repo
```

You can then compile and test go Packages like so:
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
    ],
)
```

You can define third party code using `go_get`:
```python
subinclude("///go//build_defs:go")

go_get(module = "github.com/stretchr/testify", version="vx.x.x")

# Add testify's dependencies as defined in its go.mod/go.sum
go_get(module = "github.com/davecgh/go-spew", version="vx.x.x")
go_get(module = "github.com/pmezard/go-difflib", version="vx.x.x")
go_get(module = "github.com/stretchr/objx", version="vx.x.x")
```

To compile a binary, you can use `go_binary()`:
```python
subinclude("///go//build_defs:go")

go_binary(
    name = "bin",
    srcs = ["main.go"],
)
```

Go is especially well suited to writing command line tools and utilities. Binaries can be ran with `plz run`, or used 
as a tool for other [custom rules](https://please.build/codelabs/genrule/#0). 
