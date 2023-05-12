// Package modinfo produces modinfo records for Go's importconfig which are consumed by
// `go tool link` and later accessible from runtime/debug.ReadBuildInfo()
package modinfo

import (
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"strings"
)

// WriteModInfo writes mod info to the given output file
func WriteModInfo(goTool, modulePath, pkgPath, buildMode, outputFile string) error {
	// Nab the Go version from the tool
	out, err := exec.Command(goTool, "version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to exec %s: %w", goTool, err)
	}
	bi := debug.BuildInfo{
		GoVersion: strings.TrimSpace(string(out)),
		Path:      pkgPath,
		Main: debug.Module{
			Path: modulePath,
			// We don't have a concept of a module version here
		},
		Settings: []debug.BuildSetting{
			{Key: "-buildmode", Value: buildMode},
		},
	}
	return os.WriteFile(outputFile, []byte("modinfo "+strconv.Quote(bi.String())), 0644)
}
