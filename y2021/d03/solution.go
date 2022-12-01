package busser

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 3 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	// TODO: Read the input. For example:
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var bitCount [12]int
	for _, l := range lines {
		for i, b := range l {
			if b == '1' {
				bitCount[i]++
			}
		}
	}

	var gamma, epsilon int
	value := 1
	for i := len(bitCount) - 1; i >= 0; i-- {
		if bitCount[i] > len(lines)/2 {
			gamma += value
		} else {
			epsilon += value
		}
		value *= 2
	}

	// TODO: Write the answer. For example:
	_, err = fmt.Fprintf(answer, "%d", gamma*epsilon)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 3 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	// TODO: Read the input. For example:
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	filterOxygen := func(lines []string, index int) []string {
		var bitCount int
		for _, l := range lines {
			if l[index] == '1' {
				bitCount++
			}
		}

		var newLines []string
		if bitCount >= len(lines)/2 { // most common is 1
			for _, l := range lines {
				if l[index] == '1' {
					newLines = append(newLines, l)
				}
			}
		} else { // most common is 0
			for _, l := range lines {
				if l[index] == '0' {
					newLines = append(newLines, l)
				}
			}
		}

		return newLines
	}

	filterCO2 := func(lines []string, index int) []string {
		var bitCount int
		for _, l := range lines {
			if l[index] == '1' {
				bitCount++
			}
		}

		var newLines []string
		if bitCount >= len(lines)/2 { // least common is 0
			for _, l := range lines {
				if l[index] == '0' {
					newLines = append(newLines, l)
				}
			}
		} else { // least common is 1
			for _, l := range lines {
				if l[index] == '1' {
					newLines = append(newLines, l)
				}
			}
		}

		return newLines
	}

	oxygenLines := lines
	for index := 0; index < 12; index++ {
		oxygenLines = filterOxygen(oxygenLines, index)
		if len(oxygenLines) == 1 {
			break
		}
	}

	co2Lines := lines
	for index := 0; index < 12; index++ {
		co2Lines = filterCO2(co2Lines, index)
		if len(co2Lines) == 1 {
			break
		}
	}

	helpers.Println(oxygenLines, co2Lines)

	var oxygen, co2 int
	value := 1
	for i := 11; i >= 0; i-- {
		if oxygenLines[0][i] == '1' {
			oxygen += value
		}
		if co2Lines[0][i] == '1' {
			co2 += value
		}
		value *= 2
	}

	// TODO: Write the answer. For example:
	_, err = fmt.Fprintf(answer, "%d", oxygen*co2)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}
