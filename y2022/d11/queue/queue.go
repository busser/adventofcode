package queue

type Queue[T any] struct {
	first, last int
	buffer      []T
}

func New[T any](cap int) *Queue[T] {
	return &Queue[T]{
		buffer: make([]T, cap),
	}
}

func (q *Queue[T]) Len() int {
	return (q.last - q.first + len(q.buffer)) % len(q.buffer)
}

func (q *Queue[T]) Enqueue(v T) {
	q.buffer[q.last] = v
	q.last = (q.last + 1) % len(q.buffer)
}

func (q *Queue[T]) Dequeue() T {
	v := q.buffer[q.first]
	q.first = (q.first + 1) % len(q.buffer)
	return v
}
