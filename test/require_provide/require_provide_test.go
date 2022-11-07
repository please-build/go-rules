package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequireProvide(t *testing.T) {
	assert.Equal(t, 42, WhatIsTheAnswerToLifeTheUniverseAndEverything())
}
