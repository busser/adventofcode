package busser

type cubicSimulation struct {
	current, next [][][]bool
}

func (sim *cubicSimulation) init(initialState [][]bool, margin int) {
	if margin < 0 {
		panic("margin cannot be negative")
	}

	minX, maxX, minY, maxY := bounds(initialState)

	newState := make([][][]bool, maxX-minX+2*margin+1)
	for x := range newState {
		newState[x] = make([][]bool, maxY-minY+2*margin+1)
		for y := range newState[x] {
			newState[x][y] = make([]bool, 2*margin+1)
		}
	}

	for x := range initialState {
		for y := range initialState[x] {
			newState[x+margin][y+margin][margin] = initialState[x][y]
		}
	}

	sim.current = newState
	sim.next = copy3DState(newState)
}

func (sim *cubicSimulation) iterate() {
	for x := range sim.current {
		for y := range sim.current[x] {
			for z := range sim.current[x][y] {
				activeCells := sim.activeCellsAround(x, y, z)
				switch {
				case sim.current[x][y][z] && activeCells != 2 && activeCells != 3:
					sim.next[x][y][z] = false
				case !sim.current[x][y][z] && activeCells == 3:
					sim.next[x][y][z] = true
				default:
					sim.next[x][y][z] = sim.current[x][y][z]
				}
			}
		}
	}

	sim.current, sim.next = sim.next, sim.current
}

func (sim cubicSimulation) activeCellsAround(x, y, z int) int {
	count := 0

	for xx := x - 1; xx <= x+1; xx++ {
		if xx < 0 || xx >= len(sim.current) {
			continue
		}

		for yy := y - 1; yy <= y+1; yy++ {
			if yy < 0 || yy >= len(sim.current[x]) {
				continue
			}

			for zz := z - 1; zz <= z+1; zz++ {
				if zz < 0 || zz >= len(sim.current[x][y]) {
					continue
				}
				if xx == x && yy == y && zz == z {
					continue
				}

				if sim.current[xx][yy][zz] {
					count++
				}
			}
		}
	}

	return count
}

func (sim cubicSimulation) tally() int {
	count := 0

	for x := range sim.current {
		for y := range sim.current[x] {
			for _, active := range sim.current[x][y] {
				if active {
					count++
				}
			}
		}
	}

	return count
}

func copy3DState(state [][][]bool) [][][]bool {
	new := make([][][]bool, len(state))

	for x := range state {
		new[x] = make([][]bool, len(state[x]))
		for y := range state[x] {
			new[x][y] = make([]bool, len(state[x][y]))
			copy(new[x][y], state[x][y])
		}
	}

	return new
}
