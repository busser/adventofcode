package d21

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 21 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	gardenMap, startingPos, err := gardenMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := reachableGardenPlots(gardenMap, startingPos, 64)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 21 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	gardenMap, startingPos, err := gardenMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := reachableGardenPlotsInfiniteGarden(gardenMap, startingPos, 26501365)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const (
	start      = 'S'
	gardenPlot = '.'
	rock       = '#'
)

type vector struct {
	row, col int
}

var (
	up    = vector{-1, 0}
	down  = vector{1, 0}
	left  = vector{0, -1}
	right = vector{0, 1}
)

func (v vector) add(other vector) vector {
	return vector{
		row: v.row + other.row,
		col: v.col + other.col,
	}
}

func reachableGardenPlotsInfiniteGarden(gardenMap [][]byte, startingPos vector, steps int) int {
	// This solution relies on assumptions about the input data.
	// It is not a general solution. It doesn't work on the example input.
	//
	// The input data has a peculiarity: there are no rocks on the same row and
	// column as the starting position. This means that the elf can go from the
	// starting position to the same position in a neighboring garden in a
	// number of steps equal to the size of the garden.
	//
	// From this, we can deduce that the number of garden tiles reachable in a
	// given number of steps grows quadratically:
	//   with input size S,
	//   with R(x) the number of reachable tiles in x steps,
	//   this series: R(n), R(n + S), R(n + 2S), ... is quadratic.
	// By finding the quadratic function that fits the first three values of
	// this series, we can find the number of reachable tiles in any number of
	// steps.
	//
	// We choose the first n value based on the final number of steps we want
	// to simulate:
	//   with T the number of steps to simulate,
	//   we choose n = T % S.
	//
	// So we need to compute R(T % S), R(T % S + S), R(T % S + 2S).
	// From there, we can compute R(T).

	gardenSize := len(gardenMap)

	largerGardenMap := make([][]byte, gardenSize*5)
	for row := range largerGardenMap {
		largerGardenMap[row] = make([]byte, gardenSize*5)
		for col := range largerGardenMap[row] {
			largerGardenMap[row][col] = gardenMap[row%gardenSize][col%gardenSize]
		}
	}

	startingPos = startingPos.add(vector{gardenSize * 2, gardenSize * 2})

	// Get first values of quadratic series.
	v0 := reachableGardenPlots(largerGardenMap, startingPos, steps%gardenSize+0*gardenSize)
	v1 := reachableGardenPlots(largerGardenMap, startingPos, steps%gardenSize+1*gardenSize)
	v2 := reachableGardenPlots(largerGardenMap, startingPos, steps%gardenSize+2*gardenSize)

	// Solve quadratic equation.
	a := (v2 - 2*v1 + v0) / 2
	b := v1 - v0 - a
	c := v0

	// Compute result.
	n := steps / gardenSize
	result := a*n*n + b*n + c

	return result
}

func reachableGardenPlots(gardenMap [][]byte, startingPos vector, steps int) int {
	visited := make([][]bool, len(gardenMap))
	for row := range visited {
		visited[row] = make([]bool, len(gardenMap[row]))
	}

	var (
		queue, nextQueue []vector
		count            = 0
		remainingSteps   = steps
	)

	nextQueue = append(nextQueue, startingPos)
	visited[startingPos.row][startingPos.col] = true
	if remainingSteps%2 == 0 {
		count++
	}

	for remainingSteps > 0 {
		remainingSteps--

		queue, nextQueue = nextQueue, queue[:0]

		visitPosition := func(pos vector) {
			if pos.row < 0 || pos.row >= len(gardenMap) || pos.col < 0 || pos.col >= len(gardenMap[pos.row]) {
				return
			}

			if gardenMap[pos.row][pos.col] == rock {
				return
			}

			if visited[pos.row][pos.col] {
				return
			}
			visited[pos.row][pos.col] = true

			nextQueue = append(nextQueue, pos)
			if remainingSteps%2 == 0 {
				count++
			}
		}

		for _, pos := range queue {
			visitPosition(pos.add(up))
			visitPosition(pos.add(down))
			visitPosition(pos.add(left))
			visitPosition(pos.add(right))
		}
	}

	return count
}

func gardenMapFromReader(r io.Reader) ([][]byte, vector, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, vector{}, fmt.Errorf("could not read input: %w", err)
	}

	var (
		startPos  vector
		gardenMap = make([][]byte, len(lines))
	)

	for row, line := range lines {
		gardenMap[row] = []byte(line)

		for col, char := range line {
			switch char {
			case start:
				startPos = vector{row, col}
				gardenMap[row][col] = gardenPlot
			case gardenPlot, rock:
				// do nothing
			default:
				return nil, vector{}, fmt.Errorf("invalid character %q", char)
			}
		}
	}

	for row := range gardenMap {
		if gardenMap[row][startPos.col] == rock {
			return nil, vector{}, fmt.Errorf("rock on the same column as the starting position")
		}
		if len(gardenMap[row]) != len(gardenMap) {
			return nil, vector{}, fmt.Errorf("garden map is not square")
		}
	}
	for col := range gardenMap[startPos.row] {
		if gardenMap[startPos.row][col] == rock {
			return nil, vector{}, fmt.Errorf("rock on the same row as the starting position")
		}
	}

	return gardenMap, startPos, nil
}
