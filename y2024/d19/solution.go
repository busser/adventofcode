package d19

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 19 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	towels, designs, err := towelsAndDesignsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	pc := newPossibilityChecker(towels)

	count := 0
	for _, d := range designs {
		if pc.isPossible(d) {
			count++
		}
	}

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 19 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	towels, designs, err := towelsAndDesignsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	pc := newPossibilityCounter(towels)

	total := 0
	for _, d := range designs {
		total += pc.countPossibilities(d)
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type towel []rune

type design []rune

type possibilityChecker struct {
	towels []towel
	cache  map[string]bool
}

func newPossibilityChecker(towels []towel) *possibilityChecker {
	return &possibilityChecker{
		towels: towels,
		cache:  make(map[string]bool),
	}
}

func (pc *possibilityChecker) isPossible(d design) (result bool) {
	if len(d) == 0 {
		return true
	}

	cacheKey := string(d)
	if possible, hit := pc.cache[cacheKey]; hit {
		return possible
	}
	defer func() { pc.cache[cacheKey] = result }()

	for _, t := range pc.towels {
		if !designCanStartWithTowel(d, t) {
			continue
		}

		remainder := d[len(t):]
		if pc.isPossible(remainder) {
			return true
		}
	}

	return false
}

type possibilityCounter struct {
	towels []towel
	cache  map[string]int
}

func newPossibilityCounter(towels []towel) *possibilityCounter {
	return &possibilityCounter{
		towels: towels,
		cache:  make(map[string]int),
	}
}

func (pc *possibilityCounter) countPossibilities(d design) (result int) {
	if len(d) == 0 {
		return 1
	}

	cacheKey := string(d)
	if possibilities, hit := pc.cache[cacheKey]; hit {
		return possibilities
	}
	defer func() { pc.cache[cacheKey] = result }()

	possibilities := 0
	for _, t := range pc.towels {
		if !designCanStartWithTowel(d, t) {
			continue
		}

		remainder := d[len(t):]
		possibilities += pc.countPossibilities(remainder)
	}

	return possibilities
}

func designCanStartWithTowel(d design, t towel) bool {
	if len(d) < len(t) {
		return false
	}

	for i := range t {
		if d[i] != t[i] {
			return false
		}
	}

	return true
}

func towelsAndDesignsFromReader(r io.Reader) ([]towel, []design, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read lines: %w", err)
	}

	if len(lines) < 3 {
		return nil, nil, errors.New("input must have at least 3 lines")
	}

	towels := towelsFromString(lines[0])

	designs := designsFromStrings(lines[2:])

	return towels, designs, nil
}

func towelsFromString(str string) []towel {
	rawTowels := strings.Split(str, ", ")

	towels := make([]towel, len(rawTowels))
	for i, t := range rawTowels {
		towels[i] = towel(t)
	}

	return towels
}

func designsFromStrings(strs []string) []design {
	designs := make([]design, len(strs))
	for i, str := range strs {
		designs[i] = design(str)
	}

	return designs
}
