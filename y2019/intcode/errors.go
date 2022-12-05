package intcode

import "fmt"

type unknownOpcodeError struct {
	s      *state
	opcode int
}

func (err *unknownOpcodeError) Error() string {
	return fmt.Sprintf("unknown opcode %d at position %d", err.s.read(err.s.pointer), err.s.pointer)
}

type outOfBoundsError struct {
	s *state
}

func (err *outOfBoundsError) Error() string {
	return fmt.Sprintf("position %d is out of bounds", err.s.pointer)
}
