// Package packageinfo writes information about Go packages in a JSON format.
// This is created at build time and intended to be consumed by the gopackagedriver binary.
package packageinfo

import (
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"io/fs"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// WritePackageInfo writes a series of package info files to the given file.
func WritePackageInfo(importPath string, srcRoot string, imports map[string]string, subrepo, module string, includeTests bool, w io.Writer) error {
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
	if err := filepath.WalkDir(srcRoot, walkDirFunc); err != nil {
		return fmt.Errorf("failed to read module dir: %w", err)
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
		pkg := fromBuildPackage(bpkg, subrepo, module)

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

func buildPackage(
	pkgPath string,
	pkgDir string,
) (*build.Package, error) {
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
	return bpkg, nil
}

func fromBuildPackage(pkg *build.Package, subrepo, module string) *packages.Package {
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
		if subrepo != "" {
			// this is fairly nasty... there must be a better way of getting it without the pkg/ prefix
			dir := strings.TrimPrefix(pkg.Dir, "pkg/"+runtime.GOOS+"_"+runtime.GOARCH)
			dir = strings.TrimPrefix(strings.TrimPrefix(dir, "/"), module)
			p.GoFiles[i] = filepath.Join(subrepo, dir, file)
			p.CompiledGoFiles[i] = filepath.Join(pkg.Dir, file) // Stash this here for later
		} else {
			p.GoFiles[i] = filepath.Join(pkg.Dir, file)
			p.CompiledGoFiles[i] = filepath.Join(pkg.Dir, file)
		}
	}
	for _, imp := range imports {
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

// modulePath returns the import path for a module, or the given one if the module isn't set.
func modulePath(module, importPath string) string {
	if module == "" {
		return importPath
	}
	before, _, _ := strings.Cut(module, "@")
	return before
}
