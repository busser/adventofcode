package d18

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 18 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	incomingBytes, err := incomingBytesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	steps := shortestPathAfter(incomingBytes, 1024)

	_, err = fmt.Fprintf(w, "%d", steps)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 18 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	incomingBytes, err := incomingBytesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var lastByte vector
	found := false
	for after := 1024; after < len(incomingBytes); after++ {
		steps := shortestPathAfter(incomingBytes, after)
		if steps == -1 {
			lastByte = incomingBytes[after-1]
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("there is always a path")
	}

	_, err = fmt.Fprintf(w, "%d,%d", lastByte.row, lastByte.col)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type vector struct {
	row, col int
}

func (v vector) plus(w vector) vector {
	return vector{
		row: v.row + w.row,
		col: v.col + w.col,
	}
}

var (
	up            = vector{row: -1, col: 0}
	down          = vector{row: 1, col: 0}
	left          = vector{row: 0, col: -1}
	right         = vector{row: 0, col: 1}
	allDirections = []vector{up, down, left, right}
)

const maxRow, maxCol = 70, 70

func shortestPathAfter(incomingBytes []vector, after int) int {
	corrupted := make([][]bool, maxRow+1)
	for row := range corrupted {
		corrupted[row] = make([]bool, maxCol+1)
	}
	for _, b := range incomingBytes[:after] {
		corrupted[b.row][b.col] = true
	}

	visited := make([][]bool, len(corrupted))
	for row := range visited {
		visited[row] = make([]bool, len(corrupted[row]))
	}

	toVisit := []vector{{0, 0}}
	var nextToVisit []vector

	steps := 0
	for len(toVisit) > 0 {
		for _, pos := range toVisit {
			if pos.row < 0 || pos.row > maxRow || pos.col < 0 || pos.col > maxCol {
				continue
			}
			if corrupted[pos.row][pos.col] {
				continue
			}
			if pos.row == maxRow && pos.col == maxCol {
				return steps
			}

			if visited[pos.row][pos.col] {
				continue
			}
			visited[pos.row][pos.col] = true

			for _, dir := range allDirections {
				nextToVisit = append(nextToVisit, pos.plus(dir))
			}
		}

		toVisit, nextToVisit = nextToVisit, toVisit[:0]
		steps++
	}

	return -1 // no path found
}

func incomingBytesFromReader(r io.Reader) ([]vector, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	incomingBytes := make([]vector, len(lines))
	for i, line := range lines {
		nums := helpers.IntsFromString(line)
		if len(nums) != 2 {
			return nil, fmt.Errorf("invalid vector: %q", line)
		}
		incomingBytes[i] = vector{nums[0], nums[1]}
	}

	return incomingBytes, nil
}
