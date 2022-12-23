package main

import (
	"encoding/json"
	"os"

	"github.com/peterebden/go-cli-init/v5/flags"
	"github.com/peterebden/go-cli-init/v5/logging"

	"github.com/please-build/go-rules/tools/driver/packages"
)

var log = logging.MustGetLogger()

var opts = struct {
	Usage     string
	Verbosity logging.Verbosity `short:"v" long:"verbosity" default:"warning" description:"Verbosity of output (higher number = more output)"`
	NoInput   bool              `short:"n" long:"no_input" description:"Assume a default config and don't try to read from stdin"`
	Args      struct {
		Files []string `positional-arg-name:"file" required:"true"`
	} `positional-args:"true"`
}{
	Usage: `
This is an implementation of the Go packages driver protocol for Please.
The protocol is defined by golang.org/x/tools/go/package; typically you specify this binary
in the GOPACKAGESDRIVER environment variable, and it will then be used in place of go list
to return information about Go package structure.

This tool is experimental.
`,
}

func main() {
	flags.ParseFlagsOrDie("Please Go package driver", &opts, nil)
	logging.InitLogging(opts.Verbosity)

	req := &packages.DriverRequest{}
	if !opts.NoInput {
		if err := json.NewDecoder(os.Stdin).Decode(req); err != nil {
			log.Fatalf("Failed to read request: %s", err)
		}
	}
	resp, err := packages.Load(req, opts.Args.Files)
	if err != nil {
		log.Fatalf("Failed to load packages: %s", err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(resp); err != nil {
		log.Fatalf("Failed to write packages: %s", err)
	}
}
