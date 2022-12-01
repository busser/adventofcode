package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 22 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	g, err := gameFromReader(input, newSliceDeck)
	// g, err := gameFromReader(input, newRingDeck)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	winner := g.playUntilOver()

	_, err = fmt.Fprintf(answer, "%d", winner.score())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 22 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	g, err := gameFromReader(input, newSliceDeck)
	// g, err := gameFromReader(input, newRingDeck)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	winner := g.playRecursivelyUntilOver()

	_, err = fmt.Fprintf(answer, "%d", winner.score())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func gameFromReader(r io.Reader, newDeck func() deck) (game, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return game{}, fmt.Errorf("could not read lines: %w", err)
	}

	chunks := splitSlice(lines, "")
	if len(chunks) != 2 {
		return game{}, errors.New("wrong format")
	}

	if len(chunks[0]) == 0 || len(chunks[1]) == 0 {
		return game{}, errors.New("wrong format")
	}

	deck1 := newDeck()
	deck2 := newDeck()

	if err := insertLinesIntoDeck(chunks[0][1:], deck1); err != nil {
		return game{}, fmt.Errorf("reading deck 1: %w", err)
	}

	if err := insertLinesIntoDeck(chunks[1][1:], deck2); err != nil {
		return game{}, fmt.Errorf("reading deck 2: %w", err)
	}

	return game{deck1, deck2}, nil
}

func splitSlice(slice []string, sep string) [][]string {
	var split [][]string

	start := 0

	for end := range slice {
		if slice[end] == sep {
			split = append(split, slice[start:end])
			start = end + 1
		}
	}
	split = append(split, slice[start:])

	return split
}

func insertLinesIntoDeck(lines []string, d deck) error {
	for _, l := range lines {
		v, err := strconv.Atoi(l)
		if err != nil {
			return fmt.Errorf("%q is not an integer", l)
		}
		d.insertAtBottom(uint8(v))
	}
	return nil
}
