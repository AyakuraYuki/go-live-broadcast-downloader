package cjson

import (
	"reflect"
)

func JsonArrayCompareIsEqual(a, b string) bool {
	if a == b {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	if !JSON.Valid([]byte(a)) || !JSON.Valid([]byte(b)) {
		// 有一个不是 json 就不等
		return false
	}
	var (
		mapA []map[string]interface{}
		mapB []map[string]interface{}
	)
	_ = JSON.UnmarshalFromString(a, &mapA)
	_ = JSON.UnmarshalFromString(b, &mapB)
	return reflect.DeepEqual(mapA, mapB)
}
