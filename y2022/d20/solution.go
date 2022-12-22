package d20

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// PartOne solves the first problem of day 20 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	ints, err := intsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	list := newList(ints)
	list.mix()
	coords := list.coordinates()

	_, err = fmt.Fprintf(w, "%d", coords)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 20 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	ints, err := intsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	list := newList(ints)
	list.applyDecryptionKey(811589153)
	for i := 0; i < 10; i++ {
		list.mix()
	}
	coords := list.coordinates()

	_, err = fmt.Fprintf(w, "%d", coords)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type list struct {
	nodes []*node // not in order
}

type node struct {
	prev, next *node
	list       *list
	value      int
}

func (l *list) mix() {
	for _, n := range l.nodes {
		n.shift()
	}
}

func (n *node) shift() {
	if n.value == 0 {
		return
	}

	// remove n
	n.prev.next = n.next
	n.next.prev = n.prev

	// compute shortest path
	loopSize := len(n.list.nodes) - 1
	offset := n.value % loopSize

	// find new spot
	for ; offset > 0; offset-- {
		n.prev = n.prev.next
	}
	for ; offset < 0; offset++ {
		n.prev = n.prev.prev
	}

	// insert n
	n.next = n.prev.next
	n.prev.next = n
	n.next.prev = n
}

func (l *list) applyDecryptionKey(key int) {
	for _, n := range l.nodes {
		n.value *= key
	}
}

func (l *list) coordinates() int {
	n := l.find(0)
	sum := 0
	for i := 1; i <= 3000; i++ {
		n = n.next
		if i%1000 == 0 {
			sum += n.value
		}
	}
	return sum
}

func (l *list) find(v int) *node {
	for n := l.nodes[0]; ; n = n.next {
		if n.value == v {
			return n
		}
	}
}

func newList(values []int) *list {
	l := list{
		nodes: make([]*node, len(values)),
	}

	l.nodes[0] = &node{
		list:  &l,
		value: values[0],
	}

	for i := 1; i < len(values); i++ {
		l.nodes[i] = &node{
			list:  &l,
			prev:  l.nodes[i-1],
			value: values[i],
		}
		l.nodes[i-1].next = l.nodes[i]
	}

	last := len(values) - 1
	l.nodes[last].next = l.nodes[0]
	l.nodes[0].prev = l.nodes[last]

	return &l
}

func intsFromReader(r io.Reader) ([]int, error) {
	var ints []int

	zeroCount := 0

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		s := sc.Text()
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("%q is not a number", s)
		}

		if n == 0 {
			zeroCount++
		}

		ints = append(ints, n)
	}

	if zeroCount != 1 {
		return nil, fmt.Errorf("expected a single 0, have %d", zeroCount)
	}

	return ints, nil
}
