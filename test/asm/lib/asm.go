// Package lib implements a simple test of Go assembly in a library.
package lib

import "github.com/please-build/go-rules/test/asm/lib/golib"

// add is the forward declaration of the assembly implementation.
func add(x, y int64) int64

// Add adds two numbers using assembly.
func Add() int {
	return int(add(golib.LHS, golib.RHS))
}

// subtract is the forward declaration of the assembly implementation.
func subtract(x, y int) int64

// Subtract subtracts two numbers using assembly.
func Subtract() int {
	return int(subtract(golib.LHS, golib.RHS))
}
