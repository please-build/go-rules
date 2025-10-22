module github.com/please-build/go-rules

go 1.24.0

ignore (
	plz-out
	test/import
)

require (
	github.com/bazelbuild/buildtools v0.0.0-20221110131218-762712d8ce3f
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/peterebden/go-cli-init/v5 v5.2.0
	github.com/stretchr/testify v1.7.1
	golang.org/x/mod v0.29.0
	golang.org/x/sync v0.17.0
	golang.org/x/term v0.23.0
	golang.org/x/tools v0.37.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/thought-machine/go-flags v1.6.2 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/op/go-logging.v1 v1.0.0-20160211212156-b2cb9fa56473 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/bazelbuild/buildtools => github.com/please-build/buildtools v0.0.0-20221110131218-762712d8ce3f
