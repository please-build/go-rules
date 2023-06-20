package generate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrimPath(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		target   string
		expected string
	}{
		{
			name:     "trims base path",
			base:     "third_party/go/_foo#dl",
			target:   "third_party/go/_foo#dl/foo",
			expected: "foo",
		},
		{
			name:     "returns target if base is shorter",
			base:     "foo/bar/baz",
			target:   "foo/bar",
			expected: "foo/bar",
		},
		{
			name:     "returns target if not in base",
			base:     "foo/bar",
			target:   "bar/baz",
			expected: "bar/baz",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, trimPath(test.target, test.base))
		})
	}
}

func TestBuildTarget(t *testing.T) {
	tests := []struct {
		test               string
		name, pkg, subrepo string
		expected           string
	}{
		{
			test:     "fully qualified",
			subrepo:  "subrepo",
			pkg:      "pkg",
			name:     "name",
			expected: "///subrepo//pkg:name",
		},
		{
			test:     "fully qualified local package",
			subrepo:  "",
			pkg:      "pkg",
			name:     "name",
			expected: "//pkg:name",
		},
		{
			test:     "root package",
			subrepo:  "",
			pkg:      "",
			name:     "name",
			expected: "//:name",
		},
		{
			test:     "root package via .",
			subrepo:  "",
			pkg:      ".",
			name:     "name",
			expected: "//:name",
		},
		{
			test:     "root package in subrepo",
			subrepo:  "subrepo",
			pkg:      "",
			name:     "name",
			expected: "///subrepo//:name",
		},
		{
			test:     "pkg base matches name",
			subrepo:  "",
			pkg:      "foo",
			name:     "foo",
			expected: "//foo",
		},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			assert.Equal(t, test.expected, buildTarget(test.name, test.pkg, test.subrepo))
		})
	}
}

func TestDepTarget(t *testing.T) {
	tests := []struct {
		name         string
		deps         []string
		importTarget string
		expected     string
	}{
		{
			name:         "resolves local import",
			importTarget: "github.com/this/module/foo",
			expected:     "//foo",
		},
		{
			name:         "resolves local import in base",
			importTarget: "github.com/this/module",
			expected:     "//:module",
		},
		{
			name:         "resolves import to another module",
			importTarget: "github.com/some/module/foo",
			expected:     "///third_party/go/github.com_some_module//foo",
			deps:         []string{"github.com/some/module"},
		},
		{
			name:         "resolves import to longest match",
			importTarget: "github.com/some/module/foo/bar",
			expected:     "///third_party/go/github.com_some_module_foo//bar",
			deps:         []string{"github.com/some/module", "github.com/some/module/foo"},
		},
		{
			name:         "root package matches module base",
			importTarget: "github.com/some/module",
			expected:     "///third_party/go/github.com_some_module//:module",
			deps:         []string{"github.com/some/module"},
		},
		{
			name:         "replaces all with lib locally",
			importTarget: "github.com/this/module/all",
			expected:     "//all:lib",
		},
		{
			name:         "replaces all with lib in another repo",
			importTarget: "github.com/some/module/all",
			expected:     "///third_party/go/github.com_some_module//all:lib",
			deps:         []string{"github.com/some/module"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := &Generate{
				moduleName:         "github.com/this/module",
				thirdPartyFolder:   "third_party/go",
				replace:            map[string]string{},
				knownImportTargets: map[string]string{},
				moduleDeps:         test.deps,
			}

			assert.Equal(t, test.expected, g.depTarget(test.importTarget))
		})
	}
}
