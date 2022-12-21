package packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/peterebden/go-cli-init/v5/logging"
	"golang.org/x/tools/go/packages"
)

var log = logging.MustGetLogger()

// DriverRequest is copied from go/packages; it's the format of config that we get on stdin.
type DriverRequest struct {
	Mode       packages.LoadMode `json:"mode"`
	Env        []string          `json:"env"`
	BuildFlags []string          `json:"build_flags"`
	Tests      bool              `json:"tests"`
	Overlay    map[string][]byte `json:"overlay"`
}

// DriverResponse is copied from go/packages; it's our response about packages we've loaded.
type DriverResponse struct {
	NotHandled bool
	Sizes      *types.StdSizes
	Roots      []string `json:",omitempty"`
	Packages   []*packages.Package
}

// Load reads a set of packages and returns information about them.
// Most of the request structure isn't honoured at the moment.
func Load(req *DriverRequest, files []string) (*DriverResponse, error) {
	// We need to find the plz repo that we need to be in (we might be invoked from outside it)
	// For now we're assuming they're all in the same repo (which is probably reasonable) and just
	// take the first one as indicative (which is maybe less so).
	for i, file := range files {
		file, err := filepath.Abs(file)
		if err != nil {
			return nil, err
		}
		files[i] = file
	}
	// Inputs can be either files or directories; here we turn them all into files.
	files, err := directoriesToFiles(files, ".go")
	if err != nil {
		return nil, err
	}
	if err := os.Chdir(filepath.Dir(files[0])); err != nil {
		return nil, err
	}
	reporoot := exec.Command("plz", "query", "reporoot")
	reporoot.Stderr = &bytes.Buffer{}
	root, err := reporoot.Output()
	if err != nil {
		return nil, handleSubprocessErr(reporoot, err)
	}
	rootpath := strings.TrimSpace(string(root))
	if err := os.Chdir(rootpath); err != nil {
		return nil, err
	}
	// Now we have to make the filepaths relative, so plz understands them
	// When https://github.com/thought-machine/please/issues/2618 is resolved, we won't need to do this.
	relFiles := make([]string, len(files))
	for i, file := range files {
		file, err = filepath.Rel(rootpath, file)
		if err != nil {
			return nil, err
		}
		relFiles[i] = file
	}
	// Find the GOROOT for stdlib imports.
	// TODO(peterebden): This is not really correct, plz could be supplying it.
	goroot := exec.Command("go", "env", "GOROOT")
	goroot.Stderr = &bytes.Buffer{}
	gorootPath, err := goroot.Output()
	if err != nil {
		return nil, handleSubprocessErr(goroot, err)
	}
	// We run Please as a subprocess to get the answers here.
	// A cooler way of handling this in future would be to do this in-process; for that we'd
	// need to define the SDK we keep talking about as a supported programmatic interface.
	whatinputs := exec.Command("plz", append([]string{"query", "whatinputs"}, relFiles...)...)
	whatinputs.Stderr = &bytes.Buffer{}
	targetList, err := whatinputs.Output()
	if err != nil {
		return nil, handleSubprocessErr(whatinputs, err)
	}
	// We can't pipe whatinputs directly into plz build because we need the target list to use again.
	build := exec.Command("plz", "build", "-")
	build.Stdin = bytes.NewReader(targetList)
	build.Stderr = &bytes.Buffer{}
	if err := build.Run(); err != nil {
		return nil, handleSubprocessErr(build, err)
	}
	// TODO(peterebden): It would be nicer to connect these two together with a pipe. That seems to be
	//                   hanging when I try it; I'm not sure what is wrong but we should be able to figure it out.
	deps := exec.Command("plz", "query", "deps", "-")
	deps.Stdin = bytes.NewReader(targetList)
	deps.Stderr = &bytes.Buffer{}
	depList, err := deps.Output()
	if err != nil {
		return nil, handleSubprocessErr(deps, err)
	}
	print := exec.Command("plz", "query", "print", "--json", "-")
	print.Stdin = bytes.NewReader(depList)
	print.Stderr = &bytes.Buffer{}
	output, err := print.Output()
	if err != nil {
		return nil, handleSubprocessErr(print, err)
	}
	targets := map[string]*buildTarget{}
	if err := json.Unmarshal(output, &targets); err != nil {
		return nil, err
	}
	for l, t := range targets {
		pkg, name, _ := strings.Cut(l, ":")
		t.Package = strings.TrimPrefix(pkg, "//")
		t.Name = name
		// Hardcoding plz-out/gen obviously isn't generally correct, but it works for all the
		// cases we care about.
		t.OutputDir = filepath.Join(rootpath, "plz-out/gen", t.Package)
		t.InputDir = filepath.Join(rootpath, t.Package)
	}
	m := map[string]struct{}{}
	for _, file := range files {
		m[file] = struct{}{}
	}
	return toResponse(targets, m, strings.TrimSpace(string(gorootPath))), nil
}

// buildTarget is a minimal version of the target output structure from `plz query print`
type buildTarget struct {
	Deps      []string    `json:"deps"`
	Labels    []string    `json:"labels"`
	Srcs      interface{} `json:"srcs"`
	Outs      []string    `json:"outs"`
	Package   string      // Not in the output, but useful to have on here for later
	InputDir  string      // Absolute path to input directory
	OutputDir string      // Absolute path to output directory
	Name      string
}

func handleSubprocessErr(cmd *exec.Cmd, err error) error {
	return fmt.Errorf("%s Stdout:\n%s", err, cmd.Stderr.(*bytes.Buffer).String())
}

// toResponse converts `plz query print` output to a DriverResponse
func toResponse(targets map[string]*buildTarget, originalFiles map[string]struct{}, goroot string) *DriverResponse {
	resp := &DriverResponse{
		Sizes: &types.StdSizes{
			// These are obviously hardcoded. To worry about later.
			WordSize: 8,
			MaxAlign: 8,
		},
	}
	m := map[string]*packages.Package{}
	for label, target := range targets {
		for _, pkg := range target.ToPackages() {
			// TODO(peterebden): We should set pkg.ID to the plz label, but that requires all the go_module
			//                   equivalents to be broken out to unique targets to be sensible.
			pkg.ID = label
			pkg.ID = pkg.PkgPath
			resp.Packages = append(resp.Packages, pkg)
			for _, src := range pkg.GoFiles {
				if _, present := originalFiles[src]; present {
					resp.Roots = append(resp.Roots, pkg.ID)
					break
				}
			}
			m[pkg.PkgPath] = pkg
		}
	}
	// Now we have all the packages, we need to populate the Imports properties
	for _, pkg := range resp.Packages {
		for _, file := range pkg.Syntax {
			for _, imp := range file.Imports {
				importPath := strings.Trim(imp.Path.Value, `"`)
				if p, present := m[importPath]; present {
					pkg.Imports[importPath] = p
				} else if !strings.Contains(importPath, ".") {
					// Looks like a third-party package we _should_ know it but won't if there are missing dependencies or whatever.
					log.Warning("Failed to map import path %s", importPath)
				} else if p := createStdlibImport(importPath, goroot); p != nil {
					pkg.Imports[importPath] = p
					resp.Packages = append(resp.Packages, p)
					m[importPath] = p
				}
			}
		}
	}
	return resp
}

// ToPackages converts this build target to the Go packages it represents.
// There can be an arbitrary number of packages in a target:
//   - 0 if it is not a Go package at all
//   - 1 if it's a go_library
//   - Many if it's a go_module
func (t *buildTarget) ToPackages() []*packages.Package {
	for _, l := range t.Labels {
		if strings.HasPrefix(l, "go_module_path:") {
			return t.toPackages(strings.TrimPrefix(l, "go_module_path:"))
		}
	}
	// Not a module, might be a go_library or whatever
	return t.toPackages("")
}

func (t *buildTarget) toPackages(modulePath string) []*packages.Package {
	ret := []*packages.Package{}
	for _, l := range t.Labels {
		if strings.HasPrefix(l, "go_package:") {
			if pkg := strings.TrimPrefix(l, "go_package:"); strings.HasSuffix(pkg, "/...") {
				// We specify everything under this directory, so we have to discover them all.
				pkg = strings.TrimSuffix(pkg, "/...")
				outDir := filepath.Join(t.OutputDir, t.Name, strings.TrimPrefix(pkg, modulePath))
				for _, dir := range allDirsUnder(outDir) {
					// This isn't efficient, we could have filtered this before and we will call it again, but just deal with it for now
					if len(allFilesInDir(filepath.Join(outDir, dir), ".go")) > 0 {
						ret = append(ret, t.toPackage(modulePath, filepath.Join(pkg, dir)))
					}
				}
			} else {
				ret = append(ret, t.toPackage(modulePath, pkg))
			}
		}
	}
	return ret
}

func (t *buildTarget) toPackage(modulePath, packagePath string) *packages.Package {
	// TODO(peterebden): Work out what we are meant to do with CompiledGoFiles
	pkg := &packages.Package{
		PkgPath: packagePath,
		Imports: map[string]*packages.Package{},
	}
	// This is a hack, the package name isn't necessarily the path name, but we haven't parsed it to know that.
	pkg.Name = packagePath[strings.LastIndexByte(packagePath, '/')+1:]
	// From here on we start making a lot of assumptions about exactly how these things are structured.
	if modulePath != "" {
		// This is part of a go_module.
		outDir := filepath.Join(t.OutputDir, t.Name, strings.TrimPrefix(packagePath, modulePath))
		for _, f := range allFilesInDir(outDir, ".go") {
			pkg.GoFiles = append(pkg.GoFiles, filepath.Join(outDir, f))
		}
		pkg.CompiledGoFiles = pkg.GoFiles
		// TODO(peterebden): Other file types too (EmbedFiles, OtherFiles etc)
		parseFiles(pkg)
		return pkg
	}
	// This is probably a go_library. Deal with its srcs (which are a faff because they can be named or not named, but we know go_library rules will be named)
	// TODO(peterebden): This does not work with sources that are themselves generated...
	if srcs, ok := t.Srcs.(map[string]interface{}); ok {
		if goSrcsI, present := srcs["go"]; present {
			if goSrcs, ok := goSrcsI.([]interface{}); ok {
				for _, src := range goSrcs {
					pkg.GoFiles = append(pkg.GoFiles, filepath.Join(t.InputDir, src.(string)))
				}
			}
		}
	}
	pkg.CompiledGoFiles = pkg.GoFiles
	parseFiles(pkg)
	return pkg
}

// allFilesInDir returns all the files ending in a particular suffix in a given directory.
func allFilesInDir(dirname, suffix string) []string {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		log.Error("Failed to read directory %s: %s", dirname, err)
		return nil
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if name := entry.Name(); strings.HasSuffix(name, suffix) {
			files = append(files, name)
		}
	}
	return files
}

// allDirsUnder returns all directories underneath the given one
func allDirsUnder(dirname string) (dirs []string) {
	if err := filepath.WalkDir(dirname, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if name := d.Name(); d.IsDir() && !strings.HasPrefix(name, ".") {
			dirs = append(dirs, strings.TrimPrefix(strings.TrimPrefix(path, dirname), "/"))
		}
		return nil
	}); err != nil {
		log.Error("failed to read output dir %s", dirname)
	}
	return dirs
}

// parseFiles parses all the Go sources of a package and populates appropriate fields.
func parseFiles(pkg *packages.Package) {
	pkg.Fset = token.NewFileSet()
	for _, src := range pkg.CompiledGoFiles {
		f, err := parser.ParseFile(pkg.Fset, src, nil, parser.SkipObjectResolution|parser.ParseComments)
		if err != nil {
			log.Error("Failed to parse file %s: %s", src, err)
		}
		pkg.Syntax = append(pkg.Syntax, f)
	}
}

// directoriesToFiles expands any directories in the given list to files in that directory.
func directoriesToFiles(in []string, suffix string) ([]string, error) {
	files := make([]string, 0, len(in))
	for _, x := range in {
		if info, err := os.Stat(x); err != nil {
			return nil, err
		} else if info.IsDir() {
			for _, f := range allFilesInDir(x, suffix) {
				files = append(files, filepath.Join(x, f))
			}
		} else {
			files = append(files, x)
		}
	}
	return files, nil
}

// createStdlibImport attempts to create an import for one of the stdlib packages.
func createStdlibImport(path, goroot string) *packages.Package {
	if path == "C" {
		return nil // Cgo isn't a real import, of course.
	}

	pkg := &packages.Package{
		PkgPath: path,
		Imports: map[string]*packages.Package{},
	}
	dir := filepath.Join(goroot, "src", path)
	for _, f := range allFilesInDir(dir, ".go") {
		pkg.GoFiles = append(pkg.GoFiles, filepath.Join(dir, f))
	}
	if len(pkg.GoFiles) == 0 {
		return nil // Assume this wasn't actually a stdlib package
	}
	pkg.CompiledGoFiles = pkg.GoFiles
	parseFiles(pkg)
	return pkg
}
