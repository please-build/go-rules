module github.com/please-build/go-rules

go 1.26.0

ignore (
	plz-out
	test/import
)

require (
	github.com/bazelbuild/buildtools v0.0.0-20260904073137-eaa4d125b423
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/peterebden/go-cli-init/v5 v5.2.1
	github.com/stretchr/testify v1.7.1
	golang.org/x/mod v0.40.0
	golang.org/x/sync v0.22.0
	golang.org/x/term v0.45.0
	golang.org/x/tools v0.49.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/thought-machine/go-flags v1.7.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/op/go-logging.v1 v1.0.0-20160211212156-b2cb9fa56473 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/bazelbuild/buildtools => github.com/please-build/buildtools v0.0.0-20221110131218-762712d8ce3f
