package d15

import (
	"bytes"
	"fmt"
	"io"
)

// PartOne solves the first problem of day 15 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	warehouse, movements, err := warehouseAndMovementsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	state, err := newState(warehouse)
	if err != nil {
		return fmt.Errorf("could not create state: %w", err)
	}

	state.attemptMoves(movements)

	sum := state.warehouse.sumBoxGPSCoordinates()

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 15 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	warehouse, movements, err := warehouseAndMovementsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	warehouse = widenWarehouse(warehouse)

	state, err := newState(warehouse)
	if err != nil {
		return fmt.Errorf("could not create state: %w", err)
	}

	state.attemptMoves(movements)

	sum := state.warehouse.sumBoxGPSCoordinates()

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type vector struct {
	row, col int
}

func (v vector) add(w vector) vector {
	return vector{
		row: v.row + w.row,
		col: v.col + w.col,
	}
}

var (
	up    = vector{row: -1, col: 0}
	down  = vector{row: 1, col: 0}
	left  = vector{row: 0, col: -1}
	right = vector{row: 0, col: 1}
)

type warehouseMap [][]byte

const (
	wall     = '#'
	empty    = '.'
	robot    = '@'
	box      = 'O'
	boxLeft  = '['
	boxRight = ']'
)

func (w warehouseMap) get(pos vector) byte {
	return w[pos.row][pos.col]
}

func (w warehouseMap) set(pos vector, value byte) {
	w[pos.row][pos.col] = value
}

func (w warehouseMap) canMove(pos, direction vector) bool {
	next := pos.add(direction)

	switch w.get(next) {
	case robot, box:
		return w.canMove(next, direction)
	case boxLeft:
		switch direction {
		case up, down:
			return w.canMove(next, direction) && w.canMove(next.add(right), direction)
		case left:
			return w.canMove(next, direction)
		case right:
			return w.canMove(next.add(right), direction)
		}
	case boxRight:
		// delegate to the left side of the box
		return w.canMove(pos.add(left), direction)
	case wall:
		return false
	case empty:
		return true
	}

	panic("unreachable")
}

func (w warehouseMap) move(pos, direction vector) {
	next := pos.add(direction)

	switch w.get(pos) {
	case robot, box:
		w.move(next, direction)
		w.set(next, w.get(pos))
		w.set(pos, empty)

	case boxLeft:
		switch direction {
		case up, down:
			w.move(next, direction)
			w.move(next.add(right), direction)
			w.set(next, boxLeft)
			w.set(next.add(right), boxRight)
			w.set(pos, empty)
			w.set(pos.add(right), empty)
		case left:
			w.move(next, direction)
			w.set(next, boxLeft)
			w.set(pos, boxRight)
			w.set(pos.add(right), empty)
		case right:
			w.move(next.add(right), direction)
			w.set(next.add(right), boxRight)
			w.set(next, boxLeft)
			w.set(pos, empty)
		}

	case boxRight:
		// delegate to the left side of the box
		w.move(pos.add(left), direction)

	case wall:
		panic("wall should not be moved")

	case empty:
		// do nothing
	}
}

func (w warehouseMap) sumBoxGPSCoordinates() int {
	sum := 0
	for row := range w {
		for col := range w[row] {
			c := w.get(vector{row, col})
			if c == box || c == boxLeft {
				sum += 100*row + col
			}
		}
	}
	return sum
}

type state struct {
	warehouse     warehouseMap
	robotPosition vector
}

func (s *state) attemptMove(direction vector) {
	if s.warehouse.canMove(s.robotPosition, direction) {
		s.warehouse.move(s.robotPosition, direction)
		s.robotPosition = s.robotPosition.add(direction)
	}
}

func (s *state) attemptMoves(movements []vector) {
	for _, m := range movements {
		s.attemptMove(m)
	}
}

func newState(warehouse warehouseMap) (state, error) {
	var robotPosition vector

	found := false
	for row := range warehouse {
		for col, char := range warehouse[row] {
			if char == robot {
				if found {
					return state{}, fmt.Errorf("multiple robots found")
				}
				robotPosition = vector{row, col}
				found = true
			}
		}
	}

	if !found {
		return state{}, fmt.Errorf("no robot found")
	}

	return state{
		warehouse:     warehouse,
		robotPosition: robotPosition,
	}, nil
}

func widenWarehouse(warehouse warehouseMap) warehouseMap {
	widened := make(warehouseMap, len(warehouse))
	for row := range warehouse {
		widened[row] = make([]byte, len(warehouse[row])*2)
		for col := range warehouse[row] {
			c := warehouse[row][col]
			switch c {
			case box:
				widened[row][col*2] = boxLeft
				widened[row][col*2+1] = boxRight
			case robot:
				widened[row][col*2] = robot
				widened[row][col*2+1] = empty
			case wall, empty:
				widened[row][col*2] = c
				widened[row][col*2+1] = c
			default:
				panic("unknown character")
			}
		}
	}
	return widened
}

func warehouseAndMovementsFromReader(r io.Reader) (warehouseMap, []vector, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return warehouseMap{}, nil, fmt.Errorf("could not read input: %w", err)
	}

	chunks := bytes.Split(input, []byte("\n\n"))
	if len(chunks) != 2 {
		return warehouseMap{}, nil, fmt.Errorf("invalid input")
	}

	warehouse, err := warehouseFromBytes(chunks[0])
	if err != nil {
		return warehouseMap{}, nil, fmt.Errorf("invalid warehouse: %w", err)
	}

	movements, err := movementsFromBytes(chunks[1])
	if err != nil {
		return warehouseMap{}, nil, fmt.Errorf("invalid movements: %w", err)
	}

	return warehouse, movements, nil

}

func warehouseFromBytes(b []byte) (warehouseMap, error) {
	lines := bytes.Split(b, []byte("\n"))

	warehouse := make(warehouseMap, len(lines))
	for row, line := range lines {
		warehouse[row] = make([]byte, len(line))
		copy(warehouse[row], line)
	}

	return warehouse, nil
}

func movementsFromBytes(b []byte) ([]vector, error) {
	movements := make([]vector, 0, len(b))

	for _, c := range b {
		switch c {
		case '^':
			movements = append(movements, up)
		case 'v':
			movements = append(movements, down)
		case '<':
			movements = append(movements, left)
		case '>':
			movements = append(movements, right)
		case '\n':
			// discard
		default:
			return nil, fmt.Errorf("unknown movement %q", c)
		}
	}

	return movements, nil
}
