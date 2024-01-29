// Package packageinfo writes information about Go packages in a JSON format.
// This is created at build time and intended to be consumed by the gopackagedriver binary.
package packageinfo

import (
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// WritePackageInfo writes a series of package info files to the given file.
func WritePackageInfo(importPath string, srcRoot, importconfig string, imports map[string]string, installPkgs []string, subrepo, module string, w io.Writer) error {
	// Discover all Go files in the module
	goFiles := map[string][]string{}
	module = modulePath(module, importPath)

	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if name := d.Name(); name == "testdata" {
			return filepath.SkipDir // Don't descend into testdata
		} else if strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			dir := filepath.Dir(path)
			goFiles[dir] = append(goFiles[dir], path)
		}
		return nil
	}
	// Check install packages first
	for _, pkg := range installPkgs {
		if strings.Contains(pkg, "...") {
			pkg = strings.TrimSuffix(pkg, "...")
			if err := filepath.WalkDir(filepath.Join(srcRoot, pkg), walkDirFunc); err != nil {
				return fmt.Errorf("failed to read module dir: %w", err)
			}
		} else {
			dir := filepath.Join(srcRoot, pkg)
			goFiles[dir] = append(goFiles[dir], filepath.Join(srcRoot, pkg))
		}
	}
	if len(installPkgs) == 0 {
		if err := filepath.WalkDir(srcRoot, walkDirFunc); err != nil {
			return fmt.Errorf("failed to read module dir: %w", err)
		}
	}
	if importconfig != "" {
		m, err := loadImportConfig(importconfig)
		if err != nil {
			return fmt.Errorf("failed to read importconfig: %w", err)
		}
		imports = m
	}
	pkgs := make([]*packages.Package, 0, len(goFiles))
	for dir := range goFiles {
		pkgDir := strings.TrimPrefix(strings.TrimPrefix(dir, srcRoot), "/")
		pkg, err := createPackage(filepath.Join(importPath, pkgDir), dir, subrepo, module)
		if _, ok := err.(*build.NoGoError); ok {
			continue // Don't really care, this happens sometimes for modules
		} else if err != nil {
			return fmt.Errorf("failed to import directory %s: %w", dir, err)
		}
		if subrepo != "" {
			_, pkgPath, ok := strings.Cut(imports[pkg.PkgPath], pkg.PkgPath)
			if !ok {
				return fmt.Errorf("Cannot determine export file path for package %s from %s", pkg.PkgPath, imports[pkg.PkgPath])
			}
			// This is a really gross hack to sneak both paths through the one field.
			pkg.ExportFile = filepath.Join(subrepo, pkgPath) + "|" + imports[pkg.PkgPath]
		} else {
			pkg.ExportFile = imports[pkg.PkgPath]
		}
		pkgs = append(pkgs, pkg)
	}
	// Ensure output is deterministic
	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].ID < pkgs[j].ID
	})
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(pkgs)
}

func createPackage(pkgPath, pkgDir, subrepo, module string) (*packages.Package, error) {
	if pkgDir == "" || pkgDir == "." {
		// This happens when we're in the repo root, ImportDir refuses to read it for some reason.
		path, err := filepath.Abs(pkgDir)
		if err != nil {
			return nil, err
		}
		pkgDir = path
	}
	bpkg, err := build.ImportDir(pkgDir, build.ImportComment)
	if err != nil {
		return nil, err
	}
	bpkg.ImportPath = pkgPath
	return FromBuildPackage(bpkg, subrepo, module), nil
}

// FromBuildPackage creates a packages Package from a build Package.
func FromBuildPackage(pkg *build.Package, subrepo, module string) *packages.Package {
	p := &packages.Package{
		ID:              pkg.ImportPath,
		Name:            pkg.Name,
		PkgPath:         pkg.ImportPath,
		GoFiles:         make([]string, len(pkg.GoFiles)),
		CompiledGoFiles: make([]string, len(pkg.GoFiles)),
		OtherFiles:      mappend(pkg.CFiles, pkg.CXXFiles, pkg.MFiles, pkg.HFiles, pkg.SFiles, pkg.SwigFiles, pkg.SwigCXXFiles, pkg.SysoFiles),
		EmbedPatterns:   pkg.EmbedPatterns,
		Imports:         make(map[string]*packages.Package, len(pkg.Imports)),
	}
	for i, file := range pkg.GoFiles {
		if subrepo != "" {
			// this is fairly nasty... there must be a better way of getting it without the pkg/ prefix
			log.Printf("here %s | %s | %s | %s", subrepo, pkg.Dir, file, module)
			dir := strings.TrimPrefix(pkg.Dir, "pkg/"+runtime.GOOS+"_"+runtime.GOARCH)
			dir = strings.TrimPrefix(strings.TrimPrefix(dir, "/"), module)
			p.GoFiles[i] = filepath.Join(subrepo, dir, file)
			p.CompiledGoFiles[i] = filepath.Join(pkg.Dir, file) // Stash this here for later
		} else {
			p.GoFiles[i] = filepath.Join(pkg.Dir, file)
			p.CompiledGoFiles[i] = filepath.Join(pkg.Dir, file)
		}
	}
	for _, imp := range pkg.Imports {
		p.Imports[imp] = &packages.Package{ID: imp, PkgPath: imp}
	}
	return p
}

// mappend appends multiple slices together.
func mappend(s []string, args ...[]string) []string {
	for _, arg := range args {
		s = append(s, arg...)
	}
	return s
}

// loadImportConfig reads the given importconfig file and produces a map of package name -> export path
func loadImportConfig(filename string) (map[string]string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	m := make(map[string]string, len(lines))
	for _, line := range lines {
		if strings.HasPrefix(line, "packagefile ") {
			pkg, exportFile, found := strings.Cut(strings.TrimPrefix(line, "packagefile "), "=")
			if !found {
				return nil, fmt.Errorf("unknown syntax for line: %s", line)
			}
			m[pkg] = exportFile
		}
	}
	return m, nil
}

// modulePath returns the import path for a module, or the given one if the module isn't set.
func modulePath(module, importPath string) string {
	if module == "" {
		return importPath
	}
	before, _, _ := strings.Cut(module, "@")
	return before
}
