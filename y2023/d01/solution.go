package d01

import (
	"fmt"
	"io"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 1 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for _, l := range lines {
		sum += parseValueDigitsOnly(l)
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for _, l := range lines {
		sum += parseValueWithWords(l)
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func parseValueDigitsOnly(s string) int {
	firstDigit, lastDigit := -1, -1

	for _, r := range s {
		if r >= '0' && r <= '9' {
			if firstDigit == -1 {
				firstDigit = int(r - '0')
			}
			lastDigit = int(r - '0')
		}
	}
	return firstDigit*10 + lastDigit
}

var wordsToDigits = map[string]int{
	"zero":  0,
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func parseValueWithWords(s string) int {
	firstDigit, lastDigit := -1, -1

	for i, r := range s {
		if r >= '0' && r <= '9' {
			if firstDigit == -1 {
				firstDigit = int(r - '0')
			}
			lastDigit = int(r - '0')
			continue
		}

		for word, digit := range wordsToDigits {
			if strings.HasPrefix(s[i:], word) {
				if firstDigit == -1 {
					firstDigit = digit
				}
				lastDigit = digit
			}
		}
	}

	return firstDigit*10 + lastDigit
}
