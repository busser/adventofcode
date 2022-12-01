package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 19 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	rules, messages, err := rulesAndMessagesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0
	for _, msg := range messages {
		if lengths := rules.ruleMatches(0, msg); contains(lengths, len(msg)) {
			count++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 19 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	rules, messages, err := rulesAndMessagesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	rules[8] = rule{
		kind: ruleKindComplex,
		subRules: [][]int{
			{42},
			{42, 8},
		},
	}
	rules[11] = rule{
		kind: ruleKindComplex,
		subRules: [][]int{
			{42, 31},
			{42, 11, 31},
		},
	}

	count := 0
	for _, msg := range messages {
		if lengths := rules.ruleMatches(0, msg); contains(lengths, len(msg)) {
			count++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

/*
	Solution model:

	Rules are stored in a RuleSet, which is a slice of Rules where Rule 0 is at
	index 0, Rule 1 is at index 1, etc.

	There are 2 kinds of Rules:
	- SimpleRules
	- ComplexRules

	A SimpleRule matches a single character.

	A ComplexRule is a list of SubRules.

	A SubRule is a list of references to a Rule.

	Let's take this input as an example:

	0: 1 2
	1: "a"
	2: 1 3 | 3 1
	3: "b"

	This is modeled this way in memory:

	RuleSet{
		Rule{
			kind: Complex,
			SubRules: []SubRule{
				{1, 2},
			},
		},
		Rule{
			kind: Simple,
			character: 'a',
		},
		Rule{
			kind: Complex,
			SubRules: []SubRule{
				{1, 3},
				{3, 1},
			},
		},
		Rule{
			kind: Simple,
			character: 'b',
		},
	}
*/

type ruleSet []rule

type rule struct {
	kind      ruleKind
	character rune    // For simple rules
	subRules  [][]int // For complex rules
}

type ruleKind uint8

const (
	ruleKindSimple ruleKind = iota
	ruleKindComplex
)

func (set ruleSet) ruleMatches(ruleRef int, msg []rune) []int {
	rule := set[ruleRef]

	switch rule.kind {
	case ruleKindSimple:
		if len(msg) == 0 {
			return nil
		}
		if msg[0] == rule.character {
			return []int{1}
		}
		return nil

	case ruleKindComplex:
		var lengths []int
		for _, subRule := range rule.subRules {
			lengths = append(lengths, set.subRuleMatches(subRule, msg)...)
		}
		return lengths

	default:
		panic("unknown rule kind")
	}
}

func (set ruleSet) subRuleMatches(ruleRefs []int, msg []rune) []int {
	if len(ruleRefs) == 0 {
		return nil
	}

	if len(ruleRefs) == 1 {
		return set.ruleMatches(ruleRefs[0], msg)
	}

	var lengths []int
	for _, firstRuleLength := range set.ruleMatches(ruleRefs[0], msg) {
		for _, l := range set.subRuleMatches(ruleRefs[1:], msg[firstRuleLength:]) {
			lengths = append(lengths, l+firstRuleLength)
		}
	}
	return lengths
}

func rulesAndMessagesFromReader(r io.Reader) (ruleSet, [][]rune, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, nil, fmt.Errorf("reading lines: %w", err)
	}

	chunks := splitSlice(lines, "")
	if len(chunks) != 2 {
		return nil, nil, fmt.Errorf("wrong format")
	}

	rules, err := rulesFromLines(chunks[0])
	if err != nil {
		return nil, nil, fmt.Errorf("parsing rules: %w", err)
	}

	messages := messagesFromLines(chunks[1])

	return rules, messages, nil
}

func rulesFromLines(lines []string) (ruleSet, error) {
	rules := make(ruleSet, len(lines))

	for i, line := range lines {
		id, r, err := ruleFromString(line)
		if err != nil {
			return nil, fmt.Errorf("parsingrule #%d: %w", i, err)
		}

		rules[id] = r
	}

	return rules, nil
}

func ruleFromString(s string) (int, rule, error) {
	pieces := strings.Split(s, ": ")
	if len(pieces) != 2 {
		return 0, rule{}, errors.New("wrong format")
	}

	id, err := strconv.Atoi(pieces[0])
	if err != nil {
		return 0, rule{}, fmt.Errorf("invalid ID %q: %w", pieces[0], err)
	}

	if len(pieces[1]) == 3 && pieces[1][0] == '"' && pieces[1][2] == '"' {
		r := rule{
			kind:      ruleKindSimple,
			character: rune(pieces[1][1]),
		}
		return id, r, nil
	}

	rawSubRules := strings.Split(pieces[1], " | ")
	if len(pieces) == 0 {
		return 0, rule{}, errors.New("wrong format")
	}

	subRules := make([][]int, len(rawSubRules))
	for i, rawSubRule := range rawSubRules {
		subRule, err := helpers.IntsFromString(rawSubRule, " ")
		if err != nil {
			return 0, rule{}, fmt.Errorf("parsing sub-rule: %w", err)
		}
		subRules[i] = subRule
	}

	r := rule{
		kind:     ruleKindComplex,
		subRules: subRules,
	}

	return id, r, nil
}

func messagesFromLines(lines []string) [][]rune {
	messages := make([][]rune, len(lines))

	for i, line := range lines {
		messages[i] = []rune(line)
	}

	return messages
}

func splitSlice(slice []string, sep string) [][]string {
	var split [][]string

	start := 0

	for end := range slice {
		if slice[end] == sep {
			split = append(split, slice[start:end])
			start = end + 1
		}
	}
	split = append(split, slice[start:])

	return split
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
