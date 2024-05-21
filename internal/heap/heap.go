package heap

import (
	"errors"
	"github.com/Jared-lu/GXT"
)

var ErrEmptyHeap = errors.New("empty heap")

type Heap[T any] struct {
	// 调用者决定如何比较两个元素之间的大小
	comparator GXT.Comparator[T]
	data       []T
}

// NewHeap 生成小根堆或大根堆
// 正常的 Comparator src < dst return -1，生成小根堆
// 利用规则反转，即Comparator src < dst return 1时，可以生成大根堆
func NewHeap[T any](data []T, comparator GXT.Comparator[T]) *Heap[T] {
	heap := &Heap[T]{data: data, comparator: comparator}
	// 通过倒序遍历来建堆，不断往上调整更大的子树
	for i := heap.parent(len(heap.data) - 1); i >= 0; i-- {
		// 堆化除叶节点以外的其他所有节点
		heap.siftDown(i)
	}
	return heap
}

func (h *Heap[T]) Push(ele T) {
	h.data = append(h.data, ele)
	h.siftUp(len(h.data) - 1)
}

func (h *Heap[T]) Pop() (T, error) {
	if len(h.data) == 0 {
		var t T
		return t, ErrEmptyHeap
	}
	val := h.data[0]
	h.data[0], h.data[len(h.data)-1] = h.data[len(h.data)-1], h.data[0]
	h.data = h.data[:len(h.data)-1]
	h.siftDown(0)
	return val, nil
}

// Peek 返回堆顶元素
func (h *Heap[T]) Peek() (T, error) {
	if len(h.data) == 0 {
		var t T
		return t, ErrEmptyHeap
	}
	return h.data[0], nil
}

func (h *Heap[T]) Size() int {
	return len(h.data)
}

func (h *Heap[T]) IsEmpty() bool {
	return len(h.data) == 0
}

// left 获取左子节点的索引
func (h *Heap[T]) left(i int) int {
	return 2*i + 1
}

// right 获取右子节点的索引
func (h *Heap[T]) right(i int) int {
	return 2*i + 2
}

// parent 获取父节点的索引
func (h *Heap[T]) parent(i int) int {
	return (i - 1) / 2
}

// siftUp 自底向上堆化
// 子节点和它的父节点进行比较
func (h *Heap[T]) siftUp(i int) {
	for {
		parent := h.parent(i)
		// 与父节点进行比较，小于父节点就说明当前子树要调整
		if parent >= 0 && h.comparator(h.data[i], h.data[parent]) == -1 {
			h.data[i], h.data[parent] = h.data[parent], h.data[i]
		} else {
			break
		}
		// 循环向上堆化
		i = parent
	}
}

// siftDown 自顶向上堆化
// 根节点和它的左右子节点进行比较
func (h *Heap[T]) siftDown(i int) {
	for {
		left, right := h.left(i), h.right(i)
		var temp int
		if left < h.Size() && h.comparator(h.data[left], h.data[i]) == -1 {
			temp = left
		} else {
			temp = i
		}
		if right < h.Size() && h.comparator(h.data[right], h.data[temp]) == -1 {
			temp = right
		}
		if temp == i {
			break
		}
		h.data[i], h.data[temp] = h.data[temp], h.data[i]
		// 交换元素后向下堆化，调整交换后的子树
		i = temp
	}
}
