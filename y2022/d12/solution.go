package d12

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// PartOne solves the first problem of day 12 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	hm, err := heightmapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	distance, err := shortestPath(hm.topography, []position{hm.start}, hm.end)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", distance)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 12 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	hm, err := heightmapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var starts []position
	for row := range hm.topography {
		for col := range hm.topography[row] {
			if hm.topography[row][col] == 'a' {
				starts = append(starts, position{row, col})
			}
		}
	}

	distance, err := shortestPath(hm.topography, starts, hm.end)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", distance)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type heightmap struct {
	start, end position
	topography [][]byte
}

type position struct {
	row, col int
}

func shortestPath(topography [][]byte, from []position, to position) (int, error) {
	visited := make([][]bool, len(topography))
	for row := range visited {
		visited[row] = make([]bool, len(topography[row]))
	}

	var next []position
	toVisit := from

	for _, pos := range from {
		visited[pos.row][pos.col] = true
	}

	visit := func(from, to position) {
		if !isWithinBounds(topography, to) {
			return
		}

		if visited[to.row][to.col] {
			return
		}

		if topography[to.row][to.col] > topography[from.row][from.col]+1 {
			return
		}

		next = append(next, to)
		visited[to.row][to.col] = true
	}

	for distance := 0; ; distance++ {
		for _, pos := range toVisit {
			if pos == to {
				return distance, nil
			}

			visit(pos, position{pos.row - 1, pos.col})
			visit(pos, position{pos.row + 1, pos.col})
			visit(pos, position{pos.row, pos.col - 1})
			visit(pos, position{pos.row, pos.col + 1})
		}

		if len(next) == 0 {
			return 0, errors.New("no path")
		}

		toVisit, next = next, toVisit
		next = next[:0]
	}
}

func isWithinBounds(topography [][]byte, pos position) bool {
	return pos.row >= 0 && pos.row < len(topography) &&
		pos.col >= 0 && pos.col < len(topography[pos.row])
}

func heightmapFromReader(r io.Reader) (heightmap, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return heightmap{}, err
	}

	hm := heightmap{
		topography: bytes.Split(bytes.TrimSpace(raw), []byte("\n")),
	}

	if len(hm.topography) == 0 {
		return heightmap{}, errors.New("no data")
	}

	var startFound, endFound bool
	for row := range hm.topography {
		if len(hm.topography[row]) != len(hm.topography[0]) {
			return heightmap{}, errors.New("map is not a rectangle")
		}

		for col := range hm.topography[row] {
			if hm.topography[row][col] == 'S' {
				hm.topography[row][col] = 'a'
				hm.start = position{row, col}
				startFound = true
			}
			if hm.topography[row][col] == 'E' {
				hm.topography[row][col] = 'z'
				hm.end = position{row, col}
				endFound = true
			}
		}
	}

	if !startFound {
		return heightmap{}, errors.New("no start")
	}
	if !endFound {
		return heightmap{}, errors.New("no end")
	}

	return hm, nil
}
