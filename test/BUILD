# More grunge to load its plugin config
subinclude("///e2e//build_defs:e2e")

# This is quite grungy, but allows us to export the plugin for use in e2e tests
exports = [
    "//:config",
    "//build_defs:go",
]

genrule(
    name = "export",
    cmd = "rm .plzconfig && $TOOL export {exports} -o $TMP_DIR/_tmp && tar -czf $OUT _tmp".format(exports=" ".join(exports)),
    outs = ["go_rules.tar.gz"],
    deps = exports,
    visibility=["//test/..."],
    tools = [CONFIG.E2E.PLEASE_TOOL],
)