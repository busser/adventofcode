package d25

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
	"github.com/busser/adventofcode/y2022/d25/slices"
)

// PartOne solves the first problem of day 25 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	numbers, err := snafuNumbersFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for _, n := range numbers {
		sum += snafuToInt(n)
	}

	_, err = fmt.Fprintf(w, "%s", intToSnafu(sum))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

var snafuDigits = []byte{'=', '-', '0', '1', '2'}

const firstDigitValue = -2

func snafuDigitValue(d byte) int {
	return slices.Index(snafuDigits, d) + firstDigitValue
}

func snafuToInt(s string) int {
	n, weight := 0, 1

	for i := len(s) - 1; i >= 0; i-- {
		n += weight * snafuDigitValue(s[i])
		weight *= len(snafuDigits)
	}

	return n
}

func intToSnafu(n int) string {
	var digits []byte // in reverse order for now

	for n > 0 {
		index := (n - firstDigitValue) % len(snafuDigits)
		digit := snafuDigits[index]
		digits = append(digits, digit)
		n -= index + firstDigitValue
		n /= len(snafuDigits)
	}

	slices.Reverse(digits)
	return string(digits)
}

func snafuNumbersFromReader(r io.Reader) ([]string, error) {
	numbers, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	for _, n := range numbers {
		for i := range n {
			if !slices.Contains(snafuDigits, n[i]) {
				return nil, fmt.Errorf("%q is not a snafu number", n)
			}
		}
	}

	return numbers, nil
}
