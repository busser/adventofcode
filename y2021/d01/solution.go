package busser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// PartOne solves the first problem of day 1 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	numbers, err := numbersFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", increases(numbers, 1))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	numbers, err := numbersFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", increases(numbers, 3))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// increases counts how many times the rolling sum of window numbers is larger
// than the previous one.
func increases(numbers []int, window int) int {
	n := len(numbers)
	count := 0
	for i := window; i < n; i++ {
		if numbers[i] > numbers[i-window] {
			count++
		}
	}
	return count
}

func numbersFromReader(r io.Reader) ([]int, error) {
	var numbers []int

	s := bufio.NewScanner(r)
	for s.Scan() {
		n, err := strconv.Atoi(s.Text())
		if err != nil {
			return nil, fmt.Errorf("%q is not a number", s.Text())
		}
		numbers = append(numbers, n)
	}
	if s.Err() != nil {
		return nil, fmt.Errorf("failed to scan reader: %w", s.Err())
	}

	return numbers, nil
}
