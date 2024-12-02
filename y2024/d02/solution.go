package d02

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 2 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	reports, err := reportsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var safeReports int
	for _, r := range reports {
		if r.isSafe() {
			safeReports++
		}
	}

	_, err = fmt.Fprintf(w, "%d", safeReports)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 2 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	reports, err := reportsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var safeReports int
	for _, r := range reports {
		if r.isSafeWithProblemDampener() {
			safeReports++
		}
	}

	_, err = fmt.Fprintf(w, "%d", safeReports)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type report []int

func (r report) isSafe() bool {
	minDelta, maxDelta := 1, 3
	if r[1]-r[0] < 0 {
		minDelta, maxDelta = -maxDelta, -minDelta
	}

	for i := 1; i < len(r); i++ {
		delta := r[i] - r[i-1]
		if delta < minDelta || delta > maxDelta {
			return false
		}
	}

	return true
}

func (r report) isSafeWithProblemDampener() bool {
	if r.isSafe() {
		return true
	}

	variant := make(report, len(r)-1)

	for i := range r {
		copy(variant, r[:i])
		copy(variant[i:], r[i+1:])
		if variant.isSafe() {
			return true
		}
	}

	return false
}

func reportFromString(s string) (report, error) {
	levels := helpers.IntsFromString(s)
	if len(levels) < 2 {
		return nil, fmt.Errorf("expected at least 2 numbers, got %d", len(levels))
	}
	return levels, nil
}

func reportsFromReader(r io.Reader) ([]report, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	reports := make([]report, len(lines))
	for i, l := range lines {
		r, err := reportFromString(l)
		if err != nil {
			return nil, fmt.Errorf("could not parse report %q: %w", l, err)
		}
		reports[i] = r
	}

	return reports, nil
}
