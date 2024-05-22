package queue

import (
	"errors"
	"github.com/Jared-lu/GXT"
	"github.com/Jared-lu/GXT/internal/heap"
)

// PriorityQueue 优先级队列，支持有界与无界
type PriorityQueue[T any] struct {
	heap     *heap.Heap[T]
	capacity int
	// 有界/无界
	boundless  bool
	comparator GXT.Comparator[T]
}

// NewPriorityQueue
// 当 capacity <= 0 时，为无界队列
// 当 capacity > 0 时，为有界队列
func NewPriorityQueue[T any](capacity int, comparator GXT.Comparator[T]) *PriorityQueue[T] {
	var boundless bool
	if capacity <= 0 {
		// 无界队列初始化容量为64
		capacity = 64
		boundless = true
	}
	data := make([]T, 0, capacity)
	h := heap.NewHeap[T](data, comparator)
	return &PriorityQueue[T]{
		heap:       h,
		capacity:   capacity,
		boundless:  boundless,
		comparator: comparator,
	}
}

func (p *PriorityQueue[T]) Enqueue(val T) error {
	if !p.boundless {
		// 有界队列，在入队前要先查看是否还有空位
		if p.capacity > p.heap.Size() {
			p.heap.Push(val)
		} else {
			return ErrOutOfCapacity
		}
		return nil
	}
	p.heap.Push(val)
	return nil
}

func (p *PriorityQueue[T]) Dequeue() (T, error) {
	val, err := p.heap.Pop()
	if errors.Is(err, heap.ErrEmptyHeap) {
		return val, ErrEmptyQueue
	}
	return val, err
}

func (p *PriorityQueue[T]) Peek() (T, error) {
	val, err := p.heap.Peek()
	if errors.Is(err, heap.ErrEmptyHeap) {
		return val, ErrEmptyQueue
	}
	return val, err
}

func (p *PriorityQueue[T]) Len() int {
	return p.heap.Size()
}

// Cap 无界队列返回0
func (p *PriorityQueue[T]) Cap() int {
	if p.boundless {
		return 0
	}
	return p.capacity
}

func (p *PriorityQueue[T]) IsBoundless() bool {
	return p.boundless
}

func (p *PriorityQueue[T]) Clean() {
	data := make([]T, 0, p.capacity)
	h := heap.NewHeap[T](data, p.comparator)
	p.heap = h
}
