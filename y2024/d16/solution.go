package d16

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 16 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	maze, err := reindeerMazeFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	costs := maze.computeCosts()
	minScore := maze.findMinimumScore(costs)

	_, err = fmt.Fprintf(w, "%d", minScore)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 16 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	maze, err := reindeerMazeFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := maze.countPositionsOnBestPaths()

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

func (v vector) minus(w vector) vector {
	return vector{
		row: v.row - w.row,
		col: v.col - w.col,
	}
}

var (
	north = vector{row: -1, col: 0}
	south = vector{row: 1, col: 0}
	west  = vector{row: 0, col: -1}
	east  = vector{row: 0, col: 1}
)

var allDirections = []vector{north, south, west, east}

func (v vector) rotateClockwise() vector {
	return vector{
		row: v.col,
		col: -v.row,
	}
}

func (v vector) rotateCounterClockwise() vector {
	return vector{
		row: -v.col,
		col: v.row,
	}
}

const (
	empty = '.'
	wall  = '#'
	start = 'S'
	end   = 'E'
)

type reindeerMaze struct {
	tiles          [][]byte
	startPosition  vector
	startDirection vector
	endPosition    vector
}

func (m reindeerMaze) at(v vector) byte {
	return m.tiles[v.row][v.col]
}

func (m reindeerMaze) isWithinBounds(v vector) bool {
	return v.row >= 0 && v.row < len(m.tiles) && v.col >= 0 && v.col < len(m.tiles[v.row])
}

type stateKey struct {
	position  vector
	direction vector
}

type searchState struct {
	stateKey
	score int
}

func (s searchState) nextForward() searchState {
	return searchState{
		stateKey: stateKey{
			position:  s.position.plus(s.direction),
			direction: s.direction,
		},
		score: s.score + 1,
	}
}

func (s searchState) nextRight() searchState {
	return searchState{
		stateKey: stateKey{
			position:  s.position,
			direction: s.direction.rotateClockwise(),
		},
		score: s.score + 1000,
	}
}

func (s searchState) nextLeft() searchState {
	return searchState{
		stateKey: stateKey{
			position:  s.position,
			direction: s.direction.rotateCounterClockwise(),
		},
		score: s.score + 1000,
	}
}

func (s searchState) previousForward() searchState {
	return searchState{
		stateKey: stateKey{
			position:  s.position.minus(s.direction),
			direction: s.direction,
		},
		score: s.score - 1,
	}
}

func (s searchState) previousRight() searchState {
	return searchState{
		stateKey: stateKey{
			position:  s.position,
			direction: s.direction.rotateCounterClockwise(),
		},
		score: s.score - 1000,
	}
}

func (s searchState) previousLeft() searchState {
	return searchState{
		stateKey: stateKey{
			position:  s.position,
			direction: s.direction.rotateClockwise(),
		},
		score: s.score - 1000,
	}
}

func (m reindeerMaze) computeCosts() map[stateKey]int {
	costs := make(map[stateKey]int)

	next := helpers.NewPriorityQueue(func(a, b searchState) bool {
		return a.score < b.score
	})
	next.Push(searchState{stateKey{m.startPosition, m.startDirection}, 0})

	for next.Len() > 0 {
		current := next.Pop()

		if !m.isWithinBounds(current.position) || m.at(current.position) == wall {
			continue
		}

		if _, visited := costs[current.stateKey]; visited {
			continue
		}
		costs[current.stateKey] = current.score

		if current.position == m.endPosition {
			return costs
		}

		next.Push(current.nextForward())
		next.Push(current.nextRight())
		next.Push(current.nextLeft())
	}

	return costs
}

func (m reindeerMaze) findMinimumScore(costs map[stateKey]int) int {
	for _, dir := range allDirections {
		state := stateKey{m.endPosition, dir}
		if cost, ok := costs[state]; ok {
			return cost
		}
	}

	return -1 // No path found.
}

func (m reindeerMaze) computeReversedCosts() map[stateKey]int {
	costs := make(map[stateKey]int)
	next := helpers.NewPriorityQueue(func(a, b searchState) bool {
		return a.score > b.score
	})
	for _, dir := range allDirections {
		next.Push(searchState{stateKey{m.endPosition, dir}, 0})
	}

	for next.Len() > 0 {
		current := next.Pop()

		if !m.isWithinBounds(current.position) || m.at(current.position) == wall {
			continue
		}

		if _, visited := costs[current.stateKey]; visited {
			continue
		}
		costs[current.stateKey] = current.score

		if current.position == m.startPosition && current.direction == m.startDirection {
			return costs
		}

		next.Push(current.previousForward())
		next.Push(current.previousRight())
		next.Push(current.previousLeft())
	}

	return costs
}

func (m reindeerMaze) countPositionsOnBestPaths() int {
	forwardCosts := m.computeCosts()
	bestScore := m.findMinimumScore(forwardCosts)
	reverseCosts := m.computeReversedCosts()

	var statesOnBestPaths []stateKey
	for state, forwardCost := range forwardCosts {
		reverseCost := reverseCosts[state]
		if forwardCost-reverseCost == bestScore {
			statesOnBestPaths = append(statesOnBestPaths, state)
		}
	}

	positionsOnBestPaths := make(map[vector]struct{})
	for _, state := range statesOnBestPaths {
		positionsOnBestPaths[state.position] = struct{}{}
	}

	return len(positionsOnBestPaths)
}

func reindeerMazeFromReader(r io.Reader) (reindeerMaze, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return reindeerMaze{}, fmt.Errorf("could not read input: %w", err)
	}

	var maze reindeerMaze

	maze.tiles = make([][]byte, len(lines))
	for row, line := range lines {
		maze.tiles[row] = []byte(line)
		for col, char := range line {
			switch char {
			case empty, wall:
				// do nothing
			case start:
				maze.startPosition = vector{row, col}
				maze.startDirection = east
			case end:
				maze.endPosition = vector{row, col}
			default:
				return reindeerMaze{}, fmt.Errorf("invalid character %c at row %d, column %d", char, row, col)
			}
		}
	}

	return maze, nil
}
