package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 5 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := linesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Count how many times each point appears in a line.
	pointCount := make(map[point]int)
	for _, l := range lines {
		if !l.isHorizontal() && !l.isVertical() { // Ignore diagonal lines.
			continue
		}
		for _, p := range l.points() {
			pointCount[p]++
		}
	}

	// Count how many points appear in more that one line.
	overlaps := 0
	for _, count := range pointCount {
		if count >= 2 {
			overlaps++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", overlaps)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := linesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Count how many times each point appears in a line.
	pointCount := make(map[point]int)
	for _, l := range lines {
		for _, p := range l.points() {
			pointCount[p]++
		}
	}

	// Count how many points appear in more that one line.
	overlaps := 0
	for _, count := range pointCount {
		if count >= 2 {
			overlaps++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", overlaps)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type line struct {
	start, end point
}

type point struct {
	x, y int
}

type vector struct {
	x, y int
}

func (p point) minus(other point) vector {
	return vector{
		x: p.x - other.x,
		y: p.y - other.y,
	}
}

func (p point) plus(v vector) point {
	return point{
		x: p.x + v.x,
		y: p.y + v.y,
	}
}

// unit returns a vector parallel to v, pointing in the same direction, with
// internal values of -1, 0, or 1.
func (v vector) unit() vector {
	unit := vector{0, 0}

	switch {
	case v.x < 0:
		unit.x = -1
	case v.x > 0:
		unit.x = 1
	}

	switch {
	case v.y < 0:
		unit.y = -1
	case v.y > 0:
		unit.y = 1
	}

	return unit
}

func (l line) isVertical() bool {
	return l.start.x == l.end.x
}

func (l line) isHorizontal() bool {
	return l.start.y == l.end.y
}

func (l line) points() []point {
	delta := l.end.minus(l.start).unit()

	var points []point
	for p := l.start; p != l.end; p = p.plus(delta) {
		points = append(points, p)
	}
	points = append(points, l.end)

	return points
}

func linesFromReader(r io.Reader) ([]line, error) {
	rawLines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read: %w", err)
	}

	lines := make([]line, len(rawLines))

	for i := range rawLines {
		l, err := lineFromString(rawLines[i])
		if err != nil {
			return nil, fmt.Errorf("invalid line %q: %w", rawLines[i], err)
		}
		lines[i] = l
	}

	return lines, nil
}

func lineFromString(s string) (line, error) {
	rawPoints := strings.Split(s, " -> ")
	if len(rawPoints) != 2 {
		return line{}, errors.New("wrong format")
	}

	start, err := pointFromString(rawPoints[0])
	if err != nil {
		return line{}, fmt.Errorf("%q invalid point: %w", rawPoints[0], err)
	}

	end, err := pointFromString(rawPoints[1])
	if err != nil {
		return line{}, fmt.Errorf("%q invalid point: %w", rawPoints[0], err)
	}

	return line{start, end}, nil
}

func pointFromString(s string) (point, error) {
	rawNumbers := strings.Split(s, ",")
	if len(rawNumbers) != 2 {
		return point{}, errors.New("wrong format")
	}

	x, err := strconv.Atoi(rawNumbers[0])
	if err != nil {
		return point{}, fmt.Errorf("%q is not a whole number", rawNumbers[0])
	}

	y, err := strconv.Atoi(rawNumbers[1])
	if err != nil {
		return point{}, fmt.Errorf("%q is not a whole number", rawNumbers[0])
	}

	return point{x, y}, nil
}
