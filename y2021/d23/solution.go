package busser

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 23 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	start, err := configurationFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	end := targetConfig(2)

	finder := newMinimumCostPathFinder()
	totalCost := finder.minimumTotalCost(start, end)

	_, err = fmt.Fprintf(answer, "%d", totalCost)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 23 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	start, err := configurationFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	start = unfoldConfig(start)
	end := targetConfig(4)

	finder := newMinimumCostPathFinder()
	totalCost := finder.minimumTotalCost(start, end)

	_, err = fmt.Fprintf(answer, "%d", totalCost)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type coordinates struct {
	x, y int
}

type space byte

const (
	empty space = iota
	wall
	amphipodA
	amphipodB
	amphipodC
	amphipodD
)

type configuration struct {
	roomSize int
	hallway  [11]space
	rooms    [4][4]space
}

func unfoldConfig(c configuration) configuration {
	c.roomSize = 4

	for r := range c.rooms {
		c.rooms[r][3] = c.rooms[r][1]
	}

	c.rooms[0][1] = amphipodD
	c.rooms[0][2] = amphipodD
	c.rooms[1][1] = amphipodC
	c.rooms[1][2] = amphipodB
	c.rooms[2][1] = amphipodB
	c.rooms[2][2] = amphipodA
	c.rooms[3][1] = amphipodA
	c.rooms[3][2] = amphipodC

	return c
}

func targetConfig(roomSize int) configuration {
	var c configuration
	c.roomSize = roomSize
	for i := range c.hallway {
		c.hallway[i] = empty
	}
	for y := 0; y < c.roomSize; y++ {
		c.rooms[0][y] = amphipodA
		c.rooms[1][y] = amphipodB
		c.rooms[2][y] = amphipodC
		c.rooms[3][y] = amphipodD
	}
	for r := range c.rooms {
		for y := c.roomSize; y < len(c.rooms[r]); y++ {
			c.rooms[r][y] = empty
		}
	}
	return c
}

type minimumCostPathFinder struct {
	reachableConfigs *priorityQueue
	minimumCostFound map[configuration]bool
}

func newMinimumCostPathFinder() minimumCostPathFinder {
	var finder minimumCostPathFinder

	finder.reachableConfigs = newPriorityQueue()
	finder.minimumCostFound = make(map[configuration]bool)

	return finder
}

func (finder minimumCostPathFinder) saveMinimumCost(config configuration, cost int) {
	finder.minimumCostFound[config] = true

	// Try to move amphipods that are in the hallway.
	for x := range config.hallway {
		if config.hallway[x] != empty {
			finder.canMoveAmphipodFromHallway(config, cost, x)
		}
	}

	// Try to move amphipods that are in a room.
	for room := range config.rooms {
		for spot := 0; spot < config.roomSize; spot++ {
			if config.rooms[room][spot] != empty {
				finder.canMoveAmphipodFromRoom(config, cost, room, spot)
				break
			}
		}
	}
}

func (finder minimumCostPathFinder) canMoveAmphipodFromHallway(config configuration, previousCost, hallwayX int) {
	amphipodType := config.hallway[hallwayX]

	var targetRoom int
	switch amphipodType {
	case amphipodA:
		targetRoom = 0
	case amphipodB:
		targetRoom = 1
	case amphipodC:
		targetRoom = 2
	case amphipodD:
		targetRoom = 3
	}

	currentX := 1 + hallwayX
	targetX := 3 + 2*targetRoom

	// The hallway between the amphipod and the room must be empty.
	if targetX > currentX {
		for i := currentX + 1; i <= targetX; i++ {
			if config.hallway[i-1] != empty {
				return
			}
		}
	} else {
		for i := currentX - 1; i >= targetX; i-- {
			if config.hallway[i-1] != empty {
				return
			}
		}
	}

	// There must be no amphipods of a different type in the room.
	for y := 0; y < config.roomSize; y++ {
		s := config.rooms[targetRoom][y]
		if s != empty && s != amphipodType {
			return
		}
	}

	targetSpot := -1
	for y := 0; y < config.roomSize; y++ {
		if config.rooms[targetRoom][y] != empty {
			break
		}
		targetSpot = y
	}
	if targetSpot == -1 {
		panic("the room is full of correct amphipods but another remains")
	}

	currentY := 1
	targetY := targetSpot + 2

	newConfig := config
	newConfig.hallway[hallwayX] = empty
	newConfig.rooms[targetRoom][targetSpot] = amphipodType

	cost := previousCost + costOfMovingAmphipod(amphipodType, abs(currentX-targetX)+abs(currentY-targetY))

	finder.canReach(newConfig, cost)
}

func (finder minimumCostPathFinder) canMoveAmphipodFromRoom(config configuration, previousCost, currentRoom, spot int) {
	amphipodType := config.rooms[currentRoom][spot]

	var targetRoom int
	switch amphipodType {
	case amphipodA:
		targetRoom = 0
	case amphipodB:
		targetRoom = 1
	case amphipodC:
		targetRoom = 2
	case amphipodD:
		targetRoom = 3
	}

	if currentRoom == targetRoom {
		allAmphipodsInRoomAreOK := true
		for s := spot; s < config.roomSize; s++ {
			if config.rooms[currentRoom][s] != amphipodType {
				allAmphipodsInRoomAreOK = false
				break
			}
		}
		if allAmphipodsInRoomAreOK {
			// All amphipods in the room are in the correct room.
			// None of these amphipods should move.
			return
		}
	}

	// The amphipod can move to any spot in the hallway that is not in front of
	// a room, as long as the path is unobstructed.
	currentX := 3 + 2*currentRoom
	currentY := 2 + spot
	targetY := 1

	for targetX := currentX - 1; targetX >= 1; targetX-- {
		if targetX%2 == 1 && targetX >= 3 && targetX <= 9 {
			continue
		}
		if config.hallway[targetX-1] != empty {
			break
		}
		newConfig := config
		newConfig.rooms[currentRoom][spot] = empty
		newConfig.hallway[targetX-1] = amphipodType
		cost := previousCost + costOfMovingAmphipod(amphipodType, abs(currentX-targetX)+abs(currentY-targetY))
		finder.canReach(newConfig, cost)
	}

	for targetX := currentX + 1; targetX <= 11; targetX++ {
		if targetX%2 == 1 && targetX >= 3 && targetX <= 9 {
			continue
		}
		if config.hallway[targetX-1] != empty {
			break
		}
		newConfig := config
		newConfig.rooms[currentRoom][spot] = empty
		newConfig.hallway[targetX-1] = amphipodType
		cost := previousCost + costOfMovingAmphipod(amphipodType, abs(currentX-targetX)+abs(currentY-targetY))
		finder.canReach(newConfig, cost)
	}
}

func costOfMovingAmphipod(typ space, distance int) int {
	switch typ {
	case amphipodA:
		return distance * 1
	case amphipodB:
		return distance * 10
	case amphipodC:
		return distance * 100
	case amphipodD:
		return distance * 1000
	default:
		panic("unknown amphipod type")
	}
}

func (finder minimumCostPathFinder) canReach(config configuration, risk int) {
	if finder.minimumCostFound[config] {
		return
	}
	finder.reachableConfigs.push(config, risk)
}

func (finder minimumCostPathFinder) reachableConfigWithMinimumCost() (config configuration, cost int) {
	if finder.reachableConfigs.len() == 0 {
		panic("no reachable configurations")
	}
	return finder.reachableConfigs.pop()
}

func (finder minimumCostPathFinder) minimumTotalCost(start, end configuration) (minCost int) {
	config, cost := start, 0
	for config != end {
		finder.saveMinimumCost(config, cost)
		config, cost = finder.reachableConfigWithMinimumCost()
	}
	return cost
}

func configurationFromReader(r io.Reader) (configuration, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return configuration{}, nil
	}

	if len(lines) != 5 {
		return configuration{}, fmt.Errorf("expected 5 lines, got %d", len(lines))
	}

	var config configuration

	config.roomSize = 2
	for i := range config.hallway {
		config.hallway[i] = empty
	}

	for spot := 0; spot < config.roomSize; spot++ {
		l := lines[spot+2]
		if len(l) < 9 {
			return configuration{}, errors.New("invalid input")
		}
		for room := 0; room < 4; room++ {
			switch l[3+2*room] {
			case 'A':
				config.rooms[room][spot] = amphipodA
			case 'B':
				config.rooms[room][spot] = amphipodB
			case 'C':
				config.rooms[room][spot] = amphipodC
			case 'D':
				config.rooms[room][spot] = amphipodD
			default:
				return configuration{}, fmt.Errorf("unknown amphipod type %q", l[2+2*room])
			}
		}
	}

	return config, nil
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
