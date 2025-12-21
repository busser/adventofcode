package d01

import (
	"fmt"
	"io"
	"strconv"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 1 of Advent of Code 2025.
func PartOne(r io.Reader, w io.Writer) error {
	rotations, err := readRotations(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	password := countZeroLandings(rotations)

	_, err = fmt.Fprintf(w, "%d", password)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2025.
func PartTwo(r io.Reader, w io.Writer) error {
	rotations, err := readRotations(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	password := countZeroCrossings(rotations)

	_, err = fmt.Fprintf(w, "%d", password)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type rotation struct {
	direction rotationDirection
	amount    int
}

type rotationDirection rune

const (
	rotateRight rotationDirection = 'R'
	rotateLeft  rotationDirection = 'L'
)

func (r rotation) String() string {
	return fmt.Sprintf("%c%d", r.direction, r.amount)
}

func countZeroLandings(rotations []rotation) int {
	dial := 50

	count := 0
	for _, r := range rotations {
		switch r.direction {
		case rotateRight:
			dial = (dial + r.amount) % 100
		case rotateLeft:
			dial = (100 + (dial-r.amount)%100) % 100
		default:
			panic("unknown rotation direction")
		}

		if dial == 0 {
			count++
		}
	}

	return count
}

func countZeroCrossings(rotations []rotation) int {
	dial := 50

	count := 0
	for _, r := range rotations {
		var delta int
		switch r.direction {
		case rotateRight:
			delta = +1
		case rotateLeft:
			delta = -1
		default:
			panic("unknown rotation direction")
		}

		for i := 0; i < r.amount; i++ {
			dial = (100 + dial + delta) % 100
			if dial == 0 {
				count++
			}
		}
	}

	return count
}

func readRotations(r io.Reader) ([]rotation, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	rotations := make([]rotation, len(lines))
	for i, line := range lines {
		rotations[i], err = rotationFromString(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse rotation %d: %w", i, err)
		}
	}

	return rotations, nil
}

func rotationFromString(s string) (rotation, error) {
	if len(s) <= 1 {
		return rotation{}, fmt.Errorf("invalid rotation: %q", s)
	}

	directionLetter := s[0]
	directionAmountUnparsed := s[1:]

	var direction rotationDirection
	switch directionLetter {
	case 'R':
		direction = rotateRight
	case 'L':
		direction = rotateLeft
	default:
		return rotation{}, fmt.Errorf("invalid rotation: %q", s)
	}

	amount, err := strconv.Atoi(directionAmountUnparsed)
	if err != nil {
		return rotation{}, fmt.Errorf("invalid rotation: %q", s)
	}

	return rotation{direction, amount}, nil
}
