// This test covers the case of a third-party package that includes assembly; we use xxhash
// as something that is fairly lightweight.

package asm_test

import (
	"testing"

	"github.com/cespare/xxhash/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackageWithAssembly(t *testing.T) {
	sum := xxhash.Sum64String("what is the meaning of life, the universe, and everything?")
	assert.NotEqual(t, 42, sum) // sadly xxhash does not answer this for us :(
}
