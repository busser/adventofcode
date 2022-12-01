package busser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 8 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	entries, err := entriesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read entries: %w", err)
	}

	count := 0
	for _, e := range entries {
		for _, o := range e.output {
			switch len(o) {
			case
				len(workingPatterns[1]),
				len(workingPatterns[4]),
				len(workingPatterns[7]),
				len(workingPatterns[8]):
				count++
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 8 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	entries, err := entriesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read entries: %w", err)
	}

	sumOfOutputs := 0
	for i, e := range entries {
		output, err := outputValueFromEntry(e)
		if err != nil {
			return fmt.Errorf("problem with entry %d: %w", i+1, err)
		}
		sumOfOutputs += output
	}

	_, err = fmt.Fprintf(answer, "%d", sumOfOutputs)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const (
	patternsLen = 10
	outputLen   = 4
	numSegments = 7
	numDigits   = 10
)

type entry struct {
	patterns [patternsLen]string
	output   [outputLen]string
}

const (
	topSegment int = iota
	topLeftSegment
	topRightSegment
	middleSegment
	bottomLeftSegment
	bottomRightSegment
	bottomSegment
)

type digit [numSegments]bool

var digits = [numDigits]digit{
	{true, true, true, false, true, true, true},     // 0
	{false, false, true, false, false, true, false}, // 1
	{true, false, true, true, true, false, true},    // 2
	{true, false, true, true, false, true, true},    // 3
	{false, true, true, true, false, true, false},   // 4
	{true, true, false, true, false, true, true},    // 5
	{true, true, false, true, true, true, true},     // 6
	{true, false, true, false, false, true, false},  // 7
	{true, true, true, true, true, true, true},      // 8
	{true, true, true, true, false, true, true},     // 9
}

var workingPatterns = [numDigits]string{
	"abcefg",  // 0
	"cf",      // 1
	"acdeg",   // 2
	"acdfg",   // 3
	"bcdf",    // 4
	"abdfg",   // 5
	"abdefg",  // 6
	"acf",     // 7
	"abcdefg", // 8
	"abcdfg",  // 9
}
var workingSegmentOccurences = countSegmentOccurences(workingPatterns)

type segmentOccurences struct {
	in1478, in023569 int
}

func outputValueFromEntry(e entry) (int, error) {
	rewiredSegmentOccurences := countSegmentOccurences(e.patterns)

	workingSegmentToRewiredSegment, err := mapWorkingSegmentsToRewiredSegments(workingSegmentOccurences, rewiredSegmentOccurences)
	if err != nil {
		return 0, fmt.Errorf("could not map target segments to rewired segments: %w", err)
	}

	rewiredPatternToDigit := mapRewiredPatternsToDigits(workingSegmentToRewiredSegment)

	outputValue, multiplier := 0, 1
	for i := outputLen - 1; i >= 0; i-- {
		outputValue += rewiredPatternToDigit[e.output[i]] * multiplier
		multiplier *= 10
	}

	return outputValue, nil
}

func countSegmentOccurences(patterns [patternsLen]string) [numSegments]segmentOccurences {
	var occurences [numSegments]segmentOccurences
	for _, pattern := range patterns {
		for _, segment := range pattern {
			segmentID := int(segment - 'a')
			switch len(pattern) {
			case
				len(workingPatterns[1]),
				len(workingPatterns[4]),
				len(workingPatterns[7]),
				len(workingPatterns[8]):
				occurences[segmentID].in1478++
			default:
				occurences[segmentID].in023569++
			}
		}
	}
	return occurences
}

func mapWorkingSegmentsToRewiredSegments(workingSegmentOccurences, rewiredSegmentOccurences [numSegments]segmentOccurences) ([numSegments]int, error) {
	var mapping [numSegments]int

	for i, working := range workingSegmentOccurences {
		found := false
		for j, rewired := range rewiredSegmentOccurences {
			if working == rewired {
				mapping[i] = j
				found = true
				break
			}
		}
		if !found {
			return [numSegments]int{}, fmt.Errorf("no rewired segment matches working segment %q", byte(i)+'a')
		}
	}

	return mapping, nil
}

func mapRewiredPatternsToDigits(workingSegmentToRewiredSegment [numSegments]int) map[string]int {
	mapping := make(map[string]int)
	for digit, workingPattern := range workingPatterns {
		var rewiredPattern []byte
		for _, workingSegment := range workingPattern {
			rewiredSegment := byte(workingSegmentToRewiredSegment[workingSegment-'a'] + 'a')
			rewiredPattern = append(rewiredPattern, rewiredSegment)
		}
		sort.Sort(inOrder(rewiredPattern))
		mapping[string(rewiredPattern)] = digit
	}
	return mapping
}

func entriesFromReader(r io.Reader) ([]entry, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	entries := make([]entry, len(lines))
	for i := range lines {
		e, err := entryFromString(lines[i])
		if err != nil {
			return nil, fmt.Errorf("invalid entry at line %d: %w", i+1, err)
		}
		entries[i] = e
	}

	return entries, nil
}

func entryFromString(s string) (entry, error) {
	parts := bytes.Split([]byte(s), []byte(" | "))
	if len(parts) != 2 {
		return entry{}, errors.New("wrong format")
	}

	patterns := bytes.Split(parts[0], []byte(" "))
	if len(patterns) != patternsLen {
		return entry{}, fmt.Errorf("expected %d signal patterns, found %d", patternsLen, len(patterns))
	}
	for _, p := range patterns {
		sort.Sort(inOrder(p))
	}

	output := bytes.Split(parts[1], []byte(" "))
	if len(output) != outputLen {
		return entry{}, fmt.Errorf("expected %d digits in output, found %d", outputLen, len(output))
	}
	for _, o := range output {
		sort.Sort(inOrder(o))
	}

	var e entry
	for i := range e.patterns {
		e.patterns[i] = string(patterns[i])
	}
	for i := range e.output {
		e.output[i] = string(output[i])
	}

	return e, nil
}

type inOrder []byte

func (s inOrder) Len() int           { return len(s) }
func (s inOrder) Less(i, j int) bool { return s[i] < s[j] }
func (s inOrder) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
