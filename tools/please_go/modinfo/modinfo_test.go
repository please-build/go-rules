package modinfo_test

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"

	"github.com/please-build/go-rules/tools/please_go/modinfo"
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
			Path:    "golang.org/x/mod",
			Version: "v0.40.0",
		},
		{
			Path:    "gopkg.in/yaml.v3",
			Version: "v3.0.0-20210107192922-496545a6307b",
		},
	}, info.Deps)
}

func TestWriteModInfoGoVersion(t *testing.T) {
	dir := t.TempDir()
	goTool := filepath.Join(dir, "go")
	// `go version` names the host. GOVERSION must be what gets stamped.
	script := `#!/bin/sh
case "$1" in
env)
	if [ "$2" = GOVERSION ]; then
		echo go1.27.1
		exit 0
	fi
	;;
version)
	echo "go version go1.2.3 darwin/arm64"
	exit 0
	;;
esac
echo "unexpected: $*" >&2
exit 1
`
	require.NoError(t, os.WriteFile(goTool, []byte(script), 0o755))

	// WriteModInfo walks "." for dependency .modinfo files.
	t.Chdir(dir)
	out := filepath.Join(dir, "out.importconfig")
	err := modinfo.WriteModInfo(goTool, "example.com/mod", "example.com/mod/cmd", "exe", "0", "linux", "arm64", out)
	require.NoError(t, err)

	payload := modInfoPayload(t, out)
	assert.Regexp(t, `(?m)^go\tgo1\.27\.1$`, payload)
	assert.NotContains(t, payload, "go version")
	assert.NotContains(t, payload, "darwin/arm64")
	assert.Contains(t, payload, "build\tGOOS=linux\n")
	assert.Contains(t, payload, "build\tGOARCH=arm64\n")
}

func modInfoPayload(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	line := strings.TrimSpace(string(data))
	require.True(t, strings.HasPrefix(line, "modinfo "), line)
	unquoted, err := strconv.Unquote(strings.TrimPrefix(line, "modinfo "))
	require.NoError(t, err)
	const magic = 16
	require.Greater(t, len(unquoted), magic*2)
	return unquoted[magic : len(unquoted)-magic]
}
