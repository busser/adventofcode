package d11

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 11 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	monkeys, err := monkeysFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	business := monkeyBusiness(monkeys, 20, true)

	_, err = fmt.Fprintf(w, "%d", business)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 11 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	monkeys, err := monkeysFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	business := monkeyBusiness(monkeys, 10_000, false)

	_, err = fmt.Fprintf(w, "%d", business)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type monkey struct {
	items          []int
	operation      func(int) int
	divisor        int
	throwToIfTrue  int
	throwToIfFalse int
}

func monkeyBusiness(monkeys []monkey, rounds int, relief bool) int {
	productOfDivisors := 1
	for i := range monkeys {
		productOfDivisors *= monkeys[i].divisor
	}

	business := make([]int, len(monkeys))

	for r := 0; r < rounds; r++ {
		for i := range monkeys {
			business[i] += len(monkeys[i].items)

			for _, item := range monkeys[i].items {
				item = monkeys[i].operation(item)
				if relief {
					item /= 3
				}

				// To avoid the level of worry for a given item overflowing, all
				// the while keeping all divisability results the same, we only
				// keep the modulus of the worry level and of productOfDivisors.
				item %= productOfDivisors

				next := monkeys[i].throwToIfFalse
				if item%monkeys[i].divisor == 0 {
					next = monkeys[i].throwToIfTrue
				}

				monkeys[next].items = append(monkeys[next].items, item)
			}

			// We empty the slice but keep the backing array for future use.
			// This reduces overall memory usage.
			monkeys[i].items = monkeys[i].items[:0]
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(business)))

	return business[0] * business[1]
}

func monkeysFromReader(r io.Reader) ([]monkey, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	if len(lines)%7 != 6 {
		return nil, fmt.Errorf("unexpected number of lines")
	}

	numMonkeys := (len(lines) + 1) / 7
	if numMonkeys < 2 {
		return nil, errors.New("not enough monkey to compute monkey business")
	}

	var monkeys []monkey

	for i := 0; i < len(lines); i += 7 {
		m, err := monkeyFromLines(lines[i : i+6])
		if err != nil {
			return nil, err
		}

		if m.throwToIfFalse < 0 || m.throwToIfFalse >= numMonkeys || m.throwToIfFalse == len(monkeys) {
			return nil, fmt.Errorf("cannot throw to monkey %d", m.throwToIfFalse)
		}
		if m.throwToIfTrue < 0 || m.throwToIfTrue >= numMonkeys || m.throwToIfTrue == len(monkeys) {
			return nil, fmt.Errorf("cannot throw to monkey %d", m.throwToIfTrue)
		}

		monkeys = append(monkeys, m)
	}

	return monkeys, nil
}

func monkeyFromLines(lines []string) (monkey, error) {
	items, err := itemsFromLine(lines[1])
	if err != nil {
		return monkey{}, err
	}

	operation, err := operationFromLine(lines[2])
	if err != nil {
		return monkey{}, err
	}

	divisor, err := divisorFromLine(lines[3])
	if err != nil {
		return monkey{}, err
	}

	throwToIfTrue, err := throwToIfTrueFromLine(lines[4])
	if err != nil {
		return monkey{}, err
	}

	throwToIfFalse, err := throwToIfFalseFromLine(lines[5])
	if err != nil {
		return monkey{}, err
	}

	return monkey{
		items:          items,
		operation:      operation,
		divisor:        divisor,
		throwToIfTrue:  throwToIfTrue,
		throwToIfFalse: throwToIfFalse,
	}, nil
}

func itemsFromLine(s string) ([]int, error) {
	var items []int

	rawItems := strings.TrimPrefix(s, "  Starting items: ")
	for _, raw := range strings.Split(rawItems, ", ") {
		item, err := strconv.Atoi(string(raw))
		if err != nil {
			return nil, fmt.Errorf("%q is not a number", string(raw))
		}

		items = append(items, item)
	}

	return items, nil
}

func operationFromLine(s string) (func(int) int, error) {
	rawOperation := strings.TrimPrefix(s, "  Operation: new = old ")
	parts := strings.SplitN(rawOperation, " ", 2)
	if len(parts) != 2 {
		return nil, errors.New("wrong format")
	}

	operator := parts[0]
	rawArgument := parts[1]

	if rawArgument == "old" {
		switch operator {
		case "+":
			return func(old int) int { return old + old }, nil
		case "*":
			return func(old int) int { return old * old }, nil
		default:
			return nil, fmt.Errorf("unknown operator %q", operator)
		}
	}

	arg, err := strconv.Atoi(rawArgument)
	if err != nil {
		return nil, fmt.Errorf("%q is not a number", rawArgument)
	}

	switch operator {
	case "+":
		return func(old int) int { return old + arg }, nil
	case "*":
		return func(old int) int { return old * arg }, nil
	default:
		return nil, fmt.Errorf("unknown operator %q", operator)
	}
}

func divisorFromLine(s string) (int, error) {
	raw := strings.TrimPrefix(s, "  Test: divisible by ")

	divisor, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%q is not a number", raw)
	}

	return divisor, nil
}

func throwToIfTrueFromLine(s string) (int, error) {
	raw := strings.TrimPrefix(s, "    If true: throw to monkey ")

	throwTo, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%q is not a number", raw)
	}

	return throwTo, nil
}

func throwToIfFalseFromLine(s string) (int, error) {
	raw := strings.TrimPrefix(s, "    If false: throw to monkey ")

	throwTo, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%q is not a number", raw)
	}

	return throwTo, nil
}
