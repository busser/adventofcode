package busser

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 7 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	positions, err := positionsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read crab positions: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", minimumFuelCost(positions, naiveCost))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 7 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	positions, err := positionsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read crab positions: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", minimumFuelCost(positions, actualCost))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type costFunction func(positions []int, newPosition int) (fuelCost int)

func minimumFuelCost(positions []int, cost costFunction) int {
	minPosition, maxPosition := rangeOf(positions)

	minFuelCost := cost(positions, minPosition)
	for p := minPosition + 1; p <= maxPosition; p++ {
		fuelCost := cost(positions, p)
		if fuelCost < minFuelCost {
			minFuelCost = fuelCost
		}
	}

	return minFuelCost
}

func naiveCost(positions []int, newPosition int) int {
	cost := 0
	for _, p := range positions {
		cost += abs(newPosition - p)
	}
	return cost
}

func actualCost(positions []int, newPosition int) int {
	cost := 0
	for _, p := range positions {
		cost += sumOf1ToN(abs(newPosition - p))
	}
	return cost
}

func sumOf1ToN(n int) int {
	return n * (n + 1) / 2
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func rangeOf(nn []int) (min, max int) {
	min, max = nn[0], nn[0]
	for _, n := range nn[1:] {
		if n > max {
			max = n
		}
		if n < min {
			min = n
		}
	}
	return min, max
}

func positionsFromReader(r io.Reader) ([]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read: %w", err)
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("expected 1 line of input, got %d", len(lines))
	}

	return helpers.IntsFromString(lines[0], ",")
}
