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
	assert.Equal(t, []*debug.Module{{
		Path:    "github.com/stretchr/testify",
		Version: "v1.7.0",
	}}, info.Deps)
}
