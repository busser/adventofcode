package d11

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 11 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	stones, err := stonesFromInput(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	total := countStonesAfterBlinks(stones, 25)

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 11 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	stones, err := stonesFromInput(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	total := countStonesAfterBlinks(stones, 75)

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func countStonesAfterBlinks(stones []int, blinks int) int {
	type cacheKey struct {
		stone  int
		blinks int
	}

	cache := make(map[cacheKey]int)

	var count func(int, int) int
	count = func(stone, blinks int) int {
		if blinks == 0 {
			return 1
		}

		key := cacheKey{stone, blinks}
		if cached, hit := cache[key]; hit {
			return cached
		}

		total := 0
		switch {
		case stone == 0:
			total = count(1, blinks-1)
		case hasEvenNumberOfDigits(stone):
			left, right := splitIntInHalf(stone)
			total = count(left, blinks-1) + count(right, blinks-1)
		default:
			total = count(stone*2024, blinks-1)
		}

		cache[key] = total
		return total
	}

	total := 0
	for _, stone := range stones {
		total += count(stone, blinks)
	}

	return total
}

func numberOfDigits(n int) int {
	numberOfDigits := 0
	for n > 0 {
		numberOfDigits++
		n /= 10
	}
	return numberOfDigits
}

func hasEvenNumberOfDigits(n int) bool {
	return numberOfDigits(n)%2 == 0
}

func splitIntInHalf(n int) (int, int) {
	split := 1
	for range numberOfDigits(n) / 2 {
		split *= 10
	}
	return n / split, n % split
}

func stonesFromInput(r io.Reader) ([]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("expected 1 line, got %d", len(lines))
	}

	return helpers.IntsFromString(lines[0]), nil
}
