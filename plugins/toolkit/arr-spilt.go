package toolkit

// ArrSplit 数组切割
func ArrSplit[T any](src []T, perLen uint) [][]T {
	res := make([][]T, 0)
	if len(src) <= int(perLen) || perLen == 0 {
		res = append(res, src)
		return res
	}
	start := 0
	end := int(perLen)
	for {
		if len(src) > end {
			res = append(res, src[start:end])
			start = end
			end = end + int(perLen)
		} else {
			res = append(res, src[start:])
			break
		}
	}
	return res
}
