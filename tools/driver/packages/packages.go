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
	"sync"

	"github.com/peterebden/go-cli-init/v5/logging"
	"golang.org/x/sync"
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
	files, err := directoriesToFiles(files)
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
	pkgs, err := loadPackageInfo(relFiles)
	if err != nil {
		return nil, err
	}
	// Handle stdlib imports which are not currently done elsewhere.
	// TODO(peterebden): This is not really correct, plz could be supplying it. However it's not
	//                   trivial to stitch in sensibly there, and doesn't work in this repo...
	goroot := exec.Command("go", "env", "GOROOT")
	goroot.Stderr = &bytes.Buffer{}
	gorootPath, err := goroot.Output()
	if err != nil {
		return nil, handleSubprocessErr(goroot, err)
	}
	// Now stitch in whatever stdlib packages are needed
	m := make(map[string]*packages.Package, len(pkgs)+50)
	for _, pkg := range pkgs {
		m[pkg.PkgPath] = pkg
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			for _, imp := range file.Imports {
				importPath := strings.Trim(imp.Path.Value, `"`)
				if p, present := m[importPath]; present {
					// This isn't _really_ necessary (it gets flattened back to much of what it
					// was before) but there isn't really a good reason not to either.
					pkg.Imports[importPath] = p
				} else if strings.Contains(importPath, ".") {
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
	return &DriverResponse{
		Sizes: &types.StdSizes{
			// These are obviously hardcoded. To worry about later.
			WordSize: 8,
			MaxAlign: 8,
		},
		Packages: pkgs,
		// TODO(peterebden): Roots
	}, nil
}

// loadPackageInfo loads all the package information by executing Please.
// A cooler way of handling this in future would be to do this in-process; for that we'd
// need to define the SDK we keep talking about as a supported programmatic interface.
func loadPackageInfo(files []string) ([]*packages.Package, error) {
	r1, w1, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	// N.B. Deliberate not to defer close of the readers here, we need to close them once processes
	//      finish (otherwise the downstream ones keep trying to read them forever).
	defer w1.Close()
	r2, w2, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer w2.Close()
	whatinputs := exec.Command("plz", append([]string{"query", "whatinputs"}, files...)...)
	whatinputs.Stderr = &bytes.Buffer{}
	whatinputs.Stdout = w1
	deps := exec.Command("plz", "query", "deps", "-", "--hidden", "--include", "go_pkg_info")
	deps.Stdin = r1
	deps.Stderr = &bytes.Buffer{}
	deps.Stdout = w2
	build := exec.Command("plz", "build", "-")
	build.Stdin = r2
	build.Stderr = &bytes.Buffer{}
	build.Stdout = &bytes.Buffer{}
	if err := whatinputs.Start(); err != nil {
		return nil, err
	} else if err := deps.Start(); err != nil {
		return nil, err
	} else if err := build.Start(); err != nil {
		return nil, err
	}
	if err := whatinputs.Wait(); err != nil {
		return nil, handleSubprocessErr(whatinputs, err)
	}
	r1.Close()
	if err := deps.Wait(); err != nil {
		return nil, handleSubprocessErr(deps, err)
	}
	r2.Close()
	if err := build.Wait(); err != nil {
		return nil, handleSubprocessErr(build, err)
	}
	// Now we can read all the package info files from the build process' stdout.
	pkgs := []*packages.Package{}
	var lock sync.Mutex
	var g errgroup.Group
	g.SetLimit(8) // arbitrary limit since we're doing I/O
	for _, file := range strings.Fields(strings.TrimSpace(build.Stdout.(*bytes.Buffer).String())) {
		file := file
		g.Go(func() error {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			lpkgs := []*packages.Package{}
			if err := json.NewDecoder(f).Decode(&lpkgs); err != nil {
				return err
			}
			lock.Lock()
			defer lock.Unlock()
			pkgs = append(pkgs, lpkgs...)
			return nil
		})
	}
	return pkgs, g.Wait()
}

// allGoFilesInDir returns all the files ending in a particular suffix in a given directory.
func allGoFilesInDir(dirname string) []string {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		log.Error("Failed to read directory %s: %s", dirname, err)
		return nil
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if name := entry.Name(); strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			files = append(files, name)
		}
	}
	return files
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
	for _, f := range allGoFilesInDir(dir) {
		pkg.GoFiles = append(pkg.GoFiles, filepath.Join(dir, f))
	}
	if len(pkg.GoFiles) == 0 {
		return nil // Assume this wasn't actually a stdlib package
	}
	pkg.CompiledGoFiles = pkg.GoFiles
	parseFiles(pkg)
	return pkg
}
