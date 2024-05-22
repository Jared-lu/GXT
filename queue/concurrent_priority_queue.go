package queue

import (
	"errors"
	"github.com/Jared-lu/GXT"
	"github.com/Jared-lu/GXT/internal/heap"
	"sync"
)

// ConcurrentPriorityQueue 并发安全的优先级队列
type ConcurrentPriorityQueue[T any] struct {
	heap     *heap.Heap[T]
	capacity int
	// 有界/无界
	boundless  bool
	comparator GXT.Comparator[T]
	r          sync.RWMutex
}

// NewConcurrentPriorityQueue
// 当 capacity <= 0 时，为无界队列
// 当 capacity > 0 时，为有界队列
func NewConcurrentPriorityQueue[T any](capacity int, comparator GXT.Comparator[T]) *ConcurrentPriorityQueue[T] {
	var boundless bool
	if capacity <= 0 {
		// 无界队列初始化容量为64
		capacity = 64
		boundless = true
	}
	data := make([]T, 0, capacity)
	h := heap.NewHeap[T](data, comparator)
	return &ConcurrentPriorityQueue[T]{
		heap:       h,
		capacity:   capacity,
		boundless:  boundless,
		comparator: comparator,
	}
}

func (p *ConcurrentPriorityQueue[T]) Enqueue(val T) error {
	p.r.Lock()
	defer p.r.Unlock()
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

func (p *ConcurrentPriorityQueue[T]) Dequeue() (T, error) {
	p.r.Lock()
	defer p.r.Unlock()
	val, err := p.heap.Pop()
	if errors.Is(err, heap.ErrEmptyHeap) {
		return val, ErrEmptyQueue
	}
	return val, err
}

func (p *ConcurrentPriorityQueue[T]) Peek() (T, error) {
	p.r.RLock()
	defer p.r.RUnlock()
	val, err := p.heap.Peek()
	if errors.Is(err, heap.ErrEmptyHeap) {
		return val, ErrEmptyQueue
	}
	return val, err
}

func (p *ConcurrentPriorityQueue[T]) Len() int {
	p.r.RLock()
	defer p.r.RUnlock()
	return p.heap.Size()
}

func (p *ConcurrentPriorityQueue[T]) Cap() int {
	p.r.RLock()
	defer p.r.RUnlock()
	return p.capacity
}

func (p *ConcurrentPriorityQueue[T]) IsBoundless() bool {
	p.r.RLock()
	defer p.r.RUnlock()
	return p.boundless
}

func (p *ConcurrentPriorityQueue[T]) Clean() {
	p.r.Lock()
	defer p.r.Unlock()
	data := make([]T, 0, p.capacity)
	h := heap.NewHeap[T](data, p.comparator)
	p.heap = h
}
