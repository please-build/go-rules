package modinfo_test

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModInfo(t *testing.T) {
	info, ok := debug.ReadBuildInfo()
	require.True(t, ok)
	assert.Equal(t, "github.com/please-build/go-rules/tools/please_go/modinfo", info.Path)
	assert.Equal(t, debug.Module{Path: "github.com/please-build/go-rules"}, info.Main)
	// This lot will change if/when we update versions of these things. That's fine, we just want to assert
	// that there's _something_ sensible in this field.
	assert.Equal(t, []*debug.Module{
		{
			Path:    "github.com/davecgh/go-spew",
			Version: "v1.1.1",
		},
		{
			Path:    "github.com/pmezard/go-difflib",
			Version: "v1.0.0",
		},
		{
			Path:    "github.com/stretchr/testify",
			Version: "v1.7.1",
		},
		{
			Path:    "gopkg.in/yaml.v3",
			Version: "v3.0.0-20210107192922-496545a6307b",
		},
	}, info.Deps)
}
