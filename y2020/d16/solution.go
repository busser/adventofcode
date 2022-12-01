package busser

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 16 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	fields, _, nearbyTickets, err := fieldsAndTicketsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0

	for _, t := range nearbyTickets {
		for _, v := range t {
			if !valueIsValidForAnyField(v, fields) {
				sum += v
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 16 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	fields, yourTicket, nearbyTickets, err := fieldsAndTicketsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var validTickets []ticket
	for _, t := range nearbyTickets {
		if ticketIsValid(t, fields) {
			validTickets = append(validTickets, t)
		}
	}

	mapping, err := fieldMapping(fields, validTickets)
	if err != nil {
		return fmt.Errorf("could not order fields: %w", err)
	}

	product := 1
	for i, j := range mapping {
		if strings.HasPrefix(fields[j].name, "departure") {
			product *= yourTicket[i]
		}
	}

	_, err = fmt.Fprintf(answer, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type field struct {
	name        string
	validValues []interval
}

type interval struct {
	min, max int
}

type ticket []int

func fieldMapping(fields []field, tickets []ticket) ([]int, error) {
	possibleFields := make([][]int, len(fields))
	for position := range fields {
		for fi, fv := range fields {
			if allTicketsMatchFieldForPosition(tickets, fv, position) {
				possibleFields[position] = append(possibleFields[position], fi)
			}
		}
	}

	mapping := make([]int, len(fields))
	alreadyMapped := make([]bool, len(fields))

	var helper func(int) bool
	helper = func(position int) bool {
		if position == len(mapping) {
			return true
		}

		for _, f := range possibleFields[position] {
			if alreadyMapped[f] {
				continue
			}

			mapping[position] = f
			alreadyMapped[f] = true

			if helper(position + 1) {
				return true
			}

			mapping[position] = -1
			alreadyMapped[f] = false
		}

		return false
	}

	if helper(0) {
		return mapping, nil
	}

	return nil, errors.New("no valid mapping found")
}

func allTicketsMatchFieldForPosition(tickets []ticket, f field, position int) bool {
	for _, t := range tickets {
		if len(t) <= position || !valueIsValidForField(t[position], f) {
			return false
		}
	}
	return true
}

func ticketIsValid(t ticket, fields []field) bool {
	for _, v := range t {
		if !valueIsValidForAnyField(v, fields) {
			return false
		}
	}
	return true
}

func valueIsValidForAnyField(v int, fields []field) bool {
	for _, f := range fields {
		if valueIsValidForField(v, f) {
			return true
		}
	}
	return false
}

func valueIsValidForField(v int, f field) bool {
	for _, i := range f.validValues {
		if v >= i.min && v <= i.max {
			return true
		}
	}
	return false
}

func fieldsAndTicketsFromReader(r io.Reader) ([]field, ticket, []ticket, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reading lines: %w", err)
	}

	chunks := splitSlice(lines, "")
	if len(chunks) != 3 {
		return nil, nil, nil, fmt.Errorf("wrong format")
	}

	fields, err := fieldsFromLines(chunks[0])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing fields: %w", err)
	}

	yourTicket, err := yourTicketFromLines(chunks[1])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing your ticket: %w", err)
	}

	nearbyTickets, err := nearbyTicketsFromLines(chunks[2])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing nearby tickets: %w", err)
	}

	return fields, yourTicket, nearbyTickets, nil
}

func fieldsFromLines(lines []string) ([]field, error) {
	fields := make([]field, len(lines))

	for i, line := range lines {
		f, err := fieldFromString(line)
		if err != nil {
			return nil, fmt.Errorf("parsing field #%d: %w", i, err)
		}

		fields[i] = f
	}

	return fields, nil
}

func fieldFromString(s string) (field, error) {
	splitField := strings.Split(s, ": ")
	if len(splitField) != 2 {
		return field{}, errors.New("wrong format")
	}

	validValues, err := intervalsFromString(splitField[1])
	if err != nil {
		return field{}, fmt.Errorf("parsing valid values: %w", err)
	}

	f := field{
		name:        splitField[0],
		validValues: validValues,
	}

	return f, nil
}

func intervalsFromString(s string) ([]interval, error) {
	rawIntervals := strings.Split(s, " or ")

	intervals := make([]interval, len(rawIntervals))

	for i, rawInterval := range rawIntervals {
		numbers, err := helpers.IntsFromString(rawInterval, "-")
		if err != nil {
			return nil, fmt.Errorf("reading numbers: %w", err)
		}
		if len(numbers) != 2 {
			return nil, errors.New("wrong format")
		}

		intervals[i] = interval{numbers[0], numbers[1]}
	}

	return intervals, nil
}

func yourTicketFromLines(lines []string) (ticket, error) {
	if len(lines) != 2 {
		return nil, errors.New("wrong format")
	}

	if lines[0] != "your ticket:" {
		return nil, errors.New("wrong format")
	}

	yourTicket, err := ticketFromString(lines[1])
	if err != nil {
		return nil, fmt.Errorf("parsing ticket: %w", err)
	}

	return yourTicket, nil
}

func nearbyTicketsFromLines(lines []string) ([]ticket, error) {
	if len(lines) < 1 {
		return nil, errors.New("wrong format")
	}

	if lines[0] != "nearby tickets:" {
		return nil, errors.New("wrong format")
	}

	tickets := make([]ticket, len(lines)-1)
	for i, line := range lines[1:] {
		t, err := ticketFromString(line)
		if err != nil {
			return nil, fmt.Errorf("parsing ticket #%d: %w", i, err)
		}
		tickets[i] = t
	}

	return tickets, nil
}

func ticketFromString(s string) (ticket, error) {
	numbers, err := helpers.IntsFromString(s, ",")
	if err != nil {
		return nil, fmt.Errorf("parsing numbers: %w", err)
	}

	return ticket(numbers), nil
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
