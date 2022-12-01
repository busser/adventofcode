package busser

import (
	"fmt"
	"io"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 12 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	caves, err := caveSystemFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := countPossiblePaths(caves)

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 12 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	caves, err := caveSystemFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := countPossiblePathsWithDoubleVisit(caves)

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type caveSystem map[string][]string

func countPossiblePaths(caves caveSystem) int {
	timesVisited := make(map[string]int)

	var countPathsFrom func(cave string) int
	countPathsFrom = func(cave string) int {
		if cave == "end" {
			return 1
		}

		if isSmall(cave) && timesVisited[cave] >= 1 {
			return 0
		}

		timesVisited[cave]++
		defer func() { timesVisited[cave]-- }()

		totalPaths := 0
		for _, neighbor := range caves[cave] {
			totalPaths += countPathsFrom(neighbor)
		}

		return totalPaths
	}

	return countPathsFrom("start")
}

func countPossiblePathsWithDoubleVisit(caves caveSystem) int {
	timesVisited := make(map[string]int)
	var smallCaveVisitedTwice bool

	var countPathsFrom func(string) int
	countPathsFrom = func(cave string) int {
		if cave == "end" {
			return 1
		}

		if isSmall(cave) && (timesVisited[cave] >= 2 || (smallCaveVisitedTwice && timesVisited[cave] >= 1)) {
			return 0
		}

		timesVisited[cave]++
		defer func() { timesVisited[cave]-- }()

		if isSmall(cave) && timesVisited[cave] >= 2 {
			smallCaveVisitedTwice = true
			defer func() { smallCaveVisitedTwice = false }()
		}

		totalPaths := 0
		for _, neighbor := range caves[cave] {
			if neighbor == "start" {
				continue
			}
			totalPaths += countPathsFrom(neighbor)
		}

		return totalPaths
	}

	return countPathsFrom("start")
}

func isSmall(cave string) bool {
	return cave[0] >= 'a' && cave[0] <= 'z'
}

func caveSystemFromReader(r io.Reader) (caveSystem, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	system := make(caveSystem)
	for _, l := range lines {
		parts := strings.Split(l, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid path %q", l)
		}

		caveA, caveB := parts[0], parts[1]
		system[caveA] = append(system[caveA], caveB)
		system[caveB] = append(system[caveB], caveA)
	}

	return system, nil
}
