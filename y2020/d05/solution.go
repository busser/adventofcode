package busser

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 5 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	passes, err := boardingPassesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var highestID int
	for _, p := range passes {
		if p.seatID > highestID {
			highestID = p.seatID
		}
	}

	_, err = fmt.Fprintf(answer, "%d", highestID)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	passes, err := boardingPassesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	lowestID, highestID := passes[0].seatID, passes[0].seatID
	for _, p := range passes[1:] {
		if p.seatID > highestID {
			highestID = p.seatID
		}
		if p.seatID < lowestID {
			lowestID = p.seatID
		}
	}

	// To find missing value in sequence, compare what the sum should be to
	// what the sum actually is. The difference is the missing number.
	missingID := (highestID - lowestID + 1) * (lowestID + highestID) / 2
	for _, p := range passes {
		missingID -= p.seatID
	}

	_, err = fmt.Fprintf(answer, "%d", missingID)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type boardingPass struct {
	row, column, seatID int
}

func (bp *boardingPass) fromString(s string) error {
	if len(s) != 10 {
		return errors.New("wrong length")
	}

	bitValue := 64
	for _, c := range s[:7] {
		switch c {
		case 'F':
		case 'B':
			bp.row += bitValue
		default:
			return fmt.Errorf("unknown symbol: %q", c)
		}

		bitValue /= 2
	}

	bitValue = 4
	for _, c := range s[7:] {
		switch c {
		case 'L':
		case 'R':
			bp.column += bitValue
		default:
			return fmt.Errorf("unknown symbol: %q", c)
		}

		bitValue /= 2
	}

	bp.seatID = 8*bp.row + bp.column

	return nil
}

func boardingPassesFromReader(r io.Reader) ([]boardingPass, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading lines: %w", err)
	}

	passes := make([]boardingPass, len(lines))

	for i, l := range lines {
		if err := passes[i].fromString(l); err != nil {
			return nil, fmt.Errorf("parsing boarding pass: %w", err)
		}
	}

	return passes, nil
}
