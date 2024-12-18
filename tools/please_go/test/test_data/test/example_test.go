// This isn't a 'real' source file, it's test data for //src/build/go:write_test_main_test

package buildgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The library here is a (very) reduced version of core that only has one file in it.
var coverageVars = []string{"core.GoCover_lock_go"}

func TestReadPkgdef(t *testing.T) {
	vars, err := readPkgdef("src/build/go/test_data/core.a")
	assert.NoError(t, err)
	assert.Equal(t, coverageVars, vars)
}

func TestReadCopiedPkgdef(t *testing.T) {
	// Sanity check that this file exists.
	vars, err := readPkgdef("src/build/go/test_data/x/core.a")
	assert.NoError(t, err)
	assert.Equal(t, coverageVars, vars)
}

func readPkgdef(name string) ([]string, error) {
	return nil, nil
}
