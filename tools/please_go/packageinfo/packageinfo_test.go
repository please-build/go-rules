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

func TestWritePackageInfo_Simple(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple package "foo"
	fooDir := filepath.Join(tmpDir, "foo")
	err := os.MkdirAll(fooDir, 0755)
	require.NoError(t, err)

	fooSrc := `package foo

import "fmt"

func Hello() {
	fmt.Println("hello")
}
`
	err = os.WriteFile(filepath.Join(fooDir, "foo.go"), []byte(fooSrc), 0644)
	require.NoError(t, err)

	imports := map[string]string{
		"module/foo": "foo/foo.a",
	}

	var buf bytes.Buffer
	err = WritePackageInfo(
		"module/foo",
		fooDir,
		imports,
		"",
		"",
		false,
		&buf,
	)
	require.NoError(t, err)

	var pkgs []*packages.Package
	err = json.Unmarshal(buf.Bytes(), &pkgs)
	require.NoError(t, err)

	expected := []*packages.Package{
		{
			ID:              "module/foo",
			Name:            "foo",
			PkgPath:         "module/foo",
			GoFiles:         []string{filepath.Join(fooDir, "foo.go")},
			CompiledGoFiles: []string{filepath.Join(fooDir, "foo.go")},
			ExportFile:      "foo/foo.a",
			Imports: map[string]*packages.Package{
				"fmt": {
					ID: "fmt",
				},
			},
		},
	}

	assert.Equal(t, expected, pkgs)
}
