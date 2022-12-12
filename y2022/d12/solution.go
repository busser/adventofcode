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

	distance, err := shortestPath(
		hm.topography,
		hm.start,
		func(pos position) bool {
			return pos == hm.end
		},
		func(from, to position) bool {
			return hm.topography[to.row][to.col] <= hm.topography[from.row][from.col]+1
		},
	)
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

	distance, err := shortestPath(
		hm.topography,
		hm.end,
		func(pos position) bool {
			return hm.topography[pos.row][pos.col] == 'a'
		},
		func(from, to position) bool {
			return hm.topography[to.row][to.col] >= hm.topography[from.row][from.col]-1
		},
	)
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

func shortestPath(
	topography [][]byte,
	from position,
	isTarget func(position) bool,
	isReachable func(from, to position) bool,
) (int, error) {
	visited := make([][]bool, len(topography))
	for row := range visited {
		visited[row] = make([]bool, len(topography[row]))
	}
	visited[from.row][from.col] = true

	var toVisit, next []position
	next = append(next, from)

	visit := func(from, to position) {
		if !isWithinBounds(topography, to) {
			return
		}

		if visited[to.row][to.col] {
			return
		}

		if !isReachable(from, to) {
			return
		}

		next = append(next, to)
		visited[to.row][to.col] = true
	}

	for distance := 0; ; distance++ {
		toVisit, next = next, toVisit
		next = next[:0]

		if len(toVisit) == 0 {
			return 0, errors.New("no path")
		}

		for _, pos := range toVisit {
			if isTarget(pos) {
				return distance, nil
			}

			visit(pos, position{pos.row - 1, pos.col})
			visit(pos, position{pos.row + 1, pos.col})
			visit(pos, position{pos.row, pos.col - 1})
			visit(pos, position{pos.row, pos.col + 1})
		}
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
