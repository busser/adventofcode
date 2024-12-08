package d08

import (
	"bytes"
	"fmt"
	"io"
)

// PartOne solves the first problem of day 8 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	city, err := cityMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := countAntinodePositions(city, false)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 8 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	city, err := cityMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := countAntinodePositions(city, true)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type cityMap [][]byte

const empty = '.'

func countAntinodePositions(city cityMap, resonantHarmonics bool) int {
	nodeIndex := make(map[byte][]vector)
	for row := range city {
		for col, node := range city[row] {
			if node == empty {
				continue
			}
			position := vector{row, col}
			nodeIndex[node] = append(nodeIndex[node], position)
		}
	}

	positionHasAntinode := make(map[vector]struct{})
	for _, positions := range nodeIndex {
		for _, positionA := range positions {
			for _, positionB := range positions {
				if positionA == positionB {
					continue
				}

				delta := positionA.minus(positionB)

				if !resonantHarmonics {
					antinodePosition := positionA.plus(delta)
					if withinCityBounds(city, antinodePosition) {
						positionHasAntinode[antinodePosition] = struct{}{}
					}
					continue
				}

				antinodePosition := positionA
				for withinCityBounds(city, antinodePosition) {
					positionHasAntinode[antinodePosition] = struct{}{}
					antinodePosition = antinodePosition.plus(delta)
				}
			}
		}
	}

	return len(positionHasAntinode)
}

func withinCityBounds(city cityMap, position vector) bool {
	return position.row >= 0 && position.row < len(city) &&
		position.col >= 0 && position.col < len(city[position.row])
}

type vector struct {
	row, col int
}

func (v vector) plus(w vector) vector {
	return vector{
		row: v.row + w.row,
		col: v.col + w.col,
	}
}

func (v vector) minus(w vector) vector {
	return vector{
		row: v.row - w.row,
		col: v.col - w.col,
	}
}

func (v vector) times(n int) vector {
	return vector{
		row: v.row * n,
		col: v.col * n,
	}
}

func cityMapFromReader(r io.Reader) (cityMap, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	data = bytes.TrimSpace(data)

	return bytes.Split(data, []byte("\n")), nil
}
