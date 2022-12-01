package busser

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 21 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	foods, err := foodsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_ = removeAllergens(foods)

	count := 0
	for _, f := range foods {
		count += len(f.ingredients)
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 21 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	foods, err := foodsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	allergenByIngredient := removeAllergens(foods)

	var dangerousIngredients []string
	for ing := range allergenByIngredient {
		dangerousIngredients = append(dangerousIngredients, ing)
	}
	sort.Slice(
		dangerousIngredients,
		func(i, j int) bool {
			ingA, ingB := dangerousIngredients[i], dangerousIngredients[j]
			return allergenByIngredient[ingA] < allergenByIngredient[ingB]
		},
	)

	_, err = fmt.Fprintf(answer, "%s", strings.Join(dangerousIngredients, ","))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type food struct {
	ingredients set
	allergens   set
}

func foodsFromReader(r io.Reader) ([]food, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	foods := make([]food, len(lines))

	for i := range lines {
		f, err := foodFromString(lines[i])
		if err != nil {
			return nil, fmt.Errorf("parsing line #%d: %w", i, err)
		}
		foods[i] = f
	}

	return foods, nil
}

func foodFromString(s string) (food, error) {

	parts := strings.Split(s, " (contains ")
	if len(parts) != 2 {
		return food{}, errors.New("wrong format")
	}
	rawIngredients, rawAllergens := parts[0], parts[1]
	if len(rawIngredients) == 0 || len(rawAllergens) == 0 {
		return food{}, errors.New("wrong format")
	}
	if rawAllergens[len(rawAllergens)-1] != ')' {
		return food{}, errors.New("wrong format")
	}

	f := food{
		ingredients: newSet(strings.Split(rawIngredients, " ")...),
		allergens:   newSet(strings.Split(rawAllergens[:len(rawAllergens)-1], ", ")...),
	}

	return f, nil
}

func removeAllergens(foods []food) map[string]string {
	allergenByIngredient := make(map[string]string)

	unmatchedAllergens := allAllergens(foods)

	for {
		foundAllergenicIngredient := false

		for allergen := range unmatchedAllergens {
			candidateIngredients := commonIngredients(foodsWithAllergen(foods, allergen))
			if len(candidateIngredients) == 1 {
				var ingredient string
				for ing := range candidateIngredients {
					ingredient = ing
				}

				allergenByIngredient[ingredient] = allergen

				// Remove matched ingredient and allergens from menu
				for _, f := range foods {
					f.ingredients.remove(ingredient)
					f.allergens.remove(allergen)
					unmatchedAllergens.remove(allergen)
				}

				foundAllergenicIngredient = true
				break
			}
		}

		if !foundAllergenicIngredient {
			break
		}
	}

	return allergenByIngredient
}

func allAllergens(foods []food) set {
	allergenSets := make([]set, len(foods))
	for i := range foods {
		allergenSets[i] = foods[i].allergens
	}
	all := setUnion(allergenSets...)
	return all
}

func foodsWithAllergen(allFoods []food, allergen string) []food {
	var foods []food
	for _, f := range allFoods {
		if f.allergens.contains(allergen) {
			foods = append(foods, f)
		}
	}
	return foods
}

func commonIngredients(foods []food) set {
	ingredientSets := make([]set, len(foods))
	for i := range foods {
		ingredientSets[i] = foods[i].ingredients
	}
	common := setIntersection(ingredientSets...)
	return common
}

type set map[string]struct{}

func newSet(values ...string) set {
	s := make(set)
	for _, v := range values {
		s.add(v)
	}
	return s
}

func (s set) add(v string) {
	s[v] = struct{}{}
}

func (s set) remove(v string) {
	delete(s, v)
}

func (s set) contains(v string) bool {
	_, ok := s[v]
	return ok
}

func setUnion(sets ...set) set {
	union := newSet()
	for _, s := range sets {
		for v := range s {
			union.add(v)
		}
	}
	return union
}

func setIntersection(sets ...set) set {
	if len(sets) == 0 {
		return newSet()
	}

	intersection := newSet()

	for v := range sets[0] {
		allSetsContain := true
		for _, s := range sets {
			if !s.contains(v) {
				allSetsContain = false
				break
			}
		}
		if allSetsContain {
			intersection.add(v)
		}
	}

	return intersection
}
