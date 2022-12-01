package busser

import (
	"fmt"
)

type sliceDeck []uint8

func newSliceDeck() deck {
	return &sliceDeck{}
}

func (d *sliceDeck) drawFromTop() uint8 {
	v := (*d)[0]
	*d = (*d)[1:]
	return v
}

func (d *sliceDeck) insertAtBottom(v uint8) {
	*d = append(*d, v)
}

func (d sliceDeck) size() uint8 {
	return uint8(len(d))
}

func (d sliceDeck) score() uint16 {
	total, factor := uint16(0), uint16(len(d))
	for _, v := range d {
		total += uint16(v) * factor
		factor--
	}
	return total
}

func (d sliceDeck) copyTopCards(n uint8) deck {
	if d.size() < n {
		panic(fmt.Sprintf("cannot copy %d cards from deck containing only %d cards", n, d.size()))
	}

	dd := make(sliceDeck, n)
	copy(dd[:n], d[:n])
	return &dd
}

func (d sliceDeck) hash() [MaxDeckSize]uint8 {
	var dd [MaxDeckSize]uint8
	for i, v := range d {
		dd[i] = v
	}
	return dd
}
