package test

import (
	"log"

	"github.com/please-build/go-rules/tools/please_go/install/toolchain"
)

// PleaseGoTest will generate the test main for the provided sources
func PleaseGoTest(goTool, dir, testPackage, output string, sources, exclude []string, isBenchmark, external bool) {
	if ver, err := toolchain.GoMinorVersion(goTool); err != nil {
		log.Fatalf("Couldn't determine Go version: %s", err)
	} else if err := WriteTestMain(testPackage, sources, output, dir != "", isBenchmark, ver >= 23); err != nil {
		log.Fatalf("Error writing test main: %s", err)
	}
}
