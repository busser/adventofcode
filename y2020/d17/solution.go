package busser

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 17 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	initialState, err := initialStateFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var sim cubicSimulation

	sim.init(initialState, 6)

	for i := 0; i < 6; i++ {
		sim.iterate()
	}

	_, err = fmt.Fprintf(answer, "%d", sim.tally())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 17 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	initialState, err := initialStateFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var sim hypercubicSimulation

	sim.init(initialState, 6)

	for i := 0; i < 6; i++ {
		sim.iterate()
	}

	_, err = fmt.Fprintf(answer, "%d", sim.tally())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type simulation interface {
	init(initialState [][]bool, margin int)
	iterate()
	tally() int
}

func bounds(state [][]bool) (minX, maxX, minY, maxY int) {
	for x := range state {
		for y := range state[x] {
			if !state[x][y] {
				continue
			}

			minX = min(minX, x)
			maxX = max(maxX, x)
			minY = min(minY, y)
			maxY = max(maxY, y)
		}
	}

	return minX, maxX, minY, maxY
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func initialStateFromReader(r io.Reader) ([][]bool, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading lines: %w", err)
	}

	s := make([][]bool, len(lines))
	for i, line := range lines {
		s[i] = make([]bool, len(line))
		for j, cell := range line {
			switch cell {
			case '#':
				s[i][j] = true
			case '.':
				s[i][j] = false
			default:
				return nil, fmt.Errorf("unknown symbol: %q", cell)
			}
		}
	}

	return s, nil
}
