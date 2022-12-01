package busser

import (
	"encoding/hex"
	"fmt"
	"io"
	"math"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 16 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	binary, err := binaryFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	p, err := newDecoder(binary).parsePacket()
	if err != nil {
		return fmt.Errorf("invalid encoding: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", p.versionSum())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 16 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	binary, err := binaryFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	p, err := newDecoder(binary).parsePacket()
	if err != nil {
		return fmt.Errorf("invalid encoding: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", p.value())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type decoder struct {
	data  []byte
	index int
}

type packet struct {
	version int
	typeID  int

	// for packets with a type ID of 4
	literalValue int

	// for packets with a type ID other than 4
	lengthTypeID int
	subPackets   []packet
}

func newDecoder(data []byte) *decoder {
	return &decoder{data: data}
}

func (d *decoder) readBit() bool {
	mask := byte(1 << 7 >> (d.index % 8))
	bit := d.data[d.index/8] & mask
	d.index++
	return bit != 0
}

func (d *decoder) readBits(n int) int {
	v := 0
	for i := 0; i < n; i++ {
		v *= 2
		if d.readBit() {
			v++
		}
	}
	return v
}

func (d *decoder) parsePacket() (packet, error) {
	var p packet

	p.version = d.readBits(3)
	p.typeID = d.readBits(3)

	if p.typeID == 4 {
		keepGoing := true
		for keepGoing {
			keepGoing = d.readBit()
			p.literalValue *= 1 << 4
			p.literalValue += d.readBits(4)
		}

		return p, nil
	}

	p.lengthTypeID = d.readBits(1)
	switch p.lengthTypeID {
	case 0:
		length := d.readBits(15)
		startIndex := d.index
		for d.index < startIndex+length {
			subPacket, err := d.parsePacket()
			if err != nil {
				return packet{}, fmt.Errorf("could not parse sub-packet: %w", err)
			}
			p.subPackets = append(p.subPackets, subPacket)
		}
		if d.index-startIndex != length {
			return packet{}, fmt.Errorf("expected sub-packets to take %d bits, actually took %d bits", length, d.index-startIndex)
		}
		return p, nil
	case 1:
		numSubPackets := d.readBits(11)
		for i := 0; i < numSubPackets; i++ {
			subPacket, err := d.parsePacket()
			if err != nil {
				return packet{}, fmt.Errorf("could not parse sub-packet: %w", err)
			}
			p.subPackets = append(p.subPackets, subPacket)
		}
		return p, nil
	default:
		return packet{}, fmt.Errorf("unknown length type ID %d", p.lengthTypeID)
	}
}

func (p packet) versionSum() int {
	sum := p.version
	for _, sub := range p.subPackets {
		sum += sub.versionSum()
	}
	return sum
}

func (p packet) value() int {
	switch p.typeID {
	case 0:
		sum := 0
		for _, sub := range p.subPackets {
			sum += sub.value()
		}
		return sum
	case 1:
		product := 1
		for _, sub := range p.subPackets {
			product *= sub.value()
		}
		return product
	case 2:
		min := math.MaxInt
		for _, sub := range p.subPackets {
			if v := sub.value(); v < min {
				min = v
			}
		}
		return min
	case 3:
		max := math.MinInt
		for _, sub := range p.subPackets {
			if v := sub.value(); v > max {
				max = v
			}
		}
		return max
	case 4:
		return p.literalValue
	case 5:
		if p.subPackets[0].value() > p.subPackets[1].value() {
			return 1
		} else {
			return 0
		}
	case 6:
		if p.subPackets[0].value() < p.subPackets[1].value() {
			return 1
		} else {
			return 0
		}
	case 7:
		if p.subPackets[0].value() == p.subPackets[1].value() {
			return 1
		} else {
			return 0
		}
	}
	panic("impossible packet type ID")
}

func binaryFromReader(r io.Reader) ([]byte, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("expected 1 line, got %d", len(lines))
	}

	return hex.DecodeString(lines[0])
}
