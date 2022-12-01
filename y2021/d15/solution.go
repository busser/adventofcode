package busser

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 15 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	cave, err := caveFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	finder := newMinimumRiskPathFinder(cave)

	start := coordinates{ // top-left corner
		x: 0,
		y: 0,
	}
	end := coordinates{ // bottom-right corner
		x: len(cave[len(cave)-1]) - 1,
		y: len(cave) - 1,
	}

	totalRisk := finder.minimumTotalRisk(start, end)

	_, err = fmt.Fprintf(answer, "%d", totalRisk)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 15 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	cave, err := caveFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	cave = extendCave(cave)

	finder := newMinimumRiskPathFinder(cave)

	start := coordinates{ // top-left corner
		x: 0,
		y: 0,
	}
	end := coordinates{ // bottom-right corner
		x: len(cave[len(cave)-1]) - 1,
		y: len(cave) - 1,
	}

	totalRisk := finder.minimumTotalRisk(start, end)

	_, err = fmt.Fprintf(answer, "%d", totalRisk)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type coordinates struct {
	x, y int
}

type minimumRiskPathFinder struct {
	cave               [][]int
	reachablePositions *priorityQueue
	minimumRiskFound   [][]bool
}

func newMinimumRiskPathFinder(cave [][]int) minimumRiskPathFinder {
	var finder minimumRiskPathFinder

	finder.cave = cave
	finder.reachablePositions = newPriorityQueue()
	finder.minimumRiskFound = make([][]bool, len(cave))
	for y := range cave {
		finder.minimumRiskFound[y] = make([]bool, len(cave[y]))
	}

	return finder
}

func extendCave(cave [][]int) [][]int {
	const multiplier = 5
	newCave := make([][]int, multiplier*len(cave))
	for i := 0; i < multiplier; i++ {
		for y := 0; y < len(cave); y++ {
			newCave[i*len(cave)+y] = make([]int, multiplier*len(cave[y]))
			for j := 0; j < multiplier; j++ {
				for x := 0; x < len(cave[y]); x++ {
					value := cave[y][x] + i + j
					for value > 9 {
						value -= 9
					}
					newCave[i*len(cave)+y][j*len(cave[y])+x] = value
				}
			}
		}
	}
	return newCave
}

func (finder minimumRiskPathFinder) saveMinimumRisk(pos coordinates, risk int) {
	finder.minimumRiskFound[pos.y][pos.x] = true

	finder.canReach(pos.x-1, pos.y, risk)
	finder.canReach(pos.x+1, pos.y, risk)
	finder.canReach(pos.x, pos.y-1, risk)
	finder.canReach(pos.x, pos.y+1, risk)
}

func (finder minimumRiskPathFinder) canReach(x, y, parentRisk int) {
	if y < 0 || y >= len(finder.cave) || x < 0 || x >= len(finder.cave[y]) {
		return
	}
	if finder.minimumRiskFound[y][x] {
		return
	}

	risk := parentRisk + finder.cave[y][x]
	finder.reachablePositions.push(coordinates{x, y}, risk)
}

func (finder minimumRiskPathFinder) reachablePositionWithMinimumRisk() (position coordinates, risk int) {
	if finder.reachablePositions.len() == 0 {
		panic("no reachable positions")
	}

	return finder.reachablePositions.pop()
}

func (finder minimumRiskPathFinder) minimumTotalRisk(start, end coordinates) (minRisk int) {
	position, risk := start, 0
	for position != end {
		finder.saveMinimumRisk(position, risk)
		position, risk = finder.reachablePositionWithMinimumRisk()
	}
	finder.saveMinimumRisk(position, risk)
	return risk
}

func caveFromReader(r io.Reader) ([][]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	cave := make([][]int, len(lines))

	for y := range lines {
		cave[y] = make([]int, len(lines[y]))
		for x := range lines[y] {
			cave[y][x] = int(lines[y][x] - '0')
		}
	}

	return cave, nil
}
