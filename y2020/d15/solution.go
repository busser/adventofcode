package busser

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 15 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	startingNumbers, err := startingNumbersFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	g := newGame()

	for _, n := range startingNumbers {
		g.speak(n)
	}

	for i := len(startingNumbers); i < 2020; i++ {
		g.speak(g.nextNumber())
	}

	_, err = fmt.Fprintf(answer, "%d", g.mostRecentlySpoken)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 15 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	startingNumbers, err := startingNumbersFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	g := newGame()

	for _, n := range startingNumbers {
		g.speak(n)
	}

	for i := len(startingNumbers); i < 30000000; i++ {
		g.speak(g.nextNumber())
	}

	_, err = fmt.Fprintf(answer, "%d", g.mostRecentlySpoken)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type game struct {
	currentTurn        int
	mostRecentlySpoken int
	lastTurnByNumber   map[int]int
}

func newGame() game {
	return game{
		lastTurnByNumber: make(map[int]int),
	}
}

func (g *game) speak(number int) {
	if g.currentTurn > 0 {
		g.lastTurnByNumber[g.mostRecentlySpoken] = g.currentTurn
	}

	g.currentTurn++
	g.mostRecentlySpoken = number
}

func (g game) nextNumber() int {
	lastTurn, spoken := g.lastTurnByNumber[g.mostRecentlySpoken]
	if !spoken {
		return 0
	}
	return g.currentTurn - lastTurn
}

func startingNumbersFromReader(r io.Reader) ([]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading line: %w", err)
	}

	if len(lines) != 1 {
		return nil, errors.New("wrong format")
	}

	numbers, err := helpers.IntsFromString(lines[0], ",")
	if err != nil {
		return nil, fmt.Errorf("parsing numbers: %w", err)
	}

	return numbers, nil
}
