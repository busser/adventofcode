package d05

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 5 of Advent of Code 2025.
func PartOne(r io.Reader, w io.Writer) error {
	db, err := databaseFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := db.countFreshIngredients()

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2025.
func PartTwo(r io.Reader, w io.Writer) error {
	db, err := databaseFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := db.countPossibleFreshIngredients()

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type database struct {
	freshIngredientIDRanges []idRange
	availableIngredientIDs  []int
}

func (db *database) ingredientIsFresh(id int) bool {
	for _, r := range db.freshIngredientIDRanges {
		if r.contains(id) {
			return true
		}
	}
	return false
}

func (db *database) countFreshIngredients() int {
	count := 0
	for _, id := range db.availableIngredientIDs {
		if db.ingredientIsFresh(id) {
			count++
		}
	}
	return count
}

type rangeEdge struct {
	n     int
	begin bool // false if end
}

func sortRangeEdges(edges []rangeEdge) {
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].n < edges[j].n
	})
}

func collectRangeEdges(idRanges []idRange) []rangeEdge {
	edges := make([]rangeEdge, len(idRanges)*2)
	for i, r := range idRanges {
		edges[2*i] = rangeEdge{n: r.begin, begin: true}
		edges[2*i+1] = rangeEdge{n: r.end + 1, begin: false}
	}
	sortRangeEdges(edges)
	return edges
}

func (db *database) countPossibleFreshIngredients() int {
	count := 0

	idRangeEdges := collectRangeEdges(db.freshIngredientIDRanges)
	activeRanges, previousEdge := 0, 0
	for _, edge := range idRangeEdges {
		if activeRanges > 0 {
			count += edge.n - previousEdge
		}
		previousEdge = edge.n

		if edge.begin {
			activeRanges++
		} else {
			activeRanges--
		}

	}

	return count
}

type idRange struct {
	begin, end int
}

func (r idRange) contains(id int) bool {
	return r.begin <= id && id <= r.end
}

func databaseFromReader(r io.Reader) (*database, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	sections := strings.SplitN(string(input), "\n\n", 2)
	if len(sections) != 2 {
		return nil, fmt.Errorf("expected 2 sections, got %d", len(sections))
	}

	freshIngredientIDRanges, err := idRangesFromString(sections[0])
	if err != nil {
		return nil, err
	}

	availableIngredientIDs := helpers.IntsFromString(sections[1])

	db := database{
		freshIngredientIDRanges: freshIngredientIDRanges,
		availableIngredientIDs:  availableIngredientIDs,
	}

	return &db, nil
}

func idRangesFromString(s string) ([]idRange, error) {
	ints := helpers.IntsFromString(s)

	if len(ints)%2 != 0 {
		return nil, fmt.Errorf("wrong id range format")
	}

	idRanges := make([]idRange, len(ints)/2)
	for i := 0; i < len(ints)/2; i++ {
		idRanges[i] = idRange{begin: ints[i*2], end: ints[i*2+1]}
	}

	return idRanges, nil
}
