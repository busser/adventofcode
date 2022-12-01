package busser

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 24 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	steps, err := stepsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	largest, err := largestValidInput(steps)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(answer, "%d", largest)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 24 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	steps, err := stepsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	lowest, err := lowestValidInput(steps)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(answer, "%d", lowest)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type aluStepType uint8

const (
	aluPushStep aluStepType = iota
	aluPopStep
)

type aluStep struct {
	index  int
	typ    aluStepType
	offset int
}

func largestValidInput(steps []aluStep) (int, error) {
	pairs, err := matchPushAndPopSteps(steps)
	if err != nil {
		return 0, err
	}

	inputValues := make([]int, len(steps))

	for _, p := range pairs {
		push, pop := p[0], p[1]
		delta := push.offset + pop.offset
		switch {
		case delta <= 0:
			inputValues[push.index] = 9
			inputValues[pop.index] = 9 + delta
		case delta > 0:
			inputValues[push.index] = 9 - delta
			inputValues[pop.index] = 9
		}
	}

	inputNumber := 0
	for _, v := range inputValues {
		inputNumber *= 10
		inputNumber += v
	}

	return inputNumber, nil
}

func lowestValidInput(steps []aluStep) (int, error) {
	pairs, err := matchPushAndPopSteps(steps)
	if err != nil {
		return 0, err
	}

	inputValues := make([]int, len(steps))

	for _, p := range pairs {
		push, pop := p[0], p[1]
		delta := push.offset + pop.offset
		switch {
		case delta <= 0:
			inputValues[push.index] = 1 - delta
			inputValues[pop.index] = 1
		case delta > 0:
			inputValues[push.index] = 1
			inputValues[pop.index] = 1 + delta
		}
	}

	inputNumber := 0
	for _, v := range inputValues {
		inputNumber *= 10
		inputNumber += v
	}

	return inputNumber, nil
}

func matchPushAndPopSteps(steps []aluStep) ([][2]aluStep, error) {
	var stack []aluStep
	var pairs [][2]aluStep

	for _, s := range steps {
		switch s.typ {
		case aluPushStep:
			stack = append(stack, s)
		case aluPopStep:
			if len(stack) == 0 {
				return nil, errors.New("could not match all steps to another")
			}
			match := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			pairs = append(pairs, [2]aluStep{match, s})
		}
	}

	if len(stack) > 0 {
		return nil, errors.New("could not match all steps to another")
	}

	return pairs, nil
}

func stepsFromReader(r io.Reader) ([]aluStep, error) {
	const numSteps, stepSize = 14, 18

	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	if len(lines) != numSteps*stepSize {
		return nil, fmt.Errorf("expected %d lines, got %d", numSteps*stepSize, len(lines))
	}

	steps := make([]aluStep, numSteps)

	for i := 0; i < numSteps; i++ {
		step, err := stepFromLines(lines[i*stepSize : (i+1)*stepSize])
		if err != nil {
			return nil, fmt.Errorf("invalid step n°%d: %w", i+1, err)
		}
		step.index = i
		steps[i] = step
	}

	return steps, nil
}

func stepFromLines(lines []string) (aluStep, error) {
	var step aluStep

	var rawStepType int
	if _, err := fmt.Sscanf(lines[4], "div z %d", &rawStepType); err != nil {
		return aluStep{}, fmt.Errorf("invalid format on line %d", 4+1)
	}

	switch rawStepType {
	case 1:
		step.typ = aluPushStep
	case 26:
		step.typ = aluPopStep
	default:
		return aluStep{}, fmt.Errorf("unknown step type %d", rawStepType)
	}

	switch step.typ {
	case aluPushStep:
		if _, err := fmt.Sscanf(lines[15], "add y %d", &step.offset); err != nil {
			return aluStep{}, fmt.Errorf("invalid format on line %d", 15+1)
		}
	case aluPopStep:
		if _, err := fmt.Sscanf(lines[5], "add x %d", &step.offset); err != nil {
			return aluStep{}, fmt.Errorf("invalid format on line %d", 5+1)
		}
	}

	return step, nil
}
