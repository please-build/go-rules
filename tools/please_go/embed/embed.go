// Package embed implements parsing of embed directives in Go files.
package embed

import (
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"strings"
)

// Cfg is the structure of a Go embedcfg file.
type Cfg struct {
	Patterns map[string][]string
	Files    map[string]string
}

// WriteEmbedConfig writes the embed config to w
func WriteEmbedConfig(files []string, w io.Writer) error {
	cfg, err := Parse(files)
	if err != nil {
		return err
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// Parse parses the given files and returns the embed information in them.
func Parse(gofiles []string) (*Cfg, error) {
	cfg := &Cfg{
		Patterns: map[string][]string{},
		Files:    map[string]string{},
	}
	for _, dir := range dirs(gofiles) {
		pkg, err := build.ImportDir(dir, build.ImportComment)
		if err != nil {
			return nil, err
		}
		// We munge all patterns together at this point, if a file is in our input sources we want to know about it regardless.
		if err := cfg.AddPackage(pkg); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

// AddPackage parses a go package and adds any embed patterns to the configuration
func (cfg *Cfg) AddPackage(pkg *build.Package) error {
	for _, pattern := range append(append(pkg.EmbedPatterns, pkg.TestEmbedPatterns...), pkg.XTestEmbedPatterns...) {
		paths, err := relglob(pkg.Dir, pattern)
		if err != nil {
			return err
		}
		cfg.Patterns[pattern] = paths
		for _, p := range paths {
			cfg.Files[p] = path.Join(pkg.Dir, p)
		}
	}
	return nil
}

func dirs(files []string) []string {
	dirs := []string{}
	seen := map[string]bool{}
	for _, f := range files {
		if dir := path.Dir(f); !seen[dir] {
			dirs = append(dirs, dir)
			seen[dir] = true
		}
	}
	return dirs
}

func relglob(dir, pattern string) ([]string, error) {
	// Go allows prefixing the pattern with all: which picks up files prefixed with . or _ (by default these should be ignored)
	includeHidden := false
	if strings.HasPrefix(pattern, "all:") {
		pattern = strings.TrimPrefix(pattern, "all:")
		includeHidden = true
	}

	paths, err := filepath.Glob(path.Join(dir, pattern))
	if err == nil && len(paths) == 0 {
		return nil, fmt.Errorf("pattern %s: no matching paths found", pattern)
	}
	ret := make([]string, 0, len(paths))
	for _, p := range paths {
		if err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			} else if !d.IsDir() {
				if hidden := strings.HasPrefix(d.Name(), ".") || strings.HasPrefix(d.Name(), "_"); !hidden || includeHidden {
					ret = append(ret, strings.TrimLeft(strings.TrimPrefix(path, dir), string(filepath.Separator)))
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return ret, err
}
