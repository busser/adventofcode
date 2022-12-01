package busser

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 11 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	grid, err := octopusGridFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	totalFlashes := 0
	for i := 0; i < 100; i++ {
		totalFlashes += grid.step()
	}

	_, err = fmt.Fprintf(answer, "%d", totalFlashes)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 11 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	grid, err := octopusGridFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var flashes, step int
	for ; flashes != gridSize*gridSize; step++ {
		flashes = grid.step()
	}

	_, err = fmt.Fprintf(answer, "%d", step)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type octopus struct {
	energyLevel int
	flashing    bool
}

const gridSize = 10

type octopusGrid [gridSize][gridSize]octopus

func (grid *octopusGrid) step() int {
	for i := range grid {
		for j := range grid[i] {
			grid.addEnergy(i, j)
		}
	}

	flashingCount := 0
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j].flashing {
				flashingCount++
			}
		}
	}

	for i := range grid {
		for j := range grid[i] {
			if grid[i][j].flashing {
				grid[i][j].energyLevel = 0
				grid[i][j].flashing = false
			}
		}
	}

	return flashingCount
}

func (grid *octopusGrid) addEnergy(i, j int) {
	if i < 0 || j < 0 || i >= gridSize || j >= gridSize {
		return
	}
	if grid[i][j].flashing {
		return
	}

	grid[i][j].energyLevel++
	if grid[i][j].energyLevel > 9 {
		grid[i][j].flashing = true

		grid.addEnergy(i-1, j-1)
		grid.addEnergy(i-1, j)
		grid.addEnergy(i-1, j+1)
		grid.addEnergy(i, j-1)
		grid.addEnergy(i, j+1)
		grid.addEnergy(i+1, j-1)
		grid.addEnergy(i+1, j)
		grid.addEnergy(i+1, j+1)
	}
}

func octopusGridFromReader(r io.Reader) (octopusGrid, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return octopusGrid{}, err
	}
	if len(lines) != gridSize {
		return octopusGrid{}, fmt.Errorf("expected %d lines, found %d", gridSize, len(lines))
	}

	var grid octopusGrid

	for i, l := range lines {
		if len(l) != gridSize {
			return octopusGrid{}, fmt.Errorf("expected line %d to have %d characters, found %d", i+1, gridSize, len(l))
		}
		for j, c := range l {
			if c < '0' || c > '9' {
				return octopusGrid{}, fmt.Errorf("%q is not a digit", c)
			}
			grid[i][j].energyLevel = int(c - '0')
		}
	}

	return grid, nil
}
