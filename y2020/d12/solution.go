package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 12 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	instructions, err := instructionsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	b := boat{
		direction: 'E',
		pos:       position{0, 0},
	}

	if err := b.followInstructions(instructions); err != nil {
		return fmt.Errorf("following instructions: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", abs(b.pos.east)+abs(b.pos.north))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 12 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	instructions, err := instructionsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	b := boat{
		pos:      position{0, 0},
		waypoint: position{10, 1},
	}

	if err := b.followInstructionsCorrectly(instructions); err != nil {
		return fmt.Errorf("following instructions: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", abs(b.pos.east)+abs(b.pos.north))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func abs(n int) int {
	if n > 0 {
		return n
	}
	return -n
}

type instruction struct {
	action byte
	value  int
}

type position struct {
	east, north int
}

type boat struct {
	direction byte
	pos       position
	waypoint  position
}

func (b *boat) turnLeft(degrees int) {
	for ; degrees > 0; degrees -= 90 {
		switch b.direction {
		case 'E':
			b.direction = 'N'
		case 'N':
			b.direction = 'W'
		case 'W':
			b.direction = 'S'
		case 'S':
			b.direction = 'E'
		default:
			panic(fmt.Sprintf("unknown direction %q", b.direction))
		}
	}
}

func (b *boat) turnRight(degrees int) {
	for ; degrees > 0; degrees -= 90 {
		switch b.direction {
		case 'E':
			b.direction = 'S'
		case 'S':
			b.direction = 'W'
		case 'W':
			b.direction = 'N'
		case 'N':
			b.direction = 'E'
		default:
			panic(fmt.Sprintf("unknown direction %q", b.direction))
		}
	}
}

func (b *boat) moveInDirection(direction byte, distance int) {
	switch direction {
	case 'N':
		b.pos.north += distance
	case 'S':
		b.pos.north -= distance
	case 'E':
		b.pos.east += distance
	case 'W':
		b.pos.east -= distance
	default:
		panic(fmt.Sprintf("unknown direction %q", b.direction))
	}
}

func (b *boat) moveForward(distance int) {
	b.moveInDirection(b.direction, distance)
}

func (b *boat) followInstructions(instructions []instruction) error {
	for _, inst := range instructions {
		switch inst.action {
		case 'N', 'S', 'E', 'W':
			b.moveInDirection(inst.action, inst.value)
		case 'L':
			b.turnLeft(inst.value)
		case 'R':
			b.turnRight(inst.value)
		case 'F':
			b.moveForward(inst.value)
		default:
			return fmt.Errorf("unknown action %q", inst.action)
		}
	}

	return nil
}

func (b *boat) rotateWaypointLeft(degrees int) {
	for ; degrees > 0; degrees -= 90 {
		b.waypoint.east, b.waypoint.north = -b.waypoint.north, b.waypoint.east
	}
}

func (b *boat) rotateWaypointRight(degrees int) {
	for ; degrees > 0; degrees -= 90 {
		b.waypoint.east, b.waypoint.north = b.waypoint.north, -b.waypoint.east
	}
}

func (b *boat) moveWaypointInDirection(direction byte, distance int) {
	switch direction {
	case 'N':
		b.waypoint.north += distance
	case 'S':
		b.waypoint.north -= distance
	case 'E':
		b.waypoint.east += distance
	case 'W':
		b.waypoint.east -= distance
	default:
		panic(fmt.Sprintf("unknown direction %q", b.direction))
	}
}

func (b *boat) moveTowardsWaypoint(times int) {
	b.pos.east += b.waypoint.east * times
	b.pos.north += b.waypoint.north * times
}

func (b *boat) followInstructionsCorrectly(instructions []instruction) error {
	for _, inst := range instructions {
		switch inst.action {
		case 'N', 'S', 'E', 'W':
			b.moveWaypointInDirection(inst.action, inst.value)
		case 'L':
			b.rotateWaypointLeft(inst.value)
		case 'R':
			b.rotateWaypointRight(inst.value)
		case 'F':
			b.moveTowardsWaypoint(inst.value)
		default:
			return fmt.Errorf("unknown action %q", inst.action)
		}
	}

	return nil
}

func instructionsFromReader(r io.Reader) ([]instruction, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading lines: %w", err)
	}

	instructions := make([]instruction, len(lines))
	for i, line := range lines {
		if len(line) < 2 {
			return nil, errors.New("wrong format")
		}

		action := line[0]

		value, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, fmt.Errorf("wrong format: not an int: %w", err)
		}

		if (action == 'L' || action == 'R') && value%90 != 0 {
			return nil, errors.New("when action is a turn, value must be multiple of 90")
		}

		instructions[i].action = action
		instructions[i].value = value
	}

	return instructions, nil
}
