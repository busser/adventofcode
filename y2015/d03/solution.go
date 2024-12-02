package d03

import (
	"fmt"
	"io"
	"slices"
)

// PartOne solves the first problem of day 3 of Advent of Code 2015.
func PartOne(r io.Reader, w io.Writer) error {
	moves, err := movesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := processMoves(moves)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 3 of Advent of Code 2015.
func PartTwo(r io.Reader, w io.Writer) error {
	moves, err := movesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := processMovesWithRoboSanta(moves)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func processMoves(moves []vector) int {
	visited := make([]vector, len(moves)+1)

	position := vector{0, 0}
	visited[0] = position

	for i, move := range moves {
		position = position.plus(move)
		visited[i+1] = position
	}

	return len(uniqueVectors(visited))
}

func processMovesWithRoboSanta(moves []vector) int {
	visited := make([]vector, len(moves)+1)

	position := vector{0, 0}
	visited[0] = position

	for i := 0; i < len(moves); i += 2 {
		position = position.plus(moves[i])
		visited[i+1] = position
	}

	position = vector{0, 0}
	visited[1] = position

	for i := 1; i < len(moves); i += 2 {
		position = position.plus(moves[i])
		visited[i+1] = position
	}

	return len(uniqueVectors(visited))
}

func uniqueVectors(vectors []vector) []vector {
	slices.SortFunc(
		vectors,
		func(a, b vector) int {
			if a.x == b.x {
				return a.y - b.y
			}
			return a.x - b.x
		},
	)

	unique := slices.Compact(vectors)

	return unique
}

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

func movesFromReader(r io.Reader) ([]vector, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	moves := make([]vector, len(data))
	for i, c := range data {
		switch c {
		case '^':
			moves[i] = up
		case 'v':
			moves[i] = down
		case '<':
			moves[i] = left
		case '>':
			moves[i] = right
		default:
			return nil, fmt.Errorf("unknown symbol %q", c)
		}
	}

	return moves, nil
}
