package d08

import (
	"fmt"
	"io"
	"slices"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 8 of Advent of Code 2025.
func PartOne(r io.Reader, w io.Writer) error {
	p, err := playgroundFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	p.initializeCircuits()
	p.assemblePairs()

	for range 1000 {
		_ = p.makeConnection()
	}

	sizes := p.circuitSizes()
	product := sizes[0] * sizes[1] * sizes[2]

	_, err = fmt.Fprintf(w, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 8 of Advent of Code 2025.
func PartTwo(r io.Reader, w io.Writer) error {
	p, err := playgroundFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	p.initializeCircuits()
	p.assemblePairs()

	var lastConnection boxPair
	for len(p.circuitSizes()) > 1 {
		lastConnection = p.makeConnection()
	}

	product := lastConnection.boxA.position.x * lastConnection.boxB.position.x

	_, err = fmt.Fprintf(w, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type playground struct {
	junctionBoxes []*junctionBox
	boxPairs      []boxPair
}

type junctionBox struct {
	position  vector
	circuitID int
}

type vector struct {
	x, y, z int
}

func (v vector) distanceSquared(w vector) int {
	dx := v.x - w.x
	dy := v.y - w.y
	dz := v.z - w.z

	return dx*dx + dy*dy + dz*dz
}

type boxPair struct {
	boxA, boxB *junctionBox // index
}

func (bp boxPair) distanceSquared() int {
	return bp.boxA.position.distanceSquared(bp.boxB.position)
}

func (p *playground) initializeCircuits() {
	for i, box := range p.junctionBoxes {
		box.circuitID = i
	}
}

func (p *playground) mergeCircuits(a, b int) {
	for _, box := range p.junctionBoxes {
		if box.circuitID == a {
			box.circuitID = b
		}
	}
}

func (p *playground) circuitSizes() []int {
	sizeByCircuitID := make(map[int]int)
	for _, box := range p.junctionBoxes {
		sizeByCircuitID[box.circuitID]++
	}

	sizes := make([]int, 0, len(sizeByCircuitID))
	for _, size := range sizeByCircuitID {
		sizes = append(sizes, size)
	}

	slices.Sort(sizes)
	slices.Reverse(sizes)

	return sizes
}

func (p *playground) makeConnection() boxPair {
	closestPair := p.boxPairs[0]
	p.boxPairs = p.boxPairs[1:]
	p.connectBoxes(closestPair.boxA, closestPair.boxB)
	return closestPair
}

func (p *playground) connectBoxes(boxA, boxB *junctionBox) {
	p.mergeCircuits(boxA.circuitID, boxB.circuitID)
}

func (p *playground) assemblePairs() {
	boxPairs := make([]boxPair, 0, len(p.junctionBoxes)*(len(p.junctionBoxes)+1)/2)

	for i := range p.junctionBoxes {
		for j := i + 1; j < len(p.junctionBoxes); j++ {
			pair := boxPair{
				boxA: p.junctionBoxes[i],
				boxB: p.junctionBoxes[j],
			}
			boxPairs = append(boxPairs, pair)
		}
	}

	slices.SortFunc(boxPairs, func(a, b boxPair) int {
		return a.distanceSquared() - b.distanceSquared()
	})

	p.boxPairs = boxPairs
}

func junctionBoxFromString(s string) (*junctionBox, error) {
	ints := helpers.IntsFromString(s)
	if len(ints) != 3 {
		return nil, fmt.Errorf("expected 3 ints, found %d", len(ints))
	}

	position := vector{x: ints[0], y: ints[1], z: ints[2]}
	box := junctionBox{position: position}

	return &box, nil
}

func playgroundFromReader(r io.Reader) (playground, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return playground{}, fmt.Errorf("could not read input: %w", err)
	}

	boxes := make([]*junctionBox, len(lines))
	for i := range lines {
		box, err := junctionBoxFromString(lines[i])
		if err != nil {
			return playground{}, err
		}
		boxes[i] = box
	}

	p := playground{junctionBoxes: boxes}

	return p, nil
}
