package d04

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 4 of Advent of Code 2025.
func PartOne(r io.Reader, w io.Writer) error {
	diagram, err := diagramFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	removableRolls := findRemovableRolls(diagram)

	_, err = fmt.Fprintf(w, "%d", len(removableRolls))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 4 of Advent of Code 2025.
func PartTwo(r io.Reader, w io.Writer) error {
	diagram, err := diagramFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := countTotalRemovableRolls(diagram)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func countTotalRemovableRolls(diagram [][]rune) int {
	count := 0

	for {
		removableRolls := findRemovableRolls(diagram)
		if len(removableRolls) == 0 {
			break
		}
		count += len(removableRolls)

		removeRolls(diagram, removableRolls)
	}

	return count
}

func findRemovableRolls(diagram [][]rune) []coordinate {
	var removableRolls []coordinate

	for row := range diagram {
		for col := range diagram[row] {
			if diagram[row][col] != roll {
				continue
			}

			position := coordinate{row, col}

			adjacentRolls := countAdjacentRolls(diagram, position)
			if adjacentRolls < 4 {
				removableRolls = append(removableRolls, position)
			}
		}
	}

	return removableRolls
}

func removeRolls(diagram [][]rune, positions []coordinate) {
	for _, position := range positions {
		diagram[position.row][position.col] = empty
	}
}

func countAdjacentRolls(diagram [][]rune, position coordinate) int {
	count := 0

	for _, neighbor := range position.neighbors() {
		if neighbor.row < 0 || neighbor.row >= len(diagram) {
			continue
		}
		if neighbor.col < 0 || neighbor.col >= len(diagram[neighbor.row]) {
			continue
		}

		if diagram[neighbor.row][neighbor.col] == roll {
			count++
		}
	}

	return count
}

type coordinate struct {
	row, col int
}

func (c coordinate) neighbors() []coordinate {
	neighbors := make([]coordinate, 0, 8)

	for dRow := -1; dRow <= 1; dRow++ {
		for dCol := -1; dCol <= 1; dCol++ {
			if dRow == 0 && dCol == 0 {
				continue
			}

			neighbor := coordinate{row: c.row + dRow, col: c.col + dCol}
			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors
}

const (
	empty = '.'
	roll  = '@'
)

func diagramFromReader(r io.Reader) ([][]rune, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	diagram := make([][]rune, len(lines))
	for row := range lines {
		diagram[row] = []rune(lines[row])
		for col := range diagram[row] {
			cell := diagram[row][col]
			if cell != empty && cell != roll {
				return nil, fmt.Errorf("invalid cell at row %d, col %d", row, col)
			}
		}
	}

	return diagram, nil
}
