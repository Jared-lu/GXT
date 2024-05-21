package GXT

// Comparator 比较两个元素的大小
// src < dst, return -1
// src == dst, return 0
// src > dst, return 1
type Comparator[T any] func(src T, dst T) int
