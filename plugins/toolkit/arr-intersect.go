package toolkit

// SortedArrIntersect 已经按升序排好序的两个数组 取交集部分
func SortedArrIntersect(sortNums1 []int64, sortNums2 []int64) []int64 {
	length1, length2 := len(sortNums1), len(sortNums2)
	index1, index2 := 0, 0

	intersection := make([]int64, 0)
	for index1 < length1 && index2 < length2 {
		if sortNums1[index1] < sortNums2[index2] {
			index1++
		} else if sortNums1[index1] > sortNums2[index2] {
			index2++
		} else {
			intersection = append(intersection, sortNums1[index1])
			index1++
			index2++
		}
	}
	return intersection
}
