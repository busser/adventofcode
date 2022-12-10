package d10

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 10 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	changes, err := registerChangesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	values := registerValuesOverTime(changes)

	sum := 0
	for cycle := 20; cycle <= 240; cycle += 40 {
		sum += cycle * values[cycle-1]
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 10 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	changes, err := registerChangesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	values := registerValuesOverTime(changes)
	image := crtImage(values)

	for row := range image {
		_, err = fmt.Fprintf(w, "%s\n", image[row])
		if err != nil {
			return fmt.Errorf("could not write answer: %w", err)
		}
	}

	return nil
}

const (
	crtRows    = 6
	crtColumns = 40
	pixelCount = crtRows * crtColumns
)

type registerChange struct {
	cycle int
	delta int
}

func registerValuesOverTime(changes []registerChange) []int {
	values := make([]int, pixelCount)
	values[0] = 1

	for cycle := 1; cycle < len(values); cycle++ {
		values[cycle] = values[cycle-1]

		if len(changes) == 0 {
			continue
		}

		if changes[0].cycle == cycle {
			values[cycle] += changes[0].delta
			changes = changes[1:]
		}
	}

	return values
}

func crtImage(values []int) [][]byte {
	crt := make([][]byte, crtRows)
	for row := range crt {
		crt[row] = make([]byte, crtColumns)
	}

	for row := range crt {
		for col := range crt[row] {
			value := values[row*crtColumns+col]
			if abs(value-col) <= 1 {
				crt[row][col] = '#'
			} else {
				crt[row][col] = '.'
			}
		}
	}

	return crt
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func registerChangesFromReader(r io.Reader) ([]registerChange, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	var changes []registerChange
	cycle := 0

	for _, l := range lines {
		switch {

		case l == "noop":
			cycle++

		case strings.HasPrefix(l, "addx "):
			rawDelta := strings.TrimPrefix(l, "addx ")
			delta, err := strconv.Atoi(rawDelta)
			if err != nil {
				return nil, fmt.Errorf("%q is not a number", rawDelta)
			}
			cycle += 2
			changes = append(changes, registerChange{cycle, delta})

		default:
			return nil, fmt.Errorf("unknown instruction %q", l)

		}
	}

	return changes, nil
}
