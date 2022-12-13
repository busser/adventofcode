package d06

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/bits"
)

// PartOneXOR solves the first problem of day 6 of Advent of Code 2022.
func PartOneXOR(r io.Reader, w io.Writer) error {
	datastream, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	startOfPacket, err := findMarkerXOR(datastream, 4)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", startOfPacket)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwoXOR solves the second problem of day 6 of Advent of Code 2022.
func PartTwoXOR(r io.Reader, w io.Writer) error {
	datastream, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	startOfMessage, err := findMarkerXOR(datastream, 14)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", startOfMessage)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func findMarkerXOR(datastream []byte, size int) (int, error) {
	// This algorithm is not mine, it comes from this great blog post:
	// https://www.mattkeeter.com/blog/2022-12-10-xor/

	datastream = bytes.TrimSpace(datastream)

	if len(datastream) < size {
		return 0, errors.New("datastream too short")
	}

	for _, b := range datastream {
		// The algorithm depends on there being 26 possible byte values.
		if b < 'a' || b > 'z' {
			return 0, fmt.Errorf("byte %q is out of range", b)
		}
	}

	var set uint32

	for i := 0; i < size; i++ {
		set ^= 1 << uint32(datastream[i]-'a')
	}

	for i := size; i < len(datastream); i++ {
		set ^= 1 << uint32(datastream[i]-'a')
		set ^= 1 << uint32(datastream[i-size]-'a')

		if bits.OnesCount32(set) == size {
			return i + 1, nil
		}
	}

	return 0, errors.New("not found")
}
