package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterebden/go-cli-init/v5/flags"

	"github.com/please-build/go-rules/tools/please_go/covervars"
	"github.com/please-build/go-rules/tools/please_go/embed"
	"github.com/please-build/go-rules/tools/please_go/filter"
	"github.com/please-build/go-rules/tools/please_go/install"
	"github.com/please-build/go-rules/tools/please_go/test"
)

var opts = struct {
	Usage string

	Install struct {
		BuildTags         []string `long:"build_tag" description:"Any build tags to apply to the build"`
		SrcRoot           string   `short:"r" long:"src_root" description:"The src root of the module to inspect" default:"."`
		ModuleName        string   `short:"n" long:"module_name" description:"The name of the module" required:"true"`
		ImportConfig      string   `short:"i" long:"importcfg" description:"The import config for the modules dependencies" required:"true"`
		LDFlags           string   `long:"ld_flags" description:"Any additional flags to apply to the C linker" env:"LDFLAGS"`
		CFlags            string   `long:"c_flags" description:"Any additional flags to apply when compiling C" env:"CFLAGS"`
		GoTool            string   `short:"g" long:"go_tool" description:"The location of the go binary" default:"go"`
		CCTool            string   `short:"c" long:"cc_tool" description:"The c compiler to use"`
		Out               string   `short:"o" long:"out" description:"The output directory to put compiled artifacts in" required:"true"`
		TrimPath          string   `short:"t" long:"trim_path" description:"Removes prefix from recorded source file paths."`
		PackageConfigTool string   `short:"p" long:"pkg_config_tool" description:"The path to the pkg config" default:"pkg-config"`
		Args              struct {
			Packages []string `positional-arg-name:"packages" description:"The packages to compile"`
		} `positional-args:"true" required:"true"`
	} `command:"install" alias:"i" description:"Compile a go module similarly to 'go install'"`
	Test struct {
		GoTool      string   `short:"g" long:"go_tool" description:"The location of the go binary" env:"TOOLS_GO" default:"go"`
		Dir         string   `short:"d" long:"dir" description:"Directory to search for Go package files for coverage"`
		Exclude     []string `short:"x" long:"exclude" default:"third_party/go" description:"Directories to exclude from search"`
		Output      string   `short:"o" long:"output" description:"Output filename" required:"true"`
		TestPackage string   `short:"t" long:"test_package" description:"The import path of the test package"`
		Benchmark   bool     `short:"b" long:"benchmark" description:"Whether to run benchmarks instead of tests"`
		External    bool     `long:"external" description:"Whether the test is external or not"`
		Args        struct {
			Sources []string `positional-arg-name:"sources" description:"Test source files" required:"true"`
		} `positional-args:"true" required:"true"`
	} `command:"testmain" alias:"t" description:"Generates a go main package to run the tests in a package."`
	CoverVars struct {
		ImportPath string `short:"i" long:"import_path" description:"The import path for the source files"`
		Args       struct {
			Sources []string `positional-arg-name:"sources" description:"Source files to generate embed config for"`
		} `positional-args:"true"`
	} `command:"covervars" description:"Generates coverage variable config for a set of go src files"`
	Filter struct {
		Tags []string `short:"t" long:"tags" description:"Additional build tags to apply"`
		Args struct {
			Sources []string `positional-arg-name:"sources" description:"Source files to filter"`
		} `positional-args:"true"`
	} `command:"filter" alias:"f" description:"Filter go sources based on the go build tag rules."`
	Embed struct {
		Args struct {
			Sources []string `positional-arg-name:"sources" description:"Source files to generate embed config for"`
		} `positional-args:"true"`
	} `command:"embed" alias:"f" description:"Filter go sources based on the go build tag rules."`
}{
	Usage: `
please-go is used by the go build rules to compile and test go modules and packages.

Unlike 'go build', this tool doesn't rely on the go path or modules to find packages. Instead it takes in
a go import config just like 'go tool compile/link -importcfg'.
`,
}

var subCommands = map[string]func() int{
	"install": func() int {
		pleaseGoInstall := install.New(
			opts.Install.BuildTags,
			opts.Install.SrcRoot,
			opts.Install.ModuleName,
			opts.Install.ImportConfig,
			opts.Install.LDFlags,
			opts.Install.CFlags,
			mustResolvePath(opts.Install.GoTool),
			mustResolvePath(opts.Install.CCTool),
			opts.Install.PackageConfigTool,
			opts.Install.Out,
			opts.Install.TrimPath,
		)
		if err := pleaseGoInstall.Install(opts.Install.Args.Packages); err != nil {
			log.Fatal(err)
		}
		return 0
	},
	"testmain": func() int {
		test.PleaseGoTest(opts.Test.GoTool, opts.Test.Dir, opts.Test.TestPackage, opts.Test.Output, opts.Test.Args.Sources, opts.Test.Exclude, opts.Test.Benchmark, opts.Test.External)
		return 0
	},
	"covervars": func() int {
		covervars.GenCoverVars(os.Stdout, opts.CoverVars.ImportPath, opts.CoverVars.Args.Sources)
		return 0
	},
	"filter": func() int {
		filter.Filter(opts.Filter.Tags, opts.Filter.Args.Sources)
		return 0
	},
	"embed": func() int {
		if err := embed.WriteEmbedConfig(opts.Embed.Args.Sources, os.Stdout); err != nil {
			log.Fatalf("failed to generate embed config: %v", err)
		}
		return 0
	},
}

func main() {
	command := flags.ParseFlagsOrDie("please-go", &opts, nil)
	os.Exit(subCommands[command]())
}

// mustResolvePath converts a relative path to absolute if it has any separators in it.
func mustResolvePath(in string) string {
	if in == "" {
		return in
	}
	if !filepath.IsAbs(in) && strings.ContainsRune(in, filepath.Separator) {
		abs, err := filepath.Abs(in)
		if err != nil {
			log.Fatalf("Failed to make %s absolute: %s", in, err)
		}
		return abs
	}
	return in
}
