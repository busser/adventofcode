package d01

import (
	"fmt"
	"io"
	"strconv"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 1 of Advent of Code 2019.
func PartOne(r io.Reader, w io.Writer) error {
	modules, err := modulesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	totalFuel := 0
	for _, mod := range modules {
		totalFuel += requiredFuel(mod)
	}

	_, err = fmt.Fprintf(w, "%d", totalFuel)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2019.
func PartTwo(r io.Reader, w io.Writer) error {
	modules, err := modulesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	totalFuel := 0
	for _, mod := range modules {
		totalFuel += recursiveRequiredFuel(mod)
	}

	_, err = fmt.Fprintf(w, "%d", totalFuel)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func modulesFromReader(r io.Reader) ([]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	var modules []int

	for _, l := range lines {
		mod, err := strconv.Atoi(l)
		if err != nil {
			return nil, fmt.Errorf("%q is not a number", l)
		}

		modules = append(modules, mod)
	}

	return modules, nil
}

func requiredFuel(mass int) int {
	return mass/3 - 2
}

func recursiveRequiredFuel(mass int) int {
	f := mass/3 - 2
	if f <= 0 {
		return 0
	}
	return f + recursiveRequiredFuel(f)
}
