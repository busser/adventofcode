package d12

import (
	"bytes"
	"fmt"
	"io"
)

// PartOne solves the first problem of day 12 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	f, err := farmFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	regions := findRegions(f)

	total := 0
	for _, r := range regions {
		total += r.fencePrice()
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 12 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	f, err := farmFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	regions := findRegions(f)

	total := 0
	for _, r := range regions {
		total += r.bulkFencePrice()
	}

	_, err = fmt.Fprintf(w, "%d", total)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type farm [][]byte

func findRegions(f farm) []region {
	visited := make([][]bool, len(f))
	for row := range visited {
		visited[row] = make([]bool, len(f[row]))
	}

	findRegion := func(pos vector) region {
		var r region
		regionType := f.at(pos)

		var visit func(vector)
		visit = func(pos vector) {
			if !f.isWithinBounds(pos) || f.at(pos) != regionType {
				return
			}

			if visited[pos.row][pos.col] {
				return
			}
			visited[pos.row][pos.col] = true

			r = append(r, pos)

			for _, neighbor := range []vector{up, down, left, right} {
				visit(pos.add(neighbor))
			}
		}

		visit(pos)

		return r
	}

	var regions []region
	for row := range f {
		for col := range f[row] {
			if !visited[row][col] {
				regions = append(regions, findRegion(vector{row, col}))
			}
		}
	}

	return regions
}

func (f farm) at(pos vector) byte {
	return f[pos.row][pos.col]
}

func (f farm) isWithinBounds(pos vector) bool {
	return pos.row >= 0 && pos.row < len(f) &&
		pos.col >= 0 && pos.col < len(f[pos.row])
}

type vector struct {
	row, col int
}

func (v vector) add(w vector) vector {
	return vector{
		row: v.row + w.row,
		col: v.col + w.col,
	}
}

var (
	up    = vector{row: -1, col: 0}
	down  = vector{row: 1, col: 0}
	left  = vector{row: 0, col: -1}
	right = vector{row: 0, col: 1}
)

type region []vector

func (r region) area() int {
	return len(r)
}

func (r region) perimeter() int {
	index := make(map[vector]bool)
	for _, v := range r {
		index[v] = true
	}

	perimeter := 0
	for _, v := range r {
		for _, neighbor := range []vector{up, down, left, right} {
			if !index[v.add(neighbor)] {
				perimeter++
			}
		}
	}

	return perimeter
}

func (r region) simplifiedPerimeter() int {
	index := make(map[vector]bool)
	for _, v := range r {
		index[v] = true
	}

	perimeter := 0
	for _, v := range r {
		directions := []vector{up, right, down, left} // ordered clockwise
		for i := range directions {
			dirA := directions[i]
			dirB := directions[(i+1)%len(directions)]

			// Detect convex corners.
			if !index[v.add(dirA)] && !index[v.add(dirB)] {
				perimeter++
			}

			// Detect concave corners.
			if index[v.add(dirA)] && index[v.add(dirB)] && !index[v.add(dirA).add(dirB)] {
				perimeter++
			}
		}
	}

	return perimeter
}

func (r region) fencePrice() int {
	return r.perimeter() * r.area()
}

func (r region) bulkFencePrice() int {
	return r.simplifiedPerimeter() * r.area()
}

func farmFromReader(r io.Reader) (farm, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	farm := bytes.Split(bytes.TrimSpace(data), []byte("\n"))

	return farm, nil
}
