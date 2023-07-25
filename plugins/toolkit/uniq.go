package toolkit

func UniqInt64s(a []int64) []int64 {
	if len(a) <= 1 {
		return a
	}
	b := make([]int64, 0)
	tmp := make(map[int64]struct{})
	for _, v := range a {
		if _, ok := tmp[v]; ok {
			continue
		}
		tmp[v] = struct{}{}
		b = append(b, v)
	}
	return b

}
