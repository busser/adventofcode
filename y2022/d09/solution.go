package d09

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 9 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	motions, err := motionsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := countUniqueTailPositions(motions, 2)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 9 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	motions, err := motionsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := countUniqueTailPositions(motions, 10)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func countUniqueTailPositions(motions []motion, ropeLength int) int {
	if ropeLength < 1 {
		return 0
	}

	rope := make([]vector, ropeLength)

	tailPositions := make(map[vector]struct{})
	tailPositions[rope[ropeLength-1]] = struct{}{}

	// Move head of rope and have the rest of the rope follow.
	for _, m := range motions {
		for step := 0; step < m.count; step++ {
			rope[0] = rope[0].plus(m.dir)

			for i := 1; i < ropeLength; i++ {
				rope[i] = tailFollowHead(rope[i], rope[i-1])
			}

			tailPositions[rope[ropeLength-1]] = struct{}{}
		}
	}

	return len(tailPositions)
}

func tailFollowHead(tail, head vector) vector {
	diff := head.minus(tail)

	if abs(diff.x) <= 1 && abs(diff.y) <= 1 {
		// The tail is close enough to the head, it does not need to move.
		return tail
	}

	return tail.plus(diff.unit())
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type vector struct {
	x, y int
}

func (v vector) unit() vector {
	var w vector

	switch {
	case v.x < 0:
		w.x = -1
	case v.x > 0:
		w.x = 1
	}

	switch {
	case v.y < 0:
		w.y = -1
	case v.y > 0:
		w.y = 1
	}

	return w
}

func (v vector) plus(w vector) vector {
	return vector{
		x: v.x + w.x,
		y: v.y + w.y,
	}
}

func (v vector) minus(w vector) vector {
	return vector{
		x: v.x - w.x,
		y: v.y - w.y,
	}
}

var (
	up    = vector{0, -1}
	down  = vector{0, 1}
	left  = vector{-1, 0}
	right = vector{1, 0}
)

type motion struct {
	dir   vector
	count int
}

func motionsFromReader(r io.Reader) ([]motion, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	motions := make([]motion, len(lines))

	for i := range lines {
		m, err := motionFromString(lines[i])
		if err != nil {
			return nil, err
		}

		motions[i] = m
	}

	return motions, nil
}

func motionFromString(s string) (motion, error) {
	parts := strings.SplitN(s, " ", 2)
	if len(parts) != 2 {
		return motion{}, errors.New("wrong format")
	}

	var m motion

	switch parts[0] {
	case "U":
		m.dir = up
	case "D":
		m.dir = down
	case "L":
		m.dir = left
	case "R":
		m.dir = right
	default:
		return motion{}, fmt.Errorf("unknown direction %q", parts[0])
	}

	n, err := strconv.Atoi(parts[1])
	if err != nil {
		return motion{}, fmt.Errorf("%q is not a number", parts[1])
	}

	m.count = n

	return m, nil
}
