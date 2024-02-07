package test_repo

// Testify has a dependency on github.com/pmezard/go-difflib which we DON'T
// have in our go.mod file but testify does.
// The test ensures that go_please generate considers both the host and the
// module go.mod and correctly generates the dependencies for the subrepo.

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertImportable(t *testing.T) {
	assert.True(t, true)
}
