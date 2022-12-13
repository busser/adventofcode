package d13

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 13 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	pairs, err := packetsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	indexSum := inOrderPairs(pairs)

	_, err = fmt.Fprintf(w, "%d", indexSum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 13 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	pairs, err := packetsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	key := decoderKey(pairs)

	_, err = fmt.Fprintf(w, "%d", key)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type packet struct {
	isList  bool
	values  []packet
	integer int
}

type pair[T any] struct {
	left, right T
}

func inOrderPairs(pairs []pair[packet]) int {
	indexSum := 0

	for i, p := range pairs {
		if diff(p.left, p.right) < 0 {
			indexSum += i + 1
		}
	}

	return indexSum
}

func decoderKey(pairs []pair[packet]) int {
	firstDivider := must(packetValueFromString("[[2]]"))
	secondDivider := must(packetValueFromString("[[6]]"))

	allPackets := []packet{firstDivider, secondDivider}

	for _, p := range pairs {
		allPackets = append(allPackets, p.left, p.right)
	}

	sort.Slice(allPackets, func(i, j int) bool {
		return diff(allPackets[i], allPackets[j]) < 0
	})

	decoderKey := 1
	for i, p := range allPackets {
		if diff(p, firstDivider) == 0 || diff(p, secondDivider) == 0 {
			decoderKey *= i + 1
		}
	}

	return decoderKey
}

func diff(left, right packet) int {
	switch {
	case left.isList && right.isList:
		for i := range left.values {
			if i >= len(right.values) {
				return 1
			}
			if d := diff(left.values[i], right.values[i]); d != 0 {
				return d
			}
		}
		if len(left.values) < len(right.values) {
			return -1
		}
		return 0

	case !left.isList && right.isList:
		newLeft := packet{
			isList: true,
			values: []packet{left},
		}
		return diff(newLeft, right)
	case left.isList && !right.isList:
		newRight := packet{
			isList: true,
			values: []packet{right},
		}
		return diff(left, newRight)
	case !left.isList && !right.isList:
		return left.integer - right.integer
	}

	panic("impossible")
}

func packetsFromReader(r io.Reader) ([]pair[packet], error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	if len(lines)%3 != 2 {
		return nil, errors.New("unexpected number of lines")
	}

	var pairs []pair[packet]

	for i := 0; i < len(lines); i += 3 {
		left, err := packetValueFromString(lines[i])
		if err != nil {
			return nil, fmt.Errorf("reading line %d: %w", i, err)
		}

		right, err := packetValueFromString(lines[i+1])
		if err != nil {
			return nil, fmt.Errorf("reading line %d: %w", i+1, err)
		}

		pairs = append(pairs, pair[packet]{left, right})
	}

	return pairs, nil
}

func must(p packet, err error) packet {
	if err != nil {
		panic(err)
	}
	return p
}

func packetValueFromString(s string) (packet, error) {
	i := 0

	var readValue func(*packet) error
	var readList func(*packet) error
	readValue = func(current *packet) error {
		for {
			switch {

			case i >= len(s):
				return errors.New("unexpected EOL")

			case s[i] == '[':
				// value is a list, recurse
				i++
				current.isList = true
				return readList(current)

			case isDigit(s[i]):
				// value is an integer
				var n int
				for ; isDigit(s[i]); i++ {
					n = 10*n + int(s[i]-'0')
				}
				current.integer = n
				return nil

			default:
				return fmt.Errorf("unexpected symbol %q in non-list value", s[i])

			}
		}
	}
	readList = func(current *packet) error {
		for {
			switch {

			case i >= len(s):
				return errors.New("unexpected EOL")

			case s[i] == ',':
				// next value in a list
				if !current.isList {
					return errors.New("unexpected comma is non-list value")
				}
				i++

			case s[i] == ']':
				// done with this list
				if !current.isList {
					return errors.New("unexpected closing bracket in non-list value")
				}
				i++
				return nil

			default:
				var child packet
				if err := readValue(&child); err != nil {
					return err
				}
				current.values = append(current.values, child)
			}
		}
	}

	var p packet
	if err := readValue(&p); err != nil {
		return packet{}, err
	}
	if i != len(s) {
		return packet{}, errors.New("invalid packet: data remains right of closing bracket")
	}

	return p, nil
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
