package d04

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 4 of Advent of Code 2015.
func PartOne(r io.Reader, w io.Writer) error {
	secretKey, err := secretKeyFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	n := mineAdventCoin(secretKey)

	_, err = fmt.Fprintf(w, "%d", n)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 4 of Advent of Code 2015.
func PartTwo(r io.Reader, w io.Writer) error {
	secretKey, err := secretKeyFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	n := mineMoreAdventCoin(secretKey)

	_, err = fmt.Fprintf(w, "%d", n)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func mineAdventCoin(secretKey string) int {
	for n := 0; ; n++ {
		hash := hashNumber(secretKey, n)
		if hash[:5] == "00000" {
			return n
		}
	}
}

func mineMoreAdventCoin(secretKey string) int {
	for n := 0; ; n++ {
		hash := hashNumber(secretKey, n)
		if hash[:6] == "000000" {
			return n
		}
	}
}

func hashNumber(secretKey string, n int) string {
	data := []byte(secretKey + strconv.Itoa(n))
	hashed := md5.Sum(data)
	encoded := hex.EncodeToString(hashed[:])
	return encoded
}

func secretKeyFromReader(r io.Reader) (string, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return "", fmt.Errorf("could not read input: %w", err)
	}

	if len(lines) != 1 {
		return "", fmt.Errorf("expected 1 line, got %d", len(lines))
	}

	return lines[0], nil
}
