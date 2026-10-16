package packages

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func writePackageInfoFile(t *testing.T, dir, name string, pkgs ...*packages.Package) string {
	t.Helper()
	b, err := json.Marshal(pkgs)
	require.NoError(t, err)
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, b, 0644))
	return path
}

func TestLoadPackageInfoFiles_DuplicatePaths(t *testing.T) {
	dir := t.TempDir()
	foo := writePackageInfoFile(t, dir, "foo_pkg_info.json", &packages.Package{ID: "foo"})
	bar := writePackageInfoFile(t, dir, "bar_pkg_info.json", &packages.Package{ID: "bar"})

	// `plz build` prints the outputs of every target it is given, so a target that is in the output of both
	// `plz query deps` and `plz query filter` has its package info file listed twice.
	pkgs, err := loadPackageInfoFiles([]string{foo, bar, foo})
	require.NoError(t, err)

	ids := make([]string, len(pkgs))
	for i, pkg := range pkgs {
		ids[i] = pkg.ID
	}
	assert.ElementsMatch(t, []string{"foo", "bar"}, ids)
}
