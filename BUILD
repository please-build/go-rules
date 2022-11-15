# Include the config file in the build graph otherwise `plz export` doesn't pick it up for the e2e tests
filegroup(
    name = "config",
    srcs = [".plzconfig"],
    visibility = ["PUBLIC"],
)