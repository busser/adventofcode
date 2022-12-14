package d14

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 14 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	vectors, err := vectorsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	units := amountOfSand(vectors, false)

	_, err = fmt.Fprintf(w, "%d", units)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 14 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	vectors, err := vectorsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	units := amountOfSand(vectors, true)

	_, err = fmt.Fprintf(w, "%d", units)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

//=== Pouring sand =============================================================

var sandStart = position{500, 0}

func amountOfSand(vectors []vector, floor bool) int {
	rock, topLeft := rockFromVectors(vectors)

	sandUnits := 0
	for addSand(rock, topLeft, floor) {
		sandUnits++
	}

	return sandUnits
}

func addSand(rock [][]bool, topLeft position, floor bool) bool {
	sand := sandStart.minus(topLeft)

	if rock[sand.x][sand.y] {
		// The source of sand is blocked.
		return false
	}

	for {
		below := sand.plus(position{0, 1})
		if !floor && below.y >= len(rock[below.x])-2 {
			// The sand has fallen into the abyss.
			return false
		}
		if !rock[below.x][below.y] {
			// The sand slid into the spot below it.
			sand = below
			continue
		}

		belowLeft := sand.plus(position{-1, 1})
		if !rock[belowLeft.x][belowLeft.y] {
			// The sand slid into the spot below it to the left.
			sand = belowLeft
			continue
		}

		belowRight := sand.plus(position{1, 1})
		if !rock[belowRight.x][belowRight.y] {
			// The sand slid into the spot below it to the right.
			sand = belowRight
			continue
		}

		// The sand has stopped moving.
		break
	}

	rock[sand.x][sand.y] = true
	return true
}

//=== Scan of existing rock ====================================================

func rockFromVectors(vectors []vector) ([][]bool, position) {
	topLeft, bottomRight := bounds(vectors)

	dx := bottomRight.minus(topLeft).x + 1
	dy := bottomRight.minus(topLeft).y + 1

	rock := make([][]bool, dx)
	for x := 0; x < dx; x++ {
		rock[x] = make([]bool, dy)
		rock[x][dy-1] = true // floor
	}

	for _, v := range vectors {
		v.foreach(func(p position) {
			p = p.minus(topLeft)
			rock[p.x][p.y] = true
		})
	}

	return rock, topLeft
}

func bounds(vectors []vector) (position, position) {
	var (
		minX, minY = sandStart.x, sandStart.y
		maxX, maxY = sandStart.x, sandStart.y
	)

	for _, v := range vectors {
		minX = min(minX, v.start.x)
		minX = min(minX, v.end.x)
		minY = min(minY, v.start.y)
		minY = min(minY, v.end.y)
		maxX = max(maxX, v.start.x)
		maxX = max(maxX, v.end.x)
		maxY = max(maxY, v.start.y)
		maxY = max(maxY, v.end.y)
	}

	maxY += 2 // Make some room for the floor.

	// Assume maximum sand spread.
	minX = min(minX, sandStart.x-(maxY-sandStart.y))
	maxX = max(maxX, sandStart.x+(maxY-sandStart.y))

	return position{minX, minY}, position{maxX, maxY}
}

//=== Position & vectors =======================================================

type position struct {
	x, y int
}

func (p position) unit() position {
	if p.x != 0 {
		p.x /= abs(p.x)
	}
	if p.y != 0 {
		p.y /= abs(p.y)
	}
	return p
}

func (p position) plus(offset position) position {
	return position{
		p.x + offset.x,
		p.y + offset.y,
	}
}

func (p position) minus(offset position) position {
	return position{
		p.x - offset.x,
		p.y - offset.y,
	}
}

type vector struct {
	start, end position
}

func (v vector) foreach(f func(position)) {
	iter := v.end.minus(v.start).unit()
	for p := v.start; p != v.end; p = p.plus(iter) {
		f(p)
	}
	f(v.end)
}

//=== Parsing ==================================================================

func vectorsFromReader(r io.Reader) ([]vector, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	var vectors []vector

	for _, l := range lines {
		v, err := vectorsFromString(l)
		if err != nil {
			return nil, err
		}

		vectors = append(vectors, v...)
	}

	if len(vectors) == 0 {
		return nil, fmt.Errorf("no vectors")
	}

	return vectors, nil
}

func vectorsFromString(s string) ([]vector, error) {
	rawPositions := strings.Split(s, " -> ")

	var positions []position
	for _, raw := range rawPositions {
		p, err := positionFromString(raw)
		if err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}

	var vectors []vector
	for i := 0; i < len(positions)-1; i++ {
		v := vector{positions[i], positions[i+1]}
		vectors = append(vectors, v)
	}

	return vectors, nil
}

func positionFromString(s string) (position, error) {
	parts := strings.SplitN(s, ",", 2)
	if len(parts) != 2 {
		return position{}, errors.New("wrong format")
	}

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return position{}, fmt.Errorf("%q is not a number", parts[0])
	}

	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return position{}, fmt.Errorf("%q is not a number", parts[1])
	}

	return position{x, y}, nil
}

//=== Math =====================================================================

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
