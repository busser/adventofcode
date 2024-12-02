package d02

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 2 of Advent of Code 2015.
func PartOne(r io.Reader, w io.Writer) error {
	presents, err := presentsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	total := 0
	for _, p := range presents {
		total += p.necessaryPaper()
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 2 of Advent of Code 2015.
func PartTwo(r io.Reader, w io.Writer) error {
	presents, err := presentsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	total := 0
	for _, p := range presents {
		total += p.necessaryRibbon()
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type present struct {
	l, w, h int
}

func (p present) surfaceArea() int {
	return 2*p.l*p.w + 2*p.w*p.h + 2*p.h*p.l
}

func (p present) smallestSide() int {
	return min(p.l*p.w, p.w*p.h, p.h*p.l)
}

func (p present) necessaryPaper() int {
	return p.surfaceArea() + p.smallestSide()
}

func (p present) smallestFace() int {
	return 2 * min(p.l+p.w, p.w+p.h, p.h+p.l)
}

func (p present) necessaryRibbon() int {
	return p.smallestFace() + p.l*p.w*p.h
}

func presentsFromReader(r io.Reader) ([]present, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	presents := make([]present, len(lines))
	for i := range lines {
		dimensions := helpers.IntsFromString(lines[i])
		if len(dimensions) != 3 {
			return nil, fmt.Errorf("invalid present %q", lines[i])
		}
		presents[i] = present{dimensions[0], dimensions[1], dimensions[2]}
	}

	return presents, nil
}
