package d13

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 13 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	patterns, err := patternsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	verticalMirrorSum := 0
	horizontalMirrorSum := 0

	for _, pattern := range patterns {
		verticalMirrors := findVerticalMirrors(pattern)
		for _, mirror := range verticalMirrors {
			verticalMirrorSum += mirror
		}

		horizontalMirrors := findHorizontalMirrors(pattern)
		for _, mirror := range horizontalMirrors {
			horizontalMirrorSum += mirror
		}
	}

	_, err = fmt.Fprintf(w, "%d", verticalMirrorSum+100*horizontalMirrorSum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 13 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	patterns, err := patternsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	verticalMirrorSum := 0
	horizontalMirrorSum := 0

	for _, pattern := range patterns {
		verticalMirrors := findUnfudgedVerticalMirrors(pattern)
		for _, mirror := range verticalMirrors {
			verticalMirrorSum += mirror
		}

		horizontalMirrors := findUnfudgedHorizontalMirrors(pattern)
		for _, mirror := range horizontalMirrors {
			horizontalMirrorSum += mirror
		}
	}

	_, err = fmt.Fprintf(w, "%d", verticalMirrorSum+100*horizontalMirrorSum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func unfudge(pattern [][]byte, row, col int) {
	switch pattern[row][col] {
	case '.':
		pattern[row][col] = '#'
	case '#':
		pattern[row][col] = '.'
	}
}

func deduplicateMirrors(mirrors []int) []int {
	var deduplicated []int

	for _, m := range mirrors {
		found := false
		for _, d := range deduplicated {
			if m == d {
				found = true
				break
			}
		}

		if !found {
			deduplicated = append(deduplicated, m)
		}
	}

	return deduplicated
}

func findUnfudgedHorizontalMirrors(pattern [][]byte) []int {
	originalMirrors := findHorizontalMirrors(pattern)
	var unfudgedMirrors []int

	for row := range pattern {
		for col := range pattern[row] {
			unfudge(pattern, row, col)
			mirrors := findHorizontalMirrors(pattern)
			unfudge(pattern, row, col)

			// Ignore original mirrors
			for _, m := range originalMirrors {
				for i, mirror := range mirrors {
					if m == mirror {
						mirrors = append(mirrors[:i], mirrors[i+1:]...)
						break
					}
				}
			}

			if len(mirrors) > 0 {
				unfudgedMirrors = append(unfudgedMirrors, mirrors...)
			}
		}
	}

	return deduplicateMirrors(unfudgedMirrors)
}

func findHorizontalMirrors(pattern [][]byte) []int {
	var mirrors []int

	for row := 1; row < len(pattern); row++ {
		if isHorizontalMirrorAtPosition(pattern, row) {
			mirrors = append(mirrors, row)
		}
	}

	return mirrors
}

func isHorizontalMirrorAtPosition(pattern [][]byte, row int) bool {
	mirrorRange := min(row, len(pattern)-row)

	for i := 0; i <= mirrorRange; i++ {
		if !rowsAreEqual(pattern, row-i, row+i-1) {
			return false
		}
	}

	return true
}

func rowsAreEqual(pattern [][]byte, rowA, rowB int) bool {
	for column := range pattern[0] {
		if pattern[rowA][column] != pattern[rowB][column] {
			return false
		}
	}

	return true
}

func findUnfudgedVerticalMirrors(pattern [][]byte) []int {
	originalMirrors := findVerticalMirrors(pattern)
	var unfudgedMirrors []int

	for row := range pattern {
		for col := range pattern[row] {
			unfudge(pattern, row, col)
			mirrors := findVerticalMirrors(pattern)
			unfudge(pattern, row, col)

			// Ignore original mirrors
			for _, m := range originalMirrors {
				for i, mirror := range mirrors {
					if m == mirror {
						mirrors = append(mirrors[:i], mirrors[i+1:]...)
						break
					}
				}
			}

			if len(mirrors) > 0 {
				unfudgedMirrors = append(unfudgedMirrors, mirrors...)
			}
		}
	}

	return deduplicateMirrors(unfudgedMirrors)
}

func findVerticalMirrors(pattern [][]byte) []int {
	var mirrors []int

	for col := 1; col < len(pattern[0]); col++ {
		if isVerticalMirrorAtPosition(pattern, col) {
			mirrors = append(mirrors, col)
		}
	}

	return mirrors
}

func isVerticalMirrorAtPosition(pattern [][]byte, col int) bool {
	mirrorRange := min(col, len(pattern[0])-col)

	for i := 0; i <= mirrorRange; i++ {
		if !columnsAreEqual(pattern, col-i, col+i-1) {
			return false
		}
	}

	return true
}

func columnsAreEqual(pattern [][]byte, colA, colB int) bool {
	for row := range pattern {
		if pattern[row][colA] != pattern[row][colB] {
			return false
		}
	}

	return true
}

func patternsFromReader(r io.Reader) ([][][]byte, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	var patterns [][][]byte
	var currentPattern [][]byte

	for _, line := range lines {
		if line == "" {
			patterns = append(patterns, currentPattern)
			currentPattern = nil
			continue
		}

		currentPattern = append(currentPattern, []byte(line))
	}
	patterns = append(patterns, currentPattern)

	for _, pattern := range patterns {
		if len(pattern) == 0 {
			return nil, fmt.Errorf("pattern is empty")
		}

		rowLength := len(pattern[0])
		for row := range pattern {
			if len(pattern[row]) != rowLength {
				return nil, fmt.Errorf("pattern is not rectangular")
			}

			for col := range pattern[row] {
				if pattern[row][col] != '.' && pattern[row][col] != '#' {
					return nil, fmt.Errorf("pattern contains invalid character")
				}
			}
		}
	}

	return patterns, nil
}
