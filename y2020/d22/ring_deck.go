package busser

import (
	"fmt"
)

const ringDeckMaxSize = MaxDeckSize + 1

type ringDeck struct {
	cards       [ringDeckMaxSize]uint8
	top, bottom uint8
}

func newRingDeck() deck {
	return new(ringDeck)
}

func (d *ringDeck) drawFromTop() uint8 {
	v := d.cards[d.top]
	d.top = (d.top + 1) % ringDeckMaxSize
	return v
}

func (d *ringDeck) insertAtBottom(v uint8) {
	d.cards[d.bottom] = v
	d.bottom = (d.bottom + 1) % ringDeckMaxSize
}

func (d ringDeck) size() uint8 {
	return (d.bottom - d.top + ringDeckMaxSize) % ringDeckMaxSize
}

func (d ringDeck) score() uint16 {
	total, factor := uint16(0), uint16(d.size())
	for i := d.top; i != d.bottom; i = (i + 1) % ringDeckMaxSize {
		total += uint16(d.cards[i]) * factor
		factor--
	}
	return total
}

func (d ringDeck) copyTopCards(n uint8) deck {
	if d.size() < n {
		panic(fmt.Sprintf("cannot copy %d cards from deck containing only %d cards", n, d.size()))
	}

	dd := new(ringDeck)

	i := d.top
	for ; n > 0; n-- {
		dd.insertAtBottom(d.cards[i])
		i = (i + 1) % ringDeckMaxSize
	}

	return dd
}

func (d ringDeck) hash() [MaxDeckSize]uint8 {
	var h [MaxDeckSize]uint8

	i := 0
	for j := d.top; j != d.bottom; j = (j + 1) % ringDeckMaxSize {
		h[i] = d.cards[j]
		i++
	}

	return h
}
