package d05

import (
	"bytes"
	"fmt"
	"io"
	"slices"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 5 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	rules, manuals, err := rulesAndPagesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	ruleSet := buildRuleSet(rules)

	total := 0
	for _, manual := range manuals {
		if manual.respectsRules(ruleSet) {
			total += manual.middlePage()
		}
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	rules, manuals, err := rulesAndPagesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	ruleSet := buildRuleSet(rules)

	total := 0
	for _, manual := range manuals {
		if !manual.respectsRules(ruleSet) {
			ruleSet.sort(manual.pages)
			total += manual.middlePage()
		}
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type orderingRule struct {
	before, after int
}

type safetyManual struct {
	pages []int
}

func (sm safetyManual) respectsRules(rules ruleSet) bool {
	for i, p1 := range sm.pages {
		for _, p2 := range sm.pages[i+1:] {
			if !rules.orderIsOK(p1, p2) {
				return false
			}
		}
	}
	return true
}

func (sm safetyManual) middlePage() int {
	return sm.pages[len(sm.pages)/2]
}

type ruleSet struct {
	aftersByBefore map[int][]int
}

func buildRuleSet(rules []orderingRule) ruleSet {
	aftersByBefore := make(map[int][]int)
	for _, rule := range rules {
		aftersByBefore[rule.before] = append(aftersByBefore[rule.before], rule.after)
	}
	return ruleSet{aftersByBefore}
}

func (rs ruleSet) orderIsOK(p1, p2 int) bool {
	return !slices.Contains(rs.aftersByBefore[p2], p1)
}

func (rs ruleSet) sort(pages []int) {
	for range len(pages) - 1 {
		for i := 0; i < len(pages)-1; i++ {
			if !rs.orderIsOK(pages[i], pages[i+1]) {
				pages[i], pages[i+1] = pages[i+1], pages[i]
			}
		}
	}
}

func pairRespectsRules(p1, p2 int, rules []orderingRule) bool {
	for _, rule := range rules {
		if rule.before == p2 && rule.after == p1 {
			return false
		}
	}
	return true
}

func orderingRuleFromString(s string) (orderingRule, error) {
	ints := helpers.IntsFromString(s)
	if len(ints) != 2 {
		return orderingRule{}, fmt.Errorf("invalid format")
	}
	return orderingRule{
		before: ints[0],
		after:  ints[1],
	}, nil
}

func safetyManualFromString(s string) (safetyManual, error) {
	pages := helpers.IntsFromString(s)
	if len(pages)%2 != 1 {
		return safetyManual{}, fmt.Errorf("manuals must have an odd number of pages")
	}

	return safetyManual{
		pages: pages,
	}, nil
}

func rulesAndPagesFromReader(r io.Reader) ([]orderingRule, []safetyManual, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read input: %w", err)
	}

	input = bytes.TrimSpace(input)

	chunks := bytes.Split(input, []byte("\n\n"))
	if len(chunks) != 2 {
		return nil, nil, fmt.Errorf("invalid input")
	}

	lines := bytes.Split(chunks[0], []byte("\n"))
	rules := make([]orderingRule, len(lines))
	for i, line := range lines {
		rule, err := orderingRuleFromString(string(line))
		if err != nil {
			return nil, nil, fmt.Errorf("invalid rule %d: %w", i, err)
		}
		rules[i] = rule
	}

	lines = bytes.Split(chunks[1], []byte("\n"))
	pages := make([]safetyManual, len(lines))
	for i, line := range lines {
		manual, err := safetyManualFromString(string(line))
		if err != nil {
			return nil, nil, fmt.Errorf("invalid manual %d: %w", i, err)
		}
		pages[i] = manual
	}

	return rules, pages, nil
}
