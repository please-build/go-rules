package testing_definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThatWeAreInFactTesting(t *testing.T) {
	assert.True(t, testing.Testing())
}
