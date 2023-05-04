// Package packageinfo writes information about Go packages in a JSON format.
// This is created at build time and intended to be consumed by the gopackagedriver binary.
package packageinfo

import (
	"encoding/json"
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/gcexportdata"
	"golang.org/x/tools/go/packages"
)

// WritePackageInfo writes a series of package info files to the given file.
func WritePackageInfo(modulePath, strip, src, exportsFile string, complete bool, w io.Writer) error {
	// Discover all Go files in the module
	goFiles := map[string][]string{}
	if err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if name := d.Name(); name == "testdata" {
			return filepath.SkipDir // Don't descend into testdata
		} else if strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			dir := filepath.Dir(path)
			goFiles[dir] = append(goFiles[dir], path)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to read module dir: %w", err)
	}
	pkgs := make([]*packages.Package, 0, len(goFiles))
	for dir := range goFiles {
		pkgDir := strings.TrimPrefix(strings.TrimPrefix(dir, strip), "/")
		pkg, err := createPackage(filepath.Join(modulePath, pkgDir), dir)
		if _, ok := err.(*build.NoGoError); ok {
			continue // Don't really care, this happens sometimes for modules
		} else if err != nil {
			return fmt.Errorf("failed to import directory %s: %w", dir, err)
		}
		pkgs = append(pkgs, pkg)
	}
	if exportsFile != "" {
		if err := writeExports(exportsFile, pkgs, complete); err != nil {
			return fmt.Errorf("failed to write export data: %w", err)
		}
	}
	// Ensure output is deterministic
	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].ID < pkgs[j].ID
	})
	return serialise(pkgs, w)
}

func createPackage(pkgPath, pkgDir string) (*packages.Package, error) {
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
	return FromBuildPackage(bpkg), nil
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

// writeExports writes the gc export data, which is needed by the package driver for at least some requests.
// This may be inefficient - we may be re-parsing files that build.ImportDir has already done, but
// it doesn't seem to be possible to get the AST info back from build.whatever :(
func writeExports(filename string, pkgs []*packages.Package, complete bool) error {
	// We have to build a parallel set of Package objects of a different flavour, because it'd be
	// too easy if all the different packages we need to use agreed on the same types.
	// We get to loop over them all several times to ensure we can populate everything fully
	// (especially imports) before writing anything out to disk.
	tpkgs := make([]*types.Package, len(pkgs))
	m := make(map[string]*types.Package, len(pkgs))
	for i, pkg := range pkgs {
		p := types.NewPackage(pkg.PkgPath, pkg.Name)
		if complete {
			p.MarkComplete()
		}
		tpkgs[i] = p
		m[pkg.PkgPath] = p
	}
	for i, pkg := range pkgs {
		imports := make([]*types.Package, 0, len(pkg.Imports))
		for name := range pkg.Imports {
			if imp, present := m[name]; present {
				imports = append(imports, imp)
			} else {
				imports = append(imports, types.NewPackage(name, filepath.Base(name)))
			}
		}
		sort.Slice(imports, func(i, j int) bool { return imports[i].Path() < imports[j].Path() })
		tpkgs[i].SetImports(imports)
	}
	for i, pkg := range pkgs {
		if err := writeExport(filename, pkg, tpkgs[i]); err != nil {
			return err
		}
	}
	return nil
}

func writeExport(dirname string, pkg *packages.Package, tpkg *types.Package) error {
	fset := token.NewFileSet()
	for _, file := range pkg.GoFiles {
		if _, err := parser.ParseFile(fset, file, nil, 0); err != nil {
			return fmt.Errorf("failed to parse %s: %w", file, err)
		}
	}
	filename := filepath.Join(dirname, pkg.PkgPath+".gc")
	if err := os.MkdirAll(filepath.Dir(filename), 0775); err != nil {
		return fmt.Errorf("failed to make directory: %w", err)
	}
	pkg.ExportFile = filename
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()
	return gcexportdata.Write(f, fset, tpkg)
}
