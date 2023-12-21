package d18

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 18 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	digPlan, err := digPlanFromReader(r, false)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	size := lagoonSize(digPlan)

	_, err = fmt.Fprintf(w, "%d", size)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 18 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	digPlan, err := digPlanFromReader(r, true)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	size := lagoonSize(digPlan)

	_, err = fmt.Fprintf(w, "%d", size)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type trench struct {
	direction vector
	length    int
}

type vector struct {
	x, y int
}

func (v vector) plus(o vector) vector {
	return vector{
		x: v.x + o.x,
		y: v.y + o.y,
	}
}

func (v vector) minus(o vector) vector {
	return vector{
		x: v.x - o.x,
		y: v.y - o.y,
	}
}

func (v vector) times(s int) vector {
	return vector{
		x: v.x * s,
		y: v.y * s,
	}
}

var (
	up    = vector{0, -1}
	down  = vector{0, 1}
	left  = vector{-1, 0}
	right = vector{1, 0}
)

func lagoonSize(digPlan []trench) int {
	position := vector{0, 0}

	var (
		area            = 0
		totalVertical   = 0
		totalHorizontal = 0
	)

	for _, trench := range digPlan {
		switch trench.direction {
		case left:
			totalHorizontal += trench.length
		case right:
		case down:
			totalVertical += trench.length
			area += trench.length * position.x
		case up:
			area -= trench.length * position.x
		}
		position = position.plus(trench.direction.times(trench.length))
	}

	area += totalVertical + totalHorizontal + 1

	return area
}

func lagoonSizeOld(digPlan []trench) (int, error) {

	// Step 1: Determine the width and height of the lagoon.

	position := vector{0, 0}

	minX, maxX := position.x, position.x
	minY, maxY := position.y, position.y

	for _, trench := range digPlan {
		position = position.plus(trench.direction.times(trench.length))
		minX = min(minX, position.x)
		maxX = max(maxX, position.x)
		minY = min(minY, position.y)
		maxY = max(maxY, position.y)
	}

	// Step 2: Create a 2D map of the lagoon.
	// We add a 1-tile border around the lagoon to make flood filling easier.

	offset := vector{minX, minY}

	lagoon := make([][]bool, maxX-minX+3)
	for x := range lagoon {
		lagoon[x] = make([]bool, maxY-minY+3)
	}

	startPostion := vector{1, 1}.minus(offset)
	position = startPostion

	for _, trench := range digPlan {
		for i := 0; i < trench.length; i++ {
			position = position.plus(trench.direction)
			lagoon[position.x][position.y] = true
		}
	}

	if position != startPostion {
		return 0, fmt.Errorf("dig plan is not a loop")
	}

	// Step 3: Flood fill outside the lagoon.

	visited := make([][]bool, len(lagoon))
	for x := range visited {
		visited[x] = make([]bool, len(lagoon[x]))
	}

	outsideArea := 0
	var next []vector

	addNext := func(position vector) {
		if position.x < 0 || position.x >= len(lagoon) || position.y < 0 || position.y >= len(lagoon[0]) {
			return
		}

		if visited[position.x][position.y] {
			return
		}
		visited[position.x][position.y] = true

		next = append(next, position)
	}

	addNext(vector{0, 0}.minus(offset))

	for len(next) > 0 {
		position := next[len(next)-1]
		next = next[:len(next)-1]

		if !lagoon[position.x][position.y] {
			outsideArea++

			addNext(position.plus(up))
			addNext(position.plus(down))
			addNext(position.plus(left))
			addNext(position.plus(right))
		}
	}

	// Step 4: Calculate the area of the lagoon.

	lagoonArea := len(lagoon)*len(lagoon[0]) - outsideArea

	return lagoonArea, nil
}

func digPlanFromReader(r io.Reader, correctInstructions bool) ([]trench, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	trenchFromString := incorrectTrenchFromString
	if correctInstructions {
		trenchFromString = correctTrenchFromString
	}

	trenches := make([]trench, len(lines))
	for i, line := range lines {
		trenches[i], err = trenchFromString(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse line %d: %w", i+1, err)
		}
	}

	return trenches, nil
}

func incorrectTrenchFromString(s string) (trench, error) {
	parts := strings.SplitN(s, " ", 3)
	if len(parts) != 3 {
		return trench{}, fmt.Errorf("invalid trench: %q", s)
	}

	var direction vector
	switch parts[0] {
	case "U":
		direction = up
	case "D":
		direction = down
	case "L":
		direction = left
	case "R":
		direction = right
	default:
		return trench{}, fmt.Errorf("invalid direction: %q", parts[0])
	}

	length, err := strconv.Atoi(parts[1])
	if err != nil {
		return trench{}, fmt.Errorf("invalid length: %q", parts[1])
	}

	return trench{
		direction: direction,
		length:    length,
	}, nil
}

func correctTrenchFromString(s string) (trench, error) {
	parts := strings.SplitN(s, " ", 3)
	if len(parts) != 3 {
		return trench{}, fmt.Errorf("invalid trench: %q", s)
	}

	if len(parts[2]) != 9 {
		return trench{}, fmt.Errorf("invalid color: %q", parts[2])
	}

	length, err := strconv.ParseInt(parts[2][2:7], 16, 64)
	if err != nil {
		return trench{}, fmt.Errorf("invalid length: %q", parts[1])
	}

	var direction vector
	switch parts[2][7] {
	case '0':
		direction = right
	case '1':
		direction = down
	case '2':
		direction = left
	case '3':
		direction = up
	}

	return trench{
		direction: direction,
		length:    int(length),
	}, nil
}
