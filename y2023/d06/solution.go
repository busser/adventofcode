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
	count := 0

	for charge := 0; charge < r.timeAllowed; charge++ {
		speed := charge
		if (r.timeAllowed-charge)*speed > r.bestDistance {
			count++
		}
	}

	return count
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
