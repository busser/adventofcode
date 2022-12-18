package d16

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

const startValve = "AA"

// PartOne solves the first problem of day 16 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	valves, err := valvesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	valves = densifyGraph(valves)

	var start valve
	for _, v := range valves {
		if v.name == startValve {
			start = v
			break
		}
	}
	if start.name == "" {
		return fmt.Errorf("no valve named %q", startValve)
	}

	released := mostPressurePossible(valves, start, 30, 1)

	_, err = fmt.Fprintf(w, "%d", released)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 16 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	valves, err := valvesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	valves = densifyGraph(valves)

	var start valve
	for _, v := range valves {
		if v.name == startValve {
			start = v
			break
		}
	}
	if start.name == "" {
		return fmt.Errorf("no valve named %q", startValve)
	}

	released := mostPressurePossible(valves, start, 26, 2)

	_, err = fmt.Fprintf(w, "%d", released)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

//=== Memoised DFS =============================================================

type stack[T any] []T

func (s *stack[T]) push(v T) {
	*s = append(*s, v)
}

func (s *stack[T]) pop() T {
	n := len(*s) - 1
	v := (*s)[n]
	*s = (*s)[:n]
	return v
}

func mostPressurePossible(valves []valve, startValve valve, timeAvailable, actors int) int {
	knownResults := make(map[stateCacheKey]int)

	var findBest func(state) int
	findBest = func(s state) int {
		if r, ok := knownResults[s.cacheKey()]; ok {
			return r
		}

		maxReleased := 0

		// Option A: open the valve
		if !s.openValves[s.location] && s.timeLeft >= 1 {
			next := s.copy()
			next.timeLeft--
			next.openValves[s.location] = true
			released := valves[next.location].flowRate * next.timeLeft
			maxReleased = max(maxReleased, released+findBest(next))
		}

		// Option B: move to another valve
		tunnels := valves[s.location].tunnels
		for _, t := range tunnels {
			if s.timeLeft < t.distance {
				continue
			}

			next := s.copy()
			next.timeLeft -= t.distance
			next.location = t.destination
			maxReleased = max(maxReleased, findBest(next))
		}

		// Option C: stop and send next actor
		if s.actor > 1 {
			next := s.copy()
			next.actor--
			next.timeLeft = timeAvailable
			next.location = startValve.id
			maxReleased = max(maxReleased, findBest(next))
		}

		knownResults[s.cacheKey()] = maxReleased
		return maxReleased
	}

	start := state{
		actor:      actors,
		timeLeft:   timeAvailable,
		location:   startValve.id,
		openValves: make([]bool, len(valves)),
	}

	return findBest(start)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type state struct {
	actor      int
	location   int
	timeLeft   int
	openValves []bool
}

type stateCacheKey struct {
	actor      int
	location   int
	timeLeft   int
	openValves string
}

func (s state) cacheKey() stateCacheKey {
	openValves := make([]byte, len(s.openValves))
	for i, b := range s.openValves {
		if b {
			openValves[i]++
		}
	}
	return stateCacheKey{
		actor:      s.actor,
		location:   s.location,
		timeLeft:   s.timeLeft,
		openValves: string(openValves),
	}
}

func (s state) copy() state {
	c := s

	c.openValves = make([]bool, len(s.openValves))
	copy(c.openValves, s.openValves)

	return c
}

//=== Graph densification ======================================================

const infinity = 1_000_000

func densifyGraph(valves []valve) []valve {

	// Step 1: compute all distances between valves with Floyd-Warshall

	distances := make([][]int, len(valves))
	for i := range distances {
		distances[i] = make([]int, len(valves))
		for j := range distances[i] {
			distances[i][j] = infinity
		}
	}

	for _, v := range valves {
		distances[v.id][v.id] = 0
		for _, t := range v.tunnels {
			distances[v.id][t.destination] = t.distance
		}
	}

	for k := range valves {
		for i := range valves {
			for j := range valves {
				if distances[i][j] > distances[i][k]+distances[k][j] {
					distances[i][j] = distances[i][k] + distances[k][j]
				}
			}
		}
	}

	// Step 2: rebuild graph without useless valves

	var newValves []valve
	var oldIDs []int

	for _, v := range valves {
		if v.flowRate == 0 && v.name != startValve {
			continue
		}

		oldIDs = append(oldIDs, v.id)
		v.id = len(newValves)
		newValves = append(newValves, v)
	}

	// Step 3: replace tunnels with new ones to make graph fully connected

	for i := range newValves {
		var newTunnels []tunnel

		for j := range newValves {
			if i == j {
				continue
			}

			t := tunnel{
				distance:    distances[oldIDs[i]][oldIDs[j]],
				destination: j,
			}

			newTunnels = append(newTunnels, t)
		}

		newValves[i].tunnels = newTunnels
	}

	return newValves
}

//=== Parsing ==================================================================

type valve struct {
	id       int
	name     string
	flowRate int
	tunnels  []tunnel
}

type tunnel struct {
	distance    int
	destination int
}

type rawValve struct {
	name      string
	flowRate  int
	neighbors []string
}

func valvesFromReader(r io.Reader) ([]valve, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	var rawValves []rawValve
	for _, l := range lines {
		v, err := rawValveFromString(l)
		if err != nil {
			return nil, err
		}

		rawValves = append(rawValves, v)
	}

	valves := make([]valve, len(rawValves))
	valveIDByName := make(map[string]int)

	for i, rv := range rawValves {
		valveIDByName[rv.name] = i
	}

	for i := range valves {
		valves[i].id = i
		valves[i].name = rawValves[i].name
		valves[i].flowRate = rawValves[i].flowRate
		valves[i].tunnels = make([]tunnel, len(rawValves[i].neighbors))
		for j, n := range rawValves[i].neighbors {
			valves[i].tunnels[j] = tunnel{1, valveIDByName[n]}
		}
	}

	return valves, nil
}

func rawValveFromString(s string) (rawValve, error) {
	parts := strings.SplitN(s, "; ", 2)
	if len(parts) != 2 {
		return rawValve{}, errors.New("wrong format")
	}

	subParts := strings.SplitN(parts[0], " has flow rate=", 2)
	if len(subParts) != 2 {
		return rawValve{}, errors.New("wrong format")
	}

	name := strings.TrimPrefix(subParts[0], "Valve ")

	flowRate, err := strconv.Atoi(subParts[1])
	if err != nil {
		return rawValve{}, fmt.Errorf("%q is not a number", subParts[1])
	}

	rawNeighbors := strings.TrimPrefix(parts[1], "tunnel leads to valve ")
	rawNeighbors = strings.TrimPrefix(rawNeighbors, "tunnels lead to valves ")
	neighbors := strings.Split(rawNeighbors, ", ")

	return rawValve{
		name:      name,
		flowRate:  flowRate,
		neighbors: neighbors,
	}, nil
}
