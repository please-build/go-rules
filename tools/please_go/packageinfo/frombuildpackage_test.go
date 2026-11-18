package packageinfo

import (
	"go/build"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

// TestFromBuildPackageForModule checks that test files, imports and the "_test" ID suffix are only
// applied when includeTests is set. Stdlib loading in the driver relies on includeTests=false so that
// packages like "bufio" (which have external tests) keep their real import path.
func TestFromBuildPackageForModule(t *testing.T) {
	newPkg := func() *build.Package {
		return &build.Package{
			Dir:          "/src/foo",
			Name:         "foo",
			ImportPath:   "example.com/foo",
			GoFiles:      []string{"foo.go"},
			TestGoFiles:  []string{"foo_internal_test.go"},
			XTestGoFiles: []string{"foo_test.go"},
			Imports:      []string{"fmt"},
			TestImports:  []string{"testing"},
			XTestImports: []string{"example.com/foo"},
		}
	}

	for _, tc := range []struct {
		name         string
		includeTests bool
		wantID       string
		wantName     string
		wantFiles    []string
		wantImports  map[string]*packages.Package
	}{
		{
			name:         "tests included",
			includeTests: true,
			wantID:       "example.com/foo_test",
			wantName:     "foo_test",
			wantFiles:    []string{"/src/foo/foo.go", "/src/foo/foo_internal_test.go", "/src/foo/foo_test.go"},
			wantImports: map[string]*packages.Package{
				"fmt":             {ID: "fmt", PkgPath: "fmt"},
				"testing":         {ID: "testing", PkgPath: "testing"},
				"example.com/foo": {ID: "example.com/foo", PkgPath: "example.com/foo"},
			},
		},
		{
			name:         "tests excluded",
			includeTests: false,
			wantID:       "example.com/foo",
			wantName:     "foo",
			wantFiles:    []string{"/src/foo/foo.go"},
			wantImports:  map[string]*packages.Package{"fmt": {ID: "fmt", PkgPath: "fmt"}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bpkg := newPkg()

			got := FromBuildPackageForModule(bpkg, "export.a", tc.includeTests)

			assert.Equal(t, tc.wantID, got.ID)
			assert.Equal(t, tc.wantID, got.PkgPath)
			assert.Equal(t, tc.wantName, got.Name)
			assert.Equal(t, tc.wantFiles, got.GoFiles)
			assert.Equal(t, tc.wantFiles, got.CompiledGoFiles)
			assert.Equal(t, tc.wantImports, got.Imports)
			assert.Equal(t, "export.a", got.ExportFile)
			// The caller's package must never be modified, whichever mode is used.
			assert.Equal(t, newPkg(), bpkg)
		})
	}
}

// TestFromBuildPackageForModule_NoExternalTests checks that a package without external tests keeps
// its plain ID even when tests are included.
func TestFromBuildPackageForModule_NoExternalTests(t *testing.T) {
	got := FromBuildPackageForModule(&build.Package{
		Dir:         "/src/foo",
		Name:        "foo",
		ImportPath:  "example.com/foo",
		GoFiles:     []string{"foo.go"},
		TestGoFiles: []string{"foo_test.go"},
	}, "", true)

	require.NotNil(t, got)
	assert.Equal(t, "example.com/foo", got.ID)
	assert.Equal(t, "foo", got.Name)
	assert.Equal(t, []string{"/src/foo/foo.go", "/src/foo/foo_test.go"}, got.GoFiles)
}
