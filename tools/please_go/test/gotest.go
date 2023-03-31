package test

import (
	"log"

	"github.com/please-build/go-rules/tools/please_go/install/toolchain"
)

// PleaseGoTest will generate the test main for the provided sources
func PleaseGoTest(goTool, dir, testPackage, output string, sources, exclude []string, isBenchmark, external bool) {
	tc := toolchain.Toolchain{GoTool: goTool}
	minor, err := tc.GoMinorVersion()
	if err != nil {
		log.Fatalf("Failed to determine Go version: %s", err)
	}
	if minor >= 20 {
		if err := WriteTestMain(testPackage, sources, output, dir != "", nil, isBenchmark, true, true); err != nil {
			log.Fatalf("Error writing test main: %s", err)
		}
		return
	}
	coverVars, err := FindCoverVars(dir, testPackage, external, exclude)
	if err != nil {
		log.Fatalf("Error scanning for coverage: %s", err)
	}
	if err := WriteTestMain(testPackage, sources, output, dir != "", coverVars, isBenchmark, minor >= 18, false); err != nil {
		log.Fatalf("Error writing test main: %s", err)
	}
}
