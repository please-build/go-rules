// Package cover instruments source files for coverage.
// It orchestrates some of the functionality of `go tool cover`.
package cover

import (
	"bytes"
	"encoding/json"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/please-build/go-rules/tools/please_go/install/toolchain"
)

// WriteCoverage writes the necessary Go coverage information for a set of sources.
func WriteCoverage(goTool, coverTool, covercfg, output, pkgConfigFile, pkg string, srcs []string) error {
	pkgName, err := packageName(srcs[0])
	if err != nil {
		return err
	}

	b, _ := json.Marshal(coverConfig{
		OutConfig:   covercfg,
		PkgPath:     pkg,
		PkgName:     pkgName,
		Granularity: "perblock",
	})
	if err := os.WriteFile(pkgConfigFile, b, 0644); err != nil {
		return err
	}
	var buf bytes.Buffer
	// 1.21 requires a cover vars file to be written into the output file list
	need121CoverVars := needs121CoverVars(goTool)
	// 1.27 needs package-relative paths (https://go.dev/issue/70478)
	needRelativePaths := needs127RelativePaths(goTool)
	if coverTool != "" || needRelativePaths || need121CoverVars {
		buf.WriteString(filepath.Join(filepath.Dir(srcs[0]), "_covervars.cover.go\n"))
	}
	for _, src := range srcs {
		buf.WriteString(strings.TrimSuffix(src, ".go") + ".cover.go\n")
	}
	if err := os.WriteFile(output, buf.Bytes(), 0644); err != nil {
		return err
	}
	if needRelativePaths {
		for i, src := range srcs {
			rel, err := filepath.Rel(pkg, src)
			if err != nil {
				log.Printf("failed to make path relative: %s", err)
				continue
			}
			// Mutate in place so that the package-relative path is passed to `go tool cover`
			srcs[i] = rel
		}
		// If we're using relative paths, we need a separate 'outfilelist' file with the relative paths in it for 'go tool cover' to read
		// while still producing the output file with the root-relative paths for Please to consume.
		output = "outfilelist.txt"
		buf.Reset()
		buf.WriteString("_covervars.cover.go\n")
		for _, src := range srcs {
			buf.WriteString(strings.TrimSuffix(src, ".go") + ".cover.go\n")
		}
		if err := os.WriteFile(filepath.Join(pkg, output), buf.Bytes(), 0644); err != nil {
			return err
		}
	}
	var cmd *exec.Cmd
	if coverTool != "" {
		cmd = exec.Command(coverTool, append([]string{"-mode=set", "-var=_plz_goCover", "-pkgcfg", pkgConfigFile, "-outfilelist", output}, srcs...)...)
	} else {
		cmd = exec.Command(goTool, append([]string{"tool", "cover", "-mode=set", "-var=_plz_goCover", "-pkgcfg", pkgConfigFile, "-outfilelist", output}, srcs...)...)
	}
	if needRelativePaths {
		cmd.Dir = pkg
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func needs121CoverVars(goTool string) bool {
	version, err := toolchain.GoMinorVersion(goTool)
	if err != nil {
		panic(err)
	}
	return version >= 21
}

func needs127RelativePaths(goTool string) bool {
	version, err := toolchain.GoMinorVersion(goTool)
	if err != nil {
		panic(err)
	}
	return version >= 27
}

// This is a copy of the one from internal/coverage (why does that need to be internal??)
type coverConfig struct {
	OutConfig   string
	PkgPath     string
	PkgName     string
	Granularity string
	ModulePath  string
}

// packageName returns the Go package for a file.
func packageName(filename string) (string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}
	return f.Name.Name, nil
}
