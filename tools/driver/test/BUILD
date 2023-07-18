subinclude("//build_defs:go")

go_library(
    name = "test",
    srcs = ["test.go"],
    deps = [
        "//tools/driver/packages",
        "//third_party/go:testify",
    ],
)

gentest(
    name = "gosec_test",
    tools = {
        "driver": ["//tools/driver:plz-gopackagesdriver"],
        "gosec": ["//third_party/binary:gosec"],
    },
    deps = [
        "//tools/driver/packages",
    ],
    needs_transitive_deps = True,
    cmd = 'export GOPACKAGESDRIVER="$TOOLS_DRIVER"; "$TOOLS_GOSEC" tools/driver/packages/packages.go',
    test_cmd = "true",
    requires = [
        "go",
        "go_src",
    ],
)
