package d07

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 7 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	equations, err := equationsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	total := 0
	for _, e := range equations {
		if e.isSolvable([]operator{opAdd, opMultiply}) {
			total += e.testValue
		}
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 7 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	equations, err := equationsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	total := 0
	for _, e := range equations {
		if e.isSolvable([]operator{opAdd, opMultiply, opConcatenate}) {
			total += e.testValue
		}
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type equation struct {
	testValue int
	numbers   []int
}

type operator uint8

const (
	opAdd operator = iota
	opMultiply
	opConcatenate
)

func (e equation) isSolvable(operators []operator) bool {
	var helper func(i, value int) bool
	helper = func(i, value int) bool {
		if i == len(e.numbers) {
			return value == e.testValue
		}

		for _, op := range operators {
			var newValue int
			switch op {
			case opAdd:
				newValue = value + e.numbers[i]
			case opMultiply:
				newValue = value * e.numbers[i]
			case opConcatenate:
				factor := 1
				for factor <= e.numbers[i] {
					factor *= 10
				}
				newValue = value*factor + e.numbers[i]
			}

			if helper(i+1, newValue) {
				return true
			}
		}

		return helper(i+1, value+e.numbers[i]) || helper(i+1, value*e.numbers[i])
	}

	return helper(1, e.numbers[0])
}

func equationFromString(s string) (equation, error) {
	numbers := helpers.IntsFromString(s)
	if len(numbers) < 2 {
		return equation{}, fmt.Errorf("invalid equation: %q", s)
	}

	return equation{
		testValue: numbers[0],
		numbers:   numbers[1:],
	}, nil
}

func equationsFromReader(r io.Reader) ([]equation, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	equations := make([]equation, len(lines))
	for i, l := range lines {
		e, err := equationFromString(l)
		if err != nil {
			return nil, fmt.Errorf("could not parse equation %q: %w", l, err)
		}
		equations[i] = e
	}

	return equations, nil
}
