package busser

import "math"

type priorityQueueItem struct {
	config configuration
	cost   int
	index  int
}

type priorityQueue struct {
	items         []*priorityQueueItem
	itemsByConfig map[configuration]*priorityQueueItem
}

func newPriorityQueue() *priorityQueue {
	var pq priorityQueue
	pq.itemsByConfig = make(map[configuration]*priorityQueueItem)
	return &pq
}

// =============================================================================
// External interface of the priority queue.
// =============================================================================

func (pq priorityQueue) len() int {
	return len(pq.items)
}

func (pq priorityQueue) peek() (config configuration, cost int) {
	if len(pq.items) == 0 {
		panic("queue underflow")
	}
	return pq.items[0].config, pq.items[0].cost
}

func (pq *priorityQueue) pop() (config configuration, cost int) {
	if len(pq.items) == 0 {
		panic("queue underflow")
	}

	n := len(pq.items)
	pq.swap(0, n-1)
	popped := pq.items[n-1]
	pq.items[n-1] = nil
	pq.items = pq.items[:n-1]
	pq.sink(0)

	delete(pq.itemsByConfig, popped.config)

	return popped.config, popped.cost
}

func (pq *priorityQueue) push(config configuration, cost int) {
	if item, ok := pq.itemsByConfig[config]; ok {
		pq.reduceCost(item.index, cost)
		return
	}

	pushed := new(priorityQueueItem)

	pq.itemsByConfig[config] = pushed
	pq.items = append(pq.items, pushed)

	pushed.config = config
	pushed.cost = math.MaxInt
	pushed.index = len(pq.items) - 1

	pq.reduceCost(pushed.index, cost)
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

func (pq priorityQueue) sink(i int) {
	l, r := pq.left(i), pq.right(i)

	lowest := i
	if l < len(pq.items) && pq.items[l].cost < pq.items[lowest].cost {
		lowest = l
	}
	if r < len(pq.items) && pq.items[r].cost < pq.items[lowest].cost {
		lowest = r
	}

	if lowest != i {
		pq.swap(i, lowest)
		pq.sink(lowest)
	}
}

func (pq priorityQueue) swim(i int) {
	for i > 0 && pq.items[pq.parent(i)].cost > pq.items[i].cost {
		pq.swap(pq.parent(i), i)
		i = pq.parent(i)
	}
}

func (pq priorityQueue) reduceCost(i, cost int) {
	if cost > pq.items[i].cost {
		return
	}
	pq.items[i].cost = cost
	pq.swim(i)
}
