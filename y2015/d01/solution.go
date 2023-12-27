package d01

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 1 of Advent of Code 2015.
func PartOne(r io.Reader, w io.Writer) error {
	instructions, err := instructionsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	floor := finalFloor(instructions)

	_, err = fmt.Fprintf(w, "%d", floor)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2015.
func PartTwo(r io.Reader, w io.Writer) error {
	instructions, err := instructionsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	index := firstBasementInstruction(instructions)

	_, err = fmt.Fprintf(w, "%d", index)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func finalFloor(instructions string) int {
	floor := 0
	for _, instruction := range instructions {
		switch instruction {
		case '(':
			floor++
		case ')':
			floor--
		}
	}
	return floor
}

func firstBasementInstruction(instructions string) int {
	floor := 0
	for i, instruction := range instructions {
		switch instruction {
		case '(':
			floor++
		case ')':
			floor--
		}

		if floor == -1 {
			return i + 1
		}
	}
	return -1
}

func instructionsFromReader(r io.Reader) (string, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return "", fmt.Errorf("could not read input: %w", err)
	}

	if len(lines) != 1 {
		return "", fmt.Errorf("expected 1 line, got %d", len(lines))
	}

	instructions := lines[0]
	for _, instruction := range instructions {
		if instruction != '(' && instruction != ')' {
			return "", fmt.Errorf("invalid instruction: %c", instruction)
		}
	}

	return lines[0], nil
}
