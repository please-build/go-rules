package example

import "fmt"

// ExampleWithOutput has an output comment, so `go test` runs it and checks the output.
func ExampleWithOutput() {
	fmt.Println("hello")
	// Output: hello
}

// ExampleEmptyOutput has an empty output comment, so `go test` runs it and checks that
// it prints nothing.
func ExampleEmptyOutput() {
	// Output:
}

// ExampleNoOutput has no output comment, so `go test` compiles it but never runs it.
// Running it would fail: it exists only to show the API and cannot succeed on its own.
func ExampleNoOutput() {
	panic("compile-only example must not be run")
}
