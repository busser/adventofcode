package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 8 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	prog, err := programFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	if err := prog.run(); err != nil {
		return fmt.Errorf("running program: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", prog.accumulator)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 8 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	prog, err := programFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	for i, inst := range prog.instructions {
		replacedOperation := inst.operation

		switch inst.operation {
		case "acc":
			continue
		case "jmp":
			prog.instructions[i].operation = "nop"
		case "nop":
			prog.instructions[i].operation = "jmp"
		default:
			return fmt.Errorf("unknown operation %q", inst.operation)
		}

		if err := prog.run(); err != nil {
			return fmt.Errorf("running program: %w", err)
		}

		if prog.index == len(prog.instructions) {
			if _, err := fmt.Fprintf(answer, "%d", prog.accumulator); err != nil {
				return fmt.Errorf("could not write answer: %w", err)
			}
			return nil
		}

		prog.instructions[i].operation = replacedOperation
		prog.reset()
	}

	return errors.New("could not find answer")
}

type program struct {
	accumulator  int
	index        int
	instructions []instruction
}

func (p *program) reset() {
	p.accumulator = 0
	p.index = 0
}

func (p *program) run() error {
	seen := make(map[int]bool)

	for {
		if p.index >= len(p.instructions) {
			break
		}

		if seen[p.index] {
			break
		}
		seen[p.index] = true

		inst := p.instructions[p.index]

		switch inst.operation {
		case "acc":
			p.accumulator += inst.argument
			p.index++
		case "jmp":
			p.index += inst.argument
		case "nop":
			p.index++
		default:
			return fmt.Errorf("unknown operation %q", inst.operation)
		}
	}

	return nil
}

type instruction struct {
	operation string
	argument  int
}

func programFromReader(r io.Reader) (*program, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading lines: %w", err)
	}

	instructions := make([]instruction, len(lines))
	for i, l := range lines {
		instructions[i].operation = l[:3]

		arg, err := strconv.Atoi(l[4:])
		if err != nil {
			return nil, errors.New("wrong format")
		}
		instructions[i].argument = arg
	}

	return &program{instructions: instructions}, nil
}
