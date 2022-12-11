package packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/types"
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
	// We run Please as a subprocess to get the answers here.
	// A cooler way of handling this in future would be to do this in-process; for that we'd
	// need to define the SDK we keep talking about as a supported programmatic interface.
	whatinputs := exec.Command("plz", append([]string{"query", "whatinputs"}, files...)...)
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
	reporoot := exec.Command("plz", "query", "reporoot")
	reporoot.Stderr = &bytes.Buffer{}
	root, err := reporoot.Output()
	if err != nil {
		return nil, handleSubprocessErr(reporoot, err)
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
		t.OutputDir = filepath.Join(root, "plz-out/gen", t.Package)
	}
	m := map[string]struct{}{}
	for _, file := range files {
		m[file] = struct{}{}
	}
	return toResponse(targets, m), nil
}

// buildTarget is a minimal version of the target output structure from `plz query print`
type buildTarget struct {
	Deps      []string    `json:"deps"`
	Labels    []string    `json:"labels"`
	Srcs      interface{} `json:"srcs"`
	Outs      []string    `json:"outs"`
	Package   string      // Not in the output, but useful to have on here for later
	OutputDir string      // Absolute path to output directory
	Name      string
}

func handleSubprocessErr(cmd *exec.Cmd, err error) error {
	return fmt.Errorf("%s Stdout:\n%s", err, cmd.Stderr.(*bytes.Buffer).String())
}

// toResponse converts `plz query print` output to a DriverResponse
func toResponse(targets map[string]*buildTarget, originalFiles map[string]struct{}) *DriverResponse {
	resp := &DriverResponse{}
	for label, target := range targets {
		for _, pkg := range target.ToPackages() {
			pkg.ID = label
			resp.Packages = append(resp.Packages, pkg)
			for _, src := range pkg.GoFiles {
				if _, present := originalFiles[src]; present {
					resp.Roots = append(resp.Roots, pkg.Name)
					break
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

func (t *buildTarget) toPackages(module_path string) []*packages.Package {
	ret := []*packages.Package{}
	for _, l := range t.Labels {
		if strings.HasPrefix(l, "go_package:") {
			ret = append(ret, t.toPackage(module_path, strings.TrimPrefix(l, "go_package:")))
		}
	}
	return ret
}

func (t *buildTarget) toPackage(module_path, package_path string) *packages.Package {
	pkg := &packages.Package{
		Path: package_path,
	}
	// This is a hack, the package name isn't necessarily the path name, but we haven't parsed it to know that.
	if idx := strings.LastIndexByte(path, '/'); idx != -1 {
		pkg.Name = path[idx+1:]
	}
	// From here on we start making a lot of assumptions about exactly how these things are structured.
	if module_path != "" {
		// This is part of a go_module.
		outDir := filepath.Join(t.OutputDir, t.Name, strings.TrimPrefix(package_path, module_path))
	}

	if len(t.Outs) == 0 {
		// Assume this is a go_module (its top level filegroup has no explicit outputs)

	}

	// srcs is a faff because they can be either named or not named, which encodes them differently.
	log.Fatalf("here %T %#v %s %s", t.Srcs, t.Srcs, pkg.Name, t.Labels)
	return pkg
}
