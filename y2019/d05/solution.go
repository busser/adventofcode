package d05

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/y2019/intcode"
)

// PartOne solves the first problem of day 5 of Advent of Code 2019.
func PartOne(r io.Reader, w io.Writer) error {
	program, err := intcode.ProgramFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	outputs, err := intcode.Run(program, []int{1})
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", outputs[len(outputs)-1])
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2019.
func PartTwo(r io.Reader, w io.Writer) error {
	program, err := intcode.ProgramFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	outputs, err := intcode.Run(program, []int{5})
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", outputs[len(outputs)-1])
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}
