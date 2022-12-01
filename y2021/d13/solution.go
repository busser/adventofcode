package busser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// PartOne solves the first problem of day 13 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	dots, foldInstructions, err := dotsAndFoldsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	if len(foldInstructions) == 0 {
		return errors.New("no fold instructions")
	}

	dots.fold(foldInstructions[0])

	_, err = fmt.Fprintf(answer, "%d", len(dots))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 13 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	dots, foldInstructions, err := dotsAndFoldsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	for _, f := range foldInstructions {
		dots.fold(f)
	}

	_, err = fmt.Fprintf(answer, "%v", &dots)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type coordinates struct {
	x, y int
}

type dots map[coordinates]struct{}

func (dd *dots) fold(line coordinates) {
	for d := range *dd {
		switch {
		case line.x != 0 && d.x > line.x:
			newDot := coordinates{line.x*2 - d.x, d.y}
			delete(*dd, d)
			(*dd)[newDot] = struct{}{}
		case line.y != 0 && d.y > line.y:
			newDot := coordinates{d.x, line.y*2 - d.y}
			delete(*dd, d)
			(*dd)[newDot] = struct{}{}
		}
	}
}

func (dd *dots) String() string {
	var maxX, maxY int
	for d := range *dd {
		if d.x > maxX {
			maxX = d.x
		}
		if d.y > maxY {
			maxY = d.y
		}
	}

	var sb strings.Builder

	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			if _, hasDot := (*dd)[coordinates{x, y}]; hasDot {
				sb.WriteByte('#')
			} else {
				sb.WriteByte('.')
			}
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

func dotsAndFoldsFromReader(r io.Reader) (dd dots, foldLines []coordinates, err error) {
	rawData, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}
	rawData = bytes.TrimSpace(rawData)

	chunks := bytes.Split(rawData, []byte("\n\n"))
	if len(chunks) != 2 {
		return nil, nil, errors.New("invalid input")
	}

	dd = make(map[coordinates]struct{})

	for _, rawDot := range bytes.Split(chunks[0], []byte("\n")) {
		parts := bytes.Split(rawDot, []byte(","))
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("invalid dot coordinates %q", rawDot)
		}

		x, err := strconv.Atoi(string(parts[0]))
		if err != nil {
			return nil, nil, fmt.Errorf("%q is not a number", parts[0])
		}
		y, err := strconv.Atoi(string(parts[1]))
		if err != nil {
			return nil, nil, fmt.Errorf("%q is not a number", parts[1])
		}

		d := coordinates{x, y}
		dd[d] = struct{}{}
	}

	for _, rawFold := range bytes.Split(chunks[1], []byte("\n")) {
		rawCoords := bytes.TrimPrefix(rawFold, []byte("fold along "))

		parts := bytes.SplitN(rawCoords, []byte("="), 2)
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("invalid fold instruction %q", rawFold)
		}

		if len(parts[0]) != 1 {
			return nil, nil, fmt.Errorf("invalid fold instruction %q", rawFold)
		}

		value, err := strconv.Atoi(string(parts[1]))
		if err != nil {
			return nil, nil, fmt.Errorf("%q is not a number", parts[1])
		}

		var line coordinates

		switch parts[0][0] {
		case 'x':
			line.x = value
		case 'y':
			line.y = value
		default:
			return nil, nil, fmt.Errorf("invalid fold instruction %q", rawFold)
		}

		foldLines = append(foldLines, line)
	}

	return dd, foldLines, nil
}
