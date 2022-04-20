package import_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/please-build/go-rules/test/import/legacy/foo"
	"github.com/please-build/go-rules/test/message"
)

func TestImportLocal(t *testing.T) {
	assert.Equal(t, "Hello there!", message.Message)
}

func TestImportFoo(t *testing.T) {
	assert.Equal(t, "Hello foo!", foo.Message)
}
