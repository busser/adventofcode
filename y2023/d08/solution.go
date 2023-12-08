package d08

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 8 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	dm, err := desertMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	steps := dm.stepsOnNormalPath()

	_, err = fmt.Fprintf(w, "%d", steps)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 8 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	dm, err := desertMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	steps := dm.stepsOnParallelGhostPath()

	_, err = fmt.Fprintf(w, "%d", steps)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type node struct {
	id          string
	left, right string
}

const (
	startNodeID  = "AAA"
	finishNodeID = "ZZZ"
)

type instruction byte

const (
	left  instruction = 'L'
	right instruction = 'R'
)

type desertMap struct {
	instructions []instruction
	nodes        map[string]*node
}

func (dm *desertMap) stepsOnNormalPath() int {
	steps := 0
	current := startNodeID

	for current != finishNodeID {
		currentNode := dm.nodes[current]
		instruction := dm.instructions[steps%len(dm.instructions)]
		switch instruction {
		case left:
			current = currentNode.left
		case right:
			current = currentNode.right
		}
		steps++
	}

	return steps
}

func (dm *desertMap) ghostStartNodes() []string {
	var nodes []string
	for id := range dm.nodes {
		if isGhostStartNode(id) {
			nodes = append(nodes, id)
		}
	}
	return nodes
}

func isGhostStartNode(id string) bool {
	return id[2] == 'A'
}

func isGhostFinishNode(id string) bool {
	return id[2] == 'Z'
}

func (dm *desertMap) stepsOnParallelGhostPath() int {
	startNodes := dm.ghostStartNodes()

	var pathLengths []int
	for _, start := range startNodes {
		pathLengths = append(pathLengths, dm.stepsOnSingleGhostPath(start))
	}

	if len(pathLengths) == 1 {
		return pathLengths[0]
	}

	parallelPathLength := pathLengths[0]
	for i := 1; i < len(pathLengths); i++ {
		parallelPathLength = lcm(parallelPathLength, pathLengths[i])
	}

	return parallelPathLength
}

func lcm(a, b int) int {
	if a == 0 && b == 0 {
		return 0
	}

	return abs(a*b) / gcd(a, b)
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (dm *desertMap) stepsOnSingleGhostPath(start string) int {
	steps := 0
	current := start

	for !isGhostFinishNode(current) {
		currentNode := dm.nodes[current]
		instruction := dm.instructions[steps%len(dm.instructions)]
		switch instruction {
		case left:
			current = currentNode.left
		case right:
			current = currentNode.right
		}
		steps++
	}

	return steps
}

func desertMapFromReader(r io.Reader) (desertMap, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return desertMap{}, fmt.Errorf("could not read input: %w", err)
	}

	if len(lines) < 3 {
		return desertMap{}, fmt.Errorf("not enought lines")
	}

	instructions, err := instructionsFromString(lines[0])
	if err != nil {
		return desertMap{}, fmt.Errorf("could not read instructions: %w", err)
	}

	nodes := make(map[string]*node)
	for _, line := range lines[2:] {
		n, err := nodeFromString(line)
		if err != nil {
			return desertMap{}, fmt.Errorf("could not read node: %w", err)
		}
		nodes[n.id] = &n
	}

	if _, ok := nodes[startNodeID]; !ok {
		return desertMap{}, fmt.Errorf("no start node")
	}
	if _, ok := nodes[finishNodeID]; !ok {
		return desertMap{}, fmt.Errorf("no finish node")
	}

	return desertMap{
		instructions: instructions,
		nodes:        nodes,
	}, nil
}

func instructionsFromString(s string) ([]instruction, error) {
	instructions := make([]instruction, len(s))
	for i, r := range s {
		switch r {
		case 'L':
			instructions[i] = left
		case 'R':
			instructions[i] = right
		default:
			return nil, fmt.Errorf("invalid instruction %q", r)
		}
	}
	return instructions, nil
}

func nodeFromString(s string) (node, error) {
	if len(s) != 16 {
		return node{}, fmt.Errorf("invalid node %q", s)
	}

	return node{
		id:    s[0:3],
		left:  s[7:10],
		right: s[12:15],
	}, nil
}
