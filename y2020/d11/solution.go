package busser

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 11 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	layout, err := layoutFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	current, next := layout, copyLayout(layout)

	for nextLayoutFromCurrent(current, next, occupiedSeatsAround, 4) {
		current, next = next, current
	}

	count := occupiedSeatsCount(current)

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 11 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	layout, err := layoutFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	current, next := layout, copyLayout(layout)

	for nextLayoutFromCurrent(current, next, occupiedSeatsInSight, 5) {
		current, next = next, current
	}

	count := occupiedSeatsCount(current)

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}
	return nil
}

type seat uint8

const (
	floor seat = iota
	empty
	occupied
)

func nextLayoutFromCurrent(current, next [][]seat, seatCount func([][]seat, int, int) int, maxTolerance int) bool {
	changeOccured := false

	for i := range current {
		for j := range current[i] {
			if current[i][j] == floor {
				next[i][j] = floor
				continue
			}

			occupiedSeats := seatCount(current, i, j)
			switch {
			case current[i][j] == empty && occupiedSeats == 0:
				next[i][j] = occupied
				changeOccured = true
			case current[i][j] == occupied && occupiedSeats >= maxTolerance:
				next[i][j] = empty
				changeOccured = true
			default:
				next[i][j] = current[i][j]
			}
		}
	}

	return changeOccured
}

func occupiedSeatsAround(layout [][]seat, i, j int) int {
	count := 0

	for ii := i - 1; ii <= i+1; ii++ {
		if ii < 0 || ii >= len(layout) {
			continue
		}

		for jj := j - 1; jj <= j+1; jj++ {
			if i == ii && j == jj {
				continue
			}
			if jj < 0 || jj >= len(layout[ii]) {
				continue
			}

			if layout[ii][jj] == occupied {
				count++
			}
		}
	}

	return count
}

func occupiedSeatsInSight(layout [][]seat, i, j int) int {
	count := 0

	directions := [...]struct {
		i, j int
	}{
		{-1, 0},  // up
		{1, 0},   // down
		{0, -1},  // left
		{0, 1},   // right
		{-1, -1}, // top left
		{-1, 1},  // top right
		{1, -1},  // bottom left
		{1, 1},   // bottom right
	}

	for _, dir := range directions {
		ii, jj := i, j

		for {
			ii, jj = ii+dir.i, jj+dir.j
			if ii < 0 || ii >= len(layout) || jj < 0 || jj >= len(layout[ii]) {
				break
			}

			s := layout[ii][jj]

			if s == occupied {
				count++
				break
			}
			if s == empty {
				break
			}
		}
	}

	return count
}

func occupiedSeatsCount(layout [][]seat) int {
	count := 0

	for _, row := range layout {
		for _, s := range row {
			if s == occupied {
				count++
			}
		}
	}

	return count
}

func copyLayout(layout [][]seat) [][]seat {
	new := make([][]seat, len(layout))

	for i := range layout {
		new[i] = make([]seat, len(layout[i]))
		copy(new[i], layout[i])
	}

	return new
}

func layoutFromReader(r io.Reader) ([][]seat, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading lines: %w", err)
	}

	layout := make([][]seat, len(lines))

	for i, line := range lines {
		layout[i] = make([]seat, len(line))
		for j, c := range line {
			switch c {
			case '.':
				layout[i][j] = floor
			case 'L':
				layout[i][j] = empty
			case '#':
				layout[i][j] = occupied
			default:
				return nil, fmt.Errorf("wrong format: unknown symbol %q", c)
			}
		}
	}

	return layout, nil
}
