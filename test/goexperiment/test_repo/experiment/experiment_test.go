package experiment_test

import (
	"iter"
	"slices"
	"testing"
)

func EverySecond[T any](x []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < len(x); i += 2 {
			if !yield(x[i]) {
				return
			}
		}
	}
}

func TestRangefunc(t *testing.T) {
	x := []int{1, 2, 3, 4, 5}
	y := []int{}
	for it := range EverySecond(x) {
		y = append(y, it)
	}
	if !slices.Equal(y, []int{1, 3, 5}) {
		t.Errorf("unexpected y: %v", y)
	}
}
