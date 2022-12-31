package d24

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 24 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	valley, blizzards, err := valleyAndBlizzardsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	entrance := vector{x: 1, y: 0}
	exit := vector{x: len(valley) - 2, y: len(valley[0]) - 1}

	time, err := minimumTime(valley, blizzards, entrance, exit, 1)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", time)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 24 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	valley, blizzards, err := valleyAndBlizzardsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	entrance := vector{x: 1, y: 0}
	exit := vector{x: len(valley) - 2, y: len(valley[0]) - 1}

	time, err := minimumTime(valley, blizzards, entrance, exit, 3)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", time)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type tile uint8

const (
	tileWall tile = iota
	tileGround
)

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
	up    = vector{0, -1}
	down  = vector{0, 1}
	left  = vector{-1, 0}
	right = vector{1, 0}
)

type blizzard struct {
	position  vector
	direction vector
}

func (b *blizzard) update(valley [][]tile) {
	b.position = b.position.plus(b.direction)

	if b.position.x == 0 {
		b.position.x = len(valley) - 2
	}
	if b.position.x == len(valley)-1 {
		b.position.x = 1
	}
	if b.position.y == 0 {
		b.position.y = len(valley[0]) - 2
	}
	if b.position.y == len(valley[0])-1 {
		b.position.y = 1
	}
}

type state struct {
	position       vector
	target         vector
	timePassed     int
	tripsCompleted int
}

type cacheKey struct {
	position       vector
	cycleIndex     int
	tripsCompleted int
}

func (s state) cacheKey(cycleLength int) cacheKey {
	return cacheKey{
		position:       s.position,
		cycleIndex:     s.timePassed % cycleLength,
		tripsCompleted: s.tripsCompleted,
	}
}

func (s state) moved(direction vector) state {
	s.position = s.position.plus(direction)
	return s
}

func minimumTime(valley [][]tile, blizzards []blizzard, entrance, exit vector, trips int) (int, error) {
	// Using two lists to reduce number of reallocations.
	var next []state
	toVisit := []state{{entrance, exit, 0, 0}}

	// Blizzard positions repeat on a cycle. We can use this to identify
	// identical situations we don't need to revisit.
	cycleLength := lcm(len(valley)-2, len(valley[0])-2)

	visited := make(map[cacheKey]bool)

	blockedByBlizzard := make([][]bool, len(valley))
	for x := range valley {
		blockedByBlizzard[x] = make([]bool, len(valley[x]))
	}

	visit := func(s state) {
		if !isWithinBounds(valley, s.position) {
			return
		}

		if valley[s.position.x][s.position.y] == tileWall {
			return
		}

		if blockedByBlizzard[s.position.x][s.position.y] {
			return
		}

		key := s.cacheKey(cycleLength)
		if visited[key] {
			return
		}
		visited[key] = true

		next = append(next, s)
	}

	for time := 0; ; time++ {
		// Compute where blizzards will be next.
		for i := range blizzards {
			blizzards[i].update(valley)
		}
		for x := range blockedByBlizzard {
			for y := range blockedByBlizzard[x] {
				blockedByBlizzard[x][y] = false
			}
		}
		for _, b := range blizzards {
			blockedByBlizzard[b.position.x][b.position.y] = true
		}

		// Add current position and neighbors to possible next positions, if no
		// blizzard prevents the move.
		for _, s := range toVisit {
			if s.position == s.target {
				s.tripsCompleted++
				if s.tripsCompleted == trips {
					return s.timePassed, nil
				}

				switch s.target {
				case entrance:
					s.target = exit
				case exit:
					s.target = entrance
				}
			}

			s.timePassed++

			visit(s)
			visit(s.moved(up))
			visit(s.moved(down))
			visit(s.moved(left))
			visit(s.moved(right))
		}

		if len(next) == 0 {
			return 0, errors.New("no path")
		}

		// Swap lists of positions before next iteration.
		toVisit, next = next, toVisit
		next = next[:0]
	}
}

func isWithinBounds(valley [][]tile, position vector) bool {
	return position.x >= 0 && position.x < len(valley) &&
		position.y >= 0 && position.y < len(valley[position.x])
}

func lcm(a, b int) int {
	if a == 0 && b == 0 {
		return 0
	}

	return abs(a*b) / gcd(a, b)
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func valleyAndBlizzardsFromReader(r io.Reader) ([][]tile, []blizzard, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, nil, err
	}

	if len(lines) == 0 {
		return nil, nil, errors.New("empty input")
	}

	width := len(lines[0])
	for _, l := range lines {
		if len(l) != width {
			return nil, nil, errors.New("valley is not a rectangle")
		}
	}

	valley := make([][]tile, width)
	for x := range valley {
		valley[x] = make([]tile, len(lines))
	}

	var blizzards []blizzard
	for y := range lines {
		for x := range lines[y] {
			switch lines[y][x] {
			case '#':
				valley[x][y] = tileWall
			case '.':
				valley[x][y] = tileGround
			case '^':
				valley[x][y] = tileGround
				blizzards = append(blizzards, blizzard{
					position:  vector{x, y},
					direction: up,
				})
			case 'v':
				valley[x][y] = tileGround
				blizzards = append(blizzards, blizzard{
					position:  vector{x, y},
					direction: down,
				})
			case '<':
				valley[x][y] = tileGround
				blizzards = append(blizzards, blizzard{
					position:  vector{x, y},
					direction: left,
				})
			case '>':
				valley[x][y] = tileGround
				blizzards = append(blizzards, blizzard{
					position:  vector{x, y},
					direction: right,
				})
			}
		}
	}

	return valley, blizzards, nil
}
