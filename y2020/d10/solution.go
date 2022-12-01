package busser

import (
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 10 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	joltages, err := joltagesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var diffsOf1, diffsOf3 int

	for i := range joltages[:len(joltages)-1] {
		diff := joltages[i+1] - joltages[i]

		switch diff {
		case 1:
			diffsOf1++
		case 3:
			diffsOf3++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", diffsOf1*diffsOf3)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 10 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	joltages, err := joltagesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	arrangements := possibleArrangements(joltages)

	_, err = fmt.Fprintf(answer, "%d", arrangements)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func possibleArrangements(joltages []int) int {
	arrangementsByStartIndex := map[int]int{
		len(joltages) - 1: 1,
	}

	var helper func(int) int
	helper = func(startIndex int) int {
		if arrangements, ok := arrangementsByStartIndex[startIndex]; ok {
			return arrangements
		}

		joltage := joltages[startIndex]

		var arrangements int
		for i := startIndex + 1; i < len(joltages) && joltages[i]-joltage <= 3; i++ {
			arrangements += helper(i)
		}

		arrangementsByStartIndex[startIndex] = arrangements
		return arrangements
	}

	return helper(0)
}

func joltagesFromReader(r io.Reader) ([]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading lines: %w", err)
	}

	joltages := make([]int, len(lines)+2)
	var max int

	for i, l := range lines {
		n, err := strconv.Atoi(l)
		if err != nil {
			return nil, fmt.Errorf("wrong format: not an int: %q", l)
		}

		joltages[i] = n

		if n > max {
			max = n
		}
	}

	joltages[len(joltages)-1] = max + 3

	sort.Ints(joltages)

	return joltages, nil
}
