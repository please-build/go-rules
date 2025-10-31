subinclude("//build_defs:go")

# For testing the go build commands, where we need to be careful about things
# being in subdirectories.
go_library(
    name = "test",
    srcs = ["test.go"],
    test_only = True,
    visibility = ["//test/..."],
)

# Tests a go 'external test', where you can have files in the same directory
# with a _test suffix.
go_test(
    name = "external_test",
    srcs = ["external_test.go"],
    definitions = {
        "github.com/please-build/go-rules/test/external.Var": "var",
        "github.com/please-build/go-rules/test/external.Var2": "var1 var2",
        "github.com/please-build/go-rules/test/external.ExecGitShow": git_show("%ct"),
        "github.com/please-build/go-rules/test/external.ExecGitState": git_state(
            clean_label = "shiny",
            dirty_label = "not shiny",
        ),
        "github.com/please-build/go-rules/test/external.ExecGitCommit": git_commit(),
        "github.com/please-build/go-rules/test/external.ExecGitBranchFull": git_branch(short = False),
        "github.com/please-build/go-rules/test/external.ExecGitBranchShort": git_branch(),
    },
    external = True,
    labels = ["no_cirrus"],  # Git branches aren't happening as expected here.
    deps = [
        ":test",
        "///third_party/go/github.com_stretchr_testify//assert",
    ],
)
