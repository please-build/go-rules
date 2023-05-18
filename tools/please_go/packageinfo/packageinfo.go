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

// WritePackageInfo writes a series of package info files to the given file.
func WritePackageInfo(modulePath, strip, src, importconfig string, imports map[string]string, installPkgs map[string]struct{}, complete bool, w io.Writer) error {
	logFile, err := os.Create("packageinfo.log")
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close()
	logFile.WriteString("FYI modulePath is " + modulePath + "\n")
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
	if importconfig != "" {
		m, err := loadImportConfig(importconfig)
		if err != nil {
			return fmt.Errorf("failed to read importconfig: %w", err)
		}
		imports = m
	}
	pkgs := make([]*packages.Package, 0, len(goFiles))
	for dir := range goFiles {
		pkgDir := strings.TrimPrefix(strings.TrimPrefix(dir, strip), "/")
		// modulePath should be golang.org/x/xerrors
		// pkgDir should be internal
		// dir should be third_party/go/xerrors/internal
		logFile.WriteString("modulePath=" + modulePath + " . it should be golang.org/x/xerrors\n")
		logFile.WriteString("pkgDir=" + pkgDir + " . it should be internal\n")
		logFile.WriteString("dir=" + dir + " . it should be third_party/go/xerrors/internal\n")
		pkg, err := createPackage(filepath.Join(modulePath, pkgDir), dir)
		if _, ok := err.(*build.NoGoError); ok {
			continue // Don't really care, this happens sometimes for modules
		} else if err != nil {
			return fmt.Errorf("failed to import directory %s: %w", dir, err)
		}
		pkg.ExportFile = imports[pkg.PkgPath]
		if pkg.ID == "golang.org/x/xerrors" {
			logFile.WriteString("pkg.ID is golang.org/x/xerrors\n")
			logFile.WriteString("pkg.ExportFile is " + pkg.ExportFile + "\n")
		}
		pkgs = append(pkgs, pkg)
	}
	// Ensure output is deterministic
	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].ID < pkgs[j].ID
	})
	return serialise(pkgs, w)
}

func createPackage(pkgPath, pkgDir string) (*packages.Package, error) {
	// open a logfile to append to
	logFile, err := os.OpenFile("packageinfo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logFile.WriteString("in createPackage with " + pkgPath + " and " + pkgDir + "\n")
	if pkgDir == "" || pkgDir == "." {
		// This happens when we're in the repo root, ImportDir refuses to read it for some reason.
		path, err := filepath.Abs(pkgDir)
		if err != nil {
			return nil, err
		}
		logFile.WriteString("got path " + path + "\n")
		pkgDir = path
	}
	logFile.WriteString("calling ImportDir with pkgDir " + pkgDir + "\n")
	bpkg, err := build.ImportDir(pkgDir, build.ImportComment)
	if err != nil {
		return nil, err
	}
	logFile.WriteString("got bpkg " + bpkg.ImportPath + "\n")
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
