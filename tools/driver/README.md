# Please Go package driver

A go package driver for Please. 

## Usage
Download the package driver from one of the releases 
[here](https://github.com/please-build/go-rules/releases?q=plz-gopackagesdriver-*&expanded=true), then set 
the `GOPACKAGESDRIVER` env var to that path.

## Features
The package driver provides a bridge between go tooling and build systems like Please. Tools like `gopls`, and many of 
the linters in `golangci-lint` respect the package driver. Anything that uses `github.com/x/tools` under the hood will 
work. 

Instead of using `go list`, these tools will instead use this package driver. The package driver will then use Please to 
resolve the packages, handling Please's semantics, including the generated sources which would otherwise be missing from 
the driver response.  
