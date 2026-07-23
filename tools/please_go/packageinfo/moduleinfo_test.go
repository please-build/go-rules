package packageinfo

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func TestWriteModuleInfo_Simple(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an importconfig file that preserves package "foo"
	icFile := filepath.Join(tmpDir, "importconfig")
	icContent := `packagefile github.com/example/foo=export/path/to/foo.a
`
	err := os.WriteFile(icFile, []byte(icContent), 0644)
	require.NoError(t, err)

	// Create a simple package "foo"
	fooDir := filepath.Join(tmpDir, "foo")
	err = os.MkdirAll(fooDir, 0755)
	require.NoError(t, err)

	fooSrc := `package foo

import "fmt"

func Hello() {
	fmt.Println("hello")
}
`
	err = os.WriteFile(filepath.Join(fooDir, "foo.go"), []byte(fooSrc), 0644)
	require.NoError(t, err)

	var buf bytes.Buffer
	err = WriteModuleInfo(
		"github.com/example",
		tmpDir,
		icFile,
		nil,
		&buf,
	)
	require.NoError(t, err)

	var pkgs []*packages.Package
	err = json.Unmarshal(buf.Bytes(), &pkgs)
	require.NoError(t, err)

	expected := []*packages.Package{
		{
			ID:              "github.com/example/foo",
			Name:            "foo",
			PkgPath:         "github.com/example/foo",
			GoFiles:         []string{filepath.Join(fooDir, "foo.go")},
			CompiledGoFiles: []string{filepath.Join(fooDir, "foo.go")},
			ExportFile:      "export/path/to/foo.a",
			Imports: map[string]*packages.Package{
				"fmt": {
					ID: "fmt",
				},
			},
		},
	}

	assert.Equal(t, expected, pkgs)
}

func TestWriteModuleInfo_Filtering(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an empty importconfig file
	icFile := filepath.Join(tmpDir, "importconfig")
	err := os.WriteFile(icFile, []byte{}, 0644)
	require.NoError(t, err)

	// Create package "foo"
	fooDir := filepath.Join(tmpDir, "foo")
	err = os.MkdirAll(fooDir, 0755)
	require.NoError(t, err)

	fooSrc := `package foo
`
	err = os.WriteFile(filepath.Join(fooDir, "foo.go"), []byte(fooSrc), 0644)
	require.NoError(t, err)

	var buf bytes.Buffer
	err = WriteModuleInfo(
		"github.com/example",
		tmpDir,
		icFile,
		nil,
		&buf,
	)
	require.NoError(t, err)

	var pkgs []*packages.Package
	err = json.Unmarshal(buf.Bytes(), &pkgs)
	require.NoError(t, err)

	// Since importconfig is empty, foo is not present in imports map,
	// and slices.DeleteFunc will filter it out.
	assert.Empty(t, pkgs)
}

func TestWriteModuleInfo_InstallPkgs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an importconfig file containing all possible packages
	icFile := filepath.Join(tmpDir, "importconfig")
	icContent := `packagefile github.com/example/foo=export/foo.a
packagefile github.com/example/foo/sub=export/sub.a
packagefile github.com/example/bar=export/bar.a
`
	err := os.WriteFile(icFile, []byte(icContent), 0644)
	require.NoError(t, err)

	// Create package "foo"
	fooDir := filepath.Join(tmpDir, "foo")
	require.NoError(t, os.MkdirAll(fooDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(fooDir, "foo.go"), []byte("package foo\n"), 0644))

	// Create package "foo/sub"
	fooSubDir := filepath.Join(fooDir, "sub")
	require.NoError(t, os.MkdirAll(fooSubDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(fooSubDir, "sub.go"), []byte("package sub\n"), 0644))

	// Create package "bar"
	barDir := filepath.Join(tmpDir, "bar")
	require.NoError(t, os.MkdirAll(barDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(barDir, "bar.go"), []byte("package bar\n"), 0644))

	// 1. Without installPkgs, walking the entire tmpDir should return all 3 packages.
	{
		var buf bytes.Buffer
		err := WriteModuleInfo(
			"github.com/example",
			tmpDir,
			icFile,
			nil, // installPkgs is nil
			&buf,
		)
		require.NoError(t, err)

		var pkgs []*packages.Package
		err = json.Unmarshal(buf.Bytes(), &pkgs)
		require.NoError(t, err)

		require.Len(t, pkgs, 3)
		assert.Equal(t, "github.com/example/bar", pkgs[0].ID)
		assert.Equal(t, "github.com/example/foo", pkgs[1].ID)
		assert.Equal(t, "github.com/example/foo/sub", pkgs[2].ID)
	}

	// 2. With installPkgs = []string{"foo"}, only "foo" should be imported directly.
	{
		var buf bytes.Buffer
		err := WriteModuleInfo(
			"github.com/example",
			tmpDir,
			icFile,
			[]string{"foo"},
			&buf,
		)
		require.NoError(t, err)

		var pkgs []*packages.Package
		err = json.Unmarshal(buf.Bytes(), &pkgs)
		require.NoError(t, err)

		require.Len(t, pkgs, 1)
		assert.Equal(t, "github.com/example/foo", pkgs[0].ID)
	}

	// 3. With installPkgs = []string{"foo/..."}, it should recursively discover "foo" and "foo/sub" but still ignore "bar".
	{
		var buf bytes.Buffer
		err := WriteModuleInfo(
			"github.com/example",
			tmpDir,
			icFile,
			[]string{"foo/..."},
			&buf,
		)
		require.NoError(t, err)

		var pkgs []*packages.Package
		err = json.Unmarshal(buf.Bytes(), &pkgs)
		require.NoError(t, err)

		require.Len(t, pkgs, 2)
		assert.Equal(t, "github.com/example/foo", pkgs[0].ID)
		assert.Equal(t, "github.com/example/foo/sub", pkgs[1].ID)
	}
}

func TestWriteModuleInfo_Vendor(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an importconfig file containing both packages
	icFile := filepath.Join(tmpDir, "importconfig")
	icContent := `packagefile foo=export/foo.a
packagefile vendor/github.com/example/lib=export/lib.a
`
	err := os.WriteFile(icFile, []byte(icContent), 0644)
	require.NoError(t, err)

	// Create package "foo" which imports "github.com/example/lib"
	fooDir := filepath.Join(tmpDir, "foo")
	require.NoError(t, os.MkdirAll(fooDir, 0755))
	fooSrc := `package foo

import "github.com/example/lib"
`
	require.NoError(t, os.WriteFile(filepath.Join(fooDir, "foo.go"), []byte(fooSrc), 0644))

	// Create vendor package "vendor/github.com/example/lib"
	libDir := filepath.Join(tmpDir, "vendor", "github.com/example", "lib")
	require.NoError(t, os.MkdirAll(libDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(libDir, "lib.go"), []byte("package lib\n"), 0644))

	var buf bytes.Buffer
	err = WriteModuleInfo(
		"", // empty importPath for vendor directories
		tmpDir,
		icFile,
		nil,
		&buf,
	)
	require.NoError(t, err)

	var pkgs []*packages.Package
	err = json.Unmarshal(buf.Bytes(), &pkgs)
	require.NoError(t, err)

	// Since we populated importconfig, both packages should survive.
	// Output should be sorted by ID: "foo" first, then "vendor/github.com/example/lib"
	require.Len(t, pkgs, 2)

	fooPkg := pkgs[0]
	libPkg := pkgs[1]

	assert.Equal(t, "foo", fooPkg.ID)
	assert.Equal(t, "vendor/github.com/example/lib", libPkg.ID)

	// The "github.com/example/lib" import in "fooPkg" should have been resolved to the vendor package.
	// Therefore, its ID should be "vendor/github.com/example/lib"
	require.Contains(t, fooPkg.Imports, "github.com/example/lib")
	assert.Equal(t, "vendor/github.com/example/lib", fooPkg.Imports["github.com/example/lib"].ID)
}
