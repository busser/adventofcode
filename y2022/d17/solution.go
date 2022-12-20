package d17

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// PartOne solves the first problem of day 17 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	pattern, err := jetPatternFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	height := towerHeight(pattern, 2022)

	_, err = fmt.Fprintf(w, "%d", height)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 17 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	pattern, err := jetPatternFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	height := towerHeight(pattern, 1_000_000_000_000)

	_, err = fmt.Fprintf(w, "%d", height)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// [0][1] is bottom row, second column
type rockType [][]byte

var rockTypes = []rockType{
	{
		[]byte("####"),
	},
	{
		[]byte(".#."),
		[]byte("###"),
		[]byte(".#."),
	},
	{
		[]byte("..#"),
		[]byte("..#"),
		[]byte("###"),
	},
	{
		[]byte("#"),
		[]byte("#"),
		[]byte("#"),
		[]byte("#"),
	},
	{
		[]byte("##"),
		[]byte("##"),
	},
}

const chamberWidth = 7

type simulation struct {
	chamber    [][chamberWidth]byte
	cache      map[cacheKey]cacheValue
	step       int
	rockID     int
	jetPattern []byte
	jetIter    int
}

type cacheKey struct {
	rockID  int
	jetIter int
}

type cacheValue struct {
	step   int
	height int
}

func (sim *simulation) String() string {
	var sb strings.Builder
	for y := len(sim.chamber) - 1; y >= 0; y-- {
		sb.WriteByte('|')
		for _, b := range sim.chamber[y] {
			if b == '#' {
				sb.WriteByte('#')
			} else {
				sb.WriteByte('.')
			}
		}
		sb.WriteByte('|')
		sb.WriteByte('\n')
	}
	sb.WriteString("+-------+")
	return sb.String()
}

func towerHeight(jetPattern []byte, numRocks int) int {
	sim := simulation{
		cache:      make(map[cacheKey]cacheValue),
		jetPattern: jetPattern,
	}

	for sim.step = 0; sim.step < numRocks; sim.step++ {
		// Check cache for a pattern.
		rockID := sim.step % len(rockTypes)
		key := cacheKey{rockID, sim.jetIter}
		if cached, hit := sim.cache[key]; hit {
			var (
				multiple  = (numRocks - sim.step) / (sim.step - cached.step)
				remainder = (numRocks - sim.step) % (sim.step - cached.step)
			)

			if remainder == 0 {
				// Repeating the pattern "multiple" times lands us at the target
				// number of rocks dropped.
				return len(sim.chamber) + (len(sim.chamber)-cached.height)*multiple
			}
		}

		// Update the cache.
		sim.cache[key] = cacheValue{sim.step, len(sim.chamber)}

		// Drop a rock.
		sim.nextRock()
	}

	return len(sim.chamber)
}

func (sim *simulation) nextRock() {

	rockID := sim.step % len(rockTypes)
	rock := rockTypes[rockID]

	x := 2
	y := len(sim.chamber) - 1 + len(rock) + 3

	for {
		// Jet of gas pushes rock
		switch sim.jetPattern[sim.jetIter] {
		case '<':
			if sim.rockCanMoveLeft(rock, x, y) {
				x--
			}
		case '>':
			if sim.rockCanMoveRight(rock, x, y) {
				x++
			}
		}
		sim.jetIter = (sim.jetIter + 1) % len(sim.jetPattern)

		// Rock falls
		if sim.rockCanMoveDown(rock, x, y) {
			y--
		} else {
			sim.setRock(rock, x, y)
			sim.rockID++
			break
		}
	}

	sim.rockID = (sim.rockID + 1) % len(rockTypes)
}

func (sim *simulation) rockCanMoveLeft(rock rockType, x, y int) bool {
	// A rock cannot move left if it already touches the left wall.
	if x <= 0 {
		return false
	}

	// A rock cannot move left if another rock is blocking the way.
	for ry := 0; ry < len(rock); ry++ {
		if y-ry >= len(sim.chamber) {
			// This part of the rock is too far above the tower to collide.
			continue
		}
		for rx := 0; rx < len(rock[ry]); rx++ {
			if rock[ry][rx] != '#' {
				continue
			}
			if sim.chamber[y-ry][x+rx-1] == '#' {
				return false
			}
		}
	}

	return true
}

func (sim *simulation) rockCanMoveRight(rock rockType, x, y int) bool {
	// A rock cannot move right if it already touches the right wall.
	if x+len(rock[0]) >= chamberWidth {
		return false
	}

	// A rock cannot move right if another rock is blocking the way.
	for ry := 0; ry < len(rock); ry++ {
		if y-ry >= len(sim.chamber) {
			// This part of the rock is too far above the tower to collide.
			continue
		}
		for rx := 0; rx < len(rock[ry]); rx++ {
			if rock[ry][rx] != '#' {
				continue
			}
			if sim.chamber[y-ry][x+rx+1] == '#' {
				return false
			}
		}
	}

	return true
}

func (sim *simulation) rockCanMoveDown(rock rockType, x, y int) bool {
	// A rock cannot move down if it is touching the floor.
	if y-len(rock) < 0 {
		return false
	}

	// A rock cannot move down if another rock is blocking the way.
	for ry := 0; ry < len(rock); ry++ {
		if y-ry-1 >= len(sim.chamber) {
			// This part of the rock is too far above the tower to collide.
			continue
		}
		for rx := 0; rx < len(rock[ry]); rx++ {
			if rock[ry][rx] != '#' {
				continue
			}
			if sim.chamber[y-ry-1][x+rx] == '#' {
				return false
			}
		}
	}

	return true
}

func (sim *simulation) setRock(rock rockType, x, y int) {
	for ry := 0; ry < len(rock); ry++ {
		// Grow the chamber enough to fit the rock
		for y-ry >= len(sim.chamber) {
			sim.chamber = append(sim.chamber, [chamberWidth]byte{})
		}

		for rx := 0; rx < len(rock[ry]); rx++ {
			if rock[ry][rx] != '#' {
				continue
			}
			sim.chamber[y-ry][x+rx] = '#'
		}
	}
}

func jetPatternFromReader(r io.Reader) ([]byte, error) {
	pattern, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(pattern), nil
}
