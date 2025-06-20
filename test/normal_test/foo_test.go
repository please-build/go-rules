package foo

import "testing"

func TestFoo(t *testing.T) {
	res := foo()

	if res != "foo" {
		t.Fail()
	}
}
