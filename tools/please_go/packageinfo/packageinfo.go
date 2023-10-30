// Package packageinfo writes information about Go packages in a JSON format.
// This is created at build time and intended to be consumed by the gopackagedriver binary.
package packageinfo

import (
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// WriteModuleInfo writes a series of package info files to the given file.
func WriteModuleInfo(srcRoot, module, importConfigPath string, pkgWildcards []string, w io.Writer) error {
	importConfig, err := loadImportConfig(importConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read importConfigPath: %w", err)
	}

	var pkgs []*packages.Package
	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		if name := d.Name(); name == "testdata" {
			return filepath.SkipDir // Don't descend into testdata
		}

		path, err = filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}

		pkgID := filepath.Join(module, path)
		pkg, err := createPackage(pkgID, filepath.Join(srcRoot, path), importConfig[pkgID])
		if err != nil {
			// Ignore this for wildcards.
			if _, ok := err.(*build.NoGoError); ok {
				return nil
			}
			return err
		}
		pkgs = append(pkgs, pkg)
		return nil
	}

	// Default to the whole module
	if len(pkgWildcards) == 0 {
		pkgWildcards = []string{"..."}
	}

	for _, pkg := range pkgWildcards {
		if strings.Contains(pkg, "...") {
			pkg = strings.TrimSuffix(pkg, "...")
			if err := filepath.WalkDir(filepath.Join(srcRoot, pkg), walkDirFunc); err != nil {
				return fmt.Errorf("failed to read module dir: %w", err)
			}
		} else {
			pkgID := filepath.Join(module, pkg)
			pkg, err := createPackage(pkgID, filepath.Join(srcRoot, pkg), importConfig[pkgID])
			if err != nil {
				return err
			}
			pkgs = append(pkgs, pkg)
		}
	}

	// Ensure output is deterministic
	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].ID < pkgs[j].ID
	})
	return serialise(pkgs, w)
}

// WritePackageInfo writes info for a single package
func WritePackageInfo(pkgDir, importPath string, exportFile string, w io.Writer) error {
	pkg, err := createPackage(importPath, pkgDir, exportFile)
	if err != nil {
		return err
	}

	return serialise([]*packages.Package{pkg}, w)
}

func createPackage(importPath, pkgDir, exportFile string) (*packages.Package, error) {
	buildPkg, err := build.ImportDir(pkgDir, build.ImportComment)
	if err != nil {
		return nil, err
	}
	buildPkg.ImportPath = importPath
	pkg := FromBuildPackage(buildPkg)
	pkg.ExportFile = exportFile
	return pkg, nil
}

func serialise(pkgs []*packages.Package, w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(pkgs)
}

// FromBuildPackage creates a packages Package from a build Package.
func FromBuildPackage(pkg *build.Package) *packages.Package {
	p := &packages.Package{
		ID:            pkg.ImportPath,
		Name:          pkg.Name,
		PkgPath:       pkg.ImportPath,
		GoFiles:       make([]string, len(pkg.GoFiles)),
		OtherFiles:    mappend(pkg.CFiles, pkg.CXXFiles, pkg.MFiles, pkg.HFiles, pkg.SFiles, pkg.SwigFiles, pkg.SwigCXXFiles, pkg.SysoFiles),
		EmbedPatterns: pkg.EmbedPatterns,
		Imports:       make(map[string]*packages.Package, len(pkg.Imports)),
	}
	for i, file := range pkg.GoFiles {
		p.GoFiles[i] = filepath.Join(pkg.Dir, file)
	}
	p.CompiledGoFiles = p.GoFiles // This seems to be important to e.g. gosec
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
