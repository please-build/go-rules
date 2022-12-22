// Package packageinfo writes information about Go packages in a JSON format.
// This is created at build time and intended to be consumed by the gopackagedriver binary.
package packageinfo

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/scanner"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/please-build/go-rules/tools/please_go/embed"
)

// WritePackageInfo writes package info to the given file.
func WritePackageInfo(importPath, plzPkg string, goSrcs []string, embedCfg string, w io.Writer) error {
	pkg, err := createPackage(filepath.Join(importPath, plzPkg), goSrcs, embedCfg)
	if err != nil {
		return err
	}
	return serialise([]*packages.Package{pkg}, w)
}

// WriteModuleInfo writes a series of package info files to the given file.
func WriteModuleInfo(modulePath, src string, w io.Writer) error {
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
	for dir, files := range goFiles {
		pkg, err := createPackage(filepath.Join(modulePath, dir), files, "")
		if err != nil {
			return err
		}
		pkgs = append(pkgs, pkg)
	}
	// Ensure output is deterministic
	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].ID < pkgs[j].ID
	})
	return serialise(pkgs, w)
}

func createPackage(pkgPath string, goSrcs []string, embedCfg string) (*packages.Package, error) {
	pkg := &packages.Package{
		// TODO(peterebden): This should really be the build label, but that won't work for go_module where it's not unique.
		ID:              pkgPath,
		PkgPath:         pkgPath,
		Fset:            token.NewFileSet(),
		GoFiles:         goSrcs,
		CompiledGoFiles: goSrcs, // For our purposes these are the same as GoFiles (I think)
		Imports:         map[string]*packages.Package{},
	}
	for _, src := range goSrcs {
		f, err := parser.ParseFile(pkg.Fset, src, nil, parser.SkipObjectResolution|parser.ParseComments)
		if err != nil {
			// Try to continue if there are just syntax errors; they are not our problem.
			if l, ok := err.(scanner.ErrorList); ok {
				for _, err := range l {
					pkg.Errors = append(pkg.Errors, packages.Error{
						Pos:  err.Pos.String(),
						Msg:  err.Msg,
						Kind: packages.ParseError,
					})
				}
			} else {
				return nil, fmt.Errorf("failed to parse %s: %w", src, err)
			}
		}
		pkg.Syntax = append(pkg.Syntax, f)
	}
	// Store its imports now, although we just enter something minimal (we can't sensibly link up package objects
	// here; it's kind of a weird structure anyway since you'd obviously not serialise as-is)
	for _, file := range pkg.Syntax {
		for _, imp := range file.Imports {
			ipath := strings.Trim(imp.Path.Value, `"`)
			pkg.Imports[ipath] = &packages.Package{ID: ipath}
		}
		pkg.Name = file.Name.Name
	}
	// If there's an embed config, load that info.
	if embedCfg != "" {
		f, err := os.Open(embedCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to read embed config: %w", err)
		}
		defer f.Close()
		cfg := embed.Cfg{}
		if err := json.NewDecoder(f).Decode(&cfg); err != nil {
			return nil, fmt.Errorf("failed to load embed config: %w", err)
		}
		for pattern := range cfg.Patterns {
			pkg.EmbedPatterns = append(pkg.EmbedPatterns, pattern)
		}
		for _, file := range cfg.Files {
			pkg.EmbedFiles = append(pkg.EmbedFiles, file)
		}
		// Both of these are nondeterminstic and need to be ordered.
		sort.Strings(pkg.EmbedPatterns)
		sort.Strings(pkg.EmbedFiles)
	}
	return pkg, nil
}

func serialise(pkgs []*packages.Package, w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(pkgs)
}
