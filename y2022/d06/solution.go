package d06

import (
	"errors"
	"fmt"
	"io"
)

// PartOne solves the first problem of day 6 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	datastream, err := datastreamFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	startOfPacket, err := findMarker(datastream, 4)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", startOfPacket)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 6 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	datastream, err := datastreamFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	startOfMessage, err := findMarker(datastream, 14)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", startOfMessage)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func findMarker(datastream []byte, size int) (int, error) {
	for i := 0; i < len(datastream)-size; i++ {
		if !hasDuplicates(datastream[i : i+size]) {
			return i + size, nil
		}
	}

	return 0, errors.New("not found")
}

func hasDuplicates(block []byte) bool {
	var seen [256]bool

	for _, b := range block {
		if seen[b] {
			return true
		}
		seen[b] = true
	}

	return false
}

func datastreamFromReader(r io.Reader) ([]byte, error) {
	datastream, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if len(datastream) < 4 {
		return nil, fmt.Errorf("datastream too short, only %d bytes", len(datastream))
	}

	return datastream, nil
}
