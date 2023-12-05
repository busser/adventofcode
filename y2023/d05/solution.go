package d05

import (
	"fmt"
	"io"
	"log"
	"math"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 5 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	almanac, err := almanacFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	lowestLocation := math.MaxInt

	for _, seed := range almanac.seeds {
		location := almanac.convert(seed)
		if location < lowestLocation {
			lowestLocation = location
		}
	}

	_, err = fmt.Fprintf(w, "%d", lowestLocation)
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

	numSeedRanges := len(almanac.seeds) / 2

	lowestLocations := make(chan int, numSeedRanges)

	for i := 0; i < numSeedRanges; i++ {
		go func(i int) {
			startSeed := almanac.seeds[i*2]
			endSeed := startSeed + almanac.seeds[i*2+1]
			log.Printf("worker %d: startSeed: %d, endSeed: %d", i, startSeed, endSeed)

			lowestLocation := math.MaxInt
			for seed := startSeed; seed < endSeed; seed++ {
				location := almanac.convert(seed)
				if location < lowestLocation {
					lowestLocation = location
				}
			}

			log.Printf("worker %d: lowestLocation: %d", i, lowestLocation)
			lowestLocations <- lowestLocation
		}(i)
	}

	lowestLocation := math.MaxInt
	for i := 0; i < numSeedRanges; i++ {
		location := <-lowestLocations
		if location < lowestLocation {
			lowestLocation = location
		}
	}

	_, err = fmt.Fprintf(w, "%d", lowestLocation)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
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
	destinationStart int
	sourceStart      int
	length           int
}

type almanacMap struct {
	sourceCategory      string
	destinationCatogory string
	ranges              []almanacRange
}

func (m almanacMap) convert(n int) int {
	for _, r := range m.ranges {
		if n >= r.sourceStart && n < r.sourceStart+r.length {
			return r.destinationStart + (n - r.sourceStart)
		}
	}
	return n
}

type almanac struct {
	seeds []int
	maps  []almanacMap
}

func (a almanac) convert(n int) int {
	for _, m := range a.maps {
		n = m.convert(n)
	}
	return n
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
		ranges = append(ranges, almanacRange{
			destinationStart: nums[0],
			sourceStart:      nums[1],
			length:           nums[2],
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
