package d21

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 21 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	monkeys, err := monkeysFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	v := rootValue(monkeys)

	_, err = fmt.Fprintf(w, "%d", v)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 21 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	monkeys, err := monkeysFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	v := humnValue(monkeys)

	_, err = fmt.Fprintf(w, "%d", v)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type monkey struct {
	id            string
	valueKnown    bool
	value         int
	op            operation
	dependsOnHumn bool
}

type operator uint8

const (
	opAdd operator = iota
	opSubstract
	opMultiply
	opDivide
	opEqual
)

type operation struct {
	left, right string // monkey IDs
	op          operator
}

const (
	humnMonkeyID = "humn"
	rootMonkeyID = "root"
)

func rootValue(monkeys []*monkey) int {
	monkeysByID := make(map[string]*monkey)
	for _, m := range monkeys {
		monkeysByID[m.id] = m
	}

	var computeValue func(m *monkey)
	computeValue = func(m *monkey) {
		if m.valueKnown {
			return
		}

		left := monkeysByID[m.op.left]
		right := monkeysByID[m.op.right]

		computeValue(left)
		computeValue(right)

		switch m.op.op {
		case opAdd:
			m.value = left.value + right.value
		case opSubstract:
			m.value = left.value - right.value
		case opMultiply:
			m.value = left.value * right.value
		case opDivide:
			m.value = left.value / right.value
		default:
			panic("unknown operation")
		}

		m.valueKnown = true
	}

	root := monkeysByID[rootMonkeyID]
	computeValue(root)

	return root.value
}

func humnValue(monkeys []*monkey) int {
	monkeysByID := make(map[string]*monkey)
	for _, m := range monkeys {
		monkeysByID[m.id] = m
	}

	monkeysByID[humnMonkeyID].dependsOnHumn = true
	monkeysByID[rootMonkeyID].op.op = opEqual

	var computeValue func(m *monkey)
	computeValue = func(m *monkey) {
		if m.valueKnown {
			return
		}

		left := monkeysByID[m.op.left]
		right := monkeysByID[m.op.right]

		computeValue(left)
		computeValue(right)

		if left.dependsOnHumn || right.dependsOnHumn {
			m.dependsOnHumn = true
			return
		}

		switch m.op.op {
		case opAdd:
			m.value = left.value + right.value
		case opSubstract:
			m.value = left.value - right.value
		case opMultiply:
			m.value = left.value * right.value
		case opDivide:
			m.value = left.value / right.value
		default:
			panic("unknown operation")
		}

		m.valueKnown = true
	}

	root := monkeysByID[rootMonkeyID]
	computeValue(root)

	left := monkeysByID[root.op.left]
	right := monkeysByID[root.op.right]

	if left.dependsOnHumn {
		root.value = right.value
	} else {
		root.value = left.value
	}

	for m := root; m.id != humnMonkeyID; {
		left := monkeysByID[m.op.left]
		right := monkeysByID[m.op.right]

		if left.dependsOnHumn {
			left.value, _ = reverseOp(m.op.op, nil, &right.value, m.value)
			m = left
		} else {
			_, right.value = reverseOp(m.op.op, &left.value, nil, m.value)
			m = right
		}
	}

	return monkeysByID[humnMonkeyID].value
}

func reverseOp(op operator, left, right *int, target int) (int, int) {
	switch op {
	case opAdd:
		if left == nil {
			return target - *right, *right
		}
		return *left, target - *left
	case opSubstract:
		if left == nil {
			return *right + target, *right
		}
		return *left, *left - target
	case opMultiply:
		if left == nil {
			return target / *right, *right
		}
		return *left, target / *left
	case opDivide:
		if left == nil {
			return target * *right, *right
		}
		return *left, *left / target
	case opEqual:
		if left == nil {
			return *right, *right
		}
		return *left, *left
	default:
		panic("unknown operation")
	}
}

func monkeysFromReader(r io.Reader) ([]*monkey, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	haveHumn := false
	haveRoot := false

	var monkeys []*monkey
	for _, l := range lines {
		m, err := monkeyFromString(l)
		if err != nil {
			return nil, err
		}

		if m.id == humnMonkeyID {
			haveHumn = true
		}
		if m.id == rootMonkeyID {
			haveRoot = true
		}

		monkeys = append(monkeys, m)
	}

	if !haveHumn {
		return nil, fmt.Errorf("no monkey with ID %q", humnMonkeyID)
	}
	if !haveRoot {
		return nil, fmt.Errorf("no monkey with ID %q", rootMonkeyID)
	}

	return monkeys, nil
}

func monkeyFromString(s string) (*monkey, error) {
	parts := strings.SplitN(s, ": ", 2)
	if len(parts) != 2 {
		return nil, errors.New("wrong format")
	}

	m := monkey{
		id: parts[0],
	}

	if len(parts[1]) == 11 { // heuristic, works for provided input
		op, err := operationFromString(parts[1])
		if err != nil {
			return nil, err
		}

		m.op = op
	} else {
		v, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("%q is not a number", parts[1])
		}

		m.value = v
		m.valueKnown = true
	}

	return &m, nil
}

func operationFromString(s string) (operation, error) {
	parts := strings.SplitN(s, " ", 3)
	if len(parts) != 3 {
		return operation{}, errors.New("wrong format")
	}

	switch parts[1] {
	case "+":
		return operation{parts[0], parts[2], opAdd}, nil
	case "-":
		return operation{parts[0], parts[2], opSubstract}, nil
	case "*":
		return operation{parts[0], parts[2], opMultiply}, nil
	case "/":
		return operation{parts[0], parts[2], opDivide}, nil
	default:
		return operation{}, fmt.Errorf("unknown operator %q", parts[1])
	}
}
