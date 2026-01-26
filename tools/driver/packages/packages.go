package packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"go/types"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/peterebden/go-cli-init/v5/logging"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
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
	// Now turn these back into the set of original directories; we use these to determine roots later
	dirs := map[string]struct{}{}
	for _, file := range files {
		dirs[filepath.Dir(file)] = struct{}{}
	}
	pkgs, err := loadPackageInfo(files, req.Mode)
	if err != nil {
		return nil, err
	}
	return packagesToResponse(rootpath, pkgs, dirs)
}

// LoadOffline is like Load but rather than querying plz to find the file to load, it just
// walks a file tree looking for pkg_info.json files.
func LoadOffline(req *DriverRequest, searchDir string, files []string) (*DriverResponse, error) {
	pkgs := []*packages.Package{}
	if err := filepath.WalkDir(searchDir, func(path string, d fs.DirEntry, err error) error {
		lpkgs := []*packages.Package{}
		if err != nil || d.IsDir() || !strings.HasSuffix(path, "pkg_info.json") {
			return err
		} else if b, err := os.ReadFile(path); err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		} else if err := json.Unmarshal(b, &lpkgs); err != nil {
			return fmt.Errorf("failed to decode %s: %w", path, err)
		}
		for _, pkg := range lpkgs {
			if _, after, found := strings.Cut(pkg.ExportFile, "|"); found {
				pkg.ExportFile = after
			}
			pkg.GoFiles = pkg.CompiledGoFiles
		}
		pkgs = append(pkgs, lpkgs...)
		return nil
	}); err != nil {
		return nil, err
	}
	dirs := map[string]struct{}{}
	for _, file := range files {
		// Inputs that are files need to get turned into their package directory
		if info, err := os.Stat(file); err == nil && !info.IsDir() {
			file = filepath.Dir(file)
		}
		dirs[file] = struct{}{}
	}
	log.Debug("Generating response for %s", dirs)
	return packagesToResponse(searchDir, pkgs, dirs)
}

func packagesToResponse(rootpath string, pkgs []*packages.Package, dirs map[string]struct{}) (*DriverResponse, error) {
	// Build the set of root packages
	seenRoots := map[string]struct{}{}
	roots := []string{}
	seenRuntime := false
	for _, pkg := range pkgs {
		if _, present := dirs[pkg.PkgPath]; present {
			seenRoots[pkg.ID] = struct{}{}
			roots = append(roots, pkg.ID)
			continue
		}
		for _, file := range pkg.GoFiles {
			if _, present := dirs[filepath.Dir(file)]; present {
				if _, present := seenRoots[pkg.ID]; !present {
					seenRoots[pkg.ID] = struct{}{}
					roots = append(roots, pkg.ID)
				}
			}
			if pkg.ID == "runtime" {
				seenRuntime = true
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
			if path := filepath.Join(rootpath, "plz-out/subrepos", file); pathExists(path) {
				pkg.GoFiles[i] = path
			} else if path := filepath.Join(rootpath, "plz-out/gen", file); pathExists(path) {
				pkg.GoFiles[i] = path
			} else {
				pkg.GoFiles[i] = filepath.Join(rootpath, file)
			}
		}
		pkg.CompiledGoFiles = pkg.GoFiles
		pkg.ExportFile = filepath.Join(rootpath, pkg.ExportFile)
	}
	if !seenRuntime {
		// Handle stdlib imports if we didn't already find them
		// TODO(peterebden): Get rid of loading stdlib here once we do v2 and require it to be
		//                   explicitly defined as a build rule
		stdlib, err := loadStdlibPackages()
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, stdlib...)
	}
	log.Debug("Built package set")
	return &DriverResponse{
		Sizes: &types.StdSizes{
			// These are obviously hardcoded. To worry about later.
			WordSize: 8,
			MaxAlign: 8,
		},
		Packages: pkgs,
		Roots:    roots,
	}, nil
}

// loadPackageInfo loads all the package information by executing Please.
// A cooler way of handling this in future would be to do this in-process; for that we'd
// need to define the SDK we keep talking about as a supported programmatic interface.
func loadPackageInfo(files []string, mode packages.LoadMode) ([]*packages.Package, error) {
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
	whatinputs := plz(append([]string{"query", "whatinputs", "--ignore_unknown" }, files...)...)
	whatinputs.Stdout = w1
	args := []string{"query", "deps", "-", "--hidden", "-i", "go_pkg_info", "-i", "go_src"}
	if (mode & packages.NeedExportFile) != 0 {
		args = append(args, "-i", "go")
	}
	deps := plz(args...)
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
	// Now we can read all the package info files from the build process' stdout.
	return loadPackageInfoFiles(strings.Fields(strings.TrimSpace(build.Stdout.(*bytes.Buffer).String())))
}

// loadPackageInfoFiles loads the given set of package info files
func loadPackageInfoFiles(paths []string) ([]*packages.Package, error) {
	log.Debug("Loading plz package info files...")
	pkgs := []*packages.Package{}
	var lock sync.Mutex
	var g errgroup.Group
	g.SetLimit(8) // arbitrary limit since we're doing I/O
	for _, file := range paths {
		file := file
		if !strings.HasSuffix(file, ".json") {
			continue // Ignore all the various Go sources etc.
		}
		log.Debug("Package file: %s", file)
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
			// Update the ExportFile paths to include the generated prefix
			for _, pkg := range lpkgs {
				// Undo the hack from packageinfo.go
				before, _, _ := strings.Cut(pkg.ExportFile, "|")
				pkg.ExportFile = filepath.Join("plz-out/gen", before)
				pkg.CompiledGoFiles = pkg.GoFiles
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
// a plz target as well (especially for go_toolchain)
func loadStdlibPackages() ([]*packages.Package, error) {
	// We just list the entire stdlib set, it's not worth trying to filter it right now.
	log.Debug("Loading stdlib packages...")
	goTool := "go"
	// This is a hack to try to closer match what various plz things do.
	// As noted above, this will go away once we move everything to go_toolchain / go_stdlib.
	if env := os.Getenv("TOOLS_GO"); env != "" {
		goTool = env
	}
	cmd := exec.Command(goTool, "list", "-json", "std")
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
		pkgs = append(pkgs, packageinfo.FromBuildPackage(pkg, "", ""))
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
		if strings.HasSuffix(x, "/...") {
			// We could turn this into a `/...` style thing for plz but we also need to know the
			// directories later to populate the roots correctly.
			if err := filepath.WalkDir(strings.TrimSuffix(x, "/..."), func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				} else if strings.HasSuffix(path, ".go") && (d.Type()&fs.ModeSymlink) == 0 {
					files = append(files, path)
				}
				return nil
			}); err != nil {
				return nil, err
			}
		} else if info, err := os.Stat(x); err != nil {
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

func pathExists(filename string) bool {
	_, err := os.Lstat(filename)
	return err == nil
}
