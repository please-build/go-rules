package import_test

import (
	"github.com/please-build/go-rules/test/message"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportLocal(t *testing.T) {
	assert.Equal(t, "Hello there!", message.Message)
}
