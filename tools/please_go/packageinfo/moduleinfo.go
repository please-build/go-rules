package packageinfo

import (
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// WriteModuleInfo writes a series of package info files to the given file.
func WriteModuleInfo(importPath string, srcRoot, importconfig string, installPkgs []string, w io.Writer) error {
	// Discover all Go files in the module
	goFiles := map[string][]string{}

	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if name := d.Name(); name == "testdata" {
			return filepath.SkipDir // Don't descend into testdata
		} else if strings.HasSuffix(name, ".go") && (!strings.HasSuffix(name, "_test.go")) {
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
	imports, err := loadImportConfig(importconfig)
	if err != nil {
		return fmt.Errorf("failed to read importconfig: %w", err)
	}

	pkgs := make([]*packages.Package, 0, len(goFiles))
	for dir := range goFiles {
		pkgDir := strings.TrimPrefix(strings.TrimPrefix(dir, srcRoot), "/")
		bpkg, err := buildPackage(filepath.Join(importPath, pkgDir), dir)
		if _, ok := err.(*build.NoGoError); ok {
			continue // Don't really care, this happens sometimes for modules
		} else if err != nil {
			return fmt.Errorf("failed to import directory %s: %w", dir, err)
		}
		pkg := FromBuildPackageForModule(bpkg)

		pkg.ExportFile = imports[pkg.PkgPath]
		pkgs = append(pkgs, pkg)
	}
	// If we're doing the stdlib, limit it to just things in the importconfig (i.e. no cmd/ packages)
	pkgs = slices.DeleteFunc(pkgs, func(pkg *packages.Package) bool {
		_, present := imports[pkg.PkgPath]
		return !present
	})

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

// FromBuildPackageForModule creates a packages Package from a build Package for a module.
func FromBuildPackageForModule(pkg *build.Package) *packages.Package {
	goFiles := slices.Concat(pkg.GoFiles, pkg.TestGoFiles, pkg.XTestGoFiles)
	imports := slices.Concat(pkg.Imports, pkg.TestImports, pkg.XTestImports)
	name := pkg.Name
	id := pkg.ImportPath
	if len(pkg.XTestGoFiles) > 0 || len(pkg.XTestImports) > 0 {
		// In please we may have an external test target and an internal test within the same please package.
		// To ensure they have different go package import paths we appending to the name and id.
		name += "_test"
		id += "_test"
	}
	p := &packages.Package{
		ID:              id,
		Name:            name,
		PkgPath:         id,
		GoFiles:         make([]string, len(goFiles)),
		CompiledGoFiles: make([]string, len(goFiles)),
		OtherFiles:      mappend(pkg.CFiles, pkg.CXXFiles, pkg.MFiles, pkg.HFiles, pkg.SFiles, pkg.SwigFiles, pkg.SwigCXXFiles, pkg.SysoFiles),
		EmbedPatterns:   pkg.EmbedPatterns,
		Imports:         make(map[string]*packages.Package, len(imports)),
	}
	for i, file := range goFiles {
		p.GoFiles[i] = filepath.Join(pkg.Dir, file)
		p.CompiledGoFiles[i] = filepath.Join(pkg.Dir, file)
	}
	for _, imp := range imports {
		p.Imports[imp] = &packages.Package{ID: imp, PkgPath: imp}
	}
	return p
}
