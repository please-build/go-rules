subinclude("//build_defs:go")

# This is mostly to check that the template doesn't error out when we don't have any tests
go_library(
    name = "no_test_lib",
    srcs = ["notest.go"],
)

go_test(
    name = "no_test_test",
    srcs = ["no_test_test.go"],
    deps = [":no_test_lib"],
)
