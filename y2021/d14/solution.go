package busser

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 14 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	poly, rules, err := polymerAndRulesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	for i := 0; i < 10; i++ {
		poly.react(rules)
	}

	minCount, maxCount := poly.stats()

	_, err = fmt.Fprintf(answer, "%d", maxCount-minCount)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 14 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	poly, rules, err := polymerAndRulesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	for i := 0; i < 40; i++ {
		poly.react(rules)
	}

	minCount, maxCount := poly.stats()

	_, err = fmt.Fprintf(answer, "%d", maxCount-minCount)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type polymer struct {
	elementPairs map[[2]rune]int
	elementCount map[rune]int
}

type insertionRules map[[2]rune]rune

func (poly polymer) react(rules insertionRules) {
	changes := make(map[[2]rune]int)
	for pair, count := range poly.elementPairs {
		if newElement, reactionHappens := rules[pair]; reactionHappens {
			changes[[2]rune{pair[0], newElement}] += count
			changes[[2]rune{newElement, pair[1]}] += count
			changes[pair] -= count

			poly.elementCount[newElement] += count
		}
	}
	for pair, diff := range changes {
		poly.elementPairs[pair] += diff
	}
}

func (poly polymer) stats() (minCount, maxCount int) {
	minCount, maxCount = math.MaxInt, 0
	for _, count := range poly.elementCount {
		if count < minCount {
			minCount = count
		}
		if count > maxCount {
			maxCount = count
		}
	}
	return minCount, maxCount
}

func polymerAndRulesFromReader(r io.Reader) (polymer, insertionRules, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return polymer{}, nil, err
	}

	if len(lines) < 2 {
		return polymer{}, nil, errors.New("not enough lines")
	}

	poly := polymerFromString(lines[0])

	rules := make(insertionRules)
	for _, l := range lines[2:] {
		pair, insert, err := ruleFromString(l)
		if err != nil {
			return polymer{}, nil, fmt.Errorf("invalid rule %q: %w", l, err)
		}
		rules[pair] = insert
	}

	return poly, rules, nil
}

func polymerFromString(s string) polymer {
	elements := []rune(s)
	poly := polymer{
		elementPairs: make(map[[2]rune]int),
		elementCount: make(map[rune]int),
	}
	for i := 1; i < len(elements); i++ {
		pair := [2]rune{elements[i-1], elements[i]}
		poly.elementPairs[pair]++
	}
	for i := 0; i < len(elements); i++ {
		poly.elementCount[elements[i]]++
	}
	return poly
}

func ruleFromString(s string) (pair [2]rune, insert rune, err error) {
	parts := strings.Split(s, " -> ")
	if len(parts) != 2 {
		return [2]rune{}, 0, errors.New("invalid format")
	}

	rawPair, rawInsert := []rune(parts[0]), []rune(parts[1])
	if len(rawPair) != 2 {
		return [2]rune{}, 0, errors.New("invalid format")
	}
	if len(rawInsert) != 1 {
		return [2]rune{}, 0, errors.New("invalid format")
	}

	return [2]rune{rawPair[0], rawPair[1]}, rawInsert[0], nil
}
