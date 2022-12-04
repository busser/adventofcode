package d02

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/y2019/intcode"
)

// PartOne solves the first problem of day 2 of Advent of Code 2019.
func PartOne(r io.Reader, w io.Writer) error {
	program, err := intcode.ProgramFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	result, err := runGravityAssist(program, 12, 2)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", result)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 2 of Advent of Code 2019.
func PartTwo(r io.Reader, w io.Writer) error {
	program, err := intcode.ProgramFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	noun, verb, err := findNounAndVerb(program, 19690720)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", 100*noun+verb)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func runGravityAssist(program []int, noun, verb int) (int, error) {
	if len(program) < 3 {
		return 0, errors.New("program too short")
	}

	program[1] = noun
	program[2] = verb

	if err := intcode.Run(program); err != nil {
		return 0, err
	}

	return program[0], nil
}

func findNounAndVerb(program []int, target int) (int, int, error) {
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			p := copyProgram(program)

			result, err := runGravityAssist(p, noun, verb)
			if err != nil {
				return 0, 0, err
			}

			if result == target {
				return noun, verb, nil
			}
		}
	}

	return 0, 0, fmt.Errorf("not found")
}

func copyProgram(program []int) []int {
	newProgram := make([]int, len(program))
	copy(newProgram, program)
	return newProgram
}
