package d03

import (
	"fmt"
	"io"
	"regexp"
)

// PartOne solves the first problem of day 3 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	corruptedMemory, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum, err := computeMuls(corruptedMemory)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 3 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	corruptedMemory, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum, err := scanCorruptedMemoryWithConditionals(corruptedMemory)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

var (
	mulInstructionPattern = regexp.MustCompile(`mul\(([0-9]+),([0-9]+)\)`)
	anyInstructionPattern = regexp.MustCompile(`mul\(([0-9]+),([0-9]+)\)|do\(\)|don't\(\)`)
)

func computeMuls(memory []byte) (int, error) {
	matches := mulInstructionPattern.FindAllSubmatch(memory, -1)

	sum := 0
	for _, match := range matches {
		if len(match) != 3 {
			return 0, fmt.Errorf("invalid match: %v", match)
		}

		a := bytesToInt(match[1])
		b := bytesToInt(match[2])

		sum += a * b
	}

	return sum, nil
}

func scanCorruptedMemoryWithConditionals(memory []byte) (int, error) {
	matches := anyInstructionPattern.FindAllSubmatch(memory, -1)

	sum := 0
	doing := true

	for _, match := range matches {
		switch string(match[0]) {
		case "do()":
			doing = true
		case "don't()":
			doing = false
		default:
			if !doing {
				continue
			}
		}

		if len(match) != 3 {
			return 0, fmt.Errorf("invalid match: %v", match)
		}

		a := bytesToInt(match[1])
		b := bytesToInt(match[2])

		sum += a * b
	}

	return sum, nil
}

func bytesToInt(digits []byte) int {
	n := 0
	for _, d := range digits {
		n *= 10
		n += int(d - '0')
	}
	return n
}
