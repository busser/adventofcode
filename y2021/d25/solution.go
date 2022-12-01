package busser

import (
	"fmt"
	"io"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 25 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	floor, err := seaFloorFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	stepCount := 1
	for moveCucumbers(&floor) {
		stepCount++
	}

	_, err = fmt.Fprintf(answer, "%d", stepCount)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type seaFloorLocation uint8

const (
	emptyLocation seaFloorLocation = iota
	southCucumber
	eastCucumber
)

type seaFloor struct {
	current [][]seaFloorLocation
	next    [][]seaFloorLocation
}

func (floor seaFloor) String() string {
	var sb strings.Builder
	for y := range floor.current {
		for x := range floor.current[y] {
			switch floor.current[y][x] {
			case emptyLocation:
				sb.WriteByte('.')
			case southCucumber:
				sb.WriteByte('v')
			case eastCucumber:
				sb.WriteByte('>')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func moveCucumbers(floor *seaFloor) (moved bool) {
	for y := range floor.current {
		copy(floor.next[y], floor.current[y])
	}

	for y := range floor.current {
		for x := range floor.current[y] {
			if floor.current[y][x] != eastCucumber {
				continue
			}

			nextX, nextY := (x+1)%len(floor.current[y]), y
			if floor.current[nextY][nextX] == emptyLocation {
				floor.next[y][x] = emptyLocation
				floor.next[nextY][nextX] = eastCucumber
				moved = true
			}
		}
	}

	floor.current, floor.next = floor.next, floor.current

	for y := range floor.current {
		copy(floor.next[y], floor.current[y])
	}

	for y := range floor.current {
		for x := range floor.current[y] {
			if floor.current[y][x] != southCucumber {
				continue
			}

			nextX, nextY := x, (y+1)%len(floor.current)
			if floor.current[nextY][nextX] == emptyLocation {
				floor.next[y][x] = emptyLocation
				floor.next[nextY][nextX] = southCucumber
				moved = true
			}
		}
	}

	floor.current, floor.next = floor.next, floor.current

	return moved
}

func seaFloorFromReader(r io.Reader) (seaFloor, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return seaFloor{}, err
	}

	var floor seaFloor
	floor.current = make([][]seaFloorLocation, len(lines))
	floor.next = make([][]seaFloorLocation, len(lines))
	for y := range lines {
		floor.current[y] = make([]seaFloorLocation, len(lines[y]))
		floor.next[y] = make([]seaFloorLocation, len(lines[y]))
		for x := range lines[y] {
			switch lines[y][x] {
			case '.':
				floor.current[y][x] = emptyLocation
			case '>':
				floor.current[y][x] = eastCucumber
			case 'v':
				floor.current[y][x] = southCucumber
			default:
				return seaFloor{}, fmt.Errorf("unknown symbol %q", lines[y][x])
			}
		}
	}

	return floor, nil
}
