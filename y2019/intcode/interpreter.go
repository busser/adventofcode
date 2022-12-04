package intcode

import (
	"fmt"
	"io"
	"os"

	"github.com/busser/adventofcode/helpers"
)

type interpreter struct {
	program []int
	pointer int
	err     error
}

const (
	add      = 1
	multiply = 2
	exit     = 99
)

func ProgramFromReader(r io.Reader) ([]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}
	if len(lines) != 1 {
		return nil, fmt.Errorf("expected 1 line of input, got %d", len(lines))
	}

	return helpers.IntsFromString(lines[0], ",")
}

var Debug = false

func Run(program []int) error {
	i := interpreter{program: program}
	return i.run()
}

func (i *interpreter) run() error {
	for {
		if Debug {
			fmt.Fprintf(os.Stderr, "\npos=%d\nerr=%v\nprogram=%v\n", i.pointer, i.err, i.program)
		}

		op := i.readAbsolute(i.pointer)
		if i.err != nil {
			return i.err
		}

		switch op {
		case add:
			i.add()
		case multiply:
			i.multiply()
		case exit:
			return nil
		default:
			return &unknownInstructionError{i}
		}

		i.pointer += 4
	}
}

func (i *interpreter) boundsCheck(position int) bool {
	if position < 0 || position >= len(i.program) {
		i.err = &outOfBoundsError{i}
		return false
	}
	return true
}

func (i *interpreter) readRelative(offset int) int {
	position := i.pointer + offset
	return i.readAbsolute(position)
}

func (i *interpreter) readAbsolute(position int) int {
	if !i.boundsCheck(position) {
		return 0
	}

	return i.program[position]
}

func (i *interpreter) writeRelative(offset, value int) {
	position := i.pointer + offset
	i.writeAbsolute(position, value)
}

func (i *interpreter) writeAbsolute(position, value int) {
	if !i.boundsCheck(position) {
		return
	}

	if Debug {
		fmt.Fprintf(os.Stderr, "writing %d to position %d\n", value, position)
	}

	i.program[position] = value
}

func (i *interpreter) add() {
	var (
		first  = i.readRelative(1)
		second = i.readRelative(2)
		third  = i.readRelative(3)
	)

	sum := i.readAbsolute(first) + i.readAbsolute(second)
	i.writeAbsolute(third, sum)
}

func (i *interpreter) multiply() {
	var (
		first  = i.readRelative(1)
		second = i.readRelative(2)
		third  = i.readRelative(3)
	)

	product := i.readAbsolute(first) * i.readAbsolute(second)
	i.writeAbsolute(third, product)
}
