package busser

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 13 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	start, timetable, err := startAndTimetableFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	earliestBus, earliestDeparture := 0, math.MaxInt64
	for _, bus := range timetable {
		if bus == 0 {
			continue
		}

		departure := nextDeparture(bus, start)
		if departure < earliestDeparture {
			earliestBus, earliestDeparture = bus, departure
		}
	}

	waitTime := earliestDeparture - start

	_, err = fmt.Fprintf(answer, "%d", earliestBus*waitTime)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 13 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	_, timetable, err := startAndTimetableFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	timestamp, jump := 0, 1
	for i, bus := range timetable {
		if bus == 0 {
			continue
		}

		for (timestamp+i)%bus != 0 {
			timestamp += jump
		}
		jump = lcm(jump, bus) // If all bus numbers are prime, could be `jump *= bus`
	}

	_, err = fmt.Fprintf(answer, "%d", timestamp)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func nextDeparture(bus, earliest int) int {
	mod := earliest % bus
	if mod == 0 {
		return earliest
	}
	return earliest - mod + bus
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func lcm(a, b int) int {
	if a == 0 || b == 0 {
		panic("can't compute LCM of zero")
	}
	return a * b / gcd(a, b)
}

func startAndTimetableFromReader(r io.Reader) (int, []int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return 0, nil, fmt.Errorf("reading lines: %w", err)
	}

	if len(lines) != 2 {
		return 0, nil, errors.New("wrong format")
	}

	start, err := strconv.Atoi(lines[0])
	if err != nil {
		return 0, nil, fmt.Errorf("wrong format: not an int: %w", err)
	}

	rawTimetable := strings.Split(lines[1], ",")
	timetable := make([]int, len(rawTimetable))

	for i, s := range rawTimetable {
		if s == "x" {
			continue
		}

		t, err := strconv.Atoi(s)
		if err != nil {
			return 0, nil, fmt.Errorf("wrong format: not an int: %w", err)
		}

		timetable[i] = t
	}

	return start, timetable, nil
}
