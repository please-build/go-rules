# Please Go package driver

A go package driver for Please. 

## Usage
Download the package driver from one of the releases 
[here](https://github.com/please-build/go-rules/releases?q=plz-gopackagesdriver-*&expanded=true), then set 
the `GOPACKAGESDRIVER` env var to that path.

## Features
The package driver provides a bridge between tools like `gopls`, and many of the linter in `golangci-lint`. Anything 
that uses `github.com/x/tools` under the hood to list packages will also work. Instead of using `go list`, these tools 
will instead use this package driver. This then presents these packages as Please understands them, including the 
generated sources which would otherwise be missing from the driver response.  
