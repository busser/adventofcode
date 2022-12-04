package intcode

import "fmt"

type unknownInstructionError struct {
	i *interpreter
}

func (err *unknownInstructionError) Error() string {
	return fmt.Sprintf("unknown instruction %d at position %d", err.i.readAbsolute(err.i.pointer), err.i.pointer)
}

type outOfBoundsError struct {
	i *interpreter
}

func (err *outOfBoundsError) Error() string {
	return fmt.Sprintf("position %d is out of bounds", err.i.pointer)
}
