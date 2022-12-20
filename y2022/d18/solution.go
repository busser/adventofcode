package d18

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 18 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	cubes, err := cubesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	area := totalArea(cubes, false)

	_, err = fmt.Fprintf(w, "%d", area)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 18 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	cubes, err := cubesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	area := totalArea(cubes, true)

	_, err = fmt.Fprintf(w, "%d", area)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type vector struct {
	x, y, z int
}

func totalArea(cubes []vector, onlyExposed bool) int {
	size := bounds(cubes)

	lava := make([][][]bool, size.x+2)
	for x := range lava {
		lava[x] = make([][]bool, size.y+2)
		for y := range lava[x] {
			lava[x][y] = make([]bool, size.z+2)
		}
	}

	for _, c := range cubes {
		lava[c.x+1][c.y+1][c.z+1] = true
	}

	water := make([][][]bool, size.x+2)
	for x := range water {
		water[x] = make([][]bool, size.y+2)
		for y := range water[x] {
			water[x][y] = make([]bool, size.z+2)
		}
	}

	if !onlyExposed {
		// Shortcut: assume water is everywhere lava isn't.
		for x := range water {
			for y := range water[x] {
				for z := range water[x][y] {
					if !lava[x][y][z] {
						water[x][y][z] = true
					}
				}
			}
		}
	} else {
		// Use DFS to put water on entire outside surface of lava.
		var placeWater func(x, y, z int)
		placeWater = func(x, y, z int) {
			if water[x][y][z] || lava[x][y][z] {
				return
			}
			water[x][y][z] = true

			if x > 0 {
				placeWater(x-1, y, z)
			}
			if x < len(water)-1 {
				placeWater(x+1, y, z)
			}

			if y > 0 {
				placeWater(x, y-1, z)
			}
			if y < len(water[x])-1 {
				placeWater(x, y+1, z)
			}

			if z > 0 {
				placeWater(x, y, z-1)
			}
			if z < len(water[x][y])-1 {
				placeWater(x, y, z+1)
			}
		}

		placeWater(0, 0, 0)
	}

	area := 0
	for x := range lava {
		for y := range lava[x] {
			for z := range lava[x][y] {
				if !lava[x][y][z] {
					continue
				}

				if water[x-1][y][z] {
					area++
				}
				if water[x+1][y][z] {
					area++
				}

				if water[x][y-1][z] {
					area++
				}
				if water[x][y+1][z] {
					area++
				}

				if water[x][y][z-1] {
					area++
				}
				if water[x][y][z+1] {
					area++
				}
			}
		}
	}

	return area
}

func bounds(vv []vector) vector {
	var maxX, maxY, maxZ int

	for _, v := range vv {
		maxX = max(maxX, v.x)
		maxY = max(maxY, v.y)
		maxZ = max(maxZ, v.z)
	}

	return vector{maxX + 1, maxY + 1, maxZ + 1}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func cubesFromReader(r io.Reader) ([]vector, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	var cubes []vector
	for _, l := range lines {
		rawCube, err := helpers.IntsFromString(l, ",")
		if err != nil {
			return nil, err
		}
		if len(rawCube) != 3 {
			return nil, errors.New("wrong format")
		}
		for _, n := range rawCube {
			if n < 0 {
				return nil, errors.New("solution assumes no negative coordinates")
			}
		}
		cube := vector{
			x: rawCube[0],
			y: rawCube[1],
			z: rawCube[2],
		}
		cubes = append(cubes, cube)
	}

	return cubes, nil
}
