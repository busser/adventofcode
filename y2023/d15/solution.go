package d15

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 15 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	steps, err := stepsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	hashSum := 0
	for _, step := range steps {
		hashSum += hash(step.raw)
	}

	_, err = fmt.Fprintf(w, "%d", hashSum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 15 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	steps, err := stepsFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var ht hashTable
	for _, step := range steps {
		ht.set(step.label, step.focalLength)
	}

	focusingPower := ht.focusingPower()

	_, err = fmt.Fprintf(w, "%d", focusingPower)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const numBoxes = 256

type hashTable struct {
	boxes [numBoxes][]slot
}

type slot struct {
	label       string
	focalLength int
}

func (ht *hashTable) set(label string, focalLength int) {
	labelHash := hash(label)

	if focalLength == 0 {
		for i := range ht.boxes[labelHash] {
			if ht.boxes[labelHash][i].label == label {
				ht.boxes[labelHash] = append(ht.boxes[labelHash][:i], ht.boxes[labelHash][i+1:]...)
				return
			}
		}
		// Not found, nothing to delete
		return
	}

	// If the label already exists, update the focal length.
	for i := range ht.boxes[labelHash] {
		if ht.boxes[labelHash][i].label == label {
			ht.boxes[labelHash][i].focalLength = focalLength
			return
		}
	}

	// If the label does not exist, append it.
	ht.boxes[labelHash] = append(ht.boxes[labelHash], slot{
		label:       label,
		focalLength: focalLength,
	})
}

func (ht hashTable) focusingPower() int {
	total := 0

	for boxNumber, box := range ht.boxes {
		for slotNumber, slot := range box {
			total += (boxNumber + 1) * (slotNumber + 1) * slot.focalLength
		}
	}

	return total
}

func hash(s string) int {
	v := 0
	for i := range s {
		v = (v + int(s[i])) * 17 % numBoxes
	}
	return v
}

type step struct {
	raw         string
	label       string
	focalLength int // 0 means delete
}

func stepsFromReader(r io.Reader) ([]step, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("expected 1 line, got %d", len(lines))
	}

	rawSteps := strings.Split(lines[0], ",")

	steps := make([]step, len(rawSteps))
	for i := range steps {
		steps[i], err = stepFromString(rawSteps[i])
		if err != nil {
			return nil, fmt.Errorf("could not parse step %d: %w", i, err)
		}
	}

	return steps, nil
}

func stepFromString(s string) (step, error) {
	if len(s) < 2 {
		return step{}, fmt.Errorf("step must be at least 2 characters long")
	}

	if s[len(s)-1] == '-' {
		return step{
			raw:         s,
			label:       s[0 : len(s)-1],
			focalLength: 0,
		}, nil
	}

	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return step{}, fmt.Errorf("step must contain exactly one '='")
	}

	focalLength, err := strconv.Atoi(parts[1])
	if err != nil {
		return step{}, fmt.Errorf("could not parse focal length: %w", err)
	}

	return step{
		raw:         s,
		label:       parts[0],
		focalLength: focalLength,
	}, nil
}
