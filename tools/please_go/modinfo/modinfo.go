// Package modinfo produces modinfo records for Go's importconfig which are consumed by
// `go tool link` and later accessible from runtime/debug.ReadBuildInfo()
package modinfo

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
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
	if err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".modinfo") {
			return err
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if module, version, found := strings.Cut(strings.TrimSpace(string(contents)), "@"); found {
			bi.Deps = append(bi.Deps, &debug.Module{
				Path:    module,
				Version: version,
			})
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to walk modinfo files: %w", err)
	}
	return os.WriteFile(outputFile, []byte("modinfo "+strconv.Quote(bi.String())), 0644)
}
