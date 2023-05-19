// TODO(jpoole): better name for this
package generate

import (
	"fmt"
	"github.com/bazelbuild/buildtools/build"
	gobuild "go/build"
	"os"
	"path/filepath"
	"strings"
)

// Update updates an existing Please project. It may create new BUILD files, however it tries to respect existing build
// rules, updating them as appropriate.
func (g *Generate) Update(moduleName string, paths []string) error {
	done := map[string]struct{}{}
	g.moduleName = moduleName
	for _, path := range paths {
		if strings.HasSuffix(path, "/...") {
			path = strings.TrimSuffix(path, "/...")
			err := filepath.WalkDir(path, func(path string, info os.DirEntry, err error) error {
				if info.IsDir() {
					return nil
				}

				if g.isBuildFile(path) {
					if err := g.update(done, path); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			if err := g.update(done, path); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Generate) isBuildFile(file string) bool {
	base := filepath.Base(file)
	for _, file := range g.buildFileNames {
		if base == file {
			return true
		}
	}
	return false
}

func (g *Generate) findBuildFile(dir string) string {
	for _, name := range g.buildFileNames {
		path := filepath.Join(dir, name)
		if _, err := os.Lstat(path); err == nil {
			return path
		}
	}
	return ""
}

func (g *Generate) loadBuildFile(path string) (*build.File, error) {
	if !g.isBuildFile(path) {
		path = g.findBuildFile(path)
		if path == "" {
			return nil, fmt.Errorf("faild to find build file in %v", path)
		}
	}
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return build.ParseBuild(path, bs)
}

func (g *Generate) update(done map[string]struct{}, path string) error {
	if _, ok := done[path]; ok {
		return nil
	}
	defer func() {
		done[path] = struct{}{}
	}()

	dir := path
	if g.isBuildFile(path) {
		dir = filepath.Dir(path)
	}

	// TODO(jpoole): we should break this up and check each source file so we can split tests out across multiple targets
	pkg, err := g.importDir(dir)
	if err != nil {
		if _, ok := err.(*gobuild.NoGoError); ok {
			return nil
		}
		return err
	}

	libRule := g.libRule(pkg)
	testRule := g.testRule(pkg, libRule)

	file, err := g.loadBuildFile(path)
	if err != nil {
		return err
	}

	libDone := false
	testDone := false

	for _, stmt := range file.Stmt {
		if call, ok := stmt.(*build.CallExpr); ok {
			rule := build.NewRule(call)
			if (rule.Kind() == "go_library" && !pkg.IsCommand()) || (rule.Kind() == "go_binary" && pkg.IsCommand()) {
				if libDone {
					return fmt.Errorf("too many go_library rules in %v", path)
				}
				populateRule(rule, libRule)
				libDone = true
			}
			if rule.Kind() == "go_test" {
				if testDone {
					fmt.Fprintln(os.Stderr, "WARNING: too many go_test rules in ", path)
					continue
				}
				if rule.Attr("external") != nil {
					continue
				}
				populateRule(rule, testRule)
				testDone = true
			}
		}
	}

	if !libDone && libRule != nil {
		r := NewRule("go_library", libRule.name)
		populateRule(r, libRule)
		file.Stmt = append(file.Stmt, r.Call)
	}

	if !testDone && testRule != nil {
		r := NewRule("go_test", testRule.name)
		populateRule(r, testRule)
		file.Stmt = append(file.Stmt, r.Call)
	}
	return os.WriteFile(file.Path, build.Format(file), 0664)
}
