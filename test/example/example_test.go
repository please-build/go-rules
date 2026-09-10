package example

import "fmt"

func ExampleRun() {
	fmt.Println("ran")
	// Output: ran
}

func ExampleNotRun() {
	panic("an example without an output comment must not be run")
}
