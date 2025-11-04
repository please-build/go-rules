// Package lib_test implements an external test on compiling Go with assembly.
// It has to be external since right now we don't support Go tests with assembly sources.
package lib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	asm "github.com/please-build/go-rules/test/asm/lib"
)

func TestAssemblyAdd(t *testing.T) {
	assert.Equal(t, 14, asm.Add())
}

func TestAssemblySubtract(t *testing.T) {
	assert.Equal(t, 6, asm.Subtract())
}
