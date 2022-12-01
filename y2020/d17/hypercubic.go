package busser

type hypercubicSimulation struct {
	current, next [][][][]bool
}

func (sim *hypercubicSimulation) init(initialState [][]bool, margin int) {
	if margin < 0 {
		panic("margin cannot be negative")
	}

	minW, maxW, minX, maxX := bounds(initialState)

	newState := make([][][][]bool, maxW-minW+2*margin+1)
	for w := range newState {
		newState[w] = make([][][]bool, maxX-minX+2*margin+1)
		for x := range newState[w] {
			newState[w][x] = make([][]bool, 2*margin+1)
			for y := range newState[w][x] {
				newState[w][x][y] = make([]bool, 2*margin+1)
			}
		}
	}

	for x := range initialState {
		for y := range initialState[x] {
			newState[x+margin][y+margin][margin][margin] = initialState[x][y]
		}
	}

	sim.current = newState
	sim.next = copy4DState(newState)
}

func (sim *hypercubicSimulation) iterate() {
	for w := range sim.current {
		for x := range sim.current[w] {
			for y := range sim.current[w][x] {
				for z := range sim.current[w][x][y] {
					activeCells := sim.activeCellsAround(w, x, y, z)

					switch {
					case sim.current[w][x][y][z] && activeCells != 2 && activeCells != 3:
						sim.next[w][x][y][z] = false
					case !sim.current[w][x][y][z] && activeCells == 3:
						sim.next[w][x][y][z] = true
					default:
						sim.next[w][x][y][z] = sim.current[w][x][y][z]
					}
				}
			}
		}
	}

	sim.current, sim.next = sim.next, sim.current
}

func (sim hypercubicSimulation) activeCellsAround(w, x, y, z int) int {
	count := 0

	for ww := w - 1; ww <= w+1; ww++ {
		if ww < 0 || ww >= len(sim.current) {
			continue
		}

		for xx := x - 1; xx <= x+1; xx++ {
			if xx < 0 || xx >= len(sim.current[w]) {
				continue
			}

			for yy := y - 1; yy <= y+1; yy++ {
				if yy < 0 || yy >= len(sim.current[w][x]) {
					continue
				}

				for zz := z - 1; zz <= z+1; zz++ {
					if zz < 0 || zz >= len(sim.current[w][x][y]) {
						continue
					}
					if ww == w && xx == x && yy == y && zz == z {
						continue
					}

					if sim.current[ww][xx][yy][zz] {
						count++
					}
				}
			}
		}
	}

	return count
}

func (sim hypercubicSimulation) tally() int {
	count := 0

	for w := range sim.current {
		for x := range sim.current[w] {
			for y := range sim.current[w][x] {
				for _, active := range sim.current[w][x][y] {
					if active {
						count++
					}
				}
			}
		}
	}

	return count
}

func copy4DState(state [][][][]bool) [][][][]bool {
	new := make([][][][]bool, len(state))

	for w := range state {
		new[w] = make([][][]bool, len(state[w]))
		for x := range state[w] {
			new[w][x] = make([][]bool, len(state[w][x]))
			for y := range state[w][x] {
				new[w][x][y] = make([]bool, len(state[w][x][y]))
				copy(new[w][x][y], state[w][x][y])
			}
		}
	}

	return new
}
