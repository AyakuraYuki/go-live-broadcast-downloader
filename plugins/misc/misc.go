package misc

import "strings"

// IntersectString 字符串交集
func IntersectString(sliceA []string, sliceB []string) []string {
	var low, high []string
	if len(sliceA) < len(sliceB) {
		low = sliceA
		high = sliceB
	} else {
		low = sliceB
		high = sliceA
	}

	m := make(map[string]int)
	n := make([]string, 0)

	for _, v := range high {
		if _, ok := m[v]; ok {
			continue
		}
		m[v]++
	}

	for _, v := range low {
		times, _ := m[v]
		if times == 1 {
			n = append(n, v)
		}
	}

	return n
}

// ContainString 判断字符串是否存在于切片中
func ContainString(slice []string, item string) bool {
	if len(slice) == 0 {
		return false
	}
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// StringIsBlank 判断字符串是不是空白字符串
func StringIsBlank(str string) bool {
	return strings.TrimSpace(str) == ""
}

// DefaultStringIfEmpty 当字符串是空内容时返回默认字符串
func DefaultStringIfEmpty(str, defaultStr string) string {
	if str == "" {
		return defaultStr
	} else {
		return str
	}
}

// DefaultStringIfBlank 当字符串是空白内容时返回默认字符串
func DefaultStringIfBlank(str, defaultStr string) string {
	if StringIsBlank(str) {
		return defaultStr
	} else {
		return str
	}
}

// MaskString 遮蔽字符串
func MaskString(content string, frontDisplayLength, endDisplayLength int) string {
	if content == "" {
		return ""
	}

	if frontDisplayLength < 0 {
		frontDisplayLength = 0
	}
	if endDisplayLength < 0 {
		endDisplayLength = 0
	}

	length := len(content)
	if frontDisplayLength+endDisplayLength >= length {
		// 保留显示的字符个数 大于或等于 总长度，则直接返回
		return content
	}
	if frontDisplayLength >= length {
		// 头部显示字符个数 大于或等于 总长度，则直接返回
		return content
	}
	if endDisplayLength >= length {
		// 尾部显示字符个数 大于或等于 总长度，则直接返回
		return content
	}

	endIndex := length - endDisplayLength
	maskLength := endIndex - frontDisplayLength
	if maskLength < 0 {
		maskLength = 0
	}

	builder := strings.Builder{}
	builder.WriteString(content[:frontDisplayLength])
	builder.WriteString(strings.Repeat("*", maskLength))
	builder.WriteString(content[endIndex:])

	return builder.String()
}
