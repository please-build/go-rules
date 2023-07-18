package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/please-build/go-rules/tools/driver/packages"
)

func TestPackages(t *testing.T) {
	_, err := packages.Load(&packages.DriverRequest{}, nil)
	assert.Error(t, err)
}
