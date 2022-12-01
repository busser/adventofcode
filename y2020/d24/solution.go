package busser

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 24 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	tileIdentifiers, err := tileIdentifiersFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	blackTiles := initialBlackTiles(tileIdentifiers)

	_, err = fmt.Fprintf(answer, "%d", len(blackTiles))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 24 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	tileIdentifiers, err := tileIdentifiersFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	blackTiles := initialBlackTiles(tileIdentifiers)

	var sim simulation

	sim.init(blackTiles, 100)

	for i := 0; i < 100; i++ {
		sim.iterate()
	}

	_, err = fmt.Fprintf(answer, "%d", sim.tally())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type direction uint8

const (
	east direction = iota
	southeast
	southwest
	west
	northwest
	northeast
)

var allDirections = [...]direction{east, southeast, southwest, west, northwest, northeast}

type coordinate struct {
	x, y int
}

func (c coordinate) neighbor(dir direction) coordinate {
	newC := c
	switch dir {
	case east:
		newC.x++
	case southeast:
		newC.y--
	case southwest:
		newC.x--
		newC.y--
	case west:
		newC.x--
	case northwest:
		newC.y++
	case northeast:
		newC.x++
		newC.y++
	default:
		panic("unknown direction")
	}

	return newC
}

type simulation struct {
	current, next [][]bool
}

func (sim *simulation) init(initialState map[coordinate]struct{}, margin int) {
	if margin < 0 {
		panic("margin cannot be negative")
	}

	minX, maxX, minY, maxY := bounds(initialState)

	newState := make([][]bool, maxX-minX+2*margin+1)
	for x := range newState {
		newState[x] = make([]bool, maxY-minY+2*margin+1)
	}

	for c := range initialState {
		newState[c.x+margin][c.y+margin] = true
	}

	sim.current = newState
	sim.next = copyState(newState)
}

func (sim *simulation) iterate() {
	for x := range sim.current {
		for y := range sim.current[x] {
			blackNeighbors := sim.blackTilesAround(x, y)
			switch {
			case sim.current[x][y] && (blackNeighbors == 0 || blackNeighbors > 2):
				sim.next[x][y] = false
			case !sim.current[x][y] && blackNeighbors == 2:
				sim.next[x][y] = true
			default:
				sim.next[x][y] = sim.current[x][y]
			}
		}
	}

	sim.current, sim.next = sim.next, sim.current
}

func (sim simulation) blackTilesAround(x, y int) int {
	count := 0

	tile := coordinate{x, y}

	for _, dir := range allDirections {
		neighbor := tile.neighbor(dir)

		if neighbor.x < 0 || neighbor.x >= len(sim.current) {
			continue
		}
		if neighbor.y < 0 || neighbor.y >= len(sim.current[neighbor.x]) {
			continue
		}

		if sim.current[neighbor.x][neighbor.y] {
			count++
		}
	}

	return count
}

func (sim simulation) tally() int {
	count := 0

	for x := range sim.current {
		for y := range sim.current[x] {
			if sim.current[x][y] {
				count++
			}
		}
	}
	return count
}

func bounds(state map[coordinate]struct{}) (minX, maxX, minY, maxY int) {
	for c := range state {
		minX = min(minX, c.x)
		maxX = max(maxX, c.x)
		minY = min(minY, c.y)
		maxY = max(maxY, c.y)
	}

	return minX, maxX, minY, maxY
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

func copyState(state [][]bool) [][]bool {
	new := make([][]bool, len(state))

	for x := range state {
		new[x] = make([]bool, len(state[x]))
		copy(new[x], state[x])
	}

	return new
}

func coordinateFromIdentifier(identifier []direction) coordinate {
	var c coordinate
	for _, dir := range identifier {
		c = c.neighbor(dir)
	}
	return c
}

func initialBlackTiles(tileIdentifiers [][]direction) map[coordinate]struct{} {
	blackTiles := make(map[coordinate]struct{})
	for _, identifier := range tileIdentifiers {
		c := coordinateFromIdentifier(identifier)
		if _, flipped := blackTiles[c]; flipped {
			delete(blackTiles, c)
		} else {
			blackTiles[c] = struct{}{}
		}
	}
	return blackTiles
}

func tileIdentifiersFromReader(r io.Reader) ([][]direction, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading lines: %w", err)
	}

	identifiers := make([][]direction, len(lines))
	for i := range lines {
		identifier, err := tileIdentifierFromString(lines[i])
		if err != nil {
			return nil, fmt.Errorf("parsing line %d: %w", i, err)
		}
		identifiers[i] = identifier
	}

	return identifiers, nil
}

func tileIdentifierFromString(s string) ([]direction, error) {
	var identifier []direction

	formatErr := errors.New("wrong format")

	var previousChar rune
	for _, c := range s {
		switch c {
		case 'e':
			switch previousChar {
			case 0:
				identifier = append(identifier, east)
			case 's':
				previousChar = 0
				identifier = append(identifier, southeast)
			case 'n':
				previousChar = 0
				identifier = append(identifier, northeast)
			default:
				panic("unknown character got through")
			}
		case 'w':
			switch previousChar {
			case 0:
				identifier = append(identifier, west)
			case 's':
				previousChar = 0
				identifier = append(identifier, southwest)
			case 'n':
				previousChar = 0
				identifier = append(identifier, northwest)
			default:
				panic("unknown character got through")
			}
		case 's':
			fallthrough
		case 'n':
			if previousChar != 0 {
				return nil, formatErr
			}
			previousChar = c
		default:
			return nil, formatErr
		}
	}

	return identifier, nil
}
