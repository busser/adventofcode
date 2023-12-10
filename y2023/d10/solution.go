package d10

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 10 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	pipeMap, err := pipeMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	loop, err := findLoop(pipeMap)
	if err != nil {
		return fmt.Errorf("could not find loop: %w", err)
	}

	_, err = fmt.Fprintf(w, "%d", len(loop)/2)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 10 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	pipeMap, err := pipeMapFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	loop, err := findLoop(pipeMap)
	if err != nil {
		return fmt.Errorf("could not find loop: %w", err)
	}

	area := areaWithinLoop(pipeMap, loop)

	_, err = fmt.Fprintf(w, "%d", area)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const (
	pipeVertical      = '|'
	pipeHorizontal    = '-'
	pipeBendNorthEast = 'L'
	pipeBendNorthWest = 'J'
	pipeBendSouthWest = '7'
	pipeBendSouthEast = 'F'
	ground            = '.'
	startingPosition  = 'S'
)

type position struct {
	row, col int
}

func areaWithinLoop(pipeMap [][]byte, loop []position) int {
	// To compute the area within the loop, we compute the area outside the
	// loop and then do some simple math. To find all the positions outside the
	// loop, we use DFS starting from the outside of the loop. For this to
	// work, we need the "outside" to be a single contiguous area. To ensure
	// that this is the case, we do two things:
	//   1. Pad the map with a border of ground.
	//   2. Zoom in on the map by a factor of 2. This allows us to squeeze
	//      through pipes. When zooming in, we keep the pipes in the loop
	//      connected by adding a vertical or horizontal pipe between them.
	// Once we've identified all the positions outside the loop, we unpad and
	// unzoom. Then, by simple substraction, we determine the area within the
	// loop.

	zoomedInMap := make([][]byte, len(pipeMap)*2+1)
	for i := range zoomedInMap {
		zoomedInMap[i] = make([]byte, len(pipeMap[0])*2+1)
	}

	for row := range zoomedInMap {
		for col := range zoomedInMap[row] {
			zoomedInMap[row][col] = ground
		}
	}

	for i := range loop {
		pos, nextPos := loop[i], loop[(i+1)%len(loop)]

		zoomedInMap[pos.row*2+1][pos.col*2+1] = pipeMap[pos.row][pos.col]

		rowDelta, colDelta := nextPos.row-pos.row, nextPos.col-pos.col

		switch {
		case rowDelta == -1 && colDelta == 0:
			zoomedInMap[pos.row*2][pos.col*2+1] = pipeVertical
		case rowDelta == 1 && colDelta == 0:
			zoomedInMap[pos.row*2+2][pos.col*2+1] = pipeVertical
		case rowDelta == 0 && colDelta == -1:
			zoomedInMap[pos.row*2+1][pos.col*2] = pipeHorizontal
		case rowDelta == 0 && colDelta == 1:
			zoomedInMap[pos.row*2+1][pos.col*2+2] = pipeHorizontal
		default:
			panic("diagonal pipe")
		}
	}

	outsideArea := 0

	seen := make([][]bool, len(zoomedInMap))
	for i := range seen {
		seen[i] = make([]bool, len(zoomedInMap[i]))
	}

	var stack []position

	processPosition := func(pos position) {
		if pos.row < 0 || pos.row >= len(zoomedInMap) || pos.col < 0 || pos.col >= len(zoomedInMap[pos.row]) {
			return
		}
		if zoomedInMap[pos.row][pos.col] != ground {
			return
		}
		if seen[pos.row][pos.col] {
			return
		}

		stack = append(stack, pos)
		seen[pos.row][pos.col] = true

		// Because of the way we zoom in, original positions are those with odd
		// coordinates.
		if pos.row%2 == 1 && pos.col%2 == 1 {
			outsideArea++
		}
	}

	pos := position{0, 0} // guaranteed to be outside the loop because of padding
	stack = append(stack, pos)
	seen[pos.row][pos.col] = true

	for len(stack) > 0 {
		pos, stack = stack[len(stack)-1], stack[:len(stack)-1]

		processPosition(position{pos.row - 1, pos.col})
		processPosition(position{pos.row + 1, pos.col})
		processPosition(position{pos.row, pos.col - 1})
		processPosition(position{pos.row, pos.col + 1})
	}

	mapArea := len(pipeMap) * len(pipeMap[0])
	insideArea := mapArea - outsideArea - len(loop)

	return insideArea
}

func findLoop(pipeMap [][]byte) ([]position, error) {
	startingPosition, err := findStartingPosition(pipeMap)
	if err != nil {
		return nil, fmt.Errorf("could not find starting position: %w", err)
	}

	loop := []position{startingPosition}
	seen := map[position]bool{startingPosition: true}

	for {
		neighbors := findConnectedNeighbors(pipeMap, loop[len(loop)-1])
		if len(neighbors) != 2 {
			return nil, fmt.Errorf("stumbled upon position with %d neighbors: %#v", len(neighbors), loop[len(loop)-1])
		}

		for len(neighbors) > 0 && seen[neighbors[0]] {
			neighbors = neighbors[1:]
		}

		if len(neighbors) == 0 {
			break
		}

		loop = append(loop, neighbors[0])
		seen[neighbors[0]] = true
	}

	return loop, nil
}

func findStartingPosition(pipeMap [][]byte) (position, error) {
	for row := range pipeMap {
		for col := range pipeMap[row] {
			if pipeMap[row][col] == startingPosition {
				return position{row, col}, nil
			}
		}
	}

	return position{}, fmt.Errorf("no starting position found")
}

func findConnectedNeighbors(pipeMap [][]byte, pos position) []position {
	var neighbors []position

	shape := pipeMap[pos.row][pos.col]

	switch shape {
	case startingPosition:
		neighbors = findStartingPositionConnectedNeigbors(pipeMap, pos)
	case pipeVertical:
		if pos.row > 0 {
			neighbors = append(neighbors, position{pos.row - 1, pos.col})
		}
		if pos.row < len(pipeMap)-1 {
			neighbors = append(neighbors, position{pos.row + 1, pos.col})
		}
	case pipeHorizontal:
		if pos.col > 0 {
			neighbors = append(neighbors, position{pos.row, pos.col - 1})
		}
		if pos.col < len(pipeMap[pos.row])-1 {
			neighbors = append(neighbors, position{pos.row, pos.col + 1})
		}
	case pipeBendNorthEast:
		if pos.row > 0 {
			neighbors = append(neighbors, position{pos.row - 1, pos.col})
		}
		if pos.col < len(pipeMap[pos.row])-1 {
			neighbors = append(neighbors, position{pos.row, pos.col + 1})
		}
	case pipeBendNorthWest:
		if pos.row > 0 {
			neighbors = append(neighbors, position{pos.row - 1, pos.col})
		}
		if pos.col > 0 {
			neighbors = append(neighbors, position{pos.row, pos.col - 1})
		}
	case pipeBendSouthWest:
		if pos.row < len(pipeMap)-1 {
			neighbors = append(neighbors, position{pos.row + 1, pos.col})
		}
		if pos.col > 0 {
			neighbors = append(neighbors, position{pos.row, pos.col - 1})
		}
	case pipeBendSouthEast:
		if pos.row < len(pipeMap)-1 {
			neighbors = append(neighbors, position{pos.row + 1, pos.col})
		}
		if pos.col < len(pipeMap[pos.row])-1 {
			neighbors = append(neighbors, position{pos.row, pos.col + 1})
		}
	}

	return neighbors
}

// use this function only with the starting position as input. Its neighbors
// will tell us its shape.
func findStartingPositionConnectedNeigbors(pipeMap [][]byte, pos position) []position {
	var neighbors []position

	if pos.row > 0 && contains([]byte{pipeVertical, pipeBendSouthEast, pipeBendSouthWest}, pipeMap[pos.row-1][pos.col]) {
		neighbors = append(neighbors, position{pos.row - 1, pos.col})
	}
	if pos.row < len(pipeMap)-1 && contains([]byte{pipeVertical, pipeBendNorthEast, pipeBendNorthWest}, pipeMap[pos.row+1][pos.col]) {
		neighbors = append(neighbors, position{pos.row + 1, pos.col})
	}
	if pos.col > 0 && contains([]byte{pipeHorizontal, pipeBendNorthEast, pipeBendSouthEast}, pipeMap[pos.row][pos.col-1]) {
		neighbors = append(neighbors, position{pos.row, pos.col - 1})
	}
	if pos.col < len(pipeMap[pos.row])-1 && contains([]byte{pipeHorizontal, pipeBendNorthWest, pipeBendSouthWest}, pipeMap[pos.row][pos.col+1]) {
		neighbors = append(neighbors, position{pos.row, pos.col + 1})
	}

	return neighbors
}

func contains(values []byte, value byte) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}

func pipeMapFromReader(r io.Reader) ([][]byte, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	pipeMap := make([][]byte, len(lines))
	for i, line := range lines {
		pipeMap[i] = []byte(line)
	}

	if len(pipeMap) == 0 {
		return nil, fmt.Errorf("no input")
	}

	lineLength := len(pipeMap[0])
	for _, line := range pipeMap {
		if len(line) != lineLength {
			return nil, fmt.Errorf("input is not rectangular")
		}
	}

	for row := range pipeMap {
		for col := range pipeMap[row] {
			switch pipeMap[row][col] {
			case pipeVertical, pipeHorizontal, pipeBendNorthEast, pipeBendNorthWest, pipeBendSouthWest, pipeBendSouthEast, ground, startingPosition:
				// Valid
			default:
				return nil, fmt.Errorf("unknow character %q", pipeMap[row][col])
			}
		}
	}

	return pipeMap, nil
}
