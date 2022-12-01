package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 2 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	entries, err := databaseEntriesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := validEntriesCount(entries, oldRules)

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 2 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	entries, err := databaseEntriesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := validEntriesCount(entries, newRules)

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type databaseEntry struct {
	password []rune
	policy   corporatePolicy
}

type corporatePolicy struct {
	letter            rune
	leftNum, rightNum int
}

func validEntriesCount(entries []databaseEntry, isValid func(dbe databaseEntry) bool) int {
	count := 0
	for _, entry := range entries {
		if isValid(entry) {
			count++
		}
	}
	return count
}

func oldRules(dbe databaseEntry) bool {
	count := 0
	for _, letter := range dbe.password {
		if letter == dbe.policy.letter {
			count++
		}
	}

	if count >= dbe.policy.leftNum && count <= dbe.policy.rightNum {
		return true
	}
	return false
}

func newRules(dbe databaseEntry) bool {
	leftMatch := dbe.password[dbe.policy.leftNum-1] == dbe.policy.letter
	rightMatch := dbe.password[dbe.policy.rightNum-1] == dbe.policy.letter

	return leftMatch != rightMatch
}

func databaseEntriesFromReader(r io.Reader) ([]databaseEntry, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading lines: %w", err)
	}

	entries := make([]databaseEntry, len(lines))
	for i := range lines {
		if err := entries[i].fromString(lines[i]); err != nil {
			return nil, fmt.Errorf("parsing database entry %q: %w", lines[i], err)
		}
	}

	return entries, nil
}

func (dbe *databaseEntry) fromString(s string) error {
	parts := strings.Split(s, ": ")
	if len(parts) != 2 {
		return errors.New("invalid format")
	}

	if err := dbe.policy.fromString(parts[0]); err != nil {
		return errors.New("parsing policy")
	}

	dbe.password = []rune(parts[1])

	return nil
}

func (cp *corporatePolicy) fromString(s string) error {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return errors.New("invalid format")
	}

	bounds, letters := parts[0], []rune(parts[1])

	if len(letters) != 1 {
		return errors.New("invalid format")
	}

	parts = strings.Split(bounds, "-")
	if len(parts) != 2 {
		return errors.New("invalid format")
	}

	leftNum, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}
	rightNum, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}

	cp.letter = letters[0]
	cp.leftNum = leftNum
	cp.rightNum = rightNum

	return nil
}
