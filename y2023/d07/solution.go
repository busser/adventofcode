package d07

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 7 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	hands, err := handsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	for i := range hands {
		hands[i].evaluateCardDefault()
	}

	sort.Slice(hands, func(i, j int) bool {
		return hands[i].less(hands[j])
	})

	totalWinnings := 0
	for i, h := range hands {
		rank := i + 1
		totalWinnings += h.bid * rank
	}

	_, err = fmt.Fprintf(w, "%d", totalWinnings)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 7 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	hands, err := handsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	for i := range hands {
		hands[i].evaluateCardJoker()
	}

	sort.Slice(hands, func(i, j int) bool {
		return hands[i].less(hands[j])
	})

	totalWinnings := 0
	for i, h := range hands {
		rank := i + 1
		totalWinnings += h.bid * rank
	}

	_, err = fmt.Fprintf(w, "%d", totalWinnings)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const handSize = 5

type hand struct {
	cards         [handSize]byte
	cardStrengths [handSize]int
	bid           int
	typ           handType
}

func (h hand) less(other hand) bool {
	if h.typ != other.typ {
		return h.typ < other.typ
	}

	for i := range h.cardStrengths {
		if h.cardStrengths[i] != other.cardStrengths[i] {
			return h.cardStrengths[i] < other.cardStrengths[i]
		}
	}

	return false
}

func (h *hand) evaluateCardDefault() {
	for i, card := range h.cards {
		h.cardStrengths[i] = cardStrengthDefault(card)
	}

	var cardCounts [numCards]int
	for _, card := range h.cardStrengths {
		cardCounts[card]++
	}

	var pairs, triples, quadruples, quintuples int
	for _, count := range cardCounts {
		switch count {
		case 2:
			pairs++
		case 3:
			triples++
		case 4:
			quadruples++
		case 5:
			quintuples++
		}
	}

	switch {
	case quintuples == 1:
		h.typ = fiveOfAKind
	case quadruples == 1:
		h.typ = fourOfAKind
	case triples == 1 && pairs == 1:
		h.typ = fullHouse
	case triples == 1:
		h.typ = threeOfAKind
	case pairs == 2:
		h.typ = twoPairs
	case pairs == 1:
		h.typ = onePair
	default:
		h.typ = highCard
	}
}

func (h *hand) evaluateCardJoker() {
	for i, card := range h.cards {
		h.cardStrengths[i] = cardStrengthJoker(card)
	}

	var cardCounts [numCards]int
	for _, card := range h.cardStrengths {
		cardCounts[card]++
	}

	var jokers, pairs, triples, quadruples, quintuples int
	jokers = cardCounts[0]
	for _, count := range cardCounts[1:] {
		switch count {
		case 2:
			pairs++
		case 3:
			triples++
		case 4:
			quadruples++
		case 5:
			quintuples++
		}
	}

	var typeWithoutJokers handType
	switch {
	case quintuples == 1:
		typeWithoutJokers = fiveOfAKind
	case quadruples == 1:
		typeWithoutJokers = fourOfAKind
	case triples == 1 && pairs == 1:
		typeWithoutJokers = fullHouse
	case triples == 1:
		typeWithoutJokers = threeOfAKind
	case pairs == 2:
		typeWithoutJokers = twoPairs
	case pairs == 1:
		typeWithoutJokers = onePair
	default:
		typeWithoutJokers = highCard
	}

	switch jokers {
	case 0:
		h.typ = typeWithoutJokers
	case 1:
		switch typeWithoutJokers {
		case highCard:
			h.typ = onePair
		case onePair:
			h.typ = threeOfAKind
		case twoPairs:
			h.typ = fullHouse
		case threeOfAKind:
			h.typ = fourOfAKind
		case fullHouse:
			h.typ = fourOfAKind
		case fourOfAKind:
			h.typ = fiveOfAKind
		default:
			panic("too many cards")
		}
	case 2:
		switch typeWithoutJokers {
		case highCard:
			h.typ = threeOfAKind
		case onePair:
			h.typ = fourOfAKind
		case threeOfAKind:
			h.typ = fiveOfAKind
		default:
			panic("too many cards")
		}
	case 3:
		switch typeWithoutJokers {
		case highCard:
			h.typ = fourOfAKind
		case onePair:
			h.typ = fiveOfAKind
		default:
			panic("too many cards")
		}
	case 4:
		h.typ = fiveOfAKind
	case 5:
		h.typ = fiveOfAKind
	default:
		panic("more than 5 jokers, there's a bug")
	}
}

var (
	cardSymbols      = [...]byte{'2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K', 'A'}
	cardOrderDefault = cardSymbols
	cardOrderJoker   = [...]byte{'J', '2', '3', '4', '5', '6', '7', '8', '9', 'T', 'Q', 'K', 'A'}
)

const numCards = len(cardSymbols)

func cardStrengthDefault(c byte) int {
	for i, card := range cardOrderDefault {
		if card == c {
			return i
		}
	}
	panic(fmt.Sprintf("unknown card %q", c))
}

func cardStrengthJoker(c byte) int {
	for i, card := range cardOrderJoker {
		if card == c {
			return i
		}
	}
	return -1
}

type handType int

const (
	highCard handType = iota
	onePair
	twoPairs
	threeOfAKind
	fullHouse
	fourOfAKind
	fiveOfAKind
)

func handsFromReader(r io.Reader) ([]hand, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	hands := make([]hand, 0, len(lines))
	for _, line := range lines {
		h, err := handFromLine(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse hand: %w", err)
		}
		hands = append(hands, h)
	}

	return hands, nil
}

func handFromLine(line string) (hand, error) {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return hand{}, fmt.Errorf("invalid hand")
	}
	if len(parts[0]) != handSize {
		return hand{}, fmt.Errorf("invalid hand")
	}

	for i := range parts[0] {
		cardExists := false
		for _, known := range cardSymbols {
			if parts[0][i] == known {
				cardExists = true
				break
			}
		}
		if !cardExists {
			return hand{}, fmt.Errorf("unknown card %q", parts[0][i])
		}
	}

	var cards [handSize]byte
	copy(cards[:], parts[0])

	bid, err := strconv.Atoi(parts[1])
	if err != nil {
		return hand{}, fmt.Errorf("invalid hand")
	}

	return hand{
		cards: cards,
		bid:   bid,
	}, nil
}
