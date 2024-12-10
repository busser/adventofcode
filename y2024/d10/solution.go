package d10

import (
	"bytes"
	"fmt"
	"io"
	"iter"
)

// PartOne solves the first problem of day 10 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	topology, err := topologyMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	total := 0
	for trailhead := range topology.trailheads() {
		total += topology.score(trailhead)
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 10 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	topology, err := topologyMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	total := 0
	for trailhead := range topology.trailheads() {
		total += topology.rating(trailhead)
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type topologyMap [][]byte

func (m topologyMap) at(pos vector) byte {
	return m[pos.row][pos.col]
}

func (m topologyMap) isWithinBounds(pos vector) bool {
	return pos.row >= 0 && pos.row < len(m) &&
		pos.col >= 0 && pos.col < len(m[pos.row])
}

func (m topologyMap) positions() iter.Seq[vector] {
	return func(yield func(vector) bool) {
		for row := range m {
			for col := range m[row] {
				if !yield(vector{row, col}) {
					return
				}
			}
		}
	}
}

func (m topologyMap) trailheads() iter.Seq[vector] {
	return func(yield func(vector) bool) {
		for pos := range m.positions() {
			if m.at(pos) == '0' && !yield(pos) {
				return
			}
		}
	}
}

func (m topologyMap) neighbors(pos vector) iter.Seq[vector] {
	return func(yield func(vector) bool) {
		for _, dir := range []vector{up, down, left, right} {
			if !m.isWithinBounds(pos.add(dir)) {
				continue
			}

			if !yield(pos.add(dir)) {
				return
			}
		}
	}
}

func (m topologyMap) score(trailhead vector) int {
	visited := make([][]bool, len(m))
	for row := range visited {
		visited[row] = make([]bool, len(m[row]))
	}

	var visit func(vector) int
	visit = func(pos vector) int {
		if !m.isWithinBounds(pos) {
			return 0
		}

		if visited[pos.row][pos.col] {
			return 0
		}
		visited[pos.row][pos.col] = true

		if m.at(pos) == '9' {
			return 1
		}

		score := 0
		for neighbor := range m.uphillNeighbors(pos) {
			score += visit(neighbor)
		}

		return score
	}

	return visit(trailhead)
}

func (m topologyMap) rating(trailhead vector) int {
	visited := make([][]bool, len(m))
	for row := range visited {
		visited[row] = make([]bool, len(m[row]))
	}
	cachedResults := make([][]int, len(m))
	for row := range cachedResults {
		cachedResults[row] = make([]int, len(m[row]))
	}

	var visit func(vector) int
	visit = func(pos vector) int {
		if !m.isWithinBounds(pos) {
			return 0
		}

		if visited[pos.row][pos.col] {
			return cachedResults[pos.row][pos.col]
		}
		visited[pos.row][pos.col] = true

		if m.at(pos) == '9' {
			cachedResults[pos.row][pos.col] = 1
			return 1
		}

		rating := 0
		for neighbor := range m.uphillNeighbors(pos) {
			rating += visit(neighbor)
		}

		cachedResults[pos.row][pos.col] = rating
		return rating
	}

	return visit(trailhead)
}

func (m topologyMap) uphillNeighbors(pos vector) iter.Seq[vector] {
	return func(yield func(vector) bool) {
		for neighbor := range m.neighbors(pos) {
			if m.at(neighbor) == m.at(pos)+1 {
				if !yield(neighbor) {
					return
				}
			}
		}
	}
}

type vector struct {
	row, col int
}

var (
	up    = vector{row: -1, col: 0}
	down  = vector{row: 1, col: 0}
	left  = vector{row: 0, col: -1}
	right = vector{row: 0, col: 1}
)

func (v vector) add(w vector) vector {
	return vector{
		row: v.row + w.row,
		col: v.col + w.col,
	}
}

func topologyMapFromReader(r io.Reader) (topologyMap, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	topology := bytes.Split(bytes.TrimSpace(data), []byte("\n"))

	return topology, nil
}
