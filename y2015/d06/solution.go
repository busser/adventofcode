package d06

import (
	"fmt"
	"io"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 6 of Advent of Code 2015.
func PartOne(r io.Reader, w io.Writer) error {
	instructions, err := instructionsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	g := newBinaryGrid()
	for _, i := range instructions {
		g.applyInstruction(i)
	}

	count := g.countLights()

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 6 of Advent of Code 2015.
func PartTwo(r io.Reader, w io.Writer) error {
	instructions, err := instructionsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	g := newLeveledGrid()
	for _, i := range instructions {
		g.applyInstruction(i)
	}

	brightness := g.measureBrightness()

	_, err = fmt.Fprintf(w, "%d", brightness)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const gridSize = 1000

type binaryGrid [][]bool

func newBinaryGrid() binaryGrid {
	g := make(binaryGrid, gridSize)
	for i := range g {
		g[i] = make([]bool, gridSize)
	}
	return g
}

func (g binaryGrid) countLights() int {
	count := 0
	for _, row := range g {
		for _, light := range row {
			if light {
				count++
			}
		}
	}
	return count
}

func (g binaryGrid) applyOperationToLight(op operation, x, y int) {
	switch op {
	case turnOn:
		g[x][y] = true
	case turnOff:
		g[x][y] = false
	case toggle:
		g[x][y] = !g[x][y]
	}
}

func (g binaryGrid) applyInstruction(i instruction) {
	for x := i.start.x; x <= i.end.x; x++ {
		for y := i.start.y; y <= i.end.y; y++ {
			g.applyOperationToLight(i.operation, x, y)
		}
	}
}

type leveledGrid [][]int

func newLeveledGrid() leveledGrid {
	g := make(leveledGrid, gridSize)
	for i := range g {
		g[i] = make([]int, gridSize)
	}
	return g
}

func (g leveledGrid) measureBrightness() int {
	brightness := 0
	for _, row := range g {
		for _, light := range row {
			brightness += light
		}
	}
	return brightness
}

func (g leveledGrid) applyOperationToLight(op operation, x, y int) {
	switch op {
	case turnOn:
		g[x][y]++
	case turnOff:
		if g[x][y] > 0 {
			g[x][y]--
		}
	case toggle:
		g[x][y] += 2
	}
}

func (g leveledGrid) applyInstruction(i instruction) {
	for x := i.start.x; x <= i.end.x; x++ {
		for y := i.start.y; y <= i.end.y; y++ {
			g.applyOperationToLight(i.operation, x, y)
		}
	}
}

type operation uint8

const (
	turnOn operation = iota
	turnOff
	toggle
)

type vector struct {
	x, y int
}

type instruction struct {
	operation operation
	start     vector
	end       vector
}

func instructionFromString(s string) (instruction, error) {
	var op operation
	switch {
	case strings.HasPrefix(s, "turn on"):
		op = turnOn
	case strings.HasPrefix(s, "turn off"):
		op = turnOff
	case strings.HasPrefix(s, "toggle"):
		op = toggle
	default:
		return instruction{}, fmt.Errorf("invalid operation: %q", s)
	}

	coordinates := helpers.IntsFromString(s)
	if len(coordinates) != 4 {
		return instruction{}, fmt.Errorf("invalid instruction: %q", s)
	}

	i := instruction{
		operation: op,
		start:     vector{coordinates[0], coordinates[1]},
		end:       vector{coordinates[2], coordinates[3]},
	}

	if err := validateInstruction(i); err != nil {
		return instruction{}, fmt.Errorf("invalid instruction: %w", err)
	}

	return i, nil
}

func instructionsFromReader(r io.Reader) ([]instruction, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	instructions := make([]instruction, len(lines))
	for i, line := range lines {
		instructions[i], err = instructionFromString(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse instruction %d: %w", i, err)
		}
	}

	return instructions, nil
}

func validateInstruction(i instruction) error {
	if i.start.x < 0 || i.start.y < 0 || i.end.x < 0 || i.end.y < 0 {
		return fmt.Errorf("invalid instruction: negative coordinates")
	}

	if i.start.x > gridSize || i.start.y > gridSize || i.end.x > gridSize || i.end.y > gridSize {
		return fmt.Errorf("invalid instruction: coordinates too large")
	}

	return nil
}
