package d02

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 2 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	guide, err := guideFromReader(r, false)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	totalScore := 0
	for _, r := range guide {
		r.determineOutcome()
		totalScore += r.score()
	}

	_, err = fmt.Fprintf(w, "%d", totalScore)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 2 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	guide, err := guideFromReader(r, true)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	totalScore := 0
	for _, r := range guide {
		r.determineYourShape()
		totalScore += r.score()
	}

	_, err = fmt.Fprintf(w, "%d", totalScore)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const (
	rock    = 1
	paper   = 2
	scissor = 3
)

const (
	youLose  = 0
	itsADraw = 3
	youWin   = 6
)

type round struct {
	opponent int
	you      int
	outcome  int
}

func (r *round) determineOutcome() {
	switch r.you - r.opponent {
	case -1, 2:
		r.outcome = youLose
	case 0:
		r.outcome = itsADraw
	case -2, 1:
		r.outcome = youWin
	}
}

func (r *round) determineYourShape() {
	switch r.outcome {
	case youLose:
		r.you = r.opponent + 2
	case itsADraw:
		r.you = r.opponent
	case youWin:
		r.you = r.opponent + 1
	}

	if r.you > 3 {
		r.you -= 3
	}
}

func (r round) score() int {
	return int(r.outcome) + int(r.you)
}

func guideFromReader(r io.Reader, correct bool) ([]round, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	roundFromString := incorrectRoundFromString
	if correct {
		roundFromString = correctRoundFromString
	}

	rounds := make([]round, len(lines))
	for i, l := range lines {
		r, err := roundFromString(l)
		if err != nil {
			return nil, fmt.Errorf("round %q invalid: %w", l, err)
		}

		rounds[i] = r
	}

	return rounds, err
}

func incorrectRoundFromString(s string) (round, error) {
	if len(s) != 3 {
		return round{}, fmt.Errorf("wrong length")
	}

	var r round

	switch s[0] {
	case 'A':
		r.opponent = rock
	case 'B':
		r.opponent = paper
	case 'C':
		r.opponent = scissor
	default:
		return round{}, fmt.Errorf("unknown int %q", s[0])
	}

	switch s[2] {
	case 'X':
		r.you = rock
	case 'Y':
		r.you = paper
	case 'Z':
		r.you = scissor
	default:
		return round{}, fmt.Errorf("unknown int %q", s[2])
	}

	return r, nil
}

func correctRoundFromString(s string) (round, error) {
	if len(s) != 3 {
		return round{}, fmt.Errorf("wrong length")
	}

	var r round

	switch s[0] {
	case 'A':
		r.opponent = rock
	case 'B':
		r.opponent = paper
	case 'C':
		r.opponent = scissor
	default:
		return round{}, fmt.Errorf("unknown int %q", s[0])
	}

	switch s[2] {
	case 'X':
		r.outcome = youLose
	case 'Y':
		r.outcome = itsADraw
	case 'Z':
		r.outcome = youWin
	default:
		return round{}, fmt.Errorf("unknown int %q", s[2])
	}

	return r, nil
}
