package d19

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 19 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	workflows, parts, err := workflowsAndPartsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	accepted := func(p part) {
		for _, r := range p.ratings {
			sum += r.start
		}
	}

	if err := triageParts(workflows, parts, accepted); err != nil {
		return fmt.Errorf("could not triage parts: %w", err)
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 19 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	workflows, _, err := workflowsAndPartsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0
	accepted := func(p part) {
		product := 1
		for _, r := range p.ratings {
			product *= r.end - r.start
		}
		count += product
	}

	attrRange := interval{1, 4001}

	parts := []part{{
		ratings: [4]interval{attrRange, attrRange, attrRange, attrRange},
	}}

	if err := triageParts(workflows, parts, accepted); err != nil {
		return fmt.Errorf("could not triage parts: %w", err)
	}

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type workflow struct {
	id    string
	rules []rule
}

type rule struct {
	alwaysApplies bool

	attribute int
	operator  byte
	value     int

	next string
}

const (
	lessThan    = '<'
	greaterThan = '>'
)

type part struct {
	ratings         [4]interval
	currentWorkflow string
	currentRule     int
}

const (
	// The order of the attributes is important.
	coolLooking = iota
	musical
	aeordynamic
	shiny
)

type interval struct {
	start int // inclusive
	end   int // exclusive
}

func (r rule) process(p part) []part {
	if r.alwaysApplies {
		return []part{{p.ratings, r.next, 0}}
	}

	switch r.operator {
	case lessThan:
		switch {
		case p.ratings[r.attribute].end <= r.value:
			// The entire interval matches.
			p.currentWorkflow = r.next
			p.currentRule = 0
			return []part{p}
		case p.ratings[r.attribute].start >= r.value:
			// The entire interval does not match.
			p.currentRule++
			return []part{p}
		default:
			// The interval is partially matching.
			splitParts := []part{
				{p.ratings, r.next, 0},
				{p.ratings, p.currentWorkflow, p.currentRule + 1},
			}
			splitParts[0].ratings[r.attribute].end = r.value
			splitParts[1].ratings[r.attribute].start = r.value
			return splitParts
		}

	case greaterThan:
		switch {
		case p.ratings[r.attribute].start > r.value:
			// The entire interval matches.
			p.currentWorkflow = r.next
			p.currentRule = 0
			return []part{p}
		case p.ratings[r.attribute].end <= r.value+1:
			// The entire interval does not match.
			p.currentRule++
			return []part{p}
		default:
			// The interval is partially matching.
			splitParts := []part{
				{p.ratings, p.currentWorkflow, p.currentRule + 1},
				{p.ratings, r.next, 0},
			}
			splitParts[0].ratings[r.attribute].end = r.value + 1
			splitParts[1].ratings[r.attribute].start = r.value + 1
			return splitParts
		}

	default:
		panic("invalid operator")
	}
}

func triageParts(workflows []workflow, parts []part, accepted func(part)) error {
	workflowsByID := make(map[string]workflow, len(workflows))
	for _, w := range workflows {
		workflowsByID[w.id] = w
	}

	for i := range parts {
		parts[i].currentWorkflow = "in"
	}

	for len(parts) > 0 {
		p := parts[len(parts)-1]
		parts = parts[:len(parts)-1]

		if p.currentWorkflow == "A" {
			accepted(p)
			continue
		}

		if p.currentWorkflow == "R" {
			continue
		}

		w, ok := workflowsByID[p.currentWorkflow]
		if !ok {
			return fmt.Errorf("unknown workflow %q", p.currentWorkflow)
		}

		if p.currentRule >= len(w.rules) {
			return fmt.Errorf("invalid rule index %d for workflow %q", p.currentRule, p.currentWorkflow)
		}

		r := w.rules[p.currentRule]
		parts = append(parts, r.process(p)...)
	}

	return nil
}

func ruleFromString(s string) (rule, error) {
	parts := strings.SplitN(s, ":", 2)

	if len(parts) == 1 {
		return rule{
			alwaysApplies: true,
			next:          parts[0],
		}, nil
	}

	var attribute int
	switch parts[0][0] {
	case 'x':
		attribute = coolLooking
	case 'm':
		attribute = musical
	case 'a':
		attribute = aeordynamic
	case 's':
		attribute = shiny
	default:
		return rule{}, fmt.Errorf("invalid attribute: %q", parts[0][0])
	}

	operator := parts[0][1]
	if operator != lessThan && operator != greaterThan {
		return rule{}, fmt.Errorf("invalid operator: %q", operator)
	}

	value, err := strconv.Atoi(parts[0][2:])
	if err != nil {
		return rule{}, fmt.Errorf("invalid value: %q", parts[0][2:])
	}

	return rule{
		attribute: attribute,
		operator:  operator,
		value:     value,
		next:      parts[1],
	}, nil
}

func workflowFromString(s string) (workflow, error) {
	parts := strings.SplitN(s, "{", 2)

	id := parts[0]

	rawRules := strings.Split(strings.TrimSuffix(parts[1], "}"), ",")

	rules := make([]rule, len(rawRules))
	for i := range rawRules {
		r, err := ruleFromString(rawRules[i])
		if err != nil {
			return workflow{}, fmt.Errorf("invalid rule: %w", err)
		}

		rules[i] = r
	}

	return workflow{id, rules}, nil
}

func partFromString(s string) (part, error) {
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")

	rawRatings := strings.Split(s, ",")
	if len(rawRatings) != 4 {
		return part{}, fmt.Errorf("%d attributes instead of 4", len(rawRatings))
	}

	var p part
	for i, rawRating := range rawRatings {
		if len(rawRating) < 3 {
			return part{}, fmt.Errorf("rating too short: %q", rawRating)
		}

		v, err := strconv.Atoi(rawRating[2:])
		if err != nil {
			return part{}, fmt.Errorf("rating is not a number: %q", rawRating[2:])
		}

		p.ratings[i] = interval{start: v, end: v + 1}
	}

	return p, nil
}

func workflowsAndPartsFromReader(r io.Reader) ([]workflow, []part, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read input: %w", err)
	}

	var workflows []workflow

	for i := 0; ; i++ {
		if i >= len(lines)-1 {
			return nil, nil, fmt.Errorf("didn't find empty line")
		}

		if lines[i] == "" {
			lines = lines[i+1:]
			break
		}

		w, err := workflowFromString(lines[i])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid workflow: %w", err)
		}

		workflows = append(workflows, w)
	}

	var parts []part

	for i := range lines {
		p, err := partFromString(lines[i])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid part: %w", err)
		}

		parts = append(parts, p)
	}

	return workflows, parts, nil
}
