package d24

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 24 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	hailstones, err := hailstonesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := countCrossedPaths(hailstones, 200_000_000_000_000, 400_000_000_000_000)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 24 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	hailstones, err := hailstonesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	rockVelocity, err := findRockVelocity(hailstones)
	if err != nil {
		return fmt.Errorf("could not find rock's velocity: %w", err)
	}

	rockPosition, err := findRockPosition(hailstones, rockVelocity)
	if err != nil {
		return fmt.Errorf("could not find rock's position: %w", err)
	}

	sum := rockPosition[0] + rockPosition[1] + rockPosition[2]

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type vector[T int | float64] [3]T

func (v vector[T]) minus(w vector[T]) vector[T] {
	return vector[T]{
		v[0] - w[0],
		v[1] - w[1],
		v[2] - w[2],
	}
}

type hailstone struct {
	position vector[int]
	velocity vector[int]
}

func floatify(v vector[int]) vector[float64] {
	return vector[float64]{
		float64(v[0]),
		float64(v[1]),
		float64(v[2]),
	}
}

func (v vector[T]) plus(w vector[T]) vector[T] {
	return vector[T]{
		v[0] + w[0],
		v[1] + w[1],
		v[2] + w[2],
	}
}

func (v vector[T]) times(n T) vector[T] {
	return vector[T]{
		v[0] * n,
		v[1] * n,
		v[2] * n,
	}
}

func (h hailstone) String() string {
	return fmt.Sprintf("%d, %d, %d @ %d, %d, %d",
		h.position[0], h.position[1], h.position[2],
		h.velocity[0], h.velocity[1], h.velocity[2])
}

func findRockVelocityForDimension(hailstones []hailstone, dimension int) (int, error) {
	hailstonesByXVelocity := make(map[int][]hailstone)
	for _, h := range hailstones {
		hailstonesByXVelocity[h.velocity[dimension]] = append(hailstonesByXVelocity[h.velocity[dimension]], h)
	}

	deltasByXVelocity := make(map[int][]int)

	for vx, hs := range hailstonesByXVelocity {
		for i := 0; i < len(hs); i++ {
			for j := i + 1; j < len(hs); j++ {
				delta := hs[j].position[dimension] - hs[i].position[dimension]
				deltasByXVelocity[vx] = append(deltasByXVelocity[vx], delta)
			}
		}
	}

rockVelocity:
	for rockV := -1000; rockV <= 1000; rockV++ {
		if rockV == 0 {
			continue
		}

		for vx, deltas := range deltasByXVelocity {
			if vx == rockV {
				continue
			}

			for _, delta := range deltas {
				if delta%(rockV-vx) != 0 {
					continue rockVelocity
				}
			}
		}

		return rockV, nil
	}

	return 0, fmt.Errorf("not enough information in input")
}

func findRockVelocity(hailstones []hailstone) (vector[int], error) {
	vx, err := findRockVelocityForDimension(hailstones, 0)
	if err != nil {
		return vector[int]{}, fmt.Errorf("could not find rock's x velocity: %w", err)
	}

	vy, err := findRockVelocityForDimension(hailstones, 1)
	if err != nil {
		return vector[int]{}, fmt.Errorf("could not find rock's y velocity: %w", err)
	}

	vz, err := findRockVelocityForDimension(hailstones, 2)
	if err != nil {
		return vector[int]{}, fmt.Errorf("could not find rock's z velocity: %w", err)
	}

	return vector[int]{
		vx,
		vy,
		vz,
	}, nil
}

func findRockPosition(hailstones []hailstone, rockVelocity vector[int]) (vector[int], error) {
	// We define a shifted view of the hailstones, where the rock does not move.
	// All hailstones must at some point move through the same position: this is
	// the position of rock in this shifted view.
	//
	// We then determine how long a given hailstone takes to reach the rock.
	// Based on this time T, combined with the rock's position at time T (the
	// rock is static in this view), we can determine the rock's position at
	// time 0 in the original view.

	shiftedHailstones := make([]hailstone, len(hailstones))
	for i, h := range hailstones {
		shiftedHailstones[i] = hailstone{
			position: h.position,
			velocity: h.velocity.minus(rockVelocity),
		}
	}

	equations := make([]vector[float64], len(shiftedHailstones))
	for i, h := range shiftedHailstones {
		equations[i] = equationOf(h)
	}

	intersection, intersect := intersectionOf(equations[0], equations[1])
	if !intersect {
		return vector[int]{}, errors.New("mathematically impossible")
	}

	timeToReachRock := (int(intersection[0]) - shiftedHailstones[0].position[0]) / shiftedHailstones[0].velocity[0]

	rockPositionAtIntersect := hailstones[0].position.
		plus(hailstones[0].velocity.times(timeToReachRock))

	rockPosition := rockPositionAtIntersect.
		minus(rockVelocity.times(timeToReachRock))

	return rockPosition, nil
}

func countCrossedPaths(hailstones []hailstone, testMin, testMax float64) int {
	equations := make([]vector[float64], len(hailstones))
	for i, h := range hailstones {
		equations[i] = equationOf(h)
	}

	crossed := 0

	for i := 0; i < len(hailstones); i++ {
		for j := i + 1; j < len(hailstones); j++ {
			h1, h2 := hailstones[i], hailstones[j]

			intersection, intersect := intersectionOf(equations[i], equations[j])
			if !intersect {
				continue
			}

			if (h1.velocity[0] > 0 && float64(h1.position[0]) > intersection[0]) ||
				(h1.velocity[0] < 0 && float64(h1.position[0]) < intersection[0]) ||
				(h2.velocity[0] > 0 && float64(h2.position[0]) > intersection[0]) ||
				(h2.velocity[0] < 0 && float64(h2.position[0]) < intersection[0]) {
				// intersection happened in the past.
				continue
			}

			if intersection[0] < testMin || intersection[0] > testMax ||
				intersection[1] < testMin || intersection[1] > testMax {
				// intersection happend outside of the test area.
				continue
			}

			crossed++
		}
	}

	return crossed
}

func equationOf(h hailstone) vector[float64] {
	position := floatify(h.position)
	velocity := floatify(h.velocity)

	a := velocity[1] / velocity[0]
	b := position[1] - a*position[0]
	return vector[float64]{
		a,
		b,
		0,
	}
}

func intersectionOf(eq1, eq2 vector[float64]) (vector[float64], bool) {
	if eq1[0] == eq2[0] {
		return vector[float64]{}, false
	}

	x := (eq2[1] - eq1[1]) / (eq1[0] - eq2[0])
	y := eq1[0]*x + eq1[1]

	return vector[float64]{
		x,
		y,
		0,
	}, true
}

func hailstonesFromReader(r io.Reader) ([]hailstone, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	var hailstones []hailstone
	for _, line := range lines {
		h, err := hailstoneFromString(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse hailstone: %w", err)
		}
		hailstones = append(hailstones, h)
	}

	return hailstones, nil
}

func hailstoneFromString(s string) (hailstone, error) {
	var h hailstone
	_, err := fmt.Sscanf(s, "%d, %d, %d @ %d, %d, %d",
		&h.position[0], &h.position[1], &h.position[2],
		&h.velocity[0], &h.velocity[1], &h.velocity[2])
	if err != nil {
		return hailstone{}, fmt.Errorf("could not parse hailstone: %w", err)
	}
	return h, nil
}
