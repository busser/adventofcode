package d06

import (
	"bytes"
	"fmt"
	"io"
)

// PartOne solves the first problem of day 6 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	lab, guard, err := labAndGuardFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := countGuardPositions(lab, guard)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 6 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	lab, guard, err := labAndGuardFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := countPossibleLoops(lab, guard)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func countGuardPositions(lab labMap, guard guardInfo) int {
	pt := newPositionTracker(lab)

	pt.visit(guard.position)
	for guard.move(lab) {
		pt.visit(guard.position)
	}

	return pt.count
}

type positionTracker struct {
	visited [][]bool
	count   int
}

func newPositionTracker(lab labMap) *positionTracker {
	visited := make([][]bool, len(lab))
	for row := range visited {
		visited[row] = make([]bool, len(lab[row]))
	}

	return &positionTracker{
		visited: visited,
		count:   0,
	}
}

func (pt *positionTracker) visit(pos vector) {
	if pt.visited[pos.row][pos.col] {
		return
	}
	pt.visited[pos.row][pos.col] = true
	pt.count++
}

func countPossibleLoops(lab labMap, guard guardInfo) int {
	guardStart := guard

	pt := newPositionTracker(lab)
	pt.visit(guard.position)
	for guard.move(lab) {
		pt.visit(guard.position)
	}

	count := 0

	for row := range lab {
		for col := range lab[row] {
			if !pt.visited[row][col] {
				// The only way to create a loop is by adding an obstacle
				// somewhere on the guard's path.
				continue
			}

			pos := vector{row, col}
			if guardStart.position == pos {
				continue
			}

			lab[row][col] = wall
			if checkIfLoop(lab, guardStart) {
				count++
			}
			lab[row][col] = empty
		}
	}

	return count
}

func checkIfLoop(lab labMap, guard guardInfo) bool {
	ld := newLoopDetector()

	_ = ld.isOnLoop(guard)
	for guard.move(lab) {
		if ld.isOnLoop(guard) {
			return true
		}
	}

	return false
}

type loopDetector struct {
	seen map[guardInfo]bool
}

func newLoopDetector() *loopDetector {
	return &loopDetector{
		seen: make(map[guardInfo]bool),
	}
}

func (ld *loopDetector) isOnLoop(g guardInfo) bool {
	if ld.seen[g] {
		return true
	}
	ld.seen[g] = true

	return false
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

func (v vector) rotatedRight() vector {
	return vector{
		row: v.col,
		col: -v.row,
	}
}

var up = vector{-1, 0}

const (
	empty            = '.'
	wall             = '#'
	startingPosition = '^'
)

type labMap [][]byte

type guardInfo struct {
	position  vector
	direction vector
}

func (g *guardInfo) move(lab labMap) bool {
	for {
		next := g.position.plus(g.direction)
		if next.row < 0 || next.row >= len(lab) || next.col < 0 || next.col >= len(lab[next.row]) {
			return false
		}

		if lab[next.row][next.col] == wall {
			g.direction = g.direction.rotatedRight()
			continue
		}

		g.position = next
		return true
	}
}

func labAndGuardFromReader(r io.Reader) (labMap, guardInfo, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, guardInfo{}, err
	}

	data = bytes.TrimSpace(data)

	lab := bytes.Split(data, []byte("\n"))

	var guard guardInfo
	foundGuard := false
	for row := range lab {
		for col := range lab[row] {
			switch lab[row][col] {
			case empty, wall:
				// do nothing
			case startingPosition:
				if foundGuard {
					return nil, guardInfo{}, fmt.Errorf("multiple starting positions found")
				}
				foundGuard = true
				guard.position = vector{row, col}
				guard.direction = up
			default:
				return nil, guardInfo{}, fmt.Errorf("unknown symbol %q", lab[row][col])
			}
		}
	}

	if !foundGuard {
		return nil, guardInfo{}, fmt.Errorf("no guard found")
	}

	return lab, guard, nil
}
