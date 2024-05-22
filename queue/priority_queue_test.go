package queue

import (
	"github.com/Jared-lu/GXT"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPriorityQueue(t *testing.T) {
	testCases := []struct {
		name         string
		capacity     int
		comparator   GXT.Comparator[int]
		IsBoundless  bool
		wantCapacity int
	}{
		{
			name:     "boundless queue",
			capacity: 0,
			comparator: func(src int, dst int) int {
				if src < dst {
					return -1
				}
				if src > dst {
					return 1
				}
				return 0
			},
			IsBoundless:  true,
			wantCapacity: 64,
		},
		{
			name:     "bounded queue",
			capacity: 10,
			comparator: func(src int, dst int) int {
				if src < dst {
					return -1
				}
				if src > dst {
					return 1
				}
				return 0
			},
			IsBoundless:  false,
			wantCapacity: 10,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewPriorityQueue(tc.capacity, tc.comparator)
			assert.Equal(t, tc.IsBoundless, q.boundless)
			assert.Equal(t, tc.wantCapacity, q.capacity)
		})
	}
}

func TestPriorityQueue_Enqueue(t *testing.T) {
	testCases := []struct {
		name     string
		capacity int
		wantErr  error
	}{
		{
			name:     "boundless queue - normal",
			capacity: 0,
			wantErr:  nil,
		},
		{
			name:     "bounded queue - out of capacity",
			capacity: 3,
			wantErr:  ErrOutOfCapacity,
		},
		{
			name:     "bounded queue - normal",
			capacity: 4,
			wantErr:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 新建一个队列
			q := NewPriorityQueue(tc.capacity, func(src int, dst int) int {
				if src < dst {
					return -1
				}
				if src > dst {
					return 1
				}
				return 0
			})
			// 入队
			q.Enqueue(1)
			q.Enqueue(2)
			q.Enqueue(3)
			err := q.Enqueue(4)
			assert.Equal(t, tc.wantErr, err)

		})
	}
}

func TestPriorityQueue_Dequeue(t *testing.T) {
	// 新建一个队列
	q := NewPriorityQueue(3, func(src int, dst int) int {
		if src < dst {
			return -1
		}
		if src > dst {
			return 1
		}
		return 0
	})
	// 入队
	q.Enqueue(3)
	q.Enqueue(11)
	q.Enqueue(9)

	// 测试出队，出队的顺序应该是小根堆的样子
	val, err := q.Dequeue()
	assert.Equal(t, nil, err)
	assert.Equal(t, val, 3)

	val, err = q.Dequeue()
	assert.Equal(t, nil, err)
	assert.Equal(t, val, 9)

	val, err = q.Dequeue()
	assert.Equal(t, nil, err)
	assert.Equal(t, val, 11)

	val, err = q.Dequeue()
	assert.Equal(t, ErrEmptyQueue, err)
	assert.Equal(t, val, 0)
}

func TestPriorityQueue_Peek(t *testing.T) {
	// 新建一个队列
	q := NewPriorityQueue(3, func(src int, dst int) int {
		if src < dst {
			return -1
		}
		if src > dst {
			return 1
		}
		return 0
	})
	// 入队
	q.Enqueue(1)
	q.Enqueue(2)

	val, err := q.Peek()
	assert.Equal(t, nil, err)
	assert.Equal(t, val, 1)

	q.Dequeue()
	q.Dequeue()

	val, err = q.Peek()
	assert.Equal(t, ErrEmptyQueue, err)
	assert.Equal(t, val, 0)
}

func TestPriorityQueue_Len(t *testing.T) {
	q := NewPriorityQueue(3, func(src int, dst int) int {
		if src < dst {
			return -1
		}
		if src > dst {
			return 1
		}
		return 0
	})
	// 入队
	q.Enqueue(1)
	q.Enqueue(2)

	assert.Equal(t, 2, q.Len())
}

func TestPriorityQueue_Cap(t *testing.T) {
	testCases := []struct {
		name     string
		capacity int
		wantCap  int
	}{
		{
			name:     "boundless queue",
			capacity: 0,
			wantCap:  0,
		},
		{
			name:     "bounded queue",
			capacity: 10,
			wantCap:  10,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewPriorityQueue(tc.capacity, func(src int, dst int) int {
				if src < dst {
					return -1
				}
				if src > dst {
					return 1
				}
				return 0
			})
			assert.Equal(t, tc.wantCap, q.Cap())
		})
	}
}

func TestPriorityQueue_IsBoundless(t *testing.T) {
	testCases := []struct {
		name        string
		capacity    int
		IsBoundless bool
	}{
		{
			name:        "boundless queue",
			capacity:    0,
			IsBoundless: true,
		},
		{
			name:        "bounded queue",
			capacity:    10,
			IsBoundless: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewPriorityQueue(tc.capacity, func(src int, dst int) int {
				if src < dst {
					return -1
				}
				if src > dst {
					return 1
				}
				return 0
			})
			assert.Equal(t, tc.IsBoundless, q.IsBoundless())
		})
	}
}

func TestPriorityQueue_Clean(t *testing.T) {
	testCases := []struct {
		name     string
		capacity int
		wantCap  int
	}{
		{
			name:     "boundless queue",
			capacity: 0,
			wantCap:  64,
		}, {
			name:     "bounded queue",
			capacity: 10,
			wantCap:  10,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewPriorityQueue(tc.capacity, func(src int, dst int) int {
				if src < dst {
					return -1
				}
				if src > dst {
					return 1
				}
				return 0
			})
			q.Enqueue(1)
			q.Enqueue(2)
			assert.Equal(t, 2, q.Len())

			q.Clean()
			assert.Equal(t, 0, q.Len())
			assert.Equal(t, tc.wantCap, q.capacity)
		})
	}
}
