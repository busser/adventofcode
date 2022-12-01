package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 9 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	numbers, err := numbersFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	first, err := firstNumberWithoutProperty(numbers)
	if err != nil {
		return fmt.Errorf("finding number: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", first)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 9 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	numbers, err := numbersFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	first, err := firstNumberWithoutProperty(numbers)
	if err != nil {
		return fmt.Errorf("finding number: %w", err)
	}

	sequence, err := contiguousSum(first, numbers)
	if err != nil {
		return fmt.Errorf("contiguous sum: %w", err)
	}
	if len(sequence) == 0 {
		return errors.New("contiguous sum is empty slice")
	}

	min, max := sequence[0], sequence[0]
	for _, n := range sequence {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}

	_, err = fmt.Fprintf(answer, "%d", min+max)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func contiguousSum(sum int, numbers []int) ([]int, error) {
	i, j, slidingSum := 0, 0, 0

	for {
		switch {
		case slidingSum == sum:
			return numbers[i:j], nil
		case slidingSum > sum:
			if i == j {
				return nil, errors.New("not found")
			}
			slidingSum -= numbers[i]
			i++
		case slidingSum < sum:
			if j == len(numbers) {
				return nil, errors.New("not found")
			}
			slidingSum += numbers[j]
			j++
		}
	}
}

const preambleSize = 25

func firstNumberWithoutProperty(numbers []int) (int, error) {
	if len(numbers) < preambleSize {
		return 0, fmt.Errorf("too few numbers: only %d, require at least %d", len(numbers), preambleSize)
	}

	for i := 25; i < len(numbers); i++ {
		n := numbers[i]

		if !isValid(n, numbers[i-preambleSize:i]) {
			return n, nil
		}
	}

	return 0, errors.New("not found")
}

func isValid(n int, preamble []int) bool {
	for i := 0; i < len(preamble); i++ {
		for j := i + 1; j < len(preamble); j++ {
			if preamble[i]+preamble[j] == n {
				return true
			}
		}
	}
	return false
}

func numbersFromReader(r io.Reader) ([]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading lines: %w", err)
	}

	numbers := make([]int, len(lines))
	for i, l := range lines {
		n, err := strconv.Atoi(l)
		if err != nil {
			return nil, fmt.Errorf("wrong format: not an int: %w", err)
		}

		numbers[i] = n
	}

	return numbers, nil
}
