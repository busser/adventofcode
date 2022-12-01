package busser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 25 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	cardPublicKey, doorPublicKey, err := publicKeysFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	cardLoopSize := loopSizeFromPublicKey(cardPublicKey)
	// doorLoopSize := loopSizeFromPublicKey(doorPublicKey)

	encryptionKey := encryptionKeyFromPublicKey(doorPublicKey, cardLoopSize)

	_, err = fmt.Fprintf(answer, "%d", encryptionKey)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func loopSizeFromPublicKey(key int) int {
	iteration, value := 0, 1

	for value != key {
		value = (value * 7) % 20201227
		iteration++
	}

	return iteration
}

func encryptionKeyFromPublicKey(pubKey, loopSize int) int {
	value := 1
	for i := 0; i < loopSize; i++ {
		value = (value * pubKey) % 20201227
	}
	return value
}

func publicKeysFromReader(r io.Reader) (int, int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return 0, 0, fmt.Errorf("reading lines: %w", err)
	}

	if len(lines) != 2 {
		return 0, 0, fmt.Errorf("expected 2 lines, got %d", len(lines))
	}

	cardKey, err := strconv.Atoi(lines[0])
	if err != nil {
		return 0, 0, fmt.Errorf("%q is not an integer", lines[0])
	}

	doorKey, err := strconv.Atoi(lines[1])
	if err != nil {
		return 0, 0, fmt.Errorf("%q is not an integer", lines[1])
	}

	return cardKey, doorKey, nil
}
