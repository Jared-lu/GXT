package queue

type Queue[T any] interface {
	Enqueue(val T) error
	Dequeue() (T, error)
}
