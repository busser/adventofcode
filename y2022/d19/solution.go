package d19

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 19 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	blueprints, err := blueprintsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for _, b := range blueprints {
		maxGeodes := maxOpenedGeodes(b.robotCosts, 24)
		sum += b.id * maxGeodes
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 19 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	blueprints, err := blueprintsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	product := 1
	for i, b := range blueprints {
		if i >= 3 {
			break
		}

		maxGeodes := maxOpenedGeodes(b.robotCosts, 32)
		product *= maxGeodes
	}

	_, err = fmt.Fprintf(w, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type resource int

const (
	ore resource = iota
	clay
	obsidian
	geode
)

var resources = [...]resource{ore, clay, obsidian, geode}

type amount [len(resources)]int

func (a amount) minus(b amount) amount {
	for i := range a {
		a[i] -= b[i]
	}
	return a
}

func (a amount) canSubstract(b amount) bool {
	for i := range a {
		if a[i] < b[i] {
			return false
		}
	}
	return true
}

func (a amount) plus(b amount) amount {
	for i := range a {
		a[i] += b[i]
	}
	return a
}

type blueprint struct {
	id         int
	robotCosts [len(resources)]amount
}

type state struct {
	timeLeft  int
	resources amount
	robots    amount
}

func maxOpenedGeodes(costs [len(resources)]amount, timeAvailable int) int {
	// No need to produce more than the maximum cost of any robot.
	var maxRobots amount
	for _, cost := range costs {
		for res, count := range cost {
			maxRobots[res] = max(maxRobots[res], count)
		}
	}
	maxRobots[geode] = math.MaxInt

	maxGeodes := 0

	var searchForMax func(state, resource)
	searchForMax = func(s state, nextRobot resource) {
		// Skip searches that provision more robots than necessary
		if s.robots[nextRobot] >= maxRobots[nextRobot] {
			return
		}

		// This loop lets resources accumulate over time and branches into a
		// recursive search every time a new robot can be built
		for s.timeLeft > 0 {

			// Stop if we can't possibly beat the current maximum
			potential := s.resources[geode] + s.robots[geode]*s.timeLeft + s.timeLeft*(s.timeLeft+1)/2
			if potential < maxGeodes {
				return
			}

			// If we can afford the robot, build it and begin a recursive search
			cost := costs[nextRobot]
			if s.resources.canSubstract(cost) {
				next := s
				next.timeLeft--
				next.resources = next.resources.plus(next.robots).minus(cost)
				next.robots[nextRobot]++

				for _, nextRobot := range resources {
					searchForMax(next, nextRobot)
				}

				return
			}

			// Wait for more resources to accumulate before retrying to build
			s.timeLeft--
			s.resources = s.resources.plus(s.robots)

		}

		maxGeodes = max(maxGeodes, s.resources[geode])
	}

	start := state{
		timeLeft:  timeAvailable,
		resources: amount{0, 0, 0, 0},
		robots:    amount{1, 0, 0, 0},
	}

	for _, firstRobot := range resources {
		searchForMax(start, firstRobot)
	}

	return maxGeodes
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func blueprintsFromReader(r io.Reader) ([]blueprint, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	var blueprints []blueprint
	for _, l := range lines {
		b, err := blueprintFromString(l)
		if err != nil {
			return nil, err
		}

		blueprints = append(blueprints, b)
	}

	return blueprints, nil
}

func blueprintFromString(s string) (blueprint, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return blueprint{}, errors.New("wrong format")
	}

	rawID := strings.TrimPrefix(parts[0], "Blueprint ")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		return blueprint{}, fmt.Errorf("%q is not a number", rawID)
	}

	rawRobotCosts := strings.Split(parts[1], ".")
	var robotCosts [len(resources)]amount
	for _, rawRobotCost := range rawRobotCosts {
		if rawRobotCost == "" {
			continue
		}

		res, cost, err := robotCostFromString(rawRobotCost)
		if err != nil {
			return blueprint{}, err
		}

		robotCosts[res] = cost
	}

	b := blueprint{
		id:         id,
		robotCosts: robotCosts,
	}

	return b, nil
}

func robotCostFromString(s string) (resource, amount, error) {
	s = strings.TrimPrefix(s, " Each ")
	parts := strings.SplitN(s, " robot costs ", 2)
	if len(parts) != 2 {
		return 0, amount{}, errors.New("wrong format")
	}

	var res resource
	switch parts[0] {
	case "ore":
		res = ore
	case "clay":
		res = clay
	case "obsidian":
		res = obsidian
	case "geode":
		res = geode
	default:
		return 0, amount{}, fmt.Errorf("unknown resource %q", parts[0])
	}

	var cost amount
	for _, raw := range strings.Split(parts[1], " and ") {
		parts := strings.SplitN(raw, " ", 2)
		if len(parts) != 2 {
			return 0, amount{}, errors.New("wrong format")
		}

		count, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, amount{}, fmt.Errorf("%q is not a number", parts[0])
		}

		var res resource
		switch parts[1] {
		case "ore":
			res = ore
		case "clay":
			res = clay
		case "obsidian":
			res = obsidian
		case "geode":
			res = geode
		default:
			return 0, amount{}, fmt.Errorf("unknown resource %q", parts[1])
		}

		cost[res] = count
	}

	return res, cost, nil
}
