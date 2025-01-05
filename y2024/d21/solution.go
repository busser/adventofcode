package d21

import (
	"fmt"
	"io"
	"math"
	"slices"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 21 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	codes, err := codesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	chain := newKeypadChain(3)

	total := 0
	for _, code := range codes {
		total += chain.codeComplexity(code)
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 21 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	codes, err := codesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	chain := newKeypadChain(26)

	total := 0
	for _, code := range codes {
		total += chain.codeComplexity(code)
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type button byte

const (
	button0 button = '0'
	button1 button = '1'
	button2 button = '2'
	button3 button = '3'
	button4 button = '4'
	button5 button = '5'
	button6 button = '6'
	button7 button = '7'
	button8 button = '8'
	button9 button = '9'

	buttonLeft  button = '<'
	buttonRight button = '>'
	buttonUp    button = '^'
	buttonDown  button = 'v'

	buttonA button = 'A'

	keypadGap button = 'X'
)

type vector struct {
	row, col int
}

var (
	up    = vector{row: -1, col: 0}
	down  = vector{row: 1, col: 0}
	left  = vector{row: 0, col: -1}
	right = vector{row: 0, col: 1}
)

func (v vector) plus(w vector) vector {
	return vector{
		row: v.row + w.row,
		col: v.col + w.col,
	}
}

type keypadChain struct {
	robots int
	numpad keypad
	dirpad keypad
	cache  map[cacheKey]int
}

type cacheKey struct {
	from, to button
	robot    int
}

func newKeypadChain(robots int) keypadChain {
	return keypadChain{
		robots: robots,
		numpad: newNumericKeypad(),
		dirpad: newDirectionalKeypad(),
		cache:  make(map[cacheKey]int),
	}
}

func (kc keypadChain) codeComplexity(code []button) int {
	numericPart := codeNumericPart(code)
	codeLength := kc.shortestMetaCode(kc.numpad, code, kc.robots)

	complexity := numericPart * codeLength

	return complexity
}

func (kc keypadChain) shortestMetaCode(k keypad, code []button, robots int) int {
	if robots == 0 {
		return len(code)
	}

	current := buttonA

	total := 0
	for _, next := range code {
		total += kc.shortestCodeForButtonPress(k, current, next, robots)
		current = next
	}

	return total
}

func (kc keypadChain) shortestCodeForButtonPress(k keypad, from, to button, robots int) (length int) {
	key := cacheKey{from, to, robots}
	if cached, hit := kc.cache[key]; hit {
		return cached
	}
	defer func() { kc.cache[key] = length }()

	type searchState struct {
		position vector
		code     []button
	}

	var toVisit, nextToVisit []searchState
	toVisit = append(toVisit, searchState{k.layout[from], nil})

	targetPosition := k.layout[to]
	shortestPath := math.MaxInt

	for len(toVisit) > 0 {
		for _, state := range toVisit {
			if state.position == targetPosition {
				code := append(slices.Clone(state.code), buttonA)
				pathLength := kc.shortestMetaCode(kc.dirpad, code, robots-1)
				shortestPath = min(shortestPath, pathLength)
				continue
			}

			if state.position == k.gap {
				continue
			}

			if state.position.row < targetPosition.row {
				nextToVisit = append(nextToVisit, searchState{
					state.position.plus(down),
					append(slices.Clone(state.code), buttonDown),
				})
			}
			if state.position.row > targetPosition.row {
				nextToVisit = append(nextToVisit, searchState{
					state.position.plus(up),
					append(slices.Clone(state.code), buttonUp),
				})
			}
			if state.position.col < targetPosition.col {
				nextToVisit = append(nextToVisit, searchState{
					state.position.plus(right),
					append(slices.Clone(state.code), buttonRight),
				})
			}
			if state.position.col > targetPosition.col {
				nextToVisit = append(nextToVisit, searchState{
					state.position.plus(left),
					append(slices.Clone(state.code), buttonLeft),
				})
			}
		}

		toVisit, nextToVisit = nextToVisit, toVisit[:0]
	}

	return shortestPath
}

func codeNumericPart(code []button) int {
	numericPart := 0
	for _, b := range code {
		if b >= button0 && b <= button9 {
			numericPart *= 10
			numericPart += int(b - button0)
		}
	}
	return numericPart
}

type keypad struct {
	layout map[button]vector
	gap    vector
}

func newNumericKeypad() keypad {
	return keypad{
		layout: map[button]vector{
			button7: {0, 0}, button8: {0, 1}, button9: {0, 2},
			button4: {1, 0}, button5: {1, 1}, button6: {1, 2},
			button1: {2, 0}, button2: {2, 1}, button3: {2, 2},
			keypadGap: {3, 0}, button0: {3, 1}, buttonA: {3, 2},
		},
		gap: vector{3, 0},
	}
}

func newDirectionalKeypad() keypad {
	return keypad{
		layout: map[button]vector{
			keypadGap: {0, 0}, buttonUp: {0, 1}, buttonA: {0, 2},
			buttonLeft: {1, 0}, buttonDown: {1, 1}, buttonRight: {1, 2},
		},
		gap: vector{0, 0},
	}
}

func codesFromReader(r io.Reader) ([][]button, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	codes := make([][]button, len(lines))
	for i, line := range lines {
		codes[i] = []button(line)
	}

	return codes, nil
}
