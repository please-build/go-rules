package main

import (
	"encoding/json"
	"os"

	"github.com/peterebden/go-cli-init/v5/flags"
	"github.com/peterebden/go-cli-init/v5/logging"
	xpackages "golang.org/x/tools/go/packages"

	"github.com/please-build/go-rules/tools/driver/packages"
)

var log = logging.MustGetLogger()

var opts = struct {
	Usage      string
	Verbosity  logging.Verbosity `short:"v" long:"verbosity" env:"PLZ_GOPACKAGESDRIVER_VERBOSITY" default:"warning" description:"Verbosity of output (higher number = more output)"`
	LogFile    string `long:"log_file" env:"PLZ_GOPACKAGESDRIVER_LOG_FILE" description:"Log additionally to this file"`
	NoInput    bool              `short:"n" long:"no_input" description:"Assume a default config and don't try to read from stdin"`
	WorkingDir string            `short:"w" long:"working_dir" description:"Change to this working directory before running"`
	OutputFile string            `short:"o" long:"output_file" env:"PLZ_GOPACKAGESDRIVER_OUTPUT_FILE" description:"File to write output to (in addition to stdout)"`
	SearchDir  string            `short:"s" long:"search_dir" env:"PLZ_GOPACKAGESDRIVER_SEARCHDIR" description:"Search this directory for modinfo files, instead of querying plz"`
	Args       struct {
		Files []string `positional-arg-name:"file"`
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
	if opts.LogFile != "" {
		logging.MustInitFileLogging(opts.Verbosity, opts.Verbosity, opts.LogFile)
	}

	if opts.WorkingDir != "" {
		if err := os.Chdir(opts.WorkingDir); err != nil {
			log.Fatalf("Failed to change working directory: %s", err)
		}
	}

	req := &packages.DriverRequest{}
	if opts.NoInput {
		req.Mode = xpackages.NeedExportFile
	} else {
		if err := json.NewDecoder(os.Stdin).Decode(req); err != nil {
			log.Fatalf("Failed to read request: %s", err)
		}
		log.Debug("Received driver request: %v", req)
	}
	log.Debug("Loading packages for %s...", opts.Args.Files)
	resp, err := load(req)
	if err != nil {
		log.Fatalf("Failed to load packages: %s", err)
	}

	if opts.OutputFile != "" {
		b, _ := json.MarshalIndent(resp, "", "    ")
		if err := os.WriteFile(opts.OutputFile, b, 0644); err != nil {
			log.Fatalf("Failed to write output file: %s", err)
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(resp); err != nil {
		log.Fatalf("Failed to write packages: %s", err)
	}
}

func load(req *packages.DriverRequest) (*packages.DriverResponse, error) {
	if opts.SearchDir == "" {
		return packages.Load(req, opts.Args.Files)
	}
	return packages.LoadOffline(req, opts.SearchDir, opts.Args.Files)
}
