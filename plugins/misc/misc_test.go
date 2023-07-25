package misc

import "testing"

func TestIntersectString(t *testing.T) {
	sliceA := []string{"a", "b", "c"}
	sliceB := []string{"c", "d"}
	intersection := IntersectString(sliceA, sliceB)
	t.Logf("%+v\n", intersection)
}
