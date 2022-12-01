package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 17 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	t, err := targetFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	initialY := -t.minY - 1
	peekY := initialY * (initialY + 1) / 2

	_, err = fmt.Fprintf(answer, "%d", peekY)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 17 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	t, err := targetFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := initialVelocitiesThatReachTarget(t)

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type target struct {
	minX, maxX, minY, maxY int
}

func initialVelocitiesThatReachTarget(t target) int {
	sums := make([]int, -t.minY+1)
	for i := 1; i < len(sums); i++ {
		sums[i] = i + sums[i-1]
	}

	var minPossibleX int
	for minPossibleX = 0; sums[minPossibleX] < t.minX; minPossibleX++ {
	}
	maxPossibleX := t.maxX

	minPossibleY, maxPossibleY := t.minY, -t.minY-1

	probeReachesTarget := func(x, y int) bool {
		posX, posY := 0, 0
		for posX <= t.maxX && posY >= t.minY {
			posX, posY = posX+x, posY+y
			if x > 0 {
				x--
			}
			y--

			if posX >= t.minX && posX <= t.maxX && posY >= t.minY && posY <= t.maxY {
				return true
			}
		}
		return false
	}

	count := 0
	for x := minPossibleX; x <= maxPossibleX; x++ {
		for y := minPossibleY; y <= maxPossibleY; y++ {
			if probeReachesTarget(x, y) {
				count++
			}
		}
	}

	return count
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func targetFromReader(r io.Reader) (target, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return target{}, err
	}

	if len(lines) != 1 {
		return target{}, fmt.Errorf("extected 1 line, got %d", len(lines))
	}

	rawTarget := strings.TrimPrefix(lines[0], "target area: ")

	parts := strings.Split(rawTarget, ", ")
	if len(parts) != 2 {
		return target{}, errors.New("wrong format")
	}

	minX, maxX, err := minMaxFromRange(parts[0])
	if err != nil {
		return target{}, err
	}
	minY, maxY, err := minMaxFromRange(parts[1])
	if err != nil {
		return target{}, err
	}

	return target{minX, maxX, minY, maxY}, nil
}

func minMaxFromRange(r string) (min, max int, err error) {
	parts := strings.SplitN(r, "=", 2)
	if len(parts) != 2 {
		return 0, 0, errors.New("wrong format")
	}

	rawMinMax := strings.SplitN(parts[1], "..", 2)
	if len(rawMinMax) != 2 {
		return 0, 0, errors.New("wrong format")
	}

	min, err = strconv.Atoi(rawMinMax[0])
	if err != nil {
		return 0, 0, fmt.Errorf("%q is not an integer", rawMinMax[0])
	}

	max, err = strconv.Atoi(rawMinMax[1])
	if err != nil {
		return 0, 0, fmt.Errorf("%q is not an integer", rawMinMax[1])
	}

	return min, max, nil
}
