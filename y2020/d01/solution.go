package busser

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 1 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	expenses, err := intsFromReader(input)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	sort.Ints(expenses)

	a, b, found := pairThatSumsTo(expenses, 2020)
	if !found {
		return errors.New("no answer found")
	}

	_, err = fmt.Fprintf(answer, "%d", a*b)
	if err != nil {
		return fmt.Errorf("writing answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	expenses, err := intsFromReader(input)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	sort.Ints(expenses)

	a, b, c, found := tripletThatSumsTo(expenses, 2020)
	if !found {
		return errors.New("no answer found")
	}

	_, err = fmt.Fprintf(answer, "%d", a*b*c)
	if err != nil {
		return fmt.Errorf("writing answer: %w", err)
	}

	return nil
}

// pairThatSumsTo assumes nums is sorted.
func pairThatSumsTo(nums []int, targetSum int) (a, b int, found bool) {
	left, right := 0, len(nums)-1
	for right > left {
		sum := nums[left] + nums[right]
		switch {
		case sum > targetSum:
			right--
		case sum < targetSum:
			left++
		case sum == targetSum:
			return nums[left], nums[right], true
		}
	}

	return 0, 0, false
}

// tripletThatSumsTo assumes nums is sorted.
func tripletThatSumsTo(nums []int, targetSum int) (a, b, c int, found bool) {
	for i, a := range nums {
		b, c, found := pairThatSumsTo(nums[i+1:], targetSum-a)
		if found {
			return a, b, c, true
		}
	}

	return 0, 0, 0, false
}

func intsFromReader(r io.Reader) ([]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading lines: %w", err)
	}

	ints := make([]int, len(lines))
	for i := range lines {
		n, err := strconv.Atoi(lines[i])
		if err != nil {
			return nil, fmt.Errorf("%q not an int", lines[i])
		}
		ints[i] = n
	}

	return ints, nil
}
