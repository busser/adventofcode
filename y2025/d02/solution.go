package d02

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 2 of Advent of Code 2025.
func PartOne(r io.Reader, w io.Writer) error {
	idRanges, err := rangesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var sum int
	for _, idRange := range idRanges {
		for id := idRange.begin; id <= idRange.end; id++ {
			if idIsSimpleRepetition(id) {
				sum += id
			}
		}
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 2 of Advent of Code 2025.
func PartTwo(r io.Reader, w io.Writer) error {
	idRanges, err := rangesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var sum int
	for _, idRange := range idRanges {
		for id := idRange.begin; id <= idRange.end; id++ {
			if idIsRepetition(id) {
				sum += id
			}
		}
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func idIsSimpleRepetition(id int) bool {
	strID := strconv.Itoa(id)

	if len(strID)%2 != 0 {
		return false
	}

	return stringRepeats(strID, len(strID)/2)
}

func idIsRepetition(id int) bool {
	strID := strconv.Itoa(id)

	for n := 1; n <= len(strID)/2; n++ {
		if stringRepeats(strID, n) {
			return true
		}
	}

	return false
}

func stringRepeats(s string, n int) (r bool) {
	if n == 0 {
		return false
	}
	if len(s)%n != 0 {
		return false
	}
	if len(s)/n < 2 {
		return false
	}

	for i := n; i < len(s); i += n {
		if s[i:i+n] != s[0:n] {
			return false
		}
	}

	return true
}

type idRange struct {
	begin, end int
}

func rangeFromString(s string) (idRange, error) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return idRange{}, fmt.Errorf("bad range")
	}

	begin, err := strconv.Atoi(parts[0])
	if err != nil {
		return idRange{}, err
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return idRange{}, err
	}

	return idRange{begin: begin, end: end}, nil
}

func (r idRange) String() string {
	return fmt.Sprintf("%d-%d", r.begin, r.end)
}

func rangesFromReader(r io.Reader) ([]idRange, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("expected 1 line, got %d", len(lines))
	}

	parts := strings.Split(lines[0], ",")
	ranges := make([]idRange, len(parts))
	for i, part := range parts {
		r, err := rangeFromString(part)
		if err != nil {
			return nil, fmt.Errorf("bad range %q: %w", part, err)
		}
		ranges[i] = r
	}

	return ranges, nil
}
