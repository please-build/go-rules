package cgo

import "testing"

func TestHello(t *testing.T) {
	msg := Hello()
	if msg != "hello" {
		t.Fatalf("unexpected msg %v", msg)
	}
}
