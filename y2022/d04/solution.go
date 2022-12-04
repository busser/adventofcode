package d04

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 4 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	pairs, err := sectionPairsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0
	for _, p := range pairs {
		if sectionOverlapCompletely(p) {
			count++
		}
	}

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 4 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	pairs, err := sectionPairsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0
	for _, p := range pairs {
		if sectionOverlap(p) {
			count++
		}
	}

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type section struct {
	start, end int
}

func (s section) len() int {
	return s.end - s.start + 1
}

type pair struct {
	first, second section
}

func sectionOverlapCompletely(p pair) bool {
	minStart := min(p.first.start, p.second.start)
	maxEnd := max(p.first.end, p.second.end)

	union := section{minStart, maxEnd} // not really a union, if the sections don't overlap.

	return p.first == union || p.second == union
}

func sectionOverlap(p pair) bool {
	minStart := min(p.first.start, p.second.start)
	maxEnd := max(p.first.end, p.second.end)

	union := section{minStart, maxEnd} // not really a union, if the sections don't overlap.

	return p.first.len()+p.second.len() > union.len()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func sectionPairsFromReader(r io.Reader) ([]pair, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	var pairs []pair

	for _, l := range lines {
		p, err := sectionPairFromString(l)
		if err != nil {
			return nil, err
		}

		pairs = append(pairs, p)
	}

	return pairs, nil
}

func sectionPairFromString(s string) (pair, error) {
	parts := strings.SplitN(s, ",", 2)
	if len(parts) != 2 {
		return pair{}, errors.New("wrong format")
	}

	first, err := sectionFromString(parts[0])
	if err != nil {
		return pair{}, err
	}

	second, err := sectionFromString(parts[1])
	if err != nil {
		return pair{}, err
	}

	return pair{first, second}, nil
}

func sectionFromString(s string) (section, error) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return section{}, errors.New("wrong format")
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return section{}, fmt.Errorf("%q is not a number", parts[0])
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return section{}, fmt.Errorf("%q is not a number", parts[1])
	}

	return section{start, end}, nil
}
