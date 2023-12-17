package d16

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 16 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	layout, err := layoutFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sim := newSimulation(layout)
	sim.add(beam{
		position:  vector{row: 0, col: 0},
		direction: right,
	})
	sim.run()

	energized := sim.countEnergized()

	_, err = fmt.Fprintf(w, "%d", energized)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 16 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	layout, err := layoutFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	maxEnergized := findMaxEnergized(layout)

	_, err = fmt.Fprintf(w, "%d", maxEnergized)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const (
	empty      = '.'
	vSplit     = '|'
	hSplit     = '-'
	trblMirror = '/'
	tlbrMirror = '\\'
)

type vector struct {
	row, col int
}

var (
	up    = vector{row: -1, col: 0}
	down  = vector{row: 1, col: 0}
	left  = vector{row: 0, col: -1}
	right = vector{row: 0, col: 1}
)

func (v vector) plus(o vector) vector {
	return vector{
		row: v.row + o.row,
		col: v.col + o.col,
	}
}

type beam struct {
	position  vector
	direction vector
}

func (b beam) moved(direction vector) beam {
	return beam{
		position:  b.position.plus(direction),
		direction: direction,
	}
}

type simulation struct {
	layout    [][]byte
	beams     []beam
	seen      map[beam]bool
	energized [][]bool
}

func newSimulation(layout [][]byte) *simulation {
	s := &simulation{
		layout:    layout,
		beams:     nil,
		seen:      make(map[beam]bool),
		energized: make([][]bool, len(layout)),
	}
	for i := range s.energized {
		s.energized[i] = make([]bool, len(layout))
	}
	return s
}

func (s *simulation) run() {
	for s.step() {
	}
}

func (s *simulation) add(b beam) {
	// Drop beams that leave the layout.
	if b.position.row < 0 || b.position.row >= len(s.layout) {
		return
	}
	if b.position.col < 0 || b.position.col >= len(s.layout) {
		return
	}

	// Drop beams we've already processed, to avoid infinite loops.
	if s.seen[b] {
		return
	}
	s.seen[b] = true

	s.energized[b.position.row][b.position.col] = true

	// Add beam to backlog.
	s.beams = append(s.beams, b)
}

func (s *simulation) step() bool {
	if len(s.beams) == 0 {
		return false
	}

	b := s.beams[0]
	s.beams = s.beams[1:]

	switch s.layout[b.position.row][b.position.col] {
	case empty:
		s.add(b.moved(b.direction))
	case vSplit:
		if b.direction == up || b.direction == down {
			s.add(b.moved(b.direction))
		} else {
			s.add(b.moved(up))
			s.add(b.moved(down))
		}
	case hSplit:
		if b.direction == left || b.direction == right {
			s.add(b.moved(b.direction))
		} else {
			s.add(b.moved(left))
			s.add(b.moved(right))
		}
	case trblMirror:
		switch b.direction {
		case up:
			s.add(b.moved(right))
		case down:
			s.add(b.moved(left))
		case left:
			s.add(b.moved(down))
		case right:
			s.add(b.moved(up))
		default:
			panic("invalid direction")
		}
	case tlbrMirror:
		switch b.direction {
		case up:
			s.add(b.moved(left))
		case down:
			s.add(b.moved(right))
		case left:
			s.add(b.moved(up))
		case right:
			s.add(b.moved(down))
		default:
			panic("invalid direction")
		}
	default:
		panic("invalid symbol in layout")
	}

	return true
}

func (s *simulation) countEnergized() int {
	count := 0
	for _, row := range s.energized {
		for _, cell := range row {
			if cell {
				count++
			}
		}
	}
	return count
}

func (s *simulation) reset() {
	clear(s.beams)
	s.beams = s.beams[:0]

	clear(s.seen)

	for i := range s.energized {
		clear(s.energized[i])
	}
}

func findMaxEnergized(layout [][]byte) int {
	sim := newSimulation(layout)

	maxEnergized := 0

	for row := range layout {
		sim.add(beam{
			position:  vector{row: row, col: 0},
			direction: right,
		})
		sim.run()
		maxEnergized = max(maxEnergized, sim.countEnergized())
		sim.reset()

		sim.add(beam{
			position:  vector{row: row, col: len(layout) - 1},
			direction: left,
		})
		sim.run()
		maxEnergized = max(maxEnergized, sim.countEnergized())
		sim.reset()
	}

	for col := range layout {
		sim.add(beam{
			position:  vector{row: 0, col: col},
			direction: down,
		})
		sim.run()
		maxEnergized = max(maxEnergized, sim.countEnergized())
		sim.reset()

		sim.add(beam{
			position:  vector{row: len(layout) - 1, col: col},
			direction: up,
		})
		sim.run()
		maxEnergized = max(maxEnergized, sim.countEnergized())
		sim.reset()
	}

	return maxEnergized
}

func layoutFromReader(r io.Reader) ([][]byte, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	layout := make([][]byte, len(lines))
	for row, line := range lines {
		layout[row] = []byte(line)
		if len(layout[row]) != len(layout) {
			return nil, fmt.Errorf("layout is not square")
		}
		for col := range layout[row] {
			switch layout[row][col] {
			case empty, vSplit, hSplit, trblMirror, tlbrMirror:
				// valid symbols
			default:
				return nil, fmt.Errorf("invalid symbol: %c", layout[row][col])
			}
		}
	}

	return layout, nil
}
