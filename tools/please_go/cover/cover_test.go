package cover_test

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallersUnderCoverage(t *testing.T) {
	_, filename, line, ok := runtime.Caller(0)
	require.True(t, ok)
	assert.Equal(t, "tools/please_go/cover/cover_test.go", filename)
	assert.Equal(t, 12, line)
}
