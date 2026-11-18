package packages

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

// TestPackagesToResponse_Deduplicates checks that the same package reported more than once (e.g. by
// several pkg_info files, each of which includes its transitive deps) yields a single entry, since
// go/packages rejects responses with duplicate IDs. The first occurrence must win, and roots must
// not be repeated either. A "runtime" package with a source file is included so the stdlib isn't
// loaded via `go list`.
func TestPackagesToResponse_Deduplicates(t *testing.T) {
	for _, tc := range []struct {
		name      string
		pkgs      []*packages.Package
		dirs      map[string]struct{}
		wantIDs   []string
		wantFirst map[string]string // ID -> expected Name of the surviving package
		wantRoots []string
	}{
		{
			name: "no duplicates",
			pkgs: []*packages.Package{
				{ID: "a", Name: "a", ExportFile: "x/a.a"},
				{ID: "runtime", GoFiles: []string{"runtime/rt.go"}},
			},
			dirs:      map[string]struct{}{"x": {}},
			wantIDs:   []string{"a", "runtime"},
			wantRoots: []string{"a"},
		},
		{
			name: "duplicate root matched by export file dir",
			pkgs: []*packages.Package{
				{ID: "a", Name: "first", ExportFile: "x/a.a"},
				{ID: "a", Name: "second", ExportFile: "x/a.a"},
				{ID: "runtime", GoFiles: []string{"runtime/rt.go"}},
			},
			dirs:      map[string]struct{}{"x": {}},
			wantIDs:   []string{"a", "runtime"},
			wantFirst: map[string]string{"a": "first"},
			wantRoots: []string{"a"},
		},
		{
			name: "duplicate root matched by go file dir",
			pkgs: []*packages.Package{
				{ID: "a", Name: "first", GoFiles: []string{"x/a.go"}},
				{ID: "a", Name: "second", GoFiles: []string{"x/a.go"}},
				{ID: "runtime", GoFiles: []string{"runtime/rt.go"}},
			},
			dirs:      map[string]struct{}{"x": {}},
			wantIDs:   []string{"a", "runtime"},
			wantFirst: map[string]string{"a": "first"},
			wantRoots: []string{"a"},
		},
		{
			name: "duplicate non-root dependency",
			pkgs: []*packages.Package{
				{ID: "a", Name: "a", ExportFile: "x/a.a"},
				{ID: "dep", Name: "first", ExportFile: "y/dep.a"},
				{ID: "dep", Name: "second", ExportFile: "y/dep.a"},
				{ID: "runtime", GoFiles: []string{"runtime/rt.go"}},
			},
			dirs:      map[string]struct{}{"x": {}},
			wantIDs:   []string{"a", "dep", "runtime"},
			wantFirst: map[string]string{"dep": "first"},
			wantRoots: []string{"a"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := packagesToResponse("/root", tc.pkgs, tc.dirs)
			require.NoError(t, err)

			gotIDs := make([]string, 0, len(resp.Packages))
			byID := map[string]*packages.Package{}
			for _, pkg := range resp.Packages {
				gotIDs = append(gotIDs, pkg.ID)
				byID[pkg.ID] = pkg
			}
			assert.Equal(t, tc.wantIDs, gotIDs)
			for id, name := range tc.wantFirst {
				assert.Equal(t, name, byID[id].Name)
			}
			assert.Equal(t, tc.wantRoots, resp.Roots)
		})
	}
}
