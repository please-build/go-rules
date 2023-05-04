package packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/peterebden/go-cli-init/v5/logging"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
	"golang.org/x/tools/go/gcexportdata"
	"golang.org/x/tools/go/packages"

	"github.com/please-build/go-rules/tools/please_go/packageinfo"
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
	// If there are no files provided, do nothing.
	if len(files) == 0 {
		return &DriverResponse{NotHandled: true}, nil
	}
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
	files, err := directoriesToFiles(files, req.Tests)
	if err != nil {
		return nil, err
	} else if len(files) == 0 {
		// Not obvious that this really is an error case.
		log.Warning("No Go files found in initial query")
		return &DriverResponse{NotHandled: true}, nil
	} else if err := os.Chdir(filepath.Dir(files[0])); err != nil {
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
	// Now turn these back into the set of original directories; we use these to determine roots later
	dirs := map[string]struct{}{}
	for _, file := range relFiles {
		dirs[filepath.Dir(file)] = struct{}{}
	}
	pkgs, err := loadPackageInfo(relFiles, (req.Mode&packages.NeedTypesInfo) != 0)
	if err != nil {
		return nil, err
	}
	// Build the set of root packages
	seenRoots := map[string]struct{}{}
	roots := []string{}
	for _, pkg := range pkgs {
		for _, file := range pkg.GoFiles {
			if _, present := dirs[filepath.Dir(file)]; present {
				if _, present := seenRoots[pkg.ID]; !present {
					seenRoots[pkg.ID] = struct{}{}
					roots = append(roots, pkg.ID)
				}
			}
		}
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
		pkg.CompiledGoFiles = pkg.GoFiles
		pkg.ExportFile = filepath.Join(rootpath, "plz-out/gen", pkg.ExportFile)
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
		Roots:    roots,
	}, nil
}

// loadPackageInfo loads all the package information by executing Please.
// A cooler way of handling this in future would be to do this in-process; for that we'd
// need to define the SDK we keep talking about as a supported programmatic interface.
func loadPackageInfo(files []string, needTypes bool) ([]*packages.Package, error) {
	isTerminal := term.IsTerminal(int(os.Stderr.Fd()))
	plz := func(args ...string) *exec.Cmd {
		cmd := exec.Command("plz", args...)
		if !isTerminal {
			cmd.Stderr = &bytes.Buffer{}
		} else {
			cmd.Stderr = os.Stderr
		}
		return cmd
	}

	r1, w1, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	r2, w2, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	// N.B. deliberate not to close these here, they happen exactly when needed.
	whatinputs := plz(append([]string{"query", "whatinputs"}, files...)...)
	whatinputs.Stdout = w1
	deps := plz("query", "deps", "-", "--hidden", "--include", "go_pkg_info", "--include", "go_src")
	deps.Stdin = r1
	deps.Stdout = w2
	build := plz("build", "-")
	build.Stdin = r2
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
		if !strings.HasSuffix(file, ".json") {
			continue // Ignore all the various Go sources etc.
		}
		g.Go(func() error {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			lpkgs := []*packages.Package{}
			if err := json.NewDecoder(f).Decode(&lpkgs); err != nil {
				return fmt.Errorf("failed to decode package info from %s: %s", file, err)
			}
			// Update the ExportFile paths which are relative
			for _, pkg := range lpkgs {
				pkg.ExportFile = filepath.Join(filepath.Dir(file), pkg.ExportFile)
				if needTypes {
					// If we need type info, we need to get it from a neighbouring gc_exports file
					// This 'just knows about' the filenames we define in go.build_defs
					filename := filepath.Join(strings.TrimSuffix(file, "_pkg_info.json")+"_gc_exports", pkg.Name)
					if types, err := loadGCExportData(filename, pkg.PkgPath); err != nil {
						log.Warning("failed to load GC export data: %s", err)
					} else {
						pkg.Types = types
					}
				}
			}
			lock.Lock()
			defer lock.Unlock()
			pkgs = append(pkgs, lpkgs...)
			return nil
		})
	}
	return pkgs, g.Wait()
}

func loadGCExportData(filename, pkgPath string) (*types.Package, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fset := token.NewFileSet()
	imports := map[string]*types.Package{}
	return gcexportdata.Read(f, fset, imports, pkgPath)
}

// loadStdlibPackages returns all the packages from the Go stdlib.
// TODO(peterebden): This is very much temporary, we should ideally be able to get this from
// a plz target as well (especially for go_toolchain)
func loadStdlibPackages() ([]*packages.Package, error) {
	// We just list the entire stdlib set, it's not worth trying to filter it right now.
	log.Debug("Loading stdlib packages...")
	cmd := exec.Command("go", "list", "-json", "std")
	cmd.Stderr = &bytes.Buffer{}
	cmd.Stdout = &bytes.Buffer{}
	if err := cmd.Run(); err != nil {
		return nil, handleSubprocessErr(cmd, err)
	}
	pkgs := []*packages.Package{}
	d := json.NewDecoder(cmd.Stdout.(*bytes.Buffer))
	for {
		pkg := &build.Package{}
		if err := d.Decode(pkg); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, packageinfo.FromBuildPackage(pkg))
	}
	return pkgs, nil
}

func handleSubprocessErr(cmd *exec.Cmd, err error) error {
	if buf, ok := cmd.Stderr.(*bytes.Buffer); ok {
		return fmt.Errorf("%s Stdout:\n%s", err, buf.String())
	}
	// If it's not a buffer, it was probably stderr, so the user has already seen it.
	return err
}

// directoriesToFiles expands any directories in the given list to files in that directory.
func directoriesToFiles(in []string, includeTests bool) ([]string, error) {
	files := make([]string, 0, len(in))
	for _, x := range in {
		if info, err := os.Stat(x); err != nil {
			return nil, err
		} else if info.IsDir() {
			for _, f := range allGoFilesInDir(x, includeTests) {
				files = append(files, filepath.Join(x, f))
			}
		} else {
			files = append(files, x)
		}
	}
	return files, nil
}

// allGoFilesInDir returns all the files ending in a particular suffix in a given directory.
func allGoFilesInDir(dirname string, includeTests bool) []string {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		log.Error("Failed to read directory %s: %s", dirname, err)
		return nil
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if name := entry.Name(); strings.HasSuffix(name, ".go") && (includeTests || !strings.HasSuffix(name, "_test.go")) {
			files = append(files, name)
		}
	}
	return files
}
