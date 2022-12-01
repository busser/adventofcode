package busser

import (
	"bytes"
	"fmt"
	"io"
	"sort"
)

// PartOne solves the first problem of day 9 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	hm, err := heightmapFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read heightmap: %w", err)
	}

	lowPoints := hm.lowPoints()

	totalRiskLevel := 0
	for _, p := range lowPoints {
		riskLevel := int(p - '0' + 1)
		totalRiskLevel += riskLevel
	}

	_, err = fmt.Fprintf(answer, "%d", totalRiskLevel)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 9 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	hm, err := heightmapFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read heightmap: %w", err)
	}

	sizes := hm.basinSizes()
	if len(sizes) < 3 {
		return fmt.Errorf("found only %d basins", len(sizes))
	}

	sort.Ints(sizes)
	sizes = sizes[len(sizes)-3:]

	_, err = fmt.Fprintf(answer, "%d", sizes[0]*sizes[1]*sizes[2])
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type heightmap [][]byte

func (hm heightmap) lowPoints() []byte {
	var points []byte
	for i := range hm {
		for j := range hm[i] {
			if (i == 0 || hm[i-1][j] > hm[i][j]) &&
				(i == len(hm)-1 || hm[i+1][j] > hm[i][j]) &&
				(j == 0 || hm[i][j-1] > hm[i][j]) &&
				(j == len(hm[i])-1 || hm[i][j+1] > hm[i][j]) {
				points = append(points, hm[i][j])
			}
		}
	}
	return points
}

func (hm heightmap) basinSizes() []int {
	// Find all basins with a depth-first search that stops at 9's.

	visited := make([][]bool, len(hm))
	for i := range hm {
		visited[i] = make([]bool, len(hm[i]))
	}

	var visit func(int, int) int
	visit = func(i, j int) (size int) {
		if i < 0 || i >= len(hm) || j < 0 || j >= len(hm[i]) {
			return 0
		}

		if visited[i][j] {
			return 0
		}
		visited[i][j] = true

		if hm[i][j] == '9' {
			return 0
		}

		return 1 + visit(i-1, j) + visit(i+1, j) + visit(i, j-1) + visit(i, j+1)
	}

	var sizes []int
	for i := range hm {
		for j := range hm[i] {
			s := visit(i, j)
			if s > 0 {
				sizes = append(sizes, s)
			}
		}
	}

	return sizes
}

func heightmapFromReader(r io.Reader) (heightmap, error) {
	rawMap, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	rawMap = bytes.TrimSpace(rawMap)

	hm := heightmap(bytes.Split(rawMap, []byte("\n")))

	return hm, nil
}
