package d04

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 4 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	cards, err := cardsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for _, c := range cards {
		sum += c.pointWorth()
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 4 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	cards, err := cardsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	total := countTotalCards(cards)

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type card struct {
	id             int
	winningNumbers []int
	numbersYouHave []int
}

func (c card) countMatchingNumbers() int {
	count := 0

	for _, n := range c.winningNumbers {
		for _, m := range c.numbersYouHave {
			if n == m {
				count++
			}
		}
	}

	return count
}

func (c card) pointWorth() int {
	matchingNumbers := c.countMatchingNumbers()

	if matchingNumbers == 0 {
		return 0
	}

	points := 1
	for i := 1; i < matchingNumbers; i++ {
		points *= 2
	}

	return points
}

func countTotalCards(cards []card) int {
	cardCount := make([]int, len(cards))
	for i := range cards {
		cardCount[i] = 1
	}

	for i := range cards {
		matchingNumbers := cards[i].countMatchingNumbers()
		for j := i + 1; j < i+matchingNumbers+1; j++ {
			cardCount[j] += cardCount[i]
		}
	}

	total := 0
	for _, c := range cardCount {
		total += c
	}

	return total
}

func cardsFromReader(r io.Reader) ([]card, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	cards := make([]card, len(lines))
	for i, line := range lines {
		c, err := cardFromString(line)
		if err != nil {
			return nil, fmt.Errorf("could not read card: %w", err)
		}

		cards[i] = c
	}

	return cards, nil
}

func cardFromString(s string) (card, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return card{}, fmt.Errorf("invalid card: %q", s)
	}

	id, err := idFromHeader(parts[0])
	if err != nil {
		return card{}, fmt.Errorf("could not read header: %w", err)
	}

	numberLists := strings.SplitN(parts[1], " | ", 2)
	if len(numberLists) != 2 {
		return card{}, fmt.Errorf("invalid card: %q", s)
	}

	winningNumbers, err := numbersFromString(numberLists[0])
	if err != nil {
		return card{}, fmt.Errorf("could not read winning numbers: %w", err)
	}

	numbersYouHave, err := numbersFromString(numberLists[1])
	if err != nil {
		return card{}, fmt.Errorf("could not read numbers you have: %w", err)
	}

	return card{
		id:             id,
		winningNumbers: winningNumbers,
		numbersYouHave: numbersYouHave,
	}, nil
}

func idFromHeader(s string) (int, error) {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid header: %q", s)
	}

	return strconv.Atoi(parts[1])
}

func numbersFromString(s string) ([]int, error) {
	parts := strings.Fields(s)

	numbers := make([]int, len(parts))
	for i, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %q", part)
		}

		numbers[i] = n
	}

	return numbers, nil
}
