package buildinfo_test

import (
	"debug/buildinfo"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildInfoIsReadable(t *testing.T) {
	info, err := buildinfo.ReadFile(os.Getenv("DATA"))
	assert.NoError(t, err)
	// Don't get into detail here, but there should be _some_ module info available
	assert.Greater(t, len(info.Deps), 0)
}
