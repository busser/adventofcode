package intcode

import (
	"fmt"
	"os"
)

type state struct {
	program []int
	pointer int

	inputFn  func() int
	outputFn func(int)

	err error
}

func (s *state) run() error {
	for {
		if Debug {
			fmt.Fprintf(os.Stderr, "\npos=%d\nerr=%v\nprogram=%v\n", s.pointer, s.err, s.program)
		}

		if s.err != nil {
			return s.err
		}

		opcode := s.read(s.getPointer())

		switch opcode % 100 {
		case 1:
			s.add()
		case 2:
			s.multiply()
		case 3:
			s.input()
		case 4:
			s.output()
		case 5:
			s.jumpIfTrue()
		case 6:
			s.jumpIfFalse()
		case 7:
			s.lessThan()
		case 8:
			s.equals()
		case 99:
			return nil
		default:
			return &unknownOpcodeError{s, opcode}
		}
	}
}

func (s *state) boundsCheck(position int) bool {
	if position < 0 || position >= len(s.program) {
		s.err = &outOfBoundsError{s}
		return false
	}
	return true
}

func (s *state) read(position int) int {
	if !s.boundsCheck(position) {
		return 0
	}

	if Debug {
		fmt.Fprintf(os.Stderr, "reading position %d\n", position)
	}

	return s.program[position]
}

func (s *state) write(position, value int) {
	if !s.boundsCheck(position) {
		return
	}

	if Debug {
		fmt.Fprintf(os.Stderr, "writing %d to position %d\n", value, position)
	}

	s.program[position] = value
}

func (s *state) getPointer() int {
	return s.pointer
}

func (s *state) setPointer(position int) {
	if !s.boundsCheck(position) {
		return
	}

	if Debug {
		fmt.Fprintf(os.Stderr, "moving pointer to position %d\n", position)
	}

	s.pointer = position
}
