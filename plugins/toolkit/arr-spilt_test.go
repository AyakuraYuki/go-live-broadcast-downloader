package toolkit

import "testing"

func TestArrSplit(t *testing.T) {
	src := []string{"a", "b", "c", "d", "e"}
	res := ArrSplit[string](src, 2)
	t.Log(res)
}
