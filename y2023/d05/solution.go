package d05

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 5 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	almanac, err := almanacFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	intervals := make([]interval, len(almanac.seeds))
	for i, seed := range almanac.seeds {
		// using intervals of 1 to share code with PartTwo
		intervals[i] = interval{
			start: seed,
			end:   seed + 1,
		}
	}

	final := almanac.convert(intervals)

	_, err = fmt.Fprintf(w, "%d", final[0].start)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	almanac, err := almanacFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	intervals := make([]interval, len(almanac.seeds)/2)
	for i := range intervals {
		start := almanac.seeds[2*i]
		length := almanac.seeds[2*i+1]
		intervals[i] = interval{
			start: start,
			end:   start + length,
		}
	}

	final := almanac.convert(intervals)

	_, err = fmt.Fprintf(w, "%d", final[0].start)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type interval struct {
	start int // inclusive
	end   int // exclusive
}

func (i interval) less(j interval) bool {
	if i.start == j.start {
		return i.end < j.end
	}
	return i.start < j.start
}

func sortIntervals(intervals []interval) {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].less(intervals[j])
	})
}

func mergeIntervals(intervals []interval) []interval {
	sortIntervals(intervals)

	var merged []interval
	for _, i := range intervals {
		if len(merged) == 0 || merged[len(merged)-1].end < i.start {
			merged = append(merged, i)
		} else if merged[len(merged)-1].end < i.end {
			merged[len(merged)-1].end = i.end
		}
	}

	return merged
}

var (
	categories = [...]string{
		"seed",
		"soil",
		"fertilizer",
		"water",
		"light",
		"temperature",
		"humidity",
		"location",
	}
)

type almanacRange struct {
	match interval
	shift int
}

type almanacMap struct {
	sourceCategory      string
	destinationCatogory string
	ranges              []almanacRange
}

func (m almanacMap) convert(intervals []interval) []interval {
	var converted []interval

	// This algorithm takes advantage of the fact that we preprocessed the map
	// ranges so that there are no gaps before, between, or after them. Each
	// interval is guaranteed to fit into one or more ranges.

	// We keep track of which interval and range we are at.
	// We will move forward with either one of the indices at each step.
	i, r := 0, 0

	for i < len(intervals) {
		// The interval can't end before the current range starts. This is one
		// of the algorithm's invariants.
		if intervals[i].end <= m.ranges[r].match.start {
			panic("current interval ends before current range starts")
		}

		// The interval can't start before the current range. This is one of the
		// algorithm's invariants.
		if intervals[i].start < m.ranges[r].match.start {
			panic("current interval starts before current range starts")
		}

		// If the interval ends within the range, we map the entire interval
		// based on the range's shift value, then move on to the next interval.
		if intervals[i].end <= m.ranges[r].match.end {
			shifted := interval{
				start: intervals[i].start + m.ranges[r].shift,
				end:   intervals[i].end + m.ranges[r].shift,
			}
			converted = append(converted, shifted)
			i++
			continue
		}

		// If the interval begins after the current range ends, we move on to
		// the next range. This can happen because of gaps between intervals.
		if intervals[i].start >= m.ranges[r].match.end {
			r++
			continue
		}

		// If the interval extends beyond the range, we shift the part that is
		// within the range and move on to the next range with the remainder.
		if intervals[i].end > m.ranges[r].match.end {
			within := interval{
				start: intervals[i].start,
				end:   m.ranges[r].match.end,
			}
			remainder := interval{
				start: m.ranges[r].match.end,
				end:   intervals[i].end,
			}

			shifted := interval{
				start: within.start + m.ranges[r].shift,
				end:   within.end + m.ranges[r].shift,
			}
			converted = append(converted, shifted)

			intervals[i] = remainder
			r++
			continue
		}

		panic("unhandled case")
	}

	merged := mergeIntervals(converted)

	return merged
}

type almanac struct {
	seeds []int
	maps  []almanacMap
}

func (a almanac) convert(intervals []interval) []interval {
	intervals = mergeIntervals(intervals)
	for _, m := range a.maps {
		intervals = m.convert(intervals)
	}
	return intervals
}

func almanacFromReader(r io.Reader) (*almanac, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	chunks := splitSlice(lines, func(s string) bool {
		return s == ""
	})

	if len(chunks) != len(categories) {
		return nil, fmt.Errorf("invalid input: expected %d categories, got %d", len(categories), len(chunks))
	}
	if len(chunks[0]) != 1 {
		return nil, fmt.Errorf("invalid input: expected 1 line of seeds, got %d", len(chunks[0]))
	}

	seeds, err := seedsFromString(chunks[0][0])
	if err != nil {
		return nil, fmt.Errorf("could not parse seeds: %w", err)
	}
	if len(seeds)%2 != 0 {
		return nil, fmt.Errorf("invalid input: expected even number of seeds, got %d", len(seeds))
	}

	var maps []almanacMap
	for i, chunk := range chunks[1:] {
		m, err := almanacMapFromStrings(chunk)
		if err != nil {
			return nil, fmt.Errorf("could not parse map for category %s: %w", categories[i], err)
		}
		if m.sourceCategory != categories[i] {
			return nil, fmt.Errorf("invalid input: expected category %s, got %s", categories[i], m.sourceCategory)
		}
		if m.destinationCatogory != categories[i+1] {
			return nil, fmt.Errorf("invalid input: expected category %s, got %s", categories[i+1], m.destinationCatogory)
		}
		maps = append(maps, m)
	}

	return &almanac{
		seeds: seeds,
		maps:  maps,
	}, nil
}

func seedsFromString(s string) ([]int, error) {
	s = strings.TrimPrefix(s, "seeds: ")
	return helpers.IntsFromString(s, " ")
}

func almanacMapFromStrings(s []string) (almanacMap, error) {
	if len(s) < 2 {
		return almanacMap{}, fmt.Errorf("invalid input: expected at least 2 lines, got %d", len(s))
	}

	source, destination, err := categoriesFromMapHeader(s[0])
	if err != nil {
		return almanacMap{}, fmt.Errorf("could not parse map header: %w", err)
	}

	var ranges []almanacRange
	for _, line := range s[1:] {
		nums, err := helpers.IntsFromString(line, " ")
		if err != nil {
			return almanacMap{}, fmt.Errorf("could not parse map line: %w", err)
		}
		if len(nums) != 3 {
			return almanacMap{}, fmt.Errorf("invalid input: expected 3 numbers, got %d", len(nums))
		}

		destinationStart := nums[0]
		sourceStart := nums[1]
		length := nums[2]

		ranges = append(ranges, almanacRange{
			match: interval{
				start: sourceStart,
				end:   sourceStart + length,
			},
			shift: destinationStart - sourceStart,
		})
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].match.less(ranges[j].match)
	})

	// We fill in the gaps between the ranges with dummy ranges with a shift
	// value of zero. This allows us to implement a simpler algorithm that
	// handles the default case where a source number should be mapped to the
	// same destination number.
	//
	// NOTE: the ranges in the input data don't have gaps between them but we'll
	// assume that they might since this property isn't mentioned in the problem
	// statement.

	var fullRanges []almanacRange

	start := math.MinInt

	for i := range ranges {
		if ranges[i].match.start > start {
			fullRanges = append(fullRanges, almanacRange{
				match: interval{
					start: start,
					end:   ranges[i].match.start,
				},
				shift: 0,
			})
		}

		fullRanges = append(fullRanges, ranges[i])
		start = ranges[i].match.end
	}

	fullRanges = append(fullRanges, almanacRange{
		match: interval{
			start: start,
			end:   math.MaxInt,
		},
		shift: 0,
	})

	return almanacMap{
		sourceCategory:      source,
		destinationCatogory: destination,
		ranges:              fullRanges,
	}, nil
}

func categoriesFromMapHeader(s string) (string, string, error) {
	s = strings.TrimSuffix(s, " map:")
	parts := strings.Split(s, "-to-")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid input: expected 2 parts, got %d", len(parts))
	}
	return parts[0], parts[1], nil
}

func splitSlice(s []string, f func(string) bool) [][]string {
	var chunks [][]string
	var chunk []string
	for _, v := range s {
		if f(v) {
			chunks = append(chunks, chunk)
			chunk = nil
		} else {
			chunk = append(chunk, v)
		}
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}
