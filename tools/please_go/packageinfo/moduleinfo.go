package packageinfo

import (
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"io/fs"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// WriteModuleInfo writes a series of package info files to the given file.
func WriteModuleInfo(importPath string, srcRoot, importconfig string, imports map[string]string, installPkgs []string, subrepo, module string, includeTests bool, w io.Writer) error {
	// Discover all Go files in the module
	goFiles := map[string][]string{}
	module = modulePath(module, importPath)

	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if name := d.Name(); name == "testdata" {
			return filepath.SkipDir // Don't descend into testdata
		} else if strings.HasSuffix(name, ".go") && (includeTests || !strings.HasSuffix(name, "_test.go")) {
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
	// If we're doing the stdlib, limit it to just things in the importconfig (i.e. no cmd/ packages)
	if importconfig != "" {
		pkgs = slices.DeleteFunc(pkgs, func(pkg *packages.Package) bool {
			_, present := imports[pkg.PkgPath]
			return !present
		})
	}
	// Vendor packages. They aren't identified by the original imports but we know what they are now.
	vendorised := map[string]*packages.Package{}
	for _, pkg := range pkgs {
		if strings.HasPrefix(pkg.PkgPath, "vendor/") {
			vendorised[strings.TrimPrefix(pkg.PkgPath, "vendor/")] = pkg
		}
	}
	for _, pkg := range pkgs {
		for k := range pkg.Imports {
			if v, present := vendorised[k]; present {
				pkg.Imports[k] = v
			}
		}
	}
	// Ensure output is deterministic
	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].ID < pkgs[j].ID
	})
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(pkgs)
}
