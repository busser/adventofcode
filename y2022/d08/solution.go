package d08

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 8 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	heightMap, err := heightMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_, err = fmt.Fprintf(w, "%d", countVisibleTrees(heightMap))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 8 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	heightMap, err := heightMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_, err = fmt.Fprintf(w, "%d", bestScenicScore(heightMap))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func countVisibleTrees(heightMap [][]int) int {
	visible := make([][]bool, len(heightMap))
	for row := range heightMap {
		visible[row] = make([]bool, len(heightMap[row]))
	}

	// Find trees visible from the top.
	for col := 0; col < len(heightMap[0]); col++ {
		maxHeight := -1
		for row := 0; row < len(heightMap); row++ {
			height := heightMap[row][col]
			if height > maxHeight {
				visible[row][col] = true
				maxHeight = height
			}
		}
	}

	// Find trees visible from the bottom.
	for col := 0; col < len(heightMap[0]); col++ {
		maxHeight := -1
		for row := len(heightMap) - 1; row >= 0; row-- {
			height := heightMap[row][col]
			if height > maxHeight {
				visible[row][col] = true
				maxHeight = height
			}
		}
	}

	// Find trees visible from the left.
	for row := 0; row < len(heightMap); row++ {
		maxHeight := -1
		for col := 0; col < len(heightMap[0]); col++ {
			height := heightMap[row][col]
			if height > maxHeight {
				visible[row][col] = true
				maxHeight = height
			}
		}
	}

	// Find trees visible from the right.
	for row := 0; row < len(heightMap); row++ {
		maxHeight := -1
		for col := len(heightMap[0]) - 1; col >= 0; col-- {
			height := heightMap[row][col]
			if height > maxHeight {
				visible[row][col] = true
				maxHeight = height
			}
		}
	}

	// Count visible trees.
	count := 0
	for row := range visible {
		for col := range visible[row] {
			if visible[row][col] {
				count++
			}
		}
	}

	return count
}

func bestScenicScore(heightMap [][]int) int {
	bestScore := 0

	for row := range heightMap {
		for col := range heightMap[row] {
			score := scenicScore(heightMap, row, col)
			if score > bestScore {
				bestScore = score
			}
		}
	}

	return bestScore
}

func scenicScore(heightMap [][]int, row, col int) int {
	return visibleToRight(heightMap, row, col) *
		visibleToLeft(heightMap, row, col) *
		visibleToBottom(heightMap, row, col) *
		visibleToTop(heightMap, row, col)
}

func visibleToRight(heightMap [][]int, row, col int) int {
	count := 0
	maxHeight := heightMap[row][col]

	for col++; col < len(heightMap[row]); col++ {
		count++
		if heightMap[row][col] >= maxHeight {
			break
		}
	}

	return count
}

func visibleToLeft(heightMap [][]int, row, col int) int {
	count := 0
	maxHeight := heightMap[row][col]

	for col--; col >= 0; col-- {
		count++
		if heightMap[row][col] >= maxHeight {
			break
		}
	}

	return count
}

func visibleToBottom(heightMap [][]int, row, col int) int {
	count := 0
	maxHeight := heightMap[row][col]

	for row++; row < len(heightMap); row++ {
		count++
		if heightMap[row][col] >= maxHeight {
			break
		}
	}

	return count
}

func visibleToTop(heightMap [][]int, row, col int) int {
	count := 0
	maxHeight := heightMap[row][col]

	for row--; row >= 0; row-- {
		count++
		if heightMap[row][col] >= maxHeight {
			break
		}
	}

	return count
}

func heightMapFromReader(r io.Reader) ([][]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return nil, errors.New("empty heightmap")
	}

	heightMap := make([][]int, len(lines))

	for i, l := range lines {
		heightMap[i] = make([]int, len(l))

		for j, c := range l {
			heightMap[i][j] = int(c - '0')
		}
	}

	rowLen := len(heightMap[0])
	for _, row := range heightMap {
		if len(row) != rowLen {
			return nil, errors.New("map is not a rectangle")
		}
	}

	if rowLen == 0 {
		return nil, errors.New("empty heightmap")
	}

	return heightMap, nil
}
