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
func WritePackageInfo(importPath, pkg string, goSrcs []string, embedCfg string, w io.Writer) error {
	pkgPath := filepath.Join(importPath, pkg)
	pkg := &packages.Package{
		// TODO(peterebden): This should really be the build label, but that won't work for go_module where it's not unique.
		ID:              pkgPath,
		PkgPath:         pkgPath,
		Fset:            token.NewFileSet(),
		GoFiles:         goSrcs,
		CompiledGoFiles: goSrcs,
	}
	for _, src := range GoFiles {
		f, err := parser.ParseFile(pkg.Fset, src, nil, parser.SkipObjectResolution|parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", src, err)
		}
		pkg.Syntax = append(pkg.Syntax, f)
	}
	// Store its imports now, although we just enter nulls (we can't sensibly link up package objects
	// here; it's kind of a weird structure anyway since you'd obviously not serialise as-is)
	for _, file := range pkg.Syntax {
		for _, imp := range file.Imports {
			pkg.Imports[strings.Trim(imp.Path.Value, `"`)] = nil
		}
	}
}
