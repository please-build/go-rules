package test

import (
	"log"
)

// PleaseGoTest will generate the test main for the provided sources
func PleaseGoTest(goTool, dir, testPackage, output string, sources, exclude []string, isBenchmark, external bool) {
	if err := WriteTestMain(testPackage, sources, output, dir != "", isBenchmark, true); err != nil {
		log.Fatalf("Error writing test main: %s", err)
	}
}
