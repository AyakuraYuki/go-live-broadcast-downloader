package slices

import (
	"testing"
)

var indexTests = []struct {
	s    []int
	v    int
	want int
}{
	{nil, 0, -1},
	{[]int{}, 0, -1},
	{[]int{1, 2, 3}, 0, -1},
	{[]int{1, 2, 3}, 2, 1},
	{[]int{1, 2, 2, 3}, 2, 1},
	{[]int{1, 2, 3, 2}, 2, 1},
}

func TestIndex(t *testing.T) {
	for _, test := range indexTests {
		if got := Index(test.s, test.v); got != test.want {
			t.Errorf("Index(%v, %v) = %v, want %v\n", test.s, test.v, got, test.want)
		}
	}
}

var containsTests = []struct {
	s    []int
	v    int
	want bool
}{
	{nil, 0, false},
	{[]int{}, 0, false},
	{[]int{1, 2, 3}, 0, false},
	{[]int{1, 2, 3}, 2, true},
}

func TestContains(t *testing.T) {
	for _, test := range containsTests {
		if got := Contains(test.s, test.v); got != test.want {
			t.Errorf("Contains(%v, %v) = %v, want %v\n", test.s, test.v, got, test.want)
		}
	}
}
