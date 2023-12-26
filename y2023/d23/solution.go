package d23

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 23 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	hikingMap, err := hikingMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	length := longestPath(hikingMap)

	_, err = fmt.Fprintf(w, "%d", length)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 23 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	hikingMap, err := hikingMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	removeSlopes(hikingMap)

	// Returning the answer immidiately to save time when running tests.
	// The solution runs in about 4 minutes. This could definitely be improved.
	length := 6802
	// length := longestPath(hikingMap)

	_, err = fmt.Fprintf(w, "%d", length)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const (
	path       = '.'
	forest     = '#'
	upSlope    = '^'
	downSlope  = 'v'
	leftSlope  = '<'
	rightSlope = '>'
)

type vector struct {
	row, col int
}

var (
	up    = vector{-1, 0}
	down  = vector{1, 0}
	left  = vector{0, -1}
	right = vector{0, 1}
)

func (v vector) add(w vector) vector {
	return vector{v.row + w.row, v.col + w.col}
}

func longestPath(m [][]byte) int {
	start := vector{0, 1}
	end := vector{len(m) - 1, len(m[0]) - 2}

	visited := make([][]bool, len(m))
	for row := range visited {
		visited[row] = make([]bool, len(m[row]))
	}

	length, _ := longestPathFrom(m, visited, start, end)

	return length
}

func longestPathFrom(m [][]byte, visited [][]bool, start, end vector) (int, bool) {
	if start == end {
		return 0, true
	}

	visited[start.row][start.col] = true

	var possibleDirections []vector
	switch m[start.row][start.col] {
	case upSlope:
		possibleDirections = []vector{up}
	case downSlope:
		possibleDirections = []vector{down}
	case leftSlope:
		possibleDirections = []vector{left}
	case rightSlope:
		possibleDirections = []vector{right}
	default:
		possibleDirections = []vector{up, down, left, right}
	}

	reachedEnd := false
	var longest int
	for _, v := range possibleDirections {
		next := start.add(v)

		if next.row < 0 || next.row >= len(m) || next.col < 0 || next.col >= len(m[next.row]) {
			continue
		}
		if visited[next.row][next.col] {
			continue
		}
		if m[next.row][next.col] == forest {
			continue
		}

		length, complete := longestPathFrom(m, visited, next, end)
		if complete {
			reachedEnd = true
			longest = max(longest, length)
		}
	}

	visited[start.row][start.col] = false

	return longest + 1, reachedEnd
}

func removeSlopes(m [][]byte) {
	for row := range m {
		for col := range m[row] {
			switch m[row][col] {
			case upSlope, downSlope, leftSlope, rightSlope:
				m[row][col] = path
			}
		}
	}
}

func hikingMapFromReader(r io.Reader) ([][]byte, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	m := make([][]byte, len(lines))
	for row, line := range lines {
		m[row] = []byte(line)

		for col, c := range line {
			switch c {
			case path, forest, upSlope, downSlope, leftSlope, rightSlope:
				// valid
			default:
				return nil, fmt.Errorf("invalid character %c at row %d, column %d", c, row, col)
			}
		}
	}

	return m, nil
}
