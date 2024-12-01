package d01

import (
	"fmt"
	"io"
	"slices"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 1 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	listA, listB, err := readLists(r)
	if err != nil {
		return err
	}

	slices.Sort(listA)
	slices.Sort(listB)

	distance := distanceBetweenLists(listA, listB)

	_, err = fmt.Fprintf(w, "%d", distance)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	listA, listB, err := readLists(r)
	if err != nil {
		return err
	}

	similarity := similarityBetweenLists(listA, listB)

	_, err = fmt.Fprintf(w, "%d", similarity)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func distanceBetweenNumbers(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

// Assumes that a and b are sorted and of equal length.
func distanceBetweenLists(a, b []int) int {
	var distance int
	for i := range a {
		distance += distanceBetweenNumbers(a[i], b[i])
	}
	return distance
}

func similarityBetweenLists(a, b []int) int {
	occurences := make(map[int]int)
	for _, n := range b {
		occurences[n]++
	}

	score := 0
	for _, n := range a {
		score += n * occurences[n]
	}

	return score
}

func readLists(r io.Reader) ([]int, []int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read input: %w", err)
	}

	var (
		listA = make([]int, len(lines))
		listB = make([]int, len(lines))
	)

	for i, line := range lines {
		nums := helpers.IntsFromString(line)
		if len(nums) != 2 {
			return nil, nil, fmt.Errorf("expected 2 numbers, got %d", len(nums))
		}
		listA[i] = nums[0]
		listB[i] = nums[1]
	}

	return listA, listB, nil
}
