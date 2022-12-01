package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 7 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	rules, err := rulesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := containsBagCount(rules, "shiny gold")

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 7 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	rules, err := rulesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := bagsInsideCount(rules, "shiny gold")

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func containsBagCount(rules ruleGraph, bag string) int {
	containsBagCache := make(map[string]bool)

	var containsBag func(string, string) bool
	containsBag = func(container, containee string) bool {
		if contains, ok := containsBagCache[container]; ok {
			return contains
		}

		for _, cr := range rules[container] {
			if cr.bag == containee {
				containsBagCache[container] = true
				return true
			}

			if containsBag(cr.bag, containee) {
				containsBagCache[container] = true
				return true
			}
		}

		containsBagCache[container] = false
		return false
	}

	count := 0
	for container := range rules {
		if containsBag(container, bag) {
			count++
		}
	}

	return count
}

func bagsInsideCount(rules ruleGraph, bag string) int {
	bagsInsideCache := make(map[string]int)

	var bagsInside func(string) int
	bagsInside = func(container string) int {
		if inside, ok := bagsInsideCache[container]; ok {
			return inside
		}

		count := 0

		for _, bc := range rules[container] {
			count += bc.count * (1 + bagsInside(bc.bag))
		}

		bagsInsideCache[container] = count
		return count
	}

	return bagsInside(bag)
}

type ruleGraph map[string][]containRule

type containRule struct {
	bag   string
	count int
}

func rulesFromReader(r io.Reader) (ruleGraph, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading lines: %w", err)
	}

	rules := ruleGraph{}

	for _, l := range lines {
		splitLine := strings.Split(l, " bags contain ")
		if len(splitLine) != 2 {
			return nil, errors.New("wrong format")
		}

		bag, contains := splitLine[0], splitLine[1]

		if contains == "no other bags." {
			continue
		}

		splitContains := strings.Split(contains, ", ")

		rule := make([]containRule, len(splitContains))

		for i, sc := range splitContains {
			fields := strings.Fields(sc)
			if len(fields) != 4 {
				return nil, errors.New("wrong format")
			}

			count, err := strconv.Atoi(fields[0])
			if err != nil {
				return nil, fmt.Errorf("wrong format: not an int: %q", fields[0])
			}

			rule[i] = containRule{
				bag:   strings.Join(fields[1:3], " "),
				count: count,
			}
		}

		rules[bag] = rule
	}

	return rules, nil
}
