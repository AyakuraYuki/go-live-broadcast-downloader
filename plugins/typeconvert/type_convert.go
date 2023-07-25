package typeconvert

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/hack"
	"sort"
	"strconv"
	"strings"
)

// StringToBool : string -> bool
func StringToBool(s string) bool {
	if s == "" {
		return false
	}

	lower := strings.ToLower(s)
	if lower == "false" || lower == "no" || lower == "0" {
		return false
	} else if lower == "true" || lower == "yes" || lower == "1" {
		return true
	}

	return false
}

// StringToInt : string -> int
func StringToInt(s string, optionalNum ...int) int {
	if s == "" {
		if len(optionalNum) != 0 {
			return optionalNum[0]
		}
		return 0
	}

	temp, _ := strconv.Atoi(s)
	if temp == 0 && len(optionalNum) != 0 {
		return optionalNum[0]
	}

	return temp
}

// StringToInt16 : string -> int16
func StringToInt16(s string) int16 {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return int16(TmpInt)
}

// StringToInt32 : string -> int32
func StringToInt32(s string) int32 {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return int32(TmpInt)
}

// StringToInt64 : string -> int64
func StringToInt64(s string) int64 {
	if s == "" {
		return 0
	}

	tmp, _ := strconv.ParseInt(s, 10, 64)

	return tmp
}

// StringToUint : string -> uint
func StringToUint(s string) uint {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return uint(TmpInt)
}

// StringToUint16 : string -> uint16
func StringToUint16(s string) uint16 {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return uint16(TmpInt)
}

// StringToUint32 : string -> uint32
func StringToUint32(s string) uint32 {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return uint32(TmpInt)
}

// StringToUint64 : string -> uint64
func StringToUint64(s string) uint64 {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return uint64(TmpInt)
}

// 保留字符转换问题
var kspstrs = map[string]struct{}{
	"nan":       {},
	"infinity":  {},
	"+infinity": {},
	"-infinity": {},
	"inf":       {},
	"+inf":      {},
	"-inf":      {},
}

// StringToFloat64 : string -> float64
func StringToFloat64(s string) float64 {
	if s == "" {
		return 0
	}
	// 保留字符转换为0
	if _, ok := kspstrs[strings.ToLower(s)]; ok {
		return 0
	}
	tmp, _ := strconv.ParseFloat(s, 64)
	return tmp
}

// Int64ToString : int64 -> string
func Int64ToString(a int64) string {
	return strconv.FormatInt(a, 10)
}

// Int32ToString : int32 -> string
func Int32ToString(a int32) string {
	return strconv.FormatInt(int64(a), 10)
}

// Int16ToString : int16 -> string
func Int16ToString(a int16) string {
	return strconv.FormatInt(int64(a), 10)
}

// IntToString : int -> string
func IntToString(a int) string {
	return strconv.Itoa(a)
}

// Uint64ToString : uint64 -> string
func Uint64ToString(a uint64) string {
	return strconv.FormatUint(a, 10)
}

// Uint32ToString : uint32 -> string
func Uint32ToString(a uint32) string {
	return strconv.FormatUint(uint64(a), 10)
}

// Uint16ToString : uint16 -> string
func Uint16ToString(a uint16) string {
	return strconv.FormatUint(uint64(a), 10)
}

// Float64ToString : float64 -> string
func Float64ToString(a float64) string {
	return strconv.FormatFloat(a, 'f', -1, 64)
}

// StringToInt32Slice :
func StringToInt32Slice(s string, sep string) (ret []int32) {
	tokens := strings.Split(s, sep)
	for _, k := range tokens {
		i, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil
		}
		ret = append(ret, int32(i))
	}
	return
}

// BytesToString convert bytes to string
func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}

// ToString convert some type to string
// []string{"a","b"} => "a,b"
// []int{1,2} => "1,2"
// []string{} => ""
func ToString(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case int:
		return IntToString(x)
	case int64:
		return Int64ToString(x)
	case float64:
		return Float64ToString(x)
	case json.Number:
		return x.String()
	case []string:
		return strings.Join(x, ",")
	case []int, []int64, []float64:
		return ToString(ToStringSlice(x))
	case map[string]interface{}:
		data, _ := json.Marshal(x)
		return hack.String(data)
	default:
		return fmt.Sprint(v)
	}
}

// ToStringSlice ToStringSlice
func ToStringSlice(v interface{}) []string {
	switch x := v.(type) {
	case string:
		return strings.Split(x, ",")
	case []string:
		return x
	case []int:
		s := make([]string, len(x))
		for i := range x {
			s[i] = ToString(x[i])
		}
		return s
	case []int64:
		s := make([]string, len(x))
		for i := range x {
			s[i] = ToString(x[i])
		}
		return s
	case []float64:
		s := make([]string, len(x))
		for i := range x {
			s[i] = ToString(x[i])
		}
		return s
	case []interface{}:
		s := make([]string, len(x))
		for i := range x {
			s[i] = ToString(x[i])
		}
		return s
	default:
		return nil
	}
}

// ToStringWithSort converts some type to string with sort by acs
// []string{"b", "c", "a"} => "a,b,c"
// []int{1, 3, 2} => "1,2,3"
// []int64{21, 34, 29, 17, 33} => "17,21,29,33,34"
// []float64{3.14, 3.15, 3.02, 2.98} => "2.98,3.02,3.14,3.15"
// If `v` is `[]string`, the original slice will not be modified because of deep copy in this function.
func ToStringWithSort(v any) string {
	switch x := v.(type) {
	case []string:
		// deep copy, will not modify original slice
		s := make([]string, len(x))
		for i := range x {
			s[i] = x[i]
		}
		sort.Strings(s)
		return strings.Join(s, ",")
	case []int, []int64, []float64:
		return ToStringWithSort(ToStringSliceWithSort(x))
	default:
		return ToString(v)
	}
}

// ToStringSliceWithSort converts elements into the string and provides a sorted slice.
func ToStringSliceWithSort(v any) []string {
	switch x := v.(type) {
	case string:
		s := strings.Split(x, ",")
		sort.Strings(s)
		return s
	case []string:
		sort.Strings(x)
		return x
	case []int:
		s := make([]string, len(x))
		for i := range x {
			s[i] = ToString(x[i])
		}
		sort.Strings(s)
		return s
	case []int64:
		s := make([]string, len(x))
		for i := range x {
			s[i] = ToString(x[i])
		}
		sort.Strings(s)
		return s
	case []float64:
		s := make([]string, len(x))
		for i := range x {
			s[i] = ToString(x[i])
		}
		sort.Strings(s)
		return s
	case []interface{}:
		s := make([]string, len(x))
		for i := range x {
			s[i] = ToString(x[i])
		}
		sort.Strings(s)
		return s
	default:
		return nil
	}
}

func Int64ToBytes(n int64) []byte {
	bf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bf, binary.BigEndian, n)
	return bf.Bytes()
}

func BytesToInt64(bys []byte) int64 {
	bf := bytes.NewBuffer(bys)
	var data int64
	_ = binary.Read(bf, binary.BigEndian, &data)
	return data
}

// Uint64SliceToInterfaceSlice : []uint64 -> []interface{}
func Uint64SliceToInterfaceSlice(a []uint64) []interface{} {
	if len(a) == 0 {
		return []interface{}{}
	}
	arr := make([]interface{}, len(a))
	for i, v := range a {
		arr[i] = v
	}
	return arr
}
