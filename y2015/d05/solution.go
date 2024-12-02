package d05

import (
	"fmt"
	"io"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 5 of Advent of Code 2015.
func PartOne(r io.Reader, w io.Writer) error {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0
	for _, l := range lines {
		if stringIsNiceOldRules(l) {
			count++
		}
	}

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2015.
func PartTwo(r io.Reader, w io.Writer) error {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0
	for _, l := range lines {
		if stringIsNiceNewRules(l) {
			count++
		}
	}

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func stringIsNiceOldRules(s string) bool {
	return stringHasAtLeastThreeVowels(s) &&
		stringHasDoubleLetter(s) &&
		stringDoesNotContainBadStrings(s)
}

func stringHasAtLeastThreeVowels(s string) bool {
	vowels := 0
	for _, c := range s {
		if c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' {
			vowels++
			if vowels >= 3 {
				return true
			}
		}
	}
	return false
}

func stringHasDoubleLetter(s string) bool {
	for i := 0; i < len(s)-1; i++ {
		if s[i] == s[i+1] {
			return true
		}
	}
	return false
}

func stringDoesNotContainBadStrings(s string) bool {
	return !strings.Contains(s, "ab") &&
		!strings.Contains(s, "cd") &&
		!strings.Contains(s, "pq") &&
		!strings.Contains(s, "xy")
}

func stringIsNiceNewRules(s string) bool {
	return stringHasDoublePair(s) &&
		stringHasRepeatingLetterWithOneBetween(s)
}

func stringHasDoublePair(s string) bool {
	for i := 0; i < len(s)-3; i++ {
		if strings.Contains(s[i+2:], s[i:i+2]) {
			return true
		}
	}
	return false
}

func stringHasRepeatingLetterWithOneBetween(s string) bool {
	for i := 0; i < len(s)-2; i++ {
		if s[i] == s[i+2] {
			return true
		}
	}
	return false
}
