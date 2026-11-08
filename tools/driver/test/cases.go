package main

import (
	"fmt"
	"strings"

	"github.com/please-build/go-rules/tools/driver/packages"
)

// A testCase is a request to the driver and the checks on its response that are specific to it. The invariants in
// driver_e2e.go are checked as well.
type testCase struct {
	name  string
	files []string // The files to ask for, relative to the repo root.
	tests bool     // Whether to ask for test packages.
	check func(*packages.DriverResponse) []error
}

var cases = []testCase{
	{name: "library/embed", files: []string{"test/embed/embed.go"}},
	{name: "library/external", files: []string{"test/external/test.go"}},
	{name: "test/cgo", files: []string{"test/cgo/cgo_test.go"}, tests: true},
	{name: "test/embed", files: []string{"test/embed/embed_test.go"}, tests: true},

	// #362: the export file of an external test package was left unset. These files are all in external test packages
	// (`package foo_test`), which the "roots have export files" invariant checks the response for.
	{
		name:  "external_test/buildinfo",
		files: []string{"test/buildinfo/buildinfo_test.go"},
		tests: true,
		check: checkHasExternalTest,
	},
	{
		name:  "external_test/external",
		files: []string{"test/external/external_test.go"},
		tests: true,
		check: checkHasExternalTest,
	},
	{
		name:  "external_test/import",
		files: []string{"test/import/import_test.go"},
		tests: true,
		check: checkHasExternalTest,
	},
	{
		name:  "external_test/importpath",
		files: []string{"test/importpath/importpath_test.go"},
		tests: true,
		check: checkHasExternalTest,
	},
}

// checkHasExternalTest checks that one of the roots is an external test package, so that the case is exercising what it
// is meant to.
func checkHasExternalTest(resp *packages.DriverResponse) []error {
	for _, pkg := range rootPackages(resp) {
		if strings.HasSuffix(pkg.ID, "_test") {
			return nil
		}
	}
	return []error{fmt.Errorf("no external test package among the roots %v", resp.Roots)}
}
