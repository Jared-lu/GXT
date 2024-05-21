package heap

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
			name:     "Empty Heap",
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

func TestHeap_Push(t *testing.T) {
	tesCases := []struct {
		name       string
		nums       []int
		pushNum    int
		wantNums   []int
		comparator GXT.Comparator[int]
	}{
		{
			name:     "MinHeap",
			nums:     []int{8, 6, 11, 3, 7, 9, 5},
			pushNum:  3,
			wantNums: []int{3, 3, 5, 6, 7, 9, 11, 8},
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
			name:     "Empty Heap",
			nums:     []int{},
			pushNum:  1,
			wantNums: []int{1},
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
			pushNum:  12,
			wantNums: []int{16, 14, 10, 8, 12, 9, 3, 2, 4, 1, 7},
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
			m.Push(tc.pushNum)
			assert.Equal(t, tc.wantNums, m.data)
		})
	}
}

func TestHeap_Pop(t *testing.T) {
	tesCases := []struct {
		name       string
		nums       []int
		wantNums   []int
		wantVal    int
		wantErr    error
		comparator GXT.Comparator[int]
	}{
		{
			name:     "MinHeap",
			nums:     []int{3, 3, 5, 6, 7, 11, 8},
			wantNums: []int{3, 6, 5, 8, 7, 11},
			comparator: func(src int, dst int) int {
				if src < dst {
					return -1
				}
				if src > dst {
					return 1
				}
				return 0
			},
			wantVal: 3,
			wantErr: nil,
		},
		{
			name:     "Empty Heap",
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
			wantVal: 0,
			wantErr: ErrEmptyHeap,
		},
		{
			name:     "MaxHeap",
			nums:     []int{16, 14, 10, 8, 12, 9, 3, 2, 4, 1, 7},
			wantNums: []int{14, 12, 10, 8, 7, 9, 3, 2, 4, 1},
			comparator: func(src int, dst int) int {
				if src < dst {
					return 1
				}
				if src > dst {
					return -1
				}
				return 0
			},
			wantVal: 16,
			wantErr: nil,
		},
	}
	for _, tc := range tesCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewHeap[int](tc.nums, tc.comparator)
			val, err := m.Pop()
			assert.Equal(t, tc.wantNums, m.data)
			assert.Equal(t, tc.wantVal, val)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestHeap_Peek(t *testing.T) {
	tesCases := []struct {
		name       string
		nums       []int
		wantVal    int
		wantErr    error
		comparator GXT.Comparator[int]
	}{
		{
			name: "MinHeap",
			nums: []int{3, 3, 5, 6, 7, 11, 8},
			comparator: func(src int, dst int) int {
				if src < dst {
					return -1
				}
				if src > dst {
					return 1
				}
				return 0
			},
			wantVal: 3,
			wantErr: nil,
		},
		{
			name: "Empty Heap",
			nums: []int{},
			comparator: func(src int, dst int) int {
				if src < dst {
					return -1
				}
				if src > dst {
					return 1
				}
				return 0
			},
			wantVal: 0,
			wantErr: ErrEmptyHeap,
		},
		{
			name: "MaxHeap",
			nums: []int{16, 14, 10, 8, 12, 9, 3, 2, 4, 1, 7},
			comparator: func(src int, dst int) int {
				if src < dst {
					return 1
				}
				if src > dst {
					return -1
				}
				return 0
			},
			wantVal: 16,
			wantErr: nil,
		},
	}
	for _, tc := range tesCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewHeap[int](tc.nums, tc.comparator)
			val, err := m.Peek()
			assert.Equal(t, tc.wantVal, val)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
