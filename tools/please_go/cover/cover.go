// Package cover instruments source files for coverage.
// It orchestrates some of the functionality of `go tool cover`.
package cover

import (
	"os/exec"
)

// WriteCoverage writes the necessary Go coverage information for a set of sources.
func WriteCoverage(goTool, covercfg, output, pkg string, srcs []string) error {
	// N.B. this is only the Go 1.20+ version.
	cmd := exec.Command(goTool, append([]string{"tool", "cover", "-mode=set", "-var=goCover", "-pkgcfg=pkgcfg", "-outfilelist", output}, srcs...))
}
