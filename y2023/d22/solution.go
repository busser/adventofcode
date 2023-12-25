package d22

import (
	"fmt"
	"io"
	"sort"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 22 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	bricks, err := bricksFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sort.Sort(byLowest(bricks))
	dropBricks(bricks)

	count := 0
	for i := range bricks {
		if countHowManyFall(bricks, i) == 0 {
			count++
		}
	}

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 22 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	bricks, err := bricksFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sort.Sort(byLowest(bricks))
	dropBricks(bricks)

	sum := 0
	for i := range bricks {
		sum += countHowManyFall(bricks, i)
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type vector struct {
	x, y, z int
}

type brick struct {
	start, end vector // inclusive
}

func (b brick) String() string {
	return fmt.Sprintf("%d,%d,%d~%d,%d,%d",
		b.start.x, b.start.y, b.start.z, b.end.x, b.end.y, b.end.z)
}

// Sorted from lowest to highest, looking at each brick's lowest point.
type byLowest []brick

func (b byLowest) Len() int           { return len(b) }
func (b byLowest) Less(i, j int) bool { return b[i].start.z < b[j].start.z }
func (b byLowest) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func dropBricks(bricks []brick) {
	_, maxX, _, maxY := horizontalRange(bricks)

	supportLevel := make([][]int, maxX+1)
	for x := range supportLevel {
		supportLevel[x] = make([]int, maxY+1)
	}

	for i, b := range bricks {
		brickHeight := b.end.z - b.start.z + 1

		// Find the highest support level for this brick.
		maxSupport := 0
		for x := b.start.x; x <= b.end.x; x++ {
			for y := b.start.y; y <= b.end.y; y++ {
				maxSupport = max(maxSupport, supportLevel[x][y])
			}
		}

		// Drop the brick.
		b.start.z = maxSupport + 1
		b.end.z = b.start.z + brickHeight - 1

		// Update the support level.
		for x := b.start.x; x <= b.end.x; x++ {
			for y := b.start.y; y <= b.end.y; y++ {
				supportLevel[x][y] = b.end.z
			}
		}

		bricks[i] = b
	}
}

func countHowManyFall(bricks []brick, removed int) int {
	// Disintegrate the brick.
	withRemoved := make([]brick, len(bricks))
	copy(withRemoved, bricks)
	withRemoved = append(withRemoved[:removed], withRemoved[removed+1:]...)

	// Simulated the bricks falling.
	dropped := make([]brick, len(withRemoved))
	copy(dropped, withRemoved)
	dropBricks(dropped)

	// Count how many bricks moved.
	count := 0
	for i := range withRemoved {
		if withRemoved[i].start.z != dropped[i].start.z {
			count++
		}
	}

	return count
}

func horizontalRange(bricks []brick) (minX, maxX, minY, maxY int) {
	minX, maxX = bricks[0].start.x, bricks[0].end.x
	minY, maxY = bricks[0].start.y, bricks[0].end.y
	for _, b := range bricks {
		minX = min(minX, b.start.x)
		maxX = max(maxX, b.end.x)
		minY = min(minY, b.start.y)
		maxY = max(maxY, b.end.y)
	}
	return
}

func bricksFromReader(r io.Reader) ([]brick, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read bricks: %w", err)
	}

	bricks := make([]brick, len(lines))
	for i, line := range lines {
		bricks[i], err = brickFromString(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse brick %d: %w", i, err)
		}
	}

	return bricks, nil
}

func brickFromString(s string) (brick, error) {
	var b brick
	_, err := fmt.Sscanf(s, "%d,%d,%d~%d,%d,%d",
		&b.start.x, &b.start.y, &b.start.z, &b.end.x, &b.end.y, &b.end.z)
	if err != nil {
		return brick{}, fmt.Errorf("could not parse brick: %w", err)
	}

	// We check some assumptions about the input data. These assumptions are not
	// part of the problem statement, but they happen to be true and they make
	// the code simpler.

	if b.start.x > b.end.x {
		return brick{}, fmt.Errorf("start.x > end.x")
	}
	if b.start.y > b.end.y {
		return brick{}, fmt.Errorf("start.y > end.y")
	}
	if b.start.z > b.end.z {
		return brick{}, fmt.Errorf("start.z > end.z")
	}

	if b.start.x < 0 {
		return brick{}, fmt.Errorf("start.x < 0")
	}
	if b.start.y < 0 {
		return brick{}, fmt.Errorf("start.y < 0")
	}
	if b.start.z < 0 {
		return brick{}, fmt.Errorf("start.z < 0")
	}

	return b, nil
}
