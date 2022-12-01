package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 22 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	steps, err := rebootStepsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var boundedSteps []rebootStep
	for i := range steps {
		if outOfBounds := steps[i].cuboid.bound(-50, 50); !outOfBounds {
			boundedSteps = append(boundedSteps, steps[i])
		}
	}

	cubesOn := reboot(boundedSteps)

	_, err = fmt.Fprintf(answer, "%d", cubesOn)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 22 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	steps, err := rebootStepsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	cubesOn := reboot(steps)

	_, err = fmt.Fprintf(answer, "%d", cubesOn)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type limits struct {
	min, max int
}

func (l *limits) bound(low, high int) (outOfBounds bool) {
	if (l.min < low && l.max < low) || (l.min > high && l.max > high) {
		return true
	}
	l.min = min(max(l.min, low), high)
	l.max = min(max(l.max, low), high)
	return false
}

func (l *limits) length() int {
	return l.max - l.min + 1
}

type box struct {
	x, y, z limits
}

func (b *box) bound(low, high int) (outOfBounds bool) {
	return b.x.bound(low, high) || b.y.bound(low, high) || b.z.bound(low, high)
}

func (this box) intersectionWith(other box) box {
	var i box
	i.x.min = max(this.x.min, other.x.min)
	i.x.max = min(this.x.max, other.x.max)
	i.y.min = max(this.y.min, other.y.min)
	i.y.max = min(this.y.max, other.y.max)
	i.z.min = max(this.z.min, other.z.min)
	i.z.max = min(this.z.max, other.z.max)
	return i
}

func (b *box) hasVolume() bool {
	return b.x.length() > 0 && b.y.length() > 0 && b.z.length() > 0
}

func (b *box) volume() int {
	return b.x.length() * b.y.length() * b.z.length()
}

type rebootStep struct {
	turnOn bool
	cuboid box
}

func reboot(steps []rebootStep) (cubesOn int) {
	boxWeights := make(map[box]int)

	for _, s := range steps {
		weightChanges := make(map[box]int)
		for existingBox, weight := range boxWeights {
			i := s.cuboid.intersectionWith(existingBox)
			if i.hasVolume() {
				weightChanges[i] -= weight
			}
		}
		if s.turnOn {
			weightChanges[s.cuboid]++
		}

		for b, w := range weightChanges {
			boxWeights[b] += w
		}
	}

	cubesOn = 0
	for b, w := range boxWeights {
		cubesOn += b.volume() * w
	}
	return cubesOn
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

func rebootStepsFromReader(r io.Reader) ([]rebootStep, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	steps := make([]rebootStep, len(lines))

	for i := range lines {
		s, err := rebootStepFromString(lines[i])
		if err != nil {
			return nil, fmt.Errorf("invalid step at line %d: %w", i+1, err)
		}
		steps[i] = s
	}

	return steps, nil
}

func rebootStepFromString(str string) (rebootStep, error) {
	parts := strings.SplitN(str, " ", 2)
	if len(parts) != 2 {
		return rebootStep{}, errors.New("invalid syntax")
	}

	var step rebootStep

	switch parts[0] {
	case "on":
		step.turnOn = true
	case "off":
		step.turnOn = false
	default:
		return rebootStep{}, errors.New("invalid syntax")
	}

	cuboid, err := boxFromString(parts[1])
	if err != nil {
		return rebootStep{}, fmt.Errorf("invalid cuboid: %w", err)
	}

	step.cuboid = cuboid

	return step, nil
}

func boxFromString(str string) (box, error) {
	var b box

	parts := strings.SplitN(str, ",", 3)
	if len(parts) != 3 {
		return box{}, fmt.Errorf("expected 3 ranges, got %d", len(parts))
	}

	var err error
	b.x, err = limitsFromString(parts[0])
	if err != nil {
		return box{}, fmt.Errorf("invalid limit %q: %w", parts[0], err)
	}
	b.y, err = limitsFromString(parts[1])
	if err != nil {
		return box{}, fmt.Errorf("invalid limit %q: %w", parts[1], err)
	}
	b.z, err = limitsFromString(parts[2])
	if err != nil {
		return box{}, fmt.Errorf("invalid limit %q: %w", parts[2], err)
	}

	return b, nil
}

func limitsFromString(str string) (limits, error) {
	parts := strings.SplitN(str, "=", 2)
	if len(parts) != 2 {
		return limits{}, errors.New("invalid syntax")
	}

	parts = strings.SplitN(parts[1], "..", 2)
	if len(parts) != 2 {
		return limits{}, fmt.Errorf("expected 2 parts, got %d", len(parts))
	}

	min, err := strconv.Atoi(parts[0])
	if err != nil {
		return limits{}, fmt.Errorf("%q is not a number", parts[0])
	}
	max, err := strconv.Atoi(parts[1])
	if err != nil {
		return limits{}, fmt.Errorf("%q is not a number", parts[1])
	}

	return limits{min, max}, nil
}
