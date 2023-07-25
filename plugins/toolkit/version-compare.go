package toolkit

import (
	"go-live-broadcast-downloader/plugins/typeconvert"
	"math"
	"strings"
)

type VersionCompareType int

const (
	// VersionGt 大于
	VersionGt VersionCompareType = iota
	// VersionLt 小于
	VersionLt
	// VersionEq 等于
	VersionEq
)

// VersionCompare 版本号对比
// v1>v2  VersionGt
// v1=v2  VersionEq
// v1<v2  VersionLt
func VersionCompare(ver1, ver2 string) (c VersionCompareType, tmpV1 []int64, tmpV2 []int64) {
	v1 := strings.TrimSpace(ver1)
	v2 := strings.TrimSpace(ver2)

	tmpV1 = make([]int64, 0)
	tmpV2 = make([]int64, 0)

	// 切割
	for _, v := range strings.Split(v1, ".") {
		tmpV1 = append(tmpV1, typeconvert.StringToInt64(v))
	}
	for _, v := range strings.Split(v2, ".") {
		tmpV2 = append(tmpV2, typeconvert.StringToInt64(v))
	}

	// 保证数组长度不为0
	if len(tmpV1) == 0 {
		tmpV1 = append(tmpV1, 0)
	}
	if len(tmpV2) == 0 {
		tmpV2 = append(tmpV2, 0)
	}

	div := float64(len(tmpV1) - len(tmpV2))
	divAbs := int64(math.Abs(div))

	// 补齐短的
	if div > 0 {
		tmpV2 = append(tmpV2, make([]int64, divAbs)...)
	} else if div < 0 {
		tmpV1 = append(tmpV1, make([]int64, divAbs)...)
	}

	for i, item1 := range tmpV1 {
		if item1 > tmpV2[i] {
			return VersionGt, tmpV1, tmpV2
		} else if item1 < tmpV2[i] {
			return VersionLt, tmpV1, tmpV2
		}
	}

	return VersionEq, tmpV1, tmpV2
}
