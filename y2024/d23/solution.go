package d23

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 23 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	connections, err := connectionsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	graph := newGraph(connections)
	count := countTriomesWithT(graph)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 23 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	connections, err := connectionsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	graph := newGraph(connections)
	lan := graph.biggestLAN()
	password := lanPassword(lan)

	_, err = fmt.Fprintf(w, "%s", password)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type computerName [2]byte

func (c computerName) less(other computerName) bool {
	if c[0] != other[0] {
		return c[0] < other[0]
	}
	return c[1] < other[1]
}

type connection [2]computerName

type networkGraph map[computerName][]computerName

func newGraph(connections []connection) networkGraph {
	graph := make(map[computerName][]computerName)

	for _, c := range connections {
		graph[c[0]] = append(graph[c[0]], c[1])
		graph[c[1]] = append(graph[c[1]], c[0])
	}

	return graph
}

func (g networkGraph) neighbors(computer computerName) []computerName {
	return g[computer]
}

func (g networkGraph) areConnected(a, b computerName) bool {
	return slices.Contains(g.neighbors(a), b)
}

func (g networkGraph) biggestLAN() []computerName {
	allComputers := make([]computerName, 0, len(g))
	for computer := range g {
		allComputers = append(allComputers, computer)
	}

	var computersInLAN, biggestLAN []computerName

	var helper func(i int)
	helper = func(i int) {
		newComputer := allComputers[i]

		for _, c := range computersInLAN {
			if !g.areConnected(newComputer, c) {
				return
			}
		}

		computersInLAN = append(computersInLAN, newComputer)
		if len(computersInLAN) > len(biggestLAN) {
			biggestLAN = slices.Clone(computersInLAN)
		}

		for j := i + 1; j < len(allComputers); j++ {
			helper(j)
		}

		computersInLAN = computersInLAN[:len(computersInLAN)-1]
	}

	for i := range allComputers {
		helper(i)
	}

	return biggestLAN
}

func lanPassword(lan []computerName) string {
	slices.SortFunc(lan, func(a, b computerName) int {
		if a.less(b) {
			return -1
		}
		return 1
	})

	var sb strings.Builder
	for i, c := range lan {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte(c[0])
		sb.WriteByte(c[1])
	}

	return sb.String()
}

type triome [3]computerName

func (t *triome) sort() {
	if t[0].less(t[1]) {
		t[0], t[1] = t[1], t[0]
	}
	if t[0].less(t[2]) {
		t[0], t[2] = t[2], t[0]
	}
	if t[1].less(t[2]) {
		t[1], t[2] = t[2], t[1]
	}
}

func countTriomesWithT(graph networkGraph) int {
	triomes := make(map[[3]computerName]struct{})

	for a := range graph {
		for _, b := range graph.neighbors(a) {
			for _, c := range graph.neighbors(b) {
				if a == c {
					continue
				}
				if graph.areConnected(a, c) {
					t := triome{a, b, c}
					t.sort()
					triomes[t] = struct{}{}
				}
			}
		}
	}

	count := 0
	for t := range triomes {
		if t[0][0] == 't' || t[1][0] == 't' || t[2][0] == 't' {
			count++
		}
	}

	return count
}

func connectionsFromReader(r io.Reader) ([]connection, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	connections := make([]connection, len(lines))
	for i, line := range lines {
		connections[i], err = connectionFromString(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse connection: %w", err)
		}
	}

	return connections, nil
}

func connectionFromString(s string) (connection, error) {
	if len(s) != 5 {
		return connection{}, fmt.Errorf("invalid connection: %q", s)
	}

	return connection{
		computerName{s[0], s[1]},
		computerName{s[3], s[4]},
	}, nil
}
