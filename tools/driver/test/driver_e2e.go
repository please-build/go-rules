// driver_e2e runs the package driver over some of this repo's own Go targets and checks properties of its response that
// must always hold. It has to run outside of a Please sandbox, as the driver runs plz itself, so it is a binary to
// `plz run` rather than a test.
//
// To cover a new bug, add a case to cases.go, or an invariant here if every response should satisfy it.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	xpackages "golang.org/x/tools/go/packages"

	"github.com/please-build/go-rules/tools/driver/packages"
)

var (
	driver = flag.String("driver", "", "Path to the package driver binary")
	run    = flag.String("run", "", "Only run the cases whose name matches this regular expression")
)

// An invariant is a property that the response to every case must satisfy.
type invariant struct {
	name  string
	check func(*packages.DriverResponse) error
}

var invariants = []invariant{
	{
		// A response with no roots analyses nothing, so any check of the roots would pass vacuously.
		name: "roots are not empty",
		check: func(resp *packages.DriverResponse) error {
			if len(resp.Roots) == 0 {
				return fmt.Errorf("response has no roots")
			}
			return nil
		},
	},
	{
		// Tools that read export data, rather than compiling the sources themselves, can't load a package without one.
		// The driver builds the archives when asked for export files. When it can't work out where one is, it ends up
		// with a directory, such as `plz-out/gen`.
		name: "roots have export files",
		check: func(resp *packages.DriverResponse) error {
			var errs []error
			for _, pkg := range rootPackages(resp) {
				if info, err := os.Stat(pkg.ExportFile); err != nil || !info.Mode().IsRegular() {
					errs = append(errs, fmt.Errorf("package %s has export file %q, which is not a file", pkg.ID, pkg.ExportFile))
				}
			}
			return errors.Join(errs...)
		},
	},
}

func main() {
	flag.Parse()
	if *driver == "" {
		fmt.Fprintln(os.Stderr, "--driver is required")
		os.Exit(2)
	}
	absDriver, err := filepath.Abs(*driver)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	filter, err := regexp.Compile(*run)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid --run:", err)
		os.Exit(2)
	}
	failed := false
	for _, c := range cases {
		if !filter.MatchString(c.name) {
			continue
		}
		errs := runCase(absDriver, c)
		for _, err := range errs {
			fmt.Printf("FAIL %s: %s\n", c.name, err)
			failed = true
		}
		if len(errs) == 0 {
			fmt.Printf("PASS %s\n", c.name)
		}
	}
	if failed {
		os.Exit(1)
	}
}

// runCase loads the case's files and returns everything wrong with the response.
func runCase(driver string, c testCase) []error {
	resp, err := load(driver, c)
	if err != nil {
		return []error{err}
	}
	var errs []error
	for _, inv := range invariants {
		if err := inv.check(resp); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", inv.name, err))
		}
	}
	if c.check != nil {
		errs = append(errs, c.check(resp)...)
	}
	return errs
}

// load runs the driver as go/packages would. It asks for export files, as tools that type check from export data do,
// which makes the driver build the packages' archives.
func load(driver string, c testCase) (*packages.DriverResponse, error) {
	mode := xpackages.NeedName | xpackages.NeedFiles | xpackages.NeedImports | xpackages.NeedExportFile
	req, err := json.Marshal(&packages.DriverRequest{Mode: mode, Tests: c.tests})
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(driver, c.files...)
	cmd.Stdin = bytes.NewReader(req)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	resp := &packages.DriverResponse{}
	if err := json.Unmarshal(out, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// rootPackages returns the packages in the response that were asked for.
func rootPackages(resp *packages.DriverResponse) []*xpackages.Package {
	isRoot := make(map[string]bool, len(resp.Roots))
	for _, id := range resp.Roots {
		isRoot[id] = true
	}
	var roots []*xpackages.Package
	for _, pkg := range resp.Packages {
		if isRoot[pkg.ID] {
			roots = append(roots, pkg)
		}
	}
	return roots
}
