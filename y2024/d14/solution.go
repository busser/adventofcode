package d14

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 14 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	robots, err := robotsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	space := vector{101, 103}
	for i := range robots {
		robots[i].move(space, 100)
	}

	score := safetyScore(robots, space)

	_, err = fmt.Fprintf(w, "%d", score)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 14 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	robots, err := robotsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	space := vector{101, 103}
	seconds := 0
	for {
		for i := range robots {
			robots[i].move(space, 1)
		}
		seconds++
		if robotsFormChristmasTree(robots, space) {
			break
		}
	}

	_, err = fmt.Fprintf(w, "%d", seconds)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type vector struct {
	row, col int
}

func (v vector) plus(w vector) vector {
	return vector{
		row: v.row + w.row,
		col: v.col + w.col,
	}
}

func (v vector) times(n int) vector {
	return vector{
		row: v.row * n,
		col: v.col * n,
	}
}

func (v vector) modulo(w vector) vector {
	return vector{
		row: v.row % w.row,
		col: v.col % w.col,
	}
}

type robot struct {
	position vector
	velocity vector
}

func (r *robot) move(space vector, seconds int) {
	normalizedVelocity := r.velocity.plus(space).modulo(space)
	r.position = r.position.plus(normalizedVelocity.times(seconds)).modulo(space)
}

func safetyScore(robots []robot, space vector) int {
	var counts [4]int

	for _, r := range robots {
		p := r.position
		switch {
		case p.row < space.row/2 && p.col < space.col/2:
			counts[0]++
		case p.row < space.row/2 && p.col > space.col/2:
			counts[1]++
		case p.row > space.row/2 && p.col < space.col/2:
			counts[2]++
		case p.row > space.row/2 && p.col > space.col/2:
			counts[3]++
		}
	}

	score := 1
	for _, c := range counts {
		score *= c
	}

	return score
}

func robotsFormChristmasTree(robots []robot, space vector) bool {
	pixels := make([][]bool, space.row)
	for row := range pixels {
		pixels[row] = make([]bool, space.col)
	}

	for _, r := range robots {
		pixels[r.position.row][r.position.col] = true
	}

	for row := range pixels {
		consecutivePixels := 0
		for col := range pixels[row] {
			if pixels[row][col] {
				consecutivePixels++
			} else {
				consecutivePixels = 0
			}
			if consecutivePixels == 10 {
				return true
			}
		}
	}

	return false
}

func robotFromString(s string) (robot, error) {
	nums := helpers.IntsFromString(s)
	if len(nums) != 4 {
		return robot{}, fmt.Errorf("invalid robot: %q", s)
	}

	return robot{
		position: vector{nums[0], nums[1]},
		velocity: vector{nums[2], nums[3]},
	}, nil
}

func robotsFromReader(r io.Reader) ([]robot, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	robots := make([]robot, len(lines))
	for i, line := range lines {
		r, err := robotFromString(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse robot: %w", err)
		}
		robots[i] = r
	}

	return robots, nil
}
