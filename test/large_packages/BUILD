subinclude("//test/build_defs:e2e")

please_repo_e2e_test(
    name = "large_package_test",
    # Package known to have ~1800 source files
    plz_command = "plz --noupdate build ///third_party/go/github.com_elastic_go-elasticsearch_v9//typedapi/types",
    repo = "test_repo",
)
