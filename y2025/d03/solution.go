package d03

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 3 of Advent of Code 2025.
func PartOne(r io.Reader, w io.Writer) error {
	banks, err := banksFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for _, bank := range banks {
		sum += maxBankJoltage(bank, batteryCountPart1)
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 3 of Advent of Code 2025.
func PartTwo(r io.Reader, w io.Writer) error {
	banks, err := banksFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0
	for _, bank := range banks {
		sum += maxBankJoltage(bank, batteryCountPart2)
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const (
	batteryCountPart1 = 2
	batteryCountPart2 = 12
)

func maxBankJoltage(bank []int, batteryCount int) int {
	if batteryCount == 0 {
		return 0
	}

	batteriesAfter := batteryCount - 1

	firstBattery := bestBatteryPosition(bank[:len(bank)-batteriesAfter])
	joltage := bank[firstBattery]

	for range batteriesAfter {
		joltage *= 10
	}
	joltage += maxBankJoltage(bank[firstBattery+1:], batteriesAfter)

	return joltage
}

func bestBatteryPosition(bank []int) int {
	position := 0
	for i := 1; i < len(bank); i++ {
		if bank[i] > bank[position] {
			position = i
		}
	}
	return position
}

func banksFromReader(r io.Reader) ([][]int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	banks := make([][]int, len(lines))
	for i, line := range lines {
		if len(line) < batteryCountPart2 {
			return nil, fmt.Errorf("line %d is too short", i)
		}

		banks[i] = make([]int, len(line))
		for j, battery := range line {
			if battery < '0' || battery > '9' {
				return nil, fmt.Errorf("battery out of range: %q", battery)
			}
			banks[i][j] = int(battery - '0')
		}
	}

	return banks, nil
}
