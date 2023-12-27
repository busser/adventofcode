package d25

import (
	"fmt"
	"io"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 25 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	graph, err := graphFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	firstGroupSize, secondGroupSize, err := splitGraph(graph, 3)
	if err != nil {
		return fmt.Errorf("could not split graph: %w", err)
	}

	product := firstGroupSize * secondGroupSize

	_, err = fmt.Fprintf(w, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func splitGraph(graph map[string][]string, target int) (int, int, error) {
	inSecondGroup := make(map[string]bool)

	countBoundaryCrossings := func(node string) int {
		boundaryCrossings := 0
		for _, edge := range graph[node] {
			if inSecondGroup[edge] {
				boundaryCrossings++
			}
		}
		return boundaryCrossings
	}

	totalBoundaryCrossings := func() int {
		totalBoundaryCrossings := 0
		for node := range graph {
			if inSecondGroup[node] {
				continue
			}
			totalBoundaryCrossings += countBoundaryCrossings(node)
		}
		return totalBoundaryCrossings
	}

	nodeWithMaxBoundaryCrossings := func() string {
		maxBoundaryCrossings := -1
		var maxBoundaryCrossingsNode string

		for node := range graph {
			if inSecondGroup[node] {
				continue
			}
			boundaryCrossings := countBoundaryCrossings(node)
			if boundaryCrossings > maxBoundaryCrossings {
				maxBoundaryCrossings = boundaryCrossings
				maxBoundaryCrossingsNode = node
			}
		}

		return maxBoundaryCrossingsNode
	}

	for len(inSecondGroup) < len(graph) {
		node := nodeWithMaxBoundaryCrossings()
		inSecondGroup[node] = true

		boundaryCrossings := totalBoundaryCrossings()
		if boundaryCrossings == target {
			firstGroupSize := len(graph) - len(inSecondGroup)
			secondGroupSize := len(inSecondGroup)
			return firstGroupSize, secondGroupSize, nil
		}
	}

	return 0, 0, fmt.Errorf("could not reach goal of %d boundary crossings", target)
}

func graphFromReader(r io.Reader) (map[string][]string, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	graph := make(map[string][]string)
	for _, line := range lines {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line: %q", line)
		}

		node := parts[0]
		edges := strings.Split(parts[1], " ")

		for _, edge := range edges {
			graph[node] = append(graph[node], edge)
			graph[edge] = append(graph[edge], node)
		}
	}

	return graph, nil
}
