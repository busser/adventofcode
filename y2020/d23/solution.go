package busser

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// PartOne solves the first problem of day 23 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	circle, err := cupsFromReader(input, 9)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	for i := 0; i < 100; i++ {
		circle.moveCups()
	}

	_, err = fmt.Fprintf(answer, "%s", circle.labelsAfter(1))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 23 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	circle, err := cupsFromReader(input, 1e6)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	for i := 0; i < 1e7; i++ {
		circle.moveCups()
	}

	cupOne := circle.getCupFromIndex(1)
	product := cupOne.after.label * cupOne.after.after.label

	_, err = fmt.Fprintf(answer, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type cup struct {
	label         int
	before, after *cup
}

type cupCircle struct {
	size               int
	currentCup         *cup
	cupMinLabel        int
	cupsIndexedByLabel []*cup
}

type cupSequence struct {
	first, last *cup
}

func (seq cupSequence) containsLabel(label int) bool {
	return label == seq.first.label || label == seq.first.after.label || label == seq.last.label
}

func (circle *cupCircle) moveCups() {
	seq := popSequenceAfter(circle.currentCup)

	destinationLabel := circle.labelMinusOne(circle.currentCup.label)
	for seq.containsLabel(destinationLabel) {
		destinationLabel = circle.labelMinusOne(destinationLabel)
	}

	destinationCup := circle.getCupFromIndex(destinationLabel)
	insertSequenceAfter(destinationCup, seq)

	circle.currentCup = circle.currentCup.after
}

func (circle cupCircle) labelMinusOne(label int) int {
	label = label - 1
	if label < circle.cupMinLabel {
		label += circle.size
	}
	return label
}

func popSequenceAfter(target *cup) cupSequence {
	seq := cupSequence{
		first: target.after,
		last:  target.after.after.after,
	}

	seq.first.before.after = seq.last.after
	seq.last.after.before = seq.first.before

	seq.first.before = nil
	seq.last.after = nil

	return seq
}

func insertSequenceAfter(target *cup, seq cupSequence) {
	seq.first.before = target
	seq.last.after = target.after

	target.after = seq.first
	seq.last.after.before = seq.last
}

func (circle cupCircle) String() string {
	if circle.currentCup == nil {
		return ""
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d", circle.currentCup.label)
	for c := circle.currentCup.after; c != circle.currentCup; c = c.after {
		fmt.Fprintf(&sb, "%d", c.label)
	}
	return sb.String()
}

func (circle cupCircle) labelsAfter(label int) string {
	first := circle.getCupFromIndex(label)
	var sb strings.Builder
	for c := first.after; c != first; c = c.after {
		fmt.Fprintf(&sb, "%d", c.label)
	}
	return sb.String()
}

func (circle *cupCircle) insertCupBeforeCurrent(c *cup) {
	circle.addCupToIndex(c)

	if circle.currentCup == nil { // c is the first cup added to cc
		c.before, c.after = c, c
		circle.currentCup = c
		return
	}

	// connect the new cup to the one to its left
	c.before = circle.currentCup.before
	c.before.after = c

	// connect the new cup to the one to its right
	c.after = circle.currentCup
	c.after.before = c
}

func (circle *cupCircle) addCupToIndex(c *cup) {
	circle.cupsIndexedByLabel[c.label-1] = c
}

func (cc *cupCircle) getCupFromIndex(label int) *cup {
	return cc.cupsIndexedByLabel[label-1]
}

func cupsFromReader(r io.Reader, size int) (cupCircle, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return cupCircle{}, err
	}
	raw = bytes.TrimSpace(raw)

	for _, b := range raw {
		if b < '0' || b > '9' {
			return cupCircle{}, fmt.Errorf("expected a digit, got %q", b)
		}
	}

	circle := cupCircle{
		cupMinLabel:        1,
		size:               size,
		cupsIndexedByLabel: make([]*cup, size),
	}

	for _, b := range raw {
		label := int(b - '0')
		c := cup{label, nil, nil}
		circle.insertCupBeforeCurrent(&c)
	}

	for label := len(raw) + 1; label <= size; label++ {
		c := cup{label, nil, nil}
		circle.insertCupBeforeCurrent(&c)
	}

	return circle, nil
}
