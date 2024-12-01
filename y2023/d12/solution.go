package d12

import (
	"fmt"
	"io"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 12 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	records, err := conditionRecordsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for _, record := range records {
		sum += record.countValidArrangements()
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 12 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	records, err := conditionRecordsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for _, record := range records {
		record.unfold()
		sum += record.countValidArrangements()
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const (
	operational = '.'
	damaged     = '#'
	unknown     = '?'
)

type conditionRecord struct {
	row    []byte
	groups []int
}

func (cr *conditionRecord) unfold() {
	newRow := make([]byte, 0, len(cr.row)*5+4)
	for i := 0; i < 5; i++ {
		if i > 0 {
			newRow = append(newRow, unknown)
		}
		newRow = append(newRow, cr.row...)
	}

	newGroups := make([]int, 0, len(cr.groups)*5)
	for i := 0; i < 5; i++ {
		newGroups = append(newGroups, cr.groups...)
	}

	cr.row = newRow
	cr.groups = newGroups
}

type cacheKey struct {
	springIndex      int
	firstSpring      byte
	groupCount       int
	currentGroupSize int
}

func sum(ints []int) int {
	sum := 0
	for _, i := range ints {
		sum += i
	}
	return sum
}

func (cr conditionRecord) countValidArrangements() int {
	// Having an operational spring at the end of the row makes it easier to
	// handle the last group.
	if cr.row[len(cr.row)-1] != operational {
		cr.row = append(cr.row, operational)
	}

	cache := make(map[cacheKey]int)

	var helper func(int, int, int) int
	helper = func(springIndex, groupCount, currentGroupSize int) (result int) {
		if springIndex == len(cr.row) {
			// We reached the end of the row.

			if groupCount != len(cr.groups) {
				// We found the wrong number of groups.
				return 0
			}

			// We found a valid arrangement.
			return 1
		}

		key := cacheKey{
			springIndex:      springIndex,
			firstSpring:      cr.row[springIndex],
			groupCount:       groupCount,
			currentGroupSize: currentGroupSize,
		}
		if cachedResult, hit := cache[key]; hit {
			return cachedResult
		}
		defer func() {
			cache[key] = result
		}()

		switch {
		case cr.row[springIndex] == operational:
			if currentGroupSize > 0 {
				// We just finished a group.
				if currentGroupSize != cr.groups[groupCount] {
					// The group we found is not the right size.
					return 0
				}
				// The group we found is the right size, move on to the next group.
				groupCount++
				currentGroupSize = 0

				remainingGroups := cr.groups[groupCount:]
				minLength := sum(remainingGroups) + len(remainingGroups) - 1
				if len(cr.row)-springIndex < minLength {
					// There are not enough springs left to form the remaining groups.
					return 0
				}
			}

			// Skip all the operational springs.
			for springIndex < len(cr.row) && cr.row[springIndex] == operational {
				springIndex++
			}

			return helper(springIndex, groupCount, currentGroupSize)

		case cr.row[springIndex] == damaged:
			if currentGroupSize == 0 {
				// We are starting a new group.
				if groupCount == len(cr.groups) {
					// This new group is one too much.
					return 0
				}
			}
			currentGroupSize++
			if currentGroupSize > cr.groups[groupCount] {
				// The group we found is too large.
				return 0
			}

			// Move on to the next spring.
			return helper(springIndex+1, groupCount, currentGroupSize)

		case cr.row[springIndex] == unknown:
			// Simulate the spring being operational.
			cr.row[springIndex] = operational
			ifOperational := helper(springIndex, groupCount, currentGroupSize)

			// Simulate the spring being damaged.
			cr.row[springIndex] = damaged
			ifDamaged := helper(springIndex, groupCount, currentGroupSize)

			// Reset the spring to unknow for future calls.
			cr.row[springIndex] = unknown

			return ifOperational + ifDamaged

		default:
			panic("unhandled case")
		}
	}

	return helper(0, 0, 0)
}

func conditionRecordsFromReader(r io.Reader) ([]conditionRecord, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	records := make([]conditionRecord, len(lines))
	for i, line := range lines {
		records[i], err = conditionRecordFromString(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse line %d: %w", i, err)
		}
	}

	return records, nil
}

func conditionRecordFromString(s string) (conditionRecord, error) {
	parts := strings.SplitN(s, " ", 2)
	if len(parts) != 2 {
		return conditionRecord{}, fmt.Errorf("invalid record: %s", s)
	}

	row, err := conditionRecordRowFromString(parts[0])
	if err != nil {
		return conditionRecord{}, fmt.Errorf("could not parse row: %w", err)
	}

	groups := helpers.IntsFromString(parts[1])

	return conditionRecord{
		row:    row,
		groups: groups,
	}, nil
}

func conditionRecordRowFromString(s string) ([]byte, error) {
	row := make([]byte, len(s))
	for i, c := range s {
		switch c {
		case operational:
			row[i] = operational
		case damaged:
			row[i] = damaged
		case unknown:
			row[i] = unknown
		default:
			return nil, fmt.Errorf("invalid row: %s", s)
		}
	}

	return row, nil
}
