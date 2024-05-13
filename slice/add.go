package slice

import "errors"

var ErrIndexOutOfRange = errors.New("index out of range")

// Add insert element at src[idx]
func Add[T any](src []T, ele T, idx int) ([]T, error) {
	if idx < 0 || idx > len(src) {
		return nil, ErrIndexOutOfRange
	}
	var zeroVal T
	src = append(src, zeroVal)
	copy(src[idx+1:], src[idx:])
	src[idx] = ele
	return src, nil
}
