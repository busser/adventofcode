package d05

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 5 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	stacks, moves, err := stacksAndMovesFromReader(r)
	if err != nil {
		return err
	}

	if err := moveCrates(stacks, moves); err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s", string(topCrates(stacks)))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	stacks, moves, err := stacksAndMovesFromReader(r)
	if err != nil {
		return err
	}

	if err := moveCratesKeepingOrder(stacks, moves); err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s", string(topCrates(stacks)))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type stack[T any] []T

func (s *stack[T]) len() int {
	return len(*s)
}

func (s *stack[T]) push(v T) {
	*s = append(*s, v)
}

func (s *stack[T]) pushN(v []T) {
	*s = append(*s, v...)
}

func (s *stack[T]) pop() T {
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *stack[T]) popN(n int) []T {
	v := (*s)[len(*s)-n:]
	*s = (*s)[:len(*s)-n]
	return v
}

func (s *stack[T]) peek() T {
	return (*s)[len(*s)-1]
}

type move struct {
	amount int
	from   int
	to     int
}

func moveCrates(stacks []stack[rune], moves []move) error {
	for _, m := range moves {
		if m.amount > stacks[m.from].len() {
			return fmt.Errorf("cannot move %d crates from stack %d, only have %d", m.amount, m.from, stacks[m.from].len())
		}

		for i := 0; i < m.amount; i++ {
			crate := stacks[m.from].pop()
			stacks[m.to].push(crate)
		}
	}

	return nil
}

func moveCratesKeepingOrder(stacks []stack[rune], moves []move) error {
	for _, m := range moves {
		if m.amount > stacks[m.from].len() {
			return fmt.Errorf("cannot move %d crates from stack %d, only have %d", m.amount, m.from, stacks[m.from].len())
		}

		crates := stacks[m.from].popN(m.amount)
		stacks[m.to].pushN(crates)
	}

	return nil
}

func topCrates(stacks []stack[rune]) []rune {
	var crates []rune
	for _, s := range stacks {
		if s.len() == 0 {
			continue
		}
		crates = append(crates, s.peek())
	}
	return crates
}

func stacksAndMovesFromReader(r io.Reader) ([]stack[rune], []move, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, nil, err
	}

	var stackLines, moveLines []string
	for i := 1; i < len(lines)-1; i++ {
		if len(lines[i]) == 0 {
			stackLines = lines[:i]
			moveLines = lines[i+1:]
			break
		}
	}

	if stackLines == nil {
		return nil, nil, errors.New("wrong format")
	}

	stacks, err := stacksFromLines(stackLines)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing stacks: %w", err)
	}

	moves, err := movesFromLines(moveLines, stacks)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing moves: %w", err)
	}

	return stacks, moves, nil
}

func stacksFromLines(lines []string) ([]stack[rune], error) {
	length := len(lines[0])
	for _, l := range lines {
		if len(l) != length {
			return nil, errors.New("all lines don't have same length")
		}
	}

	if length%4 != 3 {
		return nil, fmt.Errorf("line length not divisible as expected")
	}

	numStacks := length/4 + 1
	stacks := make([]stack[rune], numStacks)

	for row := len(lines) - 2; row >= 0; row-- {
		for s := range stacks {
			col := 4*s + 1
			tag := rune(lines[row][col])
			if tag != ' ' {
				stacks[s].push(tag)
			}
		}
	}

	return stacks, nil
}

func movesFromLines(lines []string, stacks []stack[rune]) ([]move, error) {
	moves := make([]move, len(lines))

	for i, l := range lines {
		m, err := moveFromString(l)
		if err != nil {
			return nil, err
		}

		if m.from < 0 || m.from >= len(stacks) {
			return nil, fmt.Errorf("cannot move from stack %d", m.from)
		}
		if m.to < 0 || m.to >= len(stacks) {
			return nil, fmt.Errorf("cannot move to stack %d", m.to)
		}

		moves[i] = m
	}

	return moves, nil
}

func moveFromString(s string) (move, error) {
	parts := strings.SplitN(s, " ", 6)
	if len(parts) != 6 {
		return move{}, errors.New("wrong format")
	}

	amount, err := strconv.Atoi(parts[1])
	if err != nil {
		return move{}, fmt.Errorf("%q is not a number", parts[1])
	}

	from, err := strconv.Atoi(parts[3])
	if err != nil {
		return move{}, fmt.Errorf("%q is not a number", parts[3])
	}

	to, err := strconv.Atoi(parts[5])
	if err != nil {
		return move{}, fmt.Errorf("%q is not a number", parts[5])
	}

	return move{
		amount: amount,
		from:   from - 1, // make it zero-indexed
		to:     to - 1,   // make it zero-indexed
	}, nil
}
