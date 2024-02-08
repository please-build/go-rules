package gomoddeps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var hostGoModPath = "tools/please_go/generate/gomoddeps/test_data/host_go_mod"
var moduleGoModPath = "tools/please_go/generate/gomoddeps/test_data/module_go_mod"
var invalidGoModPath = "tools/please_go/generate/gomoddeps/test_data/invalid_go_mod"

func TestErrors(t *testing.T) {
	t.Run("errors if host go.mod does not exist", func(t *testing.T) {
		deps, replacements, err := GetCombinedDepsAndRequirements("/does/not/exist", "/does/not/matter")
		assert.Error(t, err)
		assert.Nil(t, deps)
		assert.Nil(t, replacements)
	})

	t.Run("does not error if module go.mod does not exist", func(t *testing.T) {
		deps, replacements, err := GetCombinedDepsAndRequirements(hostGoModPath, "/does/not/exist")
		assert.NoError(t, err)
		assert.NotNil(t, deps)
		assert.NotNil(t, replacements)
	})

	t.Run("errors if host go.mod is invalid", func(t *testing.T) {
		deps, replacements, err := GetCombinedDepsAndRequirements(invalidGoModPath, "/does/not/exist")
		assert.Error(t, err)
		assert.Nil(t, deps)
		assert.Nil(t, replacements)
	})

	t.Run("errors if module go.mod is invalid", func(t *testing.T) {
		deps, replacements, err := GetCombinedDepsAndRequirements(hostGoModPath, invalidGoModPath)
		assert.Error(t, err)
		assert.Nil(t, deps)
		assert.Nil(t, replacements)
	})
}

func TestHostOnlyDeps(t *testing.T) {
	t.Run("returns all host dependencies", func(t *testing.T) {
		deps, _, err := GetCombinedDepsAndRequirements(hostGoModPath, "/does/not/matter")
		assert.NoError(t, err)

		assert.Len(t, deps, 3)
		assert.Equal(t, []string{"example.com/foo", "example.com/bar", "example.com/bob"}, deps)
	})

	t.Run("returns all host replacements", func(t *testing.T) {
		_, replacements, err := GetCombinedDepsAndRequirements(hostGoModPath, "/does/not/matter")
		assert.NoError(t, err)

		assert.Len(t, replacements, 1)
		assert.Equal(t, map[string]string{"example.com/bob": "example.com/new-bob"}, replacements)
	})
}

func TestModuleOnlyDeps(t *testing.T) {
	t.Run("returns all host dependencies", func(t *testing.T) {
		deps, _, err := GetCombinedDepsAndRequirements("", moduleGoModPath)
		assert.NoError(t, err)

		assert.Len(t, deps, 2)
		assert.Equal(t, []string{"example.com/dob", "example.com/bab"}, deps)
	})

	t.Run("returns all host replacements", func(t *testing.T) {
		_, replacements, err := GetCombinedDepsAndRequirements("", moduleGoModPath)
		assert.NoError(t, err)

		assert.Len(t, replacements, 1)
		assert.Equal(t, map[string]string{"example.com/bab": "example.com/new-bab"}, replacements)
	})
}

func TestBothDeps(t *testing.T) {
	t.Run("returns combined dependencies", func(t *testing.T) {
		deps, _, err := GetCombinedDepsAndRequirements(hostGoModPath, moduleGoModPath)
		assert.NoError(t, err)

		assert.Len(t, deps, 5)
		assert.Equal(t, []string{"example.com/foo", "example.com/bar", "example.com/bob", "example.com/dob", "example.com/bab"}, deps)
	})

	t.Run("returns only host replacements", func(t *testing.T) {
		_, replacements, err := GetCombinedDepsAndRequirements(hostGoModPath, moduleGoModPath)
		assert.NoError(t, err)

		assert.Len(t, replacements, 1)
		assert.Equal(t, map[string]string{"example.com/bob": "example.com/new-bob"}, replacements)
	})
}
