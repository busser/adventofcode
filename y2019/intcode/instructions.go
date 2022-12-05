package intcode

import "log"

func (s *state) add() {
	var (
		opcode = s.read(s.getPointer() + 0)
		first  = s.read(s.getPointer() + 1)
		second = s.read(s.getPointer() + 2)
		third  = s.read(s.getPointer() + 3)
	)

	if parameterMode(opcode, 1) == paramModePosition {
		first = s.read(first)
	}
	if parameterMode(opcode, 2) == paramModePosition {
		second = s.read(second)
	}

	sum := first + second
	s.write(third, sum)

	s.setPointer(s.getPointer() + 4)
}

func (s *state) multiply() {
	var (
		opcode = s.read(s.getPointer() + 0)
		first  = s.read(s.getPointer() + 1)
		second = s.read(s.getPointer() + 2)
		third  = s.read(s.getPointer() + 3)
	)

	if parameterMode(opcode, 1) == paramModePosition {
		first = s.read(first)
	}
	if parameterMode(opcode, 2) == paramModePosition {
		second = s.read(second)
	}

	log.Println(opcode, first, second, third)

	product := first * second
	s.write(third, product)

	s.setPointer(s.getPointer() + 4)
}

func (s *state) input() {
	var (
		_     = s.read(s.getPointer() + 0)
		first = s.read(s.getPointer() + 1)
	)

	s.write(first, s.inputFn())

	s.setPointer(s.getPointer() + 2)
}

func (s *state) output() {
	var (
		opcode = s.read(s.getPointer() + 0)
		first  = s.read(s.getPointer() + 1)
	)

	if parameterMode(opcode, 1) == paramModePosition {
		first = s.read(first)
	}

	s.outputFn(first)

	s.setPointer(s.getPointer() + 2)
}

func (s *state) jumpIfTrue() {
	var (
		opcode = s.read(s.getPointer() + 0)
		first  = s.read(s.getPointer() + 1)
		second = s.read(s.getPointer() + 2)
	)

	if parameterMode(opcode, 1) == paramModePosition {
		first = s.read(first)
	}
	if parameterMode(opcode, 2) == paramModePosition {
		second = s.read(second)
	}

	if first != 0 {
		s.setPointer(second)
	} else {
		s.setPointer(s.getPointer() + 3)
	}

}

func (s *state) jumpIfFalse() {
	var (
		opcode = s.read(s.getPointer() + 0)
		first  = s.read(s.getPointer() + 1)
		second = s.read(s.getPointer() + 2)
	)

	if parameterMode(opcode, 1) == paramModePosition {
		first = s.read(first)
	}
	if parameterMode(opcode, 2) == paramModePosition {
		second = s.read(second)
	}

	if first == 0 {
		s.setPointer(second)
	} else {
		s.setPointer(s.getPointer() + 3)
	}

}

func (s *state) lessThan() {
	var (
		opcode = s.read(s.getPointer() + 0)
		first  = s.read(s.getPointer() + 1)
		second = s.read(s.getPointer() + 2)
		third  = s.read(s.getPointer() + 3)
	)

	if parameterMode(opcode, 1) == paramModePosition {
		first = s.read(first)
	}
	if parameterMode(opcode, 2) == paramModePosition {
		second = s.read(second)
	}

	if first < second {
		s.write(third, 1)
	} else {
		s.write(third, 0)
	}

	s.setPointer(s.getPointer() + 4)
}

func (s *state) equals() {
	var (
		opcode = s.read(s.getPointer() + 0)
		first  = s.read(s.getPointer() + 1)
		second = s.read(s.getPointer() + 2)
		third  = s.read(s.getPointer() + 3)
	)

	if parameterMode(opcode, 1) == paramModePosition {
		first = s.read(first)
	}
	if parameterMode(opcode, 2) == paramModePosition {
		second = s.read(second)
	}

	if first == second {
		s.write(third, 1)
	} else {
		s.write(third, 0)
	}

	s.setPointer(s.getPointer() + 4)
}

const (
	paramModePosition  = 0
	paramModeImmediate = 1
)

func parameterMode(opcode int, n int) int {
	opcode /= 100

	for ; n > 1; n-- {
		opcode /= 10
	}

	return opcode % 10
}
