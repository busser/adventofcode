package busser

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 18 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	expressions, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	operatorPrecedence := map[rune]int{
		'(': 0,
		')': 1,
		'*': 2,
		'+': 2,
	}

	eval := newEvaluator(operatorPrecedence)

	sum := 0

	for i, exp := range expressions {
		v, err := eval.evaluate(exp)
		if err != nil {
			return fmt.Errorf("evaluating expression #%d: %w", i, err)
		}

		sum += v
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 18 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	expressions, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	operatorPrecedence := map[rune]int{
		'(': 0,
		')': 1,
		'*': 2,
		'+': 3,
	}

	eval := newEvaluator(operatorPrecedence)

	sum := 0

	for i, exp := range expressions {
		v, err := eval.evaluate(exp)
		if err != nil {
			return fmt.Errorf("evaluating expression #%d: %w", i, err)
		}

		sum += v
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type evaluator struct {
	operators  operatorStack
	numbers    numberStack
	precedence map[rune]int
}

func newEvaluator(operatorPrecedence map[rune]int) evaluator {
	return evaluator{
		precedence: operatorPrecedence,
	}
}

func (e *evaluator) evaluate(expression string) (int, error) {
	for i, token := range expression {
		if err := e.processToken(token); err != nil {
			return 0, fmt.Errorf("processing token %q at position %d: %w", token, i, err)
		}
	}

	for e.operators.len() > 0 {
		if err := e.unstackOperator(); err != nil {
			return 0, fmt.Errorf("unstacking operator: %w", err)
		}
	}

	if e.numbers.len() != 1 {
		return 0, fmt.Errorf("invalid expression: expected 1 value left in stack, have %d", e.numbers.len())
	}

	return e.numbers.pop(), nil
}

func (e *evaluator) processToken(t rune) error {
	switch t {
	case ' ':
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		e.numbers.push(int(t - '0'))
	case '+', '*':
		p := e.precedence[t]
		for e.operators.len() > 0 && p <= e.precedence[e.operators.peek()] {
			if err := e.unstackOperator(); err != nil {
				return fmt.Errorf("unstacking operator: %w", err)
			}
		}
		e.operators.push(t)
	case '(':
		e.operators.push(t)
	case ')':
		p := e.precedence[t]
		for e.operators.len() > 0 && p <= e.precedence[e.operators.peek()] {
			if err := e.unstackOperator(); err != nil {
				return fmt.Errorf("unstacking operator: %w", err)
			}
		}
		if e.operators.len() == 0 || e.operators.peek() != '(' {
			return errors.New("closing parenthesis has no matching opening parenthesis")
		}
		e.operators.pop()
	default:
		return errors.New("unknown token")
	}

	return nil
}

func (e *evaluator) unstackOperator() error {
	if e.operators.len() == 0 {
		return errors.New("no operator to unstack")
	}

	op := e.operators.pop()

	switch op {
	case '+':
		if e.numbers.len() < 2 {
			return fmt.Errorf("need 2 numbers to perform addition, have %d", e.numbers.len())
		}
		e.numbers.push(e.numbers.pop() + e.numbers.pop())
	case '*':
		if e.numbers.len() < 2 {
			return fmt.Errorf("need 2 numbers to perform multiplication, have %d", e.numbers.len())
		}
		e.numbers.push(e.numbers.pop() * e.numbers.pop())
	case '(':
		return errors.New("opening parenthesis has no matching closing parenthesis")
	case ')':
		return errors.New("closing parenthesis has no matching opening parenthesis")
	default:
		return fmt.Errorf("unknown operator %q", op)
	}

	return nil
}
