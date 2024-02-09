package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/peterebden/go-cli-init/v5/flags"

	"github.com/please-build/go-rules/tools/please_go/cover"
	"github.com/please-build/go-rules/tools/please_go/covervars"
	"github.com/please-build/go-rules/tools/please_go/embed"
	"github.com/please-build/go-rules/tools/please_go/filter"
	"github.com/please-build/go-rules/tools/please_go/generate"
	"github.com/please-build/go-rules/tools/please_go/goget"
	"github.com/please-build/go-rules/tools/please_go/install"
	"github.com/please-build/go-rules/tools/please_go/modinfo"
	"github.com/please-build/go-rules/tools/please_go/packageinfo"
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
		PackageConfigTool string   `short:"p" long:"pkg_config_tool" env:"PKG_CONFIG_TOOL" description:"The path to the pkg config" default:"pkg-config"`
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
	Cover struct {
		GoTool      string `short:"g" long:"go" default:"go" description:"Go binary to run"`
		CoverageCfg string `short:"c" long:"covcfg" required:"true" description:"Output coveragecfg file to feed into go tool compile"`
		Output      string `short:"o" long:"output" required:"true" description:"File that will contain output names of modified files"`
		Pkg         string `long:"pkg" env:"PKG_DIR" description:"Package that we're in within the repo"`
		PkgName     string `short:"p" long:"pkg_name" hidden:"true" description:"Deprecated, has no effect"`
		Args        struct {
			Sources []string `positional-arg-name:"sources" required:"true" description:"Source files to generate embed config for"`
		} `positional-args:"true"`
	} `command:"cover" description:"Generates coverage information for a package."`
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
	} `command:"embed" alias:"e" description:"Generate embed config for a set of Go source files"`
	PackageInfo struct {
		ImportPath string            `short:"i" long:"import_path" description:"Go import path (e.g. github.com/please-build/go-rules)"`
		Pkg        string            `long:"pkg" env:"PKG_DIR" description:"Package that we're in within the repo"`
		ImportMap  map[string]string `short:"m" long:"import_map" description:"Existing map of imports"`
		Subrepo    string            `short:"s" long:"subrepo" description:"Subrepo root that this package is within"`
		Module     string            `long:"mod" description:"The module this is within, if present"`
	} `command:"package_info" alias:"p" description:"Creates an info file about a Go package"`
	ModuleInfo struct {
		ModulePath   string   `short:"m" long:"module_path" required:"true" description:"Import path of the module in question"`
		Srcs         string   `long:"srcs" env:"SRCS_SRCS" required:"true" description:"Source files of the module"`
		ImportConfig string   `long:"importconfig" env:"SRCS_IC" description:"Importconfig file for locating gc export data"`
		Packages     []string `short:"p" long:"packages" description:"Packages to include in the module"`
	} `command:"module_info" alias:"m" description:"Creates an info file about a series of packages in a go_module"`
	Generate struct {
		SrcRoot          string   `short:"r" long:"src_root" description:"The src root of the module to inspect"`
		ImportPath       string   `long:"import_path" description:"overrides the module's import path. If not set, the import path from the go.mod will be used.'"`
		ThirdPartyFolder string   `short:"t" long:"third_part_folder" description:"The folder containing the third party subrepos" default:"third_party/go"`
		ModFile          string   `long:"mod_file" description:"Path to the host repo mod file to use to resolve dependencies against (dependencies will be resolved against the module as well if it exists)"`
		Module           string   `long:"module" description:"The name of the current module"`
		Version          string   `long:"version" description:"The version of the current module"`
		Install          []string `long:"install" description:"The packages to add to the :install alias"`
		BuildTags        []string `long:"build_tag" description:"Any build tags to apply to the build"`
		Subrepo          string   `long:"subrepo" description:"The subrepo root to output into"`
		Args             struct {
			Requirements []string `positional-arg-name:"requirements" description:"Any module requirements not included in the go.mod"`
		} `positional-args:"true"`
	} `command:"generate" alias:"g" description:"Generate build targets for a Go module"`
	GoGet struct {
		ModFile string `short:"m" long:"mod_file" description:"A go.mod file to use as a set of reuirementzs"`
		Args    struct {
			Requirements []string `positional-arg-name:"requirements" description:"a set of module@version pairs"`
		} `positional-args:"true"`
	} `command:"get" description:"Generate go_get rules"`
	ModInfo struct {
		GoTool     string `short:"g" long:"go" env:"TOOLS_GO" required:"true" description:"The Go tool we'll use"`
		ModulePath string `short:"m" long:"module_path" description:"The path for the module being built"`
		Pkg        string `short:"p" long:"package" env:"PKG_DIR" description:"The package directory within the repo"`
		BuildMode  string `short:"b" long:"build_mode" default:"exe" description:"The Go build mode being used"`
		Out        string `short:"o" long:"out" env:"OUT" required:"true" description:"File to write the output to"`
		Write      bool   `short:"w" long:"write" hidden:"true" description:"Print this binary's own modinfo"`
		CgoEnabled string `short:"c" long:"cgo_enabled" env:"CGO_ENABLED" description:"Whether cgo is enabled or not"`
		GoOS       string `long:"goos" env:"OS" description:"OS we're compiling for"`
		GoArch     string `long:"goarch" env:"ARCH" description:"Architecture we're compiling for"`
	} `command:"modinfo" description:"Generates Go modinfo for the linter"`
	GenerateModuleVersion struct {
		ModulePath string `short:"m" long:"module_path" required:"true" description:"The module's path"`
		Version    string `long:"version" required:"true" description:"The module's (semantic) version number"`
		Validate   bool   `long:"validate" description:"Check validity of the given module import path and version number"`
		Out        string `short:"o" long:"out" env:"OUT" required:"true" description:"File to write the output to"`
	} `command:"generate_module_version" description:"Generates a module version file for a third-party Go module"`
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
	"cover": func() int {
		if err := cover.WriteCoverage(opts.Cover.GoTool, opts.Cover.CoverageCfg, opts.Cover.Output, opts.Cover.Pkg, opts.Cover.Args.Sources); err != nil {
			log.Fatalf("failed to write coverage: %s", err)
		}
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
	"generate": func() int {
		gen := opts.Generate
		g := generate.New(gen.SrcRoot, gen.ThirdPartyFolder, gen.ModFile, gen.Module, gen.Version, gen.Subrepo, []string{"BUILD", "BUILD.plz"}, gen.Args.Requirements, gen.Install, gen.BuildTags)
		if err := g.Generate(); err != nil {
			log.Fatalf("failed to generate go rules: %v", err)
		}
		return 0
	},
	"get": func() int {
		if opts.GoGet.ModFile != "" {
			if err := goget.GetMod(opts.GoGet.ModFile); err != nil {
				log.Fatalf("failed to generate go rules: %v", err)
			}
			return 0
		}
		if err := goget.GoGet(opts.GoGet.Args.Requirements); err != nil {
			log.Fatalf("failed to generate go rules: %v", err)
		}
		return 0
	},
	"package_info": func() int {
		pi := opts.PackageInfo
		if err := packageinfo.WritePackageInfo(pi.ImportPath, pi.Pkg, "", pi.ImportMap, nil, pi.Subrepo, pi.Module, os.Stdout); err != nil {
			log.Fatalf("failed to write package info: %s", err)
		}
		return 0
	},
	"module_info": func() int {
		mi := opts.ModuleInfo
		if err := packageinfo.WritePackageInfo(mi.ModulePath, mi.Srcs, mi.ImportConfig, nil, mi.Packages, "", "", os.Stdout); err != nil {
			log.Fatalf("failed to write module info: %s", err)
		}
		return 0
	},
	"modinfo": func() int {
		mi := opts.ModInfo
		if mi.Write {
			info, ok := debug.ReadBuildInfo()
			if !ok {
				log.Fatalf("build info not available")
			}
			os.Stderr.Write([]byte(info.String() + "\n"))
		}
		if err := modinfo.WriteModInfo(mi.GoTool, mi.ModulePath, filepath.Join(mi.ModulePath, mi.Pkg), mi.BuildMode, mi.CgoEnabled, mi.GoOS, mi.GoArch, mi.Out); err != nil {
			log.Fatalf("failed to write modinfo: %s", err)
		}
		return 0
	},
	"generate_module_version": func() int {
		mv := opts.GenerateModuleVersion
		if err := modinfo.WriteModuleVersion(mv.ModulePath, mv.Version, mv.Validate, mv.Out); err != nil {
			log.Fatalf("failed to generate module version file: %v", err)
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
