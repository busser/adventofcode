package busser

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 6 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	gen, err := firstGenerationFromReader(input)
	if err != nil {
		return fmt.Errorf("could not get first generation: %w", err)
	}

	for i := 0; i < 80; i++ {
		gen = gen.next()
	}

	total := 0
	for _, count := range gen {
		total += count
	}

	_, err = fmt.Fprintf(answer, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 6 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	gen, err := firstGenerationFromReader(input)
	if err != nil {
		return fmt.Errorf("could not get first generation: %w", err)
	}

	for i := 0; i < 256; i++ {
		gen = gen.next()
	}

	total := 0
	for _, count := range gen {
		total += count
	}

	_, err = fmt.Fprintf(answer, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const (
	defaultDaysToReproduce = 7
	maxdaysToReproduce     = defaultDaysToReproduce + 2
)

type generation [maxdaysToReproduce]int

func (gen generation) next() generation {
	var nextGen generation

	nextGen[maxdaysToReproduce-1] += gen[0]     // new fish babies
	nextGen[defaultDaysToReproduce-1] += gen[0] // fish that just reproduced

	for age := 1; age < maxdaysToReproduce; age++ {
		nextGen[age-1] += gen[age] // fish that get closer to reproducing
	}

	return nextGen
}

func firstGenerationFromReader(r io.Reader) (generation, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return generation{}, fmt.Errorf("could not read: %w", err)
	}

	if len(lines) != 1 {
		return generation{}, fmt.Errorf("expected 1 line of input, got %d", len(lines))
	}

	ages := helpers.IntsFromString(lines[0])

	var gen generation
	for _, a := range ages {
		if a >= maxdaysToReproduce {
			return generation{}, fmt.Errorf("age %d should be less than %d", a, maxdaysToReproduce)
		}
		gen[a]++
	}

	return gen, nil
}
