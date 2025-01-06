package d22

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 22 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	numbers, err := initialSecretNumbersFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	total := 0
	for _, n := range numbers {
		for range 2000 {
			n = nextNumber(n)
		}
		total += n
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 22 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	numbers, err := initialSecretNumbersFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	maxBananas := findMaxBananas(numbers)

	_, err = fmt.Fprintf(w, "%d", maxBananas)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func nextNumber(n int) int {
	n = (n ^ (n * 64)) % 16777216
	n = (n ^ (n / 32)) % 16777216
	n = (n ^ (n * 2048)) % 16777216
	return n
}

func findMaxBananas(initialNumbers []int) int {
	totalBananasBySequence := make(map[[4]int8]int)

	for _, n := range initialNumbers {
		bananasBySequence := make(map[[4]int8]int)

		var seq [4]int8
		for i := range 2000 {
			next := nextNumber(n)
			priceChange := (next % 10) - (n % 10)
			n = next

			seq = [4]int8{seq[1], seq[2], seq[3], int8(priceChange)}
			if i < 4 {
				continue
			}

			if _, seen := bananasBySequence[seq]; !seen {
				bananasBySequence[seq] = (n % 10)
			}
		}

		for seq, bananas := range bananasBySequence {
			totalBananasBySequence[seq] += bananas
		}
	}

	var maxBananas int
	for _, bananas := range totalBananasBySequence {
		maxBananas = max(maxBananas, bananas)
	}

	return maxBananas
}

func initialSecretNumbersFromReader(r io.Reader) ([]int, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	return helpers.IntsFromString(string(input)), nil
}
