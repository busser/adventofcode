package d13

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 13 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	machines, err := machinesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	total := 0
	for _, m := range machines {
		tokens, ok := m.fewestTokensToSpend()
		if ok {
			total += tokens
		}
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 13 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	machines, err := machinesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	total := 0
	for _, m := range machines {
		m.prize = m.prize.plus(vector{row: 1e13, col: 1e13})
		tokens, ok := m.fewestTokensToSpend()
		if ok {
			total += tokens
		}
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type vector struct {
	row, col int
}

func (v vector) plus(w vector) vector {
	return vector{
		row: v.row + w.row,
		col: v.col + w.col,
	}
}

type machine struct {
	buttonA vector
	buttonB vector
	prize   vector
}

func (m machine) solve() (int, int, bool) {
	determinant := m.buttonA.row*m.buttonB.col - m.buttonA.col*m.buttonB.row
	if determinant == 0 {
		return 0, 0, false
	}

	numA := m.prize.row*m.buttonB.col - m.prize.col*m.buttonB.row
	numB := m.prize.col*m.buttonA.row - m.prize.row*m.buttonA.col

	if numA%determinant != 0 || numB%determinant != 0 {
		return 0, 0, false
	}

	a := numA / determinant
	b := numB / determinant

	return a, b, true
}

func (m machine) fewestTokensToSpend() (int, bool) {
	a, b, ok := m.solve()
	if !ok || a < 0 || b < 0 {
		return 0, false
	}
	return a*3 + b, true
}

func vectorFromString(s string) (vector, error) {
	numbers := helpers.IntsFromString(s)
	if len(numbers) != 2 {
		return vector{}, fmt.Errorf("invalid vector: %q", s)
	}

	return vector{
		row: numbers[0],
		col: numbers[1],
	}, nil
}

func machineFromStrings(s1, s2, s3 string) (machine, error) {
	buttonA, err := vectorFromString(s1)
	if err != nil {
		return machine{}, err
	}

	buttonB, err := vectorFromString(s2)
	if err != nil {
		return machine{}, err
	}

	prize, err := vectorFromString(s3)
	if err != nil {
		return machine{}, err
	}

	return machine{
		buttonA: buttonA,
		buttonB: buttonB,
		prize:   prize,
	}, nil
}

func machinesFromReader(r io.Reader) ([]machine, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	if len(lines)%4 != 3 {
		return nil, fmt.Errorf("invalid number of lines")
	}

	var machines []machine
	for i := 0; i < len(lines); i += 4 {
		m, err := machineFromStrings(lines[i], lines[i+1], lines[i+2])
		if err != nil {
			return nil, err
		}

		machines = append(machines, m)
	}

	return machines, nil
}
