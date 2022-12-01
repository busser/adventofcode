package busser

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 3 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	treeMap, err := treeMapFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	s := slope{3, 1}
	count := treesOnSlope(treeMap, s)

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 3 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	treeMap, err := treeMapFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	slopes := []slope{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}

	product := 1
	for _, s := range slopes {
		product *= treesOnSlope(treeMap, s)
	}

	_, err = fmt.Fprintf(answer, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwoParallel solves the second problem of day 3 of Advent of Code 2020.
func PartTwoParallel(input io.Reader, answer io.Writer) error {
	treeMap, err := treeMapFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	slopes := []slope{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}

	treeCounts := make(chan int, len(slopes))
	defer close(treeCounts)

	for _, s := range slopes {
		go func(s slope) {
			treeCounts <- treesOnSlope(treeMap, s)
		}(s)
	}

	product := 1
	for range slopes {
		product *= <-treeCounts
	}

	_, err = fmt.Fprintf(answer, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type slope struct {
	right, down int
}

func treeMapFromReader(r io.Reader) ([][]bool, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading lines: %w", err)
	}

	treeMap := make([][]bool, len(lines))
	for i := range lines {
		treeMap[i] = make([]bool, len(lines[i]))
		for j, c := range lines[i] {
			treeMap[i][j] = c == '#'
		}
	}

	return treeMap, nil
}

func treesOnSlope(treeMap [][]bool, s slope) int {
	x, y, count := 0, 0, 0

	for y < len(treeMap) {
		if treeMap[y][x] {
			count++
		}
		x, y = (x+s.right)%len(treeMap[0]), y+s.down
	}

	return count
}
