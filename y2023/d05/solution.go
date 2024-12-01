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

	seedRanges := make([]interval, len(almanac.seeds))
	for i, seed := range almanac.seeds {
		// using intervals of 1 to share code with PartTwo
		seedRanges[i] = interval{
			start: seed,
			end:   seed + 1,
		}
	}

	final := almanac.convert(seedRanges)

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

	seedRanges := make([]interval, len(almanac.seeds)/2)
	for i := range seedRanges {
		start := almanac.seeds[2*i]
		length := almanac.seeds[2*i+1]
		seedRanges[i] = interval{
			start: start,
			end:   start + length,
		}
	}

	final := almanac.convert(seedRanges)

	_, err = fmt.Fprintf(w, "%d", final[0].start)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

//=== Intervals: sorting and merging ===========================================

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

//=== Mapping ranges to new values =============================================

func (a almanac) convert(seedRanges []interval) []interval {
	farmRanges := mergeIntervals(seedRanges)
	for _, m := range a.maps {
		farmRanges = m.convert(farmRanges)
	}
	return farmRanges
}

func (m almanacMap) convert(farmRanges []interval) []interval {
	var converted []interval

	// This algorithm preprocesses the map  ranges so that there are no gaps
	// before, between, or after them. Each farm range is guaranteed to fit into
	// one or more map ranges. We also guarantee that the farm ranges are sorted
	// and non-overlapping, so that we can process them in order without ever
	// having to backtrack.

	sort.Slice(m.ranges, func(i, j int) bool {
		return m.ranges[i].match.less(m.ranges[j].match)
	})

	// We fill in the gaps between the ranges with dummy ranges with a shift
	// value of zero. This allows us to implement a simpler algorithm that
	// handles the default case where a source number should be mapped to the
	// same destination number.
	//
	// NOTE: the ranges in the input data don't have gaps between them but we'll
	// assume that they might since this property isn't mentioned in the problem
	// statement.

	var mapRanges []almanacRange

	start := math.MinInt

	for i := range m.ranges {
		if m.ranges[i].match.start > start {
			mapRanges = append(mapRanges, almanacRange{
				match: interval{
					start: start,
					end:   m.ranges[i].match.start,
				},
				shift: 0,
			})
		}

		mapRanges = append(mapRanges, m.ranges[i])
		start = m.ranges[i].match.end
	}

	mapRanges = append(mapRanges, almanacRange{
		match: interval{
			start: start,
			end:   math.MaxInt,
		},
		shift: 0,
	})

	// We keep track of which farm range and map range we are at.
	// We will move forward with either one of the indices at each step.
	fr, mr := 0, 0

	for fr < len(farmRanges) {
		// The current farm range can't end before the current map range starts.
		// This is one of the algorithm's invariants.
		if farmRanges[fr].end <= mapRanges[mr].match.start {
			panic("current interval ends before current range starts")
		}

		// The current farm range can't start before the current map range.
		// This is one of the algorithm's invariants.
		if farmRanges[fr].start < mapRanges[mr].match.start {
			panic("current interval starts before current range starts")
		}

		// If the farm range ends within the map range, we shift the entire farm
		// range and move on the next farm range.
		if farmRanges[fr].end <= mapRanges[mr].match.end {
			shifted := interval{
				start: farmRanges[fr].start + mapRanges[mr].shift,
				end:   farmRanges[fr].end + mapRanges[mr].shift,
			}
			converted = append(converted, shifted)
			fr++
			continue
		}

		// If the farm range begins after the current map range ends, we move on
		// to the next map range. This can happen because of gaps between farm
		// ranges.
		if farmRanges[fr].start >= mapRanges[mr].match.end {
			mr++
			continue
		}

		// If the farm range extends beyond the map range, we shift the part of
		// the farm range that is within the map range and move on to the next
		// map range with what remains of the farm range.
		if farmRanges[fr].end > mapRanges[mr].match.end {
			within := interval{
				start: farmRanges[fr].start,
				end:   mapRanges[mr].match.end,
			}
			remainder := interval{
				start: mapRanges[mr].match.end,
				end:   farmRanges[fr].end,
			}

			shifted := interval{
				start: within.start + mapRanges[mr].shift,
				end:   within.end + mapRanges[mr].shift,
			}
			converted = append(converted, shifted)

			farmRanges[fr] = remainder
			mr++
			continue
		}

		panic("unhandled case")
	}

	merged := mergeIntervals(converted)

	return merged
}

//=== Almanac: definition and parsing ==========================================

type almanac struct {
	seeds []int
	maps  []almanacMap
}

type almanacRange struct {
	match interval
	shift int
}

type almanacMap struct {
	sourceCategory      string
	destinationCatogory string
	ranges              []almanacRange
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

	seeds := seedsFromString(chunks[0][0])
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

func seedsFromString(s string) []int {
	s = strings.TrimPrefix(s, "seeds: ")
	return helpers.IntsFromString(s)
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
		nums := helpers.IntsFromString(line)
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

	return almanacMap{
		sourceCategory:      source,
		destinationCatogory: destination,
		ranges:              ranges,
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
