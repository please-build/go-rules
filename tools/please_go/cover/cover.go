// Package cover instruments source files for coverage.
// It orchestrates some of the functionality of `go tool cover`.
package cover

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
)

// WriteCoverage writes the necessary Go coverage information for a set of sources.
func WriteCoverage(goTool, covercfg, output, pkg, pkgName string, srcs []string) error {
	const pkgConfigFile = "pkgcfg"
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
	for _, src := range srcs {
		buf.WriteString(strings.TrimSuffix(src, ".go") + ".cover.go\n")
	}
	if err := os.WriteFile(output, buf.Bytes(), 0644); err != nil {
		return err
	}
	cmd := exec.Command(goTool, append([]string{"tool", "cover", "-mode=set", "-var=goCover", "-pkgcfg", pkgConfigFile, "-outfilelist", output}, srcs...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// This is a copy of the one from internal/coverage (why does that need to be internal??)
type coverConfig struct {
	OutConfig   string
	PkgPath     string
	PkgName     string
	Granularity string
	ModulePath  string
}
