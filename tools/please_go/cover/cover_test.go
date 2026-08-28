package cover

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteCoveragePassesAbsoluteSources(t *testing.T) {
	tmpDir := t.TempDir()

	pkgDir := filepath.Join(tmpDir, "src", "mypkg")
	require.NoError(t, os.MkdirAll(pkgDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "mypkg.go"), []byte("package mypkg\n"), 0644))

	// A fake cover tool that records the arguments it was given.
	argsFile := filepath.Join(tmpDir, "args")
	coverTool := filepath.Join(tmpDir, "cover.sh")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"" + argsFile + "\"\n"
	require.NoError(t, os.WriteFile(coverTool, []byte(script), 0755))

	t.Chdir(tmpDir)
	src := filepath.Join("src", "mypkg", "mypkg.go")
	outfilelist := filepath.Join(tmpDir, "outfilelist")
	require.NoError(t, WriteCoverage("go", coverTool, filepath.Join(tmpDir, "covcfg"), outfilelist, "example.com/mypkg", []string{src}))

	args, err := os.ReadFile(argsFile)
	require.NoError(t, err)
	assert.Contains(t, strings.Fields(string(args)), filepath.Join(pkgDir, "mypkg.go"))
	assert.NotContains(t, strings.Fields(string(args)), src)

	list, err := os.ReadFile(outfilelist)
	require.NoError(t, err)
	assert.Equal(t, []string{
		filepath.Join(pkgDir, "_covervars.cover.go"),
		filepath.Join(pkgDir, "mypkg.cover.go"),
	}, strings.Split(strings.TrimSpace(string(list)), "\n"))
}
