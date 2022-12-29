package d23

import (
	"bytes"
	"fmt"
	"io"
)

// PartOne solves the first problem of day 23 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	positions, err := elfPositionsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_ = doProcess(positions, 10)
	area := totalArea(positions)

	_, err = fmt.Fprintf(w, "%d", area-len(positions))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 23 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	positions, err := elfPositionsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	rounds := doProcess(positions, infinity)

	_, err = fmt.Fprintf(w, "%d", rounds+1)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const infinity = 1_000_000

type vector struct {
	x, y int
}

func (v vector) plus(w vector) vector {
	return vector{
		x: v.x + w.x,
		y: v.y + w.y,
	}
}

var (
	north     = vector{0, -1}
	west      = vector{-1, 0}
	south     = vector{0, 1}
	east      = vector{1, 0}
	northWest = north.plus(west)
	southWest = south.plus(west)
	southEast = south.plus(east)
	northEast = north.plus(east)

	allDirections = []vector{
		north, west, south, east,
		northWest, northEast,
		southWest, southEast,
	}
)

type decision struct {
	needEmpty []vector
	delta     vector
}

var allDecisions = []decision{
	{[]vector{north, northEast, northWest}, north},
	{[]vector{south, southEast, southWest}, south},
	{[]vector{west, northWest, southWest}, west},
	{[]vector{east, northEast, southEast}, east},
}

func doProcess(positions []vector, maxRounds int) int {
	nextPositions := make([]vector, len(positions))

	for round := 0; round < maxRounds; round++ {
		// Build an index of occupied positions.
		positionOccupied := make(map[vector]bool, len(positions))
		for _, p := range positions {
			positionOccupied[p] = true
		}

		// Determine next position for each elf.
		for elf := range positions {
			allNeighborsEmpty := all(allDirections, func(dir vector) bool {
				return !positionOccupied[positions[elf].plus(dir)]
			})
			if allNeighborsEmpty {
				nextPositions[elf] = positions[elf]
				continue
			}

			decisionMade := false
			for i := range allDecisions {
				d := allDecisions[(i+round)%len(allDecisions)]
				allEmpty := all(d.needEmpty, func(dir vector) bool {
					return !positionOccupied[positions[elf].plus(dir)]
				})
				if allEmpty {
					decisionMade = true
					nextPositions[elf] = positions[elf].plus(d.delta)
					break
				}
			}

			if !decisionMade {
				nextPositions[elf] = positions[elf]
			}
		}

		// Count how many elves want to move to each position.
		nextPositionsCount := make(map[vector]int, len(positions))
		for _, p := range nextPositions {
			nextPositionsCount[p]++
		}

		// Move the elves.
		anElfMoved := false
		for elf := range positions {
			pos := nextPositions[elf]
			if nextPositionsCount[pos] == 1 && positions[elf] != pos {
				positions[elf] = pos
				anElfMoved = true
			}
		}

		// Stop now if no elves moved.
		if !anElfMoved {
			return round
		}
	}

	return maxRounds
}

func totalArea(positions []vector) int {
	minX, maxX := positions[0].x, positions[0].x
	minY, maxY := positions[0].y, positions[0].y

	for _, p := range positions {
		minX = min(minX, p.x)
		maxX = max(maxX, p.x)
		minY = min(minY, p.y)
		maxY = max(maxY, p.y)
	}

	return (maxX - minX + 1) * (maxY - minY + 1)
}

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

func all[T any](values []T, fn func(T) bool) bool {
	for _, v := range values {
		if !fn(v) {
			return false
		}
	}
	return true
}

func elfPositionsFromReader(r io.Reader) ([]vector, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	grove := bytes.Split(bytes.TrimSpace(raw), []byte("\n"))

	var positions []vector
	for y := range grove {
		for x := range grove[y] {
			if grove[y][x] == '#' {
				positions = append(positions, vector{x, y})
			}
		}
	}

	return positions, nil
}
