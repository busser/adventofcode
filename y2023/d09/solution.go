package d09

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 9 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	sequences, err := sequencesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for _, sequence := range sequences {
		sum += nextValue(sequence)
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 9 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	sequences, err := sequencesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for _, sequence := range sequences {
		sum += previousValue(sequence)
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func previousValue(sequence []int) int {
	if allZeroes(sequence) {
		return 0
	}

	d := diffs(sequence)
	previousDiff := previousValue(d)
	return sequence[0] - previousDiff
}

func nextValue(sequence []int) int {
	if allZeroes(sequence) {
		return 0
	}

	d := diffs(sequence)
	nextDiff := nextValue(d)
	return sequence[len(sequence)-1] + nextDiff
}

func diffs(sequence []int) []int {
	diffs := make([]int, len(sequence)-1)
	for i := 0; i < len(sequence)-1; i++ {
		diffs[i] = sequence[i+1] - sequence[i]
	}
	return diffs
}

func allZeroes(sequence []int) bool {
	for _, v := range sequence {
		if v != 0 {
			return false
		}
	}
	return true
}

func sequencesFromReader(r io.Reader) ([][]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	sequences := make([][]int, len(lines))
	for i, line := range lines {
		sequences[i], err = helpers.IntsFromString(line, " ")
		if err != nil {
			return nil, fmt.Errorf("could not parse line %d: %w", i, err)
		}
	}

	return sequences, nil
}
