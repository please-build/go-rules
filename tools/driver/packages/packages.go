package packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/peterebden/go-cli-init/v5/logging"
	"golang.org/x/sync/errgroup"
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
	// Now turn these back into the set of original directories; we use these to determine roots later
	dirs := map[string]struct{}{}
	for _, file := range files {
		dirs[filepath.Dir(file)] = struct{}{}
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
	// Make all file paths absolute. Useful absolute paths cannot exist in build actions so we
	// need to rebuild here.
	log.Debug("Checking for file locations...")
	for _, pkg := range pkgs {
		for i, file := range pkg.GoFiles {
			// This is pretty awkward; we need to try to figure out where these files exist now,
			// which isn't particularly clear to the build actions that generated them.
			if _, err := os.Lstat(file); err == nil { // file exists
				pkg.GoFiles[i] = filepath.Join(rootpath, file)
			} else {
				pkg.GoFiles[i] = filepath.Join(rootpath, "plz-out/gen", file)
			}
		}
		// TODO(pebers): Do we need to care about these? go list doesn't seem to populate its
		//               equivalent (although it isn't exactly the same structure)
		pkg.CompiledGoFiles = pkg.GoFiles
	}
	// Handle stdlib imports which are not currently done elsewhere.
	stdlib, err := loadStdlibPackages()
	if err != nil {
		return nil, err
	}
	log.Debug("Built package set")
	return &DriverResponse{
		Sizes: &types.StdSizes{
			// These are obviously hardcoded. To worry about later.
			WordSize: 8,
			MaxAlign: 8,
		},
		Packages: append(pkgs, stdlib...),
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
	r2, w2, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	// N.B. deliberate not to close these here, they happen exactly when needed.
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
	log.Debug("Waiting for plz query whatinputs...")
	if err := whatinputs.Wait(); err != nil {
		return nil, handleSubprocessErr(whatinputs, err)
	}
	w1.Close()
	r1.Close()
	log.Debug("Waiting for plz query deps...")
	if err := deps.Wait(); err != nil {
		return nil, handleSubprocessErr(deps, err)
	}
	w2.Close()
	r2.Close()
	log.Debug("Waiting for plz build...")
	if err := build.Wait(); err != nil {
		return nil, handleSubprocessErr(build, err)
	}
	log.Debug("Loading plz package info files...")
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

// loadStdlibPackages returns all the packages from the Go stdlib.
// TODO(peterebden): This is very much temporary, we should ideally be able to get this from
//                   a plz target as well (especially for go_toolchain)
func loadStdlibPackages() ([]*packages.Package, error) {
	// We just list the entire stdlib set, it's not worth trying to filter it right now.
	log.Debug("Loading stdlib packages...")
	cmd := exec.Command("go", "list", "-json", "std")
	cmd.Stderr = &bytes.Buffer{}
	cmd.Stdout = &bytes.Buffer{}
	if err := cmd.Run(); err != nil {
		return nil, handleSubprocessErr(cmd, err)
	}
	goPkgs := []*goPackage{}
	d := json.NewDecoder(cmd.Stdout.(*bytes.Buffer))
	for {
		pkg := &goPackage{}
		if err := d.Decode(pkg); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		goPkgs = append(goPkgs, pkg)
	}
	pkgs := make([]*packages.Package, len(goPkgs))
	for i, pkg := range goPkgs {
		for i, file := range pkg.GoFiles {
			pkg.GoFiles[i] = filepath.Join(pkg.Dir, file)
		}
		pkgs[i] = &packages.Package{
			ID:              pkg.ImportPath,
			Name:            pkg.Name,
			PkgPath:         pkg.ImportPath,
			GoFiles:         pkg.GoFiles,
			CompiledGoFiles: pkg.CompiledGoFiles,
			OtherFiles:      mappend(pkg.CFiles, pkg.CXXFiles, pkg.MFiles, pkg.HFiles, pkg.SFiles, pkg.SwigFiles, pkg.SwigCXXFiles, pkg.SysoFiles),
			EmbedPatterns:   pkg.EmbedPatterns,
			EmbedFiles:      pkg.EmbedFiles,
			Imports:         map[string]*packages.Package{},
		}
		for _, imp := range pkg.Imports {
			pkgs[i].Imports[imp] = &packages.Package{ID: imp, PkgPath: imp}
		}
	}
	return pkgs, nil
}

// mappend appends multiple slices together.
func mappend(s []string, args ...[]string) []string {
	for _, arg := range args {
		s = append(s, arg...)
	}
	return s
}

// goPackage is a subset of the struct that `go list` outputs (AFAIK this isn't importable)
type goPackage struct {
	Dir             string   // directory containing package sources
	ImportPath      string   // import path of package in dir
	Name            string   // package name
	Root            string   // Go root or Go path dir containing this package
	GoFiles         []string // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles        []string // .go source files that import "C"
	CompiledGoFiles []string // .go files presented to compiler (when using -compiled)
	CFiles          []string // .c source files
	CXXFiles        []string // .cc, .cxx and .cpp source files
	MFiles          []string // .m source files
	HFiles          []string // .h, .hh, .hpp and .hxx source files
	SFiles          []string // .s source files
	SwigFiles       []string // .swig files
	SwigCXXFiles    []string // .swigcxx files
	SysoFiles       []string // .syso object files to add to archive
	EmbedPatterns   []string // //go:embed patterns
	EmbedFiles      []string // files matched by EmbedPatterns
	Imports         []string // import paths used by this package
}

func handleSubprocessErr(cmd *exec.Cmd, err error) error {
	return fmt.Errorf("%s Stdout:\n%s", err, cmd.Stderr.(*bytes.Buffer).String())
}

// directoriesToFiles expands any directories in the given list to files in that directory.
func directoriesToFiles(in []string) ([]string, error) {
	files := make([]string, 0, len(in))
	for _, x := range in {
		if info, err := os.Stat(x); err != nil {
			return nil, err
		} else if info.IsDir() {
			for _, f := range allGoFilesInDir(x) {
				files = append(files, filepath.Join(x, f))
			}
		} else {
			files = append(files, x)
		}
	}
	return files, nil
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
