package busser

import "math"

type priorityQueueItem struct {
	position coordinates
	risk     int
	index    int
}

type priorityQueue struct {
	items           []*priorityQueueItem
	itemsByPosition map[coordinates]*priorityQueueItem
}

func newPriorityQueue() *priorityQueue {
	var pq priorityQueue
	pq.itemsByPosition = make(map[coordinates]*priorityQueueItem)
	return &pq
}

// =============================================================================
// External interface of the priority queue.
// =============================================================================

func (pq priorityQueue) len() int {
	return len(pq.items)
}

func (pq priorityQueue) peek() (position coordinates, risk int) {
	if len(pq.items) == 0 {
		panic("queue underflow")
	}
	return pq.items[0].position, pq.items[0].risk
}

func (pq *priorityQueue) pop() (position coordinates, risk int) {
	if len(pq.items) == 0 {
		panic("queue underflow")
	}

	n := len(pq.items)
	pq.swap(0, n-1)
	popped := pq.items[n-1]
	pq.items[n-1] = nil
	pq.items = pq.items[:n-1]
	pq.fix(0)

	delete(pq.itemsByPosition, popped.position)

	return popped.position, popped.risk
}

func (pq *priorityQueue) push(position coordinates, risk int) {
	if item, ok := pq.itemsByPosition[position]; ok {
		pq.lowerRisk(item.index, risk)
		return
	}

	pushed := new(priorityQueueItem)

	pq.itemsByPosition[position] = pushed
	pq.items = append(pq.items, pushed)

	pushed.position = position
	pushed.risk = math.MaxInt
	pushed.index = len(pq.items) - 1

	pq.lowerRisk(pushed.index, risk)
}

// =============================================================================
// Internal function of the priority queue.
// =============================================================================

func (_ priorityQueue) parent(i int) int {
	return (i - 1) / 2
}

func (_ priorityQueue) left(i int) int {
	return 2*i + 1
}

func (_ priorityQueue) right(i int) int {
	return 2*i + 2
}

func (pq priorityQueue) swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

func (pq priorityQueue) fix(i int) {
	l, r := pq.left(i), pq.right(i)

	lowest := i
	if l < len(pq.items) && pq.items[l].risk < pq.items[lowest].risk {
		lowest = l
	}
	if r < len(pq.items) && pq.items[r].risk < pq.items[lowest].risk {
		lowest = r
	}

	if lowest != i {
		pq.swap(i, lowest)
		pq.fix(lowest)
	}
}

func (pq priorityQueue) lowerRisk(i, risk int) {
	if risk > pq.items[i].risk {
		return
	}
	pq.items[i].risk = risk
	for i > 0 && pq.items[pq.parent(i)].risk > pq.items[i].risk {
		pq.swap(pq.parent(i), i)
		i = pq.parent(i)
	}
}
