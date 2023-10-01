package embed

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLibEmbed(t *testing.T) {
	assert.Equal(t, "hello", strings.TrimSpace(hello))
}

func TestEmbedDir(t *testing.T) {
	_, err := subdir.ReadFile("subdir/_test.txt")
	assert.Error(t, err)
}

func TestEmbedDirAll(t *testing.T) {
	b, err := subdirAll.ReadFile("subdir/_test.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello", strings.TrimSpace(string(b)))
}
