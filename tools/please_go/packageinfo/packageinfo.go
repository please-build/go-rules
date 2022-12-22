// Package packageinfo reads & writes information about Go packages in a JSON format.
// This is created at build time and intended to be consumed by the gopackagedriver binary.
package packageinfo

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// WritePackageInfo writes package info to the given file.
func WritePackageInfo(importPath, plzPkg string, goSrcs []string, embedCfg string, w io.Writer) error {
	pkgPath := filepath.Join(importPath, plzPkg)
	pkg := &packages.Package{
		// TODO(peterebden): This should really be the build label, but that won't work for go_module where it's not unique.
		ID:      pkgPath,
		PkgPath: pkgPath,
		Fset:    token.NewFileSet(),
		Imports: map[string]*packages.Package{},
	}
	for _, src := range goSrcs {
		f, err := parser.ParseFile(pkg.Fset, src, nil, parser.SkipObjectResolution|parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", src, err)
		}
		pkg.Syntax = append(pkg.Syntax, f)
		// Make GoFiles relative to the package's directory. These are normally absolute in a
		// *packages.Package but we can't determine useful absolute paths inside a build action.
		pkg.GoFiles = append(pkg.GoFiles, strings.TrimPrefix(strings.TrimPrefix(src, plzPkg), "/"))
	}
	// For our purposes these are the same as GoFiles (I think)
	pkg.CompiledGoFiles = pkg.GoFiles
	// Store its imports now, although we just enter something minimal (we can't sensibly link up package objects
	// here; it's kind of a weird structure anyway since you'd obviously not serialise as-is)
	for _, file := range pkg.Syntax {
		for _, imp := range file.Imports {
			ipath := strings.Trim(imp.Path.Value, `"`)
			pkg.Imports[ipath] = &packages.Package{ID: ipath}
		}
		pkg.Name = file.Name.Name
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode([]*packages.Package{pkg})
}
