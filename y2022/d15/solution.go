package d15

import (
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 15 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	sensors, err := sensorsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	possibleX := interval{math.MinInt, math.MaxInt}
	ruledOut, _ := ruledOutPositions(sensors, possibleX, 2_000_000)

	_, err = fmt.Fprintf(w, "%d", ruledOut)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 15 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	sensors, err := sensorsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	possibleX := interval{0, 4_000_000}
	possibleY := interval{0, 4_000_000}
	p, err := beakonPosition(sensors, possibleX, possibleY)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", p.x*4_000_000+p.y)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type position struct {
	x, y int
}

func (p position) distanceTo(other position) int {
	return abs(p.x-other.x) + abs(p.y-other.y)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type sensor struct {
	position
	closestBeacon position
}

// O(1)
func (s sensor) ruledOut(y int) interval {
	radius := s.distanceTo(s.closestBeacon)
	if abs(s.y-y) > radius {
		// The row is outside the radius.
		return interval{}
	}

	offset := radius - abs(s.y-y)
	start := s.x - offset
	end := s.x + offset

	return interval{start, end}
}

type interval struct {
	start, end int
}

type intervalLimit struct {
	isStart bool
	value   int
}

// O(len(sensors) * len(possibleY))
func beakonPosition(sensors []sensor, possibleX, possibleY interval) (position, error) {
	for y := possibleY.start; y <= possibleY.end; y++ {
		ruledOut, missingX := ruledOutPositions(sensors, possibleX, y)
		if ruledOut < possibleX.end-possibleX.start {
			return position{missingX, y}, nil
		}
	}

	return position{}, errors.New("not found")
}

// O(len(sensors))
func ruledOutPositions(sensors []sensor, possibleX interval, y int) (int, int) {
	var intervals []interval
	for _, s := range sensors {
		intervals = append(intervals, s.ruledOut(y))
	}

	var limits []intervalLimit
	for _, in := range intervals {
		limits = append(
			limits,
			intervalLimit{true, in.start},
			intervalLimit{false, in.end},
		)
	}

	sort.Slice(limits, func(i, j int) bool {
		if limits[i].value != limits[j].value {
			return limits[i].value < limits[j].value
		}
		return limits[i].isStart
	})

	ruledOut := 0

	lastMissingX := possibleX.start
	overlapping := 0
	var intervalStart int
	for _, limit := range limits {
		// if limit.value > possibleX.end {
		// 	break
		// }

		if limit.isStart {
			if overlapping == 0 {
				intervalStart = max(limit.value, possibleX.start)
				if intervalStart-1 > lastMissingX {
					lastMissingX = intervalStart - 1
				}
			}
			overlapping++
			continue
		}

		overlapping--
		if overlapping == 0 {
			ruledOut += min(limit.value, possibleX.end) - intervalStart
		}
	}

	return ruledOut, lastMissingX
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func sensorsFromReader(r io.Reader) ([]sensor, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	var sensors []sensor

	for _, l := range lines {
		s, err := sensorFromString(l)
		if err != nil {
			return nil, err
		}
		sensors = append(sensors, s)
	}

	if len(sensors) == 0 {
		return nil, errors.New("no sensors")
	}

	return sensors, nil
}

func sensorFromString(s string) (sensor, error) {
	parts := strings.SplitN(strings.TrimPrefix(s, "Sensor at "), ": closest beacon is at ", 2)
	if len(parts) != 2 {
		return sensor{}, errors.New("wrong format")
	}

	sp, err := positionFromString(parts[0])
	if err != nil {
		return sensor{}, err
	}

	bp, err := positionFromString(parts[1])
	if err != nil {
		return sensor{}, err
	}

	return sensor{
		position:      sp,
		closestBeacon: bp,
	}, nil
}

func positionFromString(s string) (position, error) {
	parts := strings.SplitN(s, ", ", 2)
	if len(parts) != 2 {
		return position{}, fmt.Errorf("invalid position %q", s)
	}

	x, err := strconv.Atoi(strings.TrimPrefix(parts[0], "x="))
	if err != nil {
		return position{}, fmt.Errorf("invalid x %q", parts[0])
	}

	y, err := strconv.Atoi(strings.TrimPrefix(parts[1], "y="))
	if err != nil {
		return position{}, fmt.Errorf("invalid y %q", parts[1])
	}

	return position{x, y}, nil
}
