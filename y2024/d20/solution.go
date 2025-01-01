package d20

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 20 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	track, err := raceTrackFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	track.findPath()
	track.buildCostMap()
	track.findCheats(2)

	count := track.countBestCheats(100)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 20 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	track, err := raceTrackFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	track.findPath()
	track.buildCostMap()
	track.findCheats(20)

	count := track.countBestCheats(100)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type vector struct {
	row, col int
}

func (v vector) plus(w vector) vector {
	return vector{
		row: v.row + w.row,
		col: v.col + w.col,
	}
}

func (v vector) distanceTo(w vector) int {
	return abs(v.row-w.row) + abs(v.col-w.col)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type cheat struct {
	start, end vector
}

const (
	empty = '.'
	wall  = '#'
	start = 'S'
	end   = 'E'
)

type raceTrack struct {
	tiles [][]byte
	start vector
	end   vector

	path          []vector
	costs         [][]int
	cheatDistance int
	cheats        []cheat
}

func (t *raceTrack) costAt(v vector) int {
	return t.costs[v.row][v.col]
}

func (t raceTrack) at(v vector) byte {
	return t.tiles[v.row][v.col]
}

func (t raceTrack) isWithinBounds(v vector) bool {
	return v.row >= 0 && v.row < len(t.tiles) && v.col >= 0 && v.col < len(t.tiles[v.row])
}

func (t *raceTrack) findCheats(cheatDistance int) {
	t.cheatDistance = cheatDistance

	for _, start := range t.path {
		for _, end := range t.reachablePositions(start, t.cheatDistance) {
			if t.at(end) == wall {
				continue
			}
			c := cheat{start, end}
			if t.cheatWorth(c) > 0 {
				t.cheats = append(t.cheats, cheat{start, end})
			}
		}
	}
}

func (t *raceTrack) cheatWorth(c cheat) int {
	return t.costAt(c.end) - t.costAt(c.start) - c.start.distanceTo(c.end)
}

func (t *raceTrack) countBestCheats(threshold int) int {
	count := 0
	for _, c := range t.cheats {
		if t.cheatWorth(c) >= threshold {
			count++
		}
	}

	return count
}

func (t *raceTrack) reachablePositions(start vector, distance int) []vector {
	var reachable []vector

	for drow := -distance; drow <= distance; drow++ {
		absDrow := abs(drow)
		for dcol := -(distance - absDrow); dcol <= distance-absDrow; dcol++ {
			pos := start.plus(vector{row: drow, col: dcol})
			if t.isWithinBounds(pos) && t.at(pos) != wall {
				reachable = append(reachable, pos)
			}
		}
	}

	return reachable
}

func (t *raceTrack) buildCostMap() {
	t.costs = make([][]int, len(t.tiles))
	for row := range t.costs {
		t.costs[row] = make([]int, len(t.tiles[row]))
		for col := range t.costs[row] {
			t.costs[row][col] = -1
		}
	}

	for timeToReach, pos := range t.path {
		t.costs[pos.row][pos.col] = timeToReach
	}
}

func (t *raceTrack) findPath() {
	toVisit := []vector{t.start}
	var nextToVisit []vector

	visited := make([][]bool, len(t.tiles))
	for row := range visited {
		visited[row] = make([]bool, len(t.tiles[row]))
	}

	for len(toVisit) > 0 {
		for _, pos := range toVisit {
			if pos == t.end {
				t.path = append(t.path, pos)
				return
			}

			if !t.isWithinBounds(pos) || t.at(pos) == wall {
				continue
			}

			if visited[pos.row][pos.col] {
				continue
			}
			visited[pos.row][pos.col] = true

			t.path = append(t.path, pos)

			nextToVisit = append(nextToVisit, pos.plus(vector{row: -1, col: 0}))
			nextToVisit = append(nextToVisit, pos.plus(vector{row: 1, col: 0}))
			nextToVisit = append(nextToVisit, pos.plus(vector{row: 0, col: -1}))
			nextToVisit = append(nextToVisit, pos.plus(vector{row: 0, col: 1}))
		}

		toVisit, nextToVisit = nextToVisit, toVisit[:0]
	}
}

func raceTrackFromReader(r io.Reader) (raceTrack, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return raceTrack{}, fmt.Errorf("could not read input: %w", err)
	}

	var track raceTrack

	track.tiles = make([][]byte, len(lines))
	for row, line := range lines {
		track.tiles[row] = []byte(line)
		for col, char := range line {
			switch char {
			case empty, wall:
				// do nothing
			case start:
				track.start = vector{row, col}
				track.tiles[row][col] = empty
			case end:
				track.end = vector{row, col}
				track.tiles[row][col] = empty
			default:
				return raceTrack{}, fmt.Errorf("invalid character %c at row %d, column %d", char, row, col)
			}
		}
	}

	return track, nil
}
