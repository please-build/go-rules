package sample

import (
	"fmt"
	"testing"
)

func TestTestingFlag(t *testing.T) {
	fmt.Printf("testing.Testing() = %v\n", testing.Testing())
	if !testing.Testing() {
		t.Fatalf("expected testing.Testing() to be true")
	}
}
