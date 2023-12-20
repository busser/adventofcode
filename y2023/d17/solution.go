package d17

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 17 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	heatLossMap, err := heatLossMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	heatLoss := totalHeatLoss(heatLossMap, 0, 3)

	_, err = fmt.Fprintf(w, "%d", heatLoss)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 17 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	heatLossMap, err := heatLossMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	heatLoss := totalHeatLoss(heatLossMap, 4, 10)

	_, err = fmt.Fprintf(w, "%d", heatLoss)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type vector struct {
	row, col int
}

func (v vector) add(o vector) vector {
	return vector{v.row + o.row, v.col + o.col}
}

func (v vector) rotateLeft() vector {
	return vector{-v.col, v.row}
}

func (v vector) rotateRight() vector {
	return vector{v.col, -v.row}
}

type state struct {
	postion             vector
	direction           vector
	iterationsSinceTurn int
	totalHeatLoss       int
}

type stateQueue []state

func (q stateQueue) less(i, j int) bool {
	return q[i].totalHeatLoss < q[j].totalHeatLoss
}

func (q stateQueue) swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *stateQueue) push(x state) {
	*q = append(*q, x)
	q.up(len(*q) - 1)
}

func (q *stateQueue) pop() state {
	old := *q
	n := len(old)
	x := old[0]
	old[0] = old[n-1]
	*q = old[0 : n-1]
	if n > 1 {
		q.down(0)
	}
	return x
}

func (q *stateQueue) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !q.less(j, i) {
			break
		}
		q.swap(i, j)
		j = i
	}
}

func (q *stateQueue) down(i0 int) {
	n := len(*q)
	i := i0
	for {
		left := 2*i + 1
		if left >= n || left < 0 { // left < 0 after int overflow
			break
		}
		j := left // left child
		if right := left + 1; right < n && q.less(right, left) {
			j = right // = 2*i + 2  // right child
		}
		if !q.less(j, i) {
			break
		}
		q.swap(i, j)
		i = j
	}
}

func totalHeatLoss(m [][]int, minStraightLine, maxStraightLine int) int {
	start := vector{0, 0}
	end := vector{len(m) - 1, len(m[0]) - 1}

	startState := state{
		postion:             start,
		direction:           vector{0, 1}, // arbitrary
		iterationsSinceTurn: 0,
		totalHeatLoss:       0,
	}

	history := make(map[state]bool)
	visited := func(s state) bool {
		// Ignoring totalHeatLoss in visited check because even if it's greater
		// we don't want to visit it again; the optimal both does not have
		// cycles.
		s.totalHeatLoss = 0
		return history[s]
	}
	remember := func(s state) {
		// Ignoring totalHeatLoss in visited check. See above for why.
		s.totalHeatLoss = 0
		history[s] = true
	}

	remember(startState)

	q := stateQueue{}
	q.push(startState)

	for len(q) > 0 {
		s := q.pop()

		if s.postion == end && s.iterationsSinceTurn >= minStraightLine {
			return s.totalHeatLoss
		}

		var possibleDirections []vector
		if s.iterationsSinceTurn >= minStraightLine {
			possibleDirections = append(possibleDirections, s.direction.rotateLeft(), s.direction.rotateRight())
		}
		if s.iterationsSinceTurn < maxStraightLine {
			possibleDirections = append(possibleDirections, s.direction)
		}

		for _, d := range possibleDirections {
			nextPosition := s.postion.add(d)
			if nextPosition.row < 0 || nextPosition.row >= len(m) || nextPosition.col < 0 || nextPosition.col >= len(m[nextPosition.row]) {
				continue
			}

			iterationsSinceTurn := 1
			if d == s.direction {
				iterationsSinceTurn = s.iterationsSinceTurn + 1
			}

			nextState := state{
				postion:             nextPosition,
				direction:           d,
				iterationsSinceTurn: iterationsSinceTurn,
				totalHeatLoss:       s.totalHeatLoss + m[nextPosition.row][nextPosition.col],
			}

			if visited(nextState) {
				continue
			}
			remember(nextState)

			q.push(nextState)
		}
	}

	return -1
}

func heatLossMapFromReader(r io.Reader) ([][]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	m := make([][]int, len(lines))
	for row := range m {
		m[row] = make([]int, len(lines[row]))
		for col := range m[row] {
			m[row][col] = int(lines[row][col] - '0')
			if m[row][col] <= 0 || m[row][col] > 9 {
				return nil, fmt.Errorf("row %d, col %d: invalid heat loss %d", row, col, m[row][col])
			}
		}
	}

	if len(m) == 0 {
		return nil, fmt.Errorf("no lines")
	}

	rowLength := len(m[0])
	for row := range m {
		if len(m[row]) != rowLength {
			return nil, errors.New("map is not a rectangle")
		}
	}

	return m, nil
}
