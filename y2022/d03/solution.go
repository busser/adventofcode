package d03

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 3 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	rucksacks, err := rucksacksFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for _, r := range rucksacks {
		item, found := outOfPlaceItem(r)
		if !found {
			return errors.New("no out of place item")
		}

		sum += itemPriority(item)
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 3 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	rucksacks, err := rucksacksFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for i := 0; i < len(rucksacks); i += 3 {
		item, found := commonItem(
			rucksacks[i],
			rucksacks[i+1],
			rucksacks[i+2],
		)
		if !found {
			return errors.New("no common item found for group")
		}

		sum += itemPriority(item)
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func rucksacksFromReader(r io.Reader) ([]string, error) {
	rucksacks, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	if len(rucksacks)%3 != 0 {
		return nil, fmt.Errorf("%d rucksacks cannot be divided into groups of 3", len(rucksacks))
	}

	for _, r := range rucksacks {
		if len(r)%2 != 0 {
			return nil, fmt.Errorf("rucksack %q has odd number of items", r)
		}

		for _, item := range r {
			if !((item >= 'a' && item <= 'z') || (item >= 'A' && item <= 'Z')) {
				return nil, fmt.Errorf("unknown item %q in rucksack %q", item, r)
			}
		}
	}

	return rucksacks, nil
}

func outOfPlaceItem(rucksack string) (rune, bool) {
	middle := len(rucksack) / 2

	for _, first := range rucksack[:middle] {
		for _, second := range rucksack[middle:] {
			if first == second {
				return first, true
			}
		}
	}

	return 0, false
}

func itemPriority(item rune) int {
	if item >= 'a' && item <= 'z' {
		return int(item - 'a' + 1)
	}
	if item >= 'A' && item <= 'Z' {
		return int(item - 'A' + 27)
	}

	panic("invalid item")
}

func commonItem(a, b, c string) (rune, bool) {
	for _, item := range a {
		if contains(b, item) && contains(c, item) {
			return item, true
		}
	}

	return 0, false
}

func contains(rucksack string, item rune) bool {
	for _, it := range rucksack {
		if it == item {
			return true
		}
	}

	return false
}
