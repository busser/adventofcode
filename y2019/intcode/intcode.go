package intcode

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

var Debug = true

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

func Run(program, inputs []int) ([]int, error) {
	var outputs []int

	s := state{
		program: program,
		inputFn: func() int {
			v := inputs[0]
			inputs = inputs[1:]
			return v
		},
		outputFn: func(v int) {
			outputs = append(outputs, v)
		},
	}

	if err := s.run(); err != nil {
		return nil, err
	}

	return outputs, nil
}
