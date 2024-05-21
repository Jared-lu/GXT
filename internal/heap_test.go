package internal

import (
	"github.com/Jared-lu/GXT"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeap(t *testing.T) {
	tesCases := []struct {
		name       string
		nums       []int
		wantNums   []int
		comparator GXT.Comparator[int]
	}{
		{
			name:     "MinHeap",
			nums:     []int{8, 6, 11, 3, 7, 9, 5},
			wantNums: []int{3, 6, 5, 8, 7, 9, 11},
			comparator: func(src int, dst int) int {
				if src < dst {
					return -1
				}
				if src > dst {
					return 1
				}
				return 0
			},
		},
		{
			name:     "0",
			nums:     []int{},
			wantNums: []int{},
			comparator: func(src int, dst int) int {
				if src < dst {
					return -1
				}
				if src > dst {
					return 1
				}
				return 0
			},
		},
		{
			name:     "MaxHeap",
			nums:     []int{4, 1, 3, 2, 16, 9, 10, 14, 8, 7},
			wantNums: []int{16, 14, 10, 8, 7, 9, 3, 2, 4, 1},
			comparator: func(src int, dst int) int {
				if src < dst {
					return 1
				}
				if src > dst {
					return -1
				}
				return 0
			},
		},
	}
	for _, tc := range tesCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewHeap[int](tc.nums, tc.comparator)
			assert.Equal(t, tc.wantNums, m.data)
		})
	}

}
