package slice

// Delete 删除切片特定下标的元素
func Delete[T any](src []T, idx int) ([]T, error) {
	length := len(src)
	if idx < 0 || idx >= length {
		return nil, ErrIndexOutOfRange
	}
	src = append(src[:idx], src[idx+1:]...)
	return src, nil
}
