// Package test contains utilities used by plz_go_test.
package test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// A CoverVar is just a combination of package path and variable name
// for one of the templated-in coverage variables.
type CoverVar struct {
	Dir, ImportPath, ImportName, Var, File string
}

// FindCoverVars searches the given directory recursively to find all Go files with coverage variables.
func FindCoverVars(dir string, excludedDirs []string) ([]CoverVar, error) {
	if dir == "" {
		return nil, nil
	}
	excludeMap := map[string]struct{}{}
	for _, e := range excludedDirs {
		excludeMap[e] = struct{}{}
	}
	var ret []CoverVar

	err := filepath.Walk(dir, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if _, present := excludeMap[name]; present {
			if info.IsDir() {
				return filepath.SkipDir
			}
		} else if strings.HasSuffix(name, ".cover_vars") {
			vars, err := parseCoverVars(name)
			if err != nil {
				return err
			}
			ret = append(ret, vars...)
		}
		return nil
	})
	return ret, err
}

// parseCoverVars parses the coverage variables file for all cover vars
func parseCoverVars(filepath string) ([]CoverVar, error) {
	dir := strings.TrimRight(path.Dir(filepath), "/")
	if dir == "" {
		dir = "."
	}

	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	ret := make([]CoverVar, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed cover var line in %v: %v", filepath, line)
		}

		ret = append(ret, coverVar(dir, parts[0], parts[1]))
	}

	return ret, nil
}

func coverVar(dir, importPath, v string) CoverVar {
	fmt.Fprintf(os.Stderr, "Found cover variable: %s %s %s\n", dir, importPath, v)
	f := path.Join(dir, strings.TrimPrefix(v, "GoCover_"))
	if strings.HasSuffix(f, "_go") {
		f = f[:len(f)-3] + ".go"
	}
	return CoverVar{
		Dir:        dir,
		ImportPath: importPath,
		Var:        v,
		File:       f,
	}
}

