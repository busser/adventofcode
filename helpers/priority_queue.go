package helpers

type PriorityQueue[T any] struct {
	items []T
	less  func(a, b T) bool
}

func NewPriorityQueue[T any](less func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		items: make([]T, 0),
		less:  less,
	}
}

func (pq PriorityQueue[T]) Len() int {
	return len(pq.items)
}

func (pq *PriorityQueue[T]) Push(x T) {
	pq.items = append(pq.items, x)
	pq.up(len(pq.items) - 1)
}

func (pq *PriorityQueue[T]) Pop() T {
	n := len(pq.items)
	pq.swap(0, n-1)
	popped := pq.items[n-1]
	pq.items[n-1] = *new(T)
	pq.items = pq.items[0 : n-1]
	pq.down(0)
	return popped
}

func (pq *PriorityQueue[T]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !pq.less(pq.items[j], pq.items[i]) {
			break
		}
		pq.swap(i, j)
		j = i
	}
}

func (pq *PriorityQueue[T]) down(i0 int) {
	n := len(pq.items)
	i := i0
	for {
		left := 2*i + 1
		if left >= n || left < 0 { // left < 0 after int overflow
			break
		}
		j := left // left child
		if right := left + 1; right < n && pq.less(pq.items[right], pq.items[left]) {
			j = right // = 2*i + 2  // right child
		}
		if !pq.less(pq.items[j], pq.items[i]) {
			break
		}
		pq.swap(i, j)
		i = j
	}
}

func (pq *PriorityQueue[T]) swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}
