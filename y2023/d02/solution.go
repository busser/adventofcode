package d02

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 2 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	games, err := gamesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	const (
		maxRed   = 12
		maxGreen = 13
		maxBlue  = 14
	)

	sum := 0
	for _, game := range games {
		possible := true
		for _, sample := range game.samples {
			if sample.red > maxRed || sample.green > maxGreen || sample.blue > maxBlue {
				possible = false
				break
			}
		}

		if possible {
			sum += game.id
		}
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 2 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	games, err := gamesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for _, game := range games {
		sum += game.power()
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type ballSample struct {
	blue, red, green int
}

type game struct {
	id      int
	samples []ballSample
}

func (g game) power() int {
	var maxRed, maxGreen, maxBlue int

	for _, sample := range g.samples {
		maxRed = max(maxRed, sample.red)
		maxGreen = max(maxGreen, sample.green)
		maxBlue = max(maxBlue, sample.blue)
	}

	return maxRed * maxGreen * maxBlue
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func gamesFromReader(r io.Reader) ([]game, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	games := make([]game, len(lines))
	for i, line := range lines {
		game, err := gameFromString(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse game: %w", err)
		}

		games[i] = game
	}

	return games, nil
}

func gameFromString(s string) (game, error) {
	parts := strings.Split(s, ": ")

	rawGameID := strings.TrimPrefix(parts[0], "Game ")
	gameID, err := strconv.Atoi(rawGameID)
	if err != nil {
		return game{}, fmt.Errorf("could not parse game id: %w", err)
	}

	samples, err := ballSamplesFromString(parts[1])
	if err != nil {
		return game{}, fmt.Errorf("could not parse samples: %w", err)
	}

	return game{
		id:      gameID,
		samples: samples,
	}, nil
}

func ballSamplesFromString(s string) ([]ballSample, error) {
	parts := strings.Split(s, "; ")

	samples := make([]ballSample, len(parts))
	for i, part := range parts {
		sample, err := ballSampleFromString(part)
		if err != nil {
			return nil, fmt.Errorf("could not parse sample: %w", err)
		}

		samples[i] = sample
	}

	return samples, nil
}

func ballSampleFromString(s string) (ballSample, error) {
	parts := strings.Split(s, ", ")

	var blue, red, green int

	for _, part := range parts {
		halves := strings.SplitN(part, " ", 2)
		if len(halves) != 2 {
			return ballSample{}, fmt.Errorf("invalid ball sample: %q", s)
		}

		amount, err := strconv.Atoi(halves[0])
		if err != nil {
			return ballSample{}, fmt.Errorf("could not parse amount: %w", err)
		}

		color := halves[1]
		switch color {
		case "blue":
			blue = amount
		case "red":
			red = amount
		case "green":
			green = amount
		default:
			return ballSample{}, fmt.Errorf("invalid color: %q", color)
		}
	}

	return ballSample{
		blue:  blue,
		red:   red,
		green: green,
	}, nil
}
