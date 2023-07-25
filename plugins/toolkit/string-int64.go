package toolkit

import (
	"go-live-broadcast-downloader/plugins/typeconvert"
	"strings"
)

func Strings2Int64List(str string, sep string) []int64 {
	if str == "" {
		return []int64{}
	}
	strList := strings.Split(str, sep)
	if len(strList) == 0 {
		return []int64{}
	}
	intList := make([]int64, 0)
	for _, v := range strList {
		intList = append(intList, typeconvert.StringToInt64(v))
	}
	return intList
}
