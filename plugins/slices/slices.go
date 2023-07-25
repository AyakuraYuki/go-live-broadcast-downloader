package slices

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[E comparable](s []E, v E) int {
	for i, vs := range s {
		if v == vs {
			return i
		}
	}
	return -1
}

// Contains reports whether v is present in s
func Contains[E comparable](s []E, v E) bool {
	return Index(s, v) >= 0
}

// Int64Index returns the index of the first occurrence of int64 v in s,
// or -1 if not present.
func Int64Index(s []int64, v int64) int {
	for i, vs := range s {
		if v == vs {
			return i
		}
	}
	return -1
}

// Int64Contains reports whether v is present in s
func Int64Contains(s []int64, v int64) bool {
	return Int64Index(s, v) >= 0
}
