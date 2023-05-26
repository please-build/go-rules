// Package modinfo produces modinfo records for Go's importconfig which are consumed by
// `go tool link` and later accessible from runtime/debug.ReadBuildInfo()
package modinfo

import (
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"sort"
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
	sort.Slice(bi.Deps, func(i, j int) bool { return bi.Deps[i].Path < bi.Deps[j].Path })
	return os.WriteFile(outputFile, []byte(fmt.Sprintf("modinfo %q\n", modInfoData(bi.String()))), 0644)
}

// modInfoData wraps the given string in Go's modinfo. This mimics what go build does in order
// for `go version` to be able to find this lot later on.
func modInfoData(modinfo string) string {
	// These are not exported from the stdlib (they're in cmd/go/internal/modload) so we must duplicate :(
	start, _ := hex.DecodeString("3077af0c9274080241e1c107e6d618e6")
	end, _ := hex.DecodeString("f932433186182072008242104116d8f2")
	return string(start) + modinfo + string(end)
}
