subinclude("//test/build_defs:e2e")

# TODO(jpoole): remove these sed calls once this has been fixed:
#   https://github.com/thought-machine/please/issues/2604
please_repo_e2e_test(
    name = "go_get_test",
    plz_command = "sed -i 's@<pleasegotool>'@$TOOLS_PLEASE_GO'@g' .plzconfig && sed -i 's@<gotool>'@$TOOLS_GO'@g' .plzconfig && plz build ///third_party/go/github_com_stretchr_testify//...",
    repo = "test_repo",
)
