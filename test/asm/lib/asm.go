// Package lib implements a simple test of Go assembly in a library.
package lib

import "github.com/please-build/go-rules/test/asm/lib/golib"

// add is the forward declaration of the assembly implementation.
func add(x, y int64) int64

// Add adds two numbers using assembly.
func Add() int {
	return int(add(golib.LHS, golib.RHS))
}
