package d23

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 23 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	hikingMap, err := hikingMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	start := vector{0, 1}
	end := vector{len(hikingMap) - 1, len(hikingMap[0]) - 2}

	g := graphFromMap(hikingMap, start, end)
	length := longestPath(g, start, end)

	_, err = fmt.Fprintf(w, "%d", length)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 23 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	hikingMap, err := hikingMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	start := vector{0, 1}
	end := vector{len(hikingMap) - 1, len(hikingMap[0]) - 2}

	removeSlopes(hikingMap)
	g := graphFromMap(hikingMap, start, end)
	length := longestPath(g, start, end)

	_, err = fmt.Fprintf(w, "%d", length)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const (
	path       = '.'
	forest     = '#'
	upSlope    = '^'
	downSlope  = 'v'
	leftSlope  = '<'
	rightSlope = '>'
)

type vector struct {
	row, col int
}

var (
	up    = vector{-1, 0}
	down  = vector{1, 0}
	left  = vector{0, -1}
	right = vector{0, 1}
)

func (v vector) add(w vector) vector {
	return vector{v.row + w.row, v.col + w.col}
}

type graph map[vector][]graphEdge

type graphEdge struct {
	to     vector
	length int
}

func graphFromMap(m [][]byte, start, end vector) graph {
	// This graph's nodes are only the positions we want to keep. We choose
	// these nodes to be the positions where the path can take multiple
	// directions. This ensures that the graph is a small as possible.

	nodesToKeep := make(map[vector]bool)

	nodesToKeep[start] = true
	nodesToKeep[end] = true

	for row := range m {
		for col := range m[row] {
			if m[row][col] == forest {
				continue
			}

			possibleDirections := possibleDirectionsFor(m[row][col])

			validNeighbors := 0
			for _, v := range possibleDirections {
				next := vector{row, col}.add(v)

				if next.row < 0 || next.row >= len(m) || next.col < 0 || next.col >= len(m[next.row]) {
					continue
				}
				if m[next.row][next.col] == forest {
					continue
				}

				validNeighbors++
			}

			if validNeighbors > 2 {
				nodesToKeep[vector{row, col}] = true
			}
		}
	}

	// The edges are the paths from each node to neighboring nodes. The length
	// of the edge is the number of steps it takes to get from one node to the
	// other. We don't add edges for paths that go through another node.
	//
	// Since each node represents a fork, there cannot be forks on the path
	// between two nodes. With simple DFS, we can explore from a node to all
	// neighboring nodes and add the edges to the graph. We stop exploring when
	// we reach another node.

	g := make(graph)

	for node := range nodesToKeep {
		visited := make(map[vector]bool)
		visited[node] = true

		var stack []graphEdge
		stack = append(stack, graphEdge{node, 0})

		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			possibleDirections := possibleDirectionsFor(m[current.to.row][current.to.col])

			for _, v := range possibleDirections {
				next := current.to.add(v)

				if next.row < 0 || next.row >= len(m) || next.col < 0 || next.col >= len(m[next.row]) {
					continue
				}
				if m[next.row][next.col] == forest {
					continue
				}
				if visited[next] {
					continue
				}
				if nodesToKeep[next] {
					g[node] = append(g[node], graphEdge{next, current.length + 1})
					continue
				}

				visited[next] = true
				stack = append(stack, graphEdge{next, current.length + 1})
			}
		}
	}

	return g
}

func possibleDirectionsFor(c byte) []vector {
	switch c {
	case path:
		return []vector{up, down, left, right}
	case upSlope:
		return []vector{up}
	case downSlope:
		return []vector{down}
	case leftSlope:
		return []vector{left}
	case rightSlope:
		return []vector{right}
	default:
		return nil
	}
}

func longestPath(g graph, start, end vector) int {
	visited := make(map[vector]bool)

	length, _ := longestPathFrom(g, visited, start, end)

	return length
}

func longestPathFrom(g graph, visited map[vector]bool, start, end vector) (int, bool) {
	if start == end {
		return 0, true
	}

	visited[start] = true

	var longest int
	foundPath := false
	for _, edge := range g[start] {
		if visited[edge.to] {
			continue
		}

		length, endReached := longestPathFrom(g, visited, edge.to, end)
		if endReached {
			longest = max(longest, length+edge.length)
			foundPath = true
		}
	}

	visited[start] = false

	return longest, foundPath
}

func removeSlopes(m [][]byte) {
	for row := range m {
		for col := range m[row] {
			switch m[row][col] {
			case upSlope, downSlope, leftSlope, rightSlope:
				m[row][col] = path
			}
		}
	}
}

func hikingMapFromReader(r io.Reader) ([][]byte, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	m := make([][]byte, len(lines))
	for row, line := range lines {
		m[row] = []byte(line)

		for col, c := range line {
			switch c {
			case path, forest, upSlope, downSlope, leftSlope, rightSlope:
				// valid
			default:
				return nil, fmt.Errorf("invalid character %c at row %d, column %d", c, row, col)
			}
		}
	}

	return m, nil
}
