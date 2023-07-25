package typeconvert

import (
	"reflect"
	"testing"
)

func TestToString(t *testing.T) {
	if r := ToString(1); r != "1" {
		t.Fatal("should equal", r)
	}

	if r := ToString(3.14); r != "3.14" {
		t.Fatal("should equal", r)
	}

	a1 := []string{"a", "b", "c"}
	if r := ToString(a1); r != "a,b,c" {
		t.Fatal("should equal", r)
	}

	a2 := []int64{1, 2, 3}
	if r := ToString(a2); r != "1,2,3" {
		t.Fatal("should equal", r)
	}

	a3 := []float64{3.14, 3.15}
	if r := ToString(a3); r != "3.14,3.15" {
		t.Fatal("should equal", r)
	}

	var a4 []float64
	if r := ToString(a4); r != "" {
		t.Fatal("should equal", r)
	}

	a5 := []int64{100202, 100221, 100103, 100554, 100387}
	if r := ToString(a5); r != "" {
		t.Log(r)
	} else {
		t.Fatal("should equal", r)
	}
}

func TestToStringSlice(t *testing.T) {
	if r := ToStringSlice("a,b,c"); !reflect.DeepEqual(r, []string{"a", "b", "c"}) {
		t.Fatal("should equal", r)
	}

	var x []interface{}
	x = append(x, "a")
	x = append(x, 1)

	if r := ToStringSlice(x); !reflect.DeepEqual(r, []string{"a", "1"}) {
		t.Fatal("should equal", r)
	}
}

func TestToStringWithSort(t *testing.T) {
	tests := []struct {
		arr  any
		want string
	}{
		{[]string{"b", "c", "a"}, "a,b,c"},
		{[]int{1, 2, 3}, "1,2,3"},
		{[]int{1, 3, 2}, "1,2,3"},
		{[]int64{21, 34, 29, 17, 33}, "17,21,29,33,34"},
		{[]float64{3.14, 3.15, 3.02, 2.98}, "2.98,3.02,3.14,3.15"},
	}
	for _, test := range tests {
		if r := ToStringWithSort(test.arr); r != test.want {
			t.Fatalf("should equals to %s but get %s\n", test.want, r)
		} else {
			t.Logf("passed: %v", test.arr)
		}
	}
}
