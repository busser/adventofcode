package d11

import (
	"fmt"
	"io"
	"math"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 11 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	galaxies, err := galaxiesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	expanded := expand(galaxies, 2)

	sum := sumOfDistances(expanded)

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 11 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	galaxies, err := galaxiesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	expanded := expand(galaxies, 1_000_000)

	sum := sumOfDistances(expanded)

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type position struct {
	row, col int
}

func sumOfDistances(galaxies []position) int {
	var sum int
	for i := range galaxies {
		for j := i + 1; j < len(galaxies); j++ {
			sum += distance(galaxies[i], galaxies[j])
		}
	}
	return sum
}

func distance(a, b position) int {
	return abs(a.row-b.row) + abs(a.col-b.col)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func expand(galaxies []position, factor int) []position {
	columnsWithGalaxies := make(map[int]bool)
	rowsWithGalaxies := make(map[int]bool)

	for _, galaxy := range galaxies {
		columnsWithGalaxies[galaxy.col] = true
		rowsWithGalaxies[galaxy.row] = true
	}

	minCol, maxCol := keyRange(columnsWithGalaxies)
	minRow, maxRow := keyRange(rowsWithGalaxies)

	var columnsWithoutGalaxies []int
	for col := minCol; col <= maxCol; col++ {
		if !columnsWithGalaxies[col] {
			columnsWithoutGalaxies = append(columnsWithoutGalaxies, col)
		}
	}

	var rowsWithoutGalaxies []int
	for row := minRow; row <= maxRow; row++ {
		if !rowsWithGalaxies[row] {
			rowsWithoutGalaxies = append(rowsWithoutGalaxies, row)
		}
	}

	var expanded []position

	for _, g := range galaxies {
		e := g
		for _, row := range rowsWithoutGalaxies {
			if row > g.row {
				break
			}
			e.row += factor - 1
		}
		for _, col := range columnsWithoutGalaxies {
			if col > g.col {
				break
			}
			e.col += factor - 1
		}
		expanded = append(expanded, e)
	}

	return expanded
}

func keyRange(m map[int]bool) (int, int) {
	min, max := math.MaxInt, math.MinInt
	for key := range m {
		if key < min {
			min = key
		}
		if key > max {
			max = key
		}
	}
	return min, max
}

func galaxiesFromReader(r io.Reader) ([]position, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	var galaxies []position
	for row, line := range lines {
		for col, char := range line {
			if char == '#' {
				galaxies = append(galaxies, position{row, col})
			}
		}
	}

	if len(galaxies) == 0 {
		return nil, fmt.Errorf("no galaxies found")
	}

	return galaxies, nil
}
