package d01

import (
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 1 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	inventories, err := elfInventoriesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_, err = fmt.Fprintf(w, "%d", highestCalorieCount(inventories))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	inventories, err := elfInventoriesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_, err = fmt.Fprintf(w, "%d", sumOfTop3CalorieCounts(inventories))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func elfInventoriesFromReader(r io.Reader) ([]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	var inventories []int
	currentInventory := 0

	for _, l := range lines {
		if len(l) == 0 {
			inventories = append(inventories, currentInventory)
			currentInventory = 0
			continue
		}

		calories, err := strconv.Atoi(l)
		if err != nil {
			return nil, fmt.Errorf("%q is not a number", l)
		}

		currentInventory += calories
	}
	inventories = append(inventories, currentInventory)

	if len(inventories) < 3 {
		return nil, fmt.Errorf("not enough elves (need 3, have %d)", len(inventories))
	}

	return inventories, nil
}

func highestCalorieCount(counts []int) int {
	max := 0
	for _, c := range counts {
		if c > max {
			max = c
		}
	}
	return max
}

func sumOfTop3CalorieCounts(counts []int) int {
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))

	return counts[0] + counts[1] + counts[2]
}
