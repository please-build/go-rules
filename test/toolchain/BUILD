subinclude("///e2e//build_defs:e2e")

please_repo_e2e_test(
    name = "toolchain_test",
    expected_output = {
        "out.txt": "wibble wibble wibble",
    },
    plz_command = "plz test -o please.PluginRepo:file://$TMP_DIR/$DATA_PLUGIN_REPO -o plugin.go.pleasegotool:$TOOLS_PLEASE_GO //src/foo:foo_test && plz run //src:main > out.txt",
    repo = "test_repo",
    data = {
        "PLUGIN_REPO": ["//test:export"],
    },
    tools = {
        "PLEASE_GO": ["//tools/please_go:bootstrap"],
    },
)

please_repo_e2e_test(
    name = "freebsd_arm_toolchain_test",
    expected_output = {
        "out.txt": "plz-out/bin/freebsd_arm64/src/main: ELF 64-bit LSB executable, ARM aarch64, version 1 (FreeBSD), statically linked, stripped",
    },
    plz_command = "plz build -o please.PluginRepo:file://$TMP_DIR/$DATA_PLUGIN_REPO -o plugin.go.pleasegotool:$TOOLS_PLEASE_GO //src:main --arch freebsd_arm64 && file plz-out/bin/freebsd_arm64/src/main > out.txt",
    repo = "test_repo",
    data = {
        "PLUGIN_REPO": ["//test:export"],
    },
    tools = {
        "PLEASE_GO": ["//tools/please_go:please_go"],
    },
)
