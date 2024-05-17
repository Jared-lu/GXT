package slice

// 判断是否满足条件
type matchFunc[T any] func(val T) bool

// Find 查找符合条件的元素，如果有多个符合，就返回第一个
func Find[T any](s []T, match matchFunc[T]) (T, bool) {
	for _, val := range s {
		if match(val) {
			return val, true
		}
	}
	var t T
	return t, false
}

// FindAll 返回所有满足条件的元素
func FindAll[T any](s []T, match matchFunc[T]) []T {
	// 返回的切片容量默认是原切片的 1/8
	result := make([]T, 0, len(s)>>3+1)
	for _, val := range s {
		if match(val) {
			result = append(result, val)
		}
	}
	return result
}
