package d06

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 6 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	races, err := racesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	product := 1
	for _, r := range races {
		product *= r.numberOfWaysToBreakRecord()
	}

	_, err = fmt.Fprintf(w, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 6 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	race, err := singleRaceFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	ways := race.numberOfWaysToBreakRecord()

	_, err = fmt.Fprintf(w, "%d", ways)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type race struct {
	timeAllowed  int
	bestDistance int
}

func (r race) numberOfWaysToBreakRecord() int {
	// Essentially, we are solving the following quadratic inequation:
	//   (timeAllowed - charge) * charge > bestDistance
	// and we want to know how many integer solutions there are for charge.
	// We can rewrite this as:
	//   charge^2 - timeAllowed * charge + bestDistance < 0
	// The inequation has two roots, minCharge and maxCharge:
	//   minCharge = (timeAllowed - sqrt(delta)) / 2
	//   maxCharge = (timeAllowed + sqrt(delta)) / 2
	//   delta = timeAllowed^2 - 4 * bestDistance
	// If delta is negative, we cannot beat the race's current record.
	delta := r.timeAllowed*r.timeAllowed - 4*r.bestDistance
	if delta < 0 {
		return 0
	}

	// Computing sqrt(delta) on integers is a pain, so we are going to compute
	// isqrt(delta) instead, where isqrt is the integer square root function:
	//   root^2 + remainder == delta
	root, remainder := isqrt(delta)

	// We can now rewrite the inequation roots as:
	//   minCharge = (timeAllowed - root) / 2
	//   maxCharge = (timeAllowed + root) / 2
	// We round minCharge up and maxCharge down.
	minCharge := (r.timeAllowed - root + 1) / 2
	maxCharge := (r.timeAllowed + root) / 2

	// If the remainder is 0, both minCharge and maxCharge will yield a race
	// time equal to the time allowed. This is not a valid solution: we need to
	// do better than the current record.
	if remainder == 0 {
		minCharge++
		maxCharge--
	}

	return maxCharge - minCharge + 1
}

func isqrt(n int) (int, int) {
	var (
		root               int
		remainder          int
		threshold          int
		tentativeRemainder int
	)

	threshold = 1
	for threshold <= n {
		threshold *= 4
	}

	remainder, root = n, 0
	for threshold > 1 {
		threshold /= 4
		tentativeRemainder = remainder - root - threshold
		root = root / 2
		if tentativeRemainder >= 0 {
			remainder, root = tentativeRemainder, root+threshold
		}
	}

	return root, remainder
}

func racesFromReader(r io.Reader) ([]race, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}
	if len(lines) != 2 {
		return nil, fmt.Errorf("expected 2 lines, got %d", len(lines))
	}

	times, err := intsFromLine(lines[0])
	if err != nil {
		return nil, fmt.Errorf("could not parse times: %w", err)
	}
	distances, err := intsFromLine(lines[1])
	if err != nil {
		return nil, fmt.Errorf("could not parse distances: %w", err)
	}

	if len(times) != len(distances) {
		return nil, fmt.Errorf("expected %d distances, got %d", len(times), len(distances))
	}

	races := make([]race, len(times))
	for i := range times {
		races[i] = race{
			timeAllowed:  times[i],
			bestDistance: distances[i],
		}
	}

	return races, nil
}

func intsFromLine(s string) ([]int, error) {
	// Ignoring first word; assuming it's a label.

	rawInts := strings.Fields(s)
	if len(rawInts) < 2 {
		return nil, fmt.Errorf("no values given")
	}

	ints := make([]int, len(rawInts)-1)
	for i, rawInt := range rawInts[1:] {
		var err error
		ints[i], err = strconv.Atoi(rawInt)
		if err != nil {
			return nil, fmt.Errorf("could not parse int %q: %w", rawInt, err)
		}
	}

	return ints, nil
}

func singleRaceFromReader(r io.Reader) (race, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return race{}, fmt.Errorf("could not read input: %w", err)
	}
	if len(lines) != 2 {
		return race{}, fmt.Errorf("expected 2 lines, got %d", len(lines))
	}

	time, err := mergedIntFromLine(lines[0])
	if err != nil {
		return race{}, fmt.Errorf("could not parse time: %w", err)
	}
	distance, err := mergedIntFromLine(lines[1])
	if err != nil {
		return race{}, fmt.Errorf("could not parse distance: %w", err)
	}

	return race{
		timeAllowed:  time,
		bestDistance: distance,
	}, nil
}

func mergedIntFromLine(s string) (int, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("expected 2 parts, got %d", len(parts))
	}

	rawInt := strings.ReplaceAll(parts[1], " ", "")
	return strconv.Atoi(rawInt)
}
