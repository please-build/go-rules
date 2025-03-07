subinclude("//build_defs:go")

# This is an integration test proving we can run within a test directory; without the
# package driver running, gosec will get import errors from the package and fail.
gentest(
    name = "gosec_test",
    cmd = [
        'export PLZ_GOPACKAGESDRIVER_SEARCHDIR="$TMP_DIR"',
        'export GOPACKAGESDRIVER="$TOOLS_DRIVER"',
        '"$TOOLS_GOSEC" --exclude G104,G304,G307 tools/driver/packages',
    ],
    needs_transitive_deps = True,
    no_test_output = True,
    requires = [
        "go",
        "go_src",
    ],
    test_cmd = "true",
    tools = {
        "driver": ["//tools/driver:plz-gopackagesdriver"],
        "gosec": ["//third_party/binary:gosec"],
    },
    deps = [
        "//tools/driver/packages",
    ],
)
