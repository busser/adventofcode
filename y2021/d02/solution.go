package busser

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

func PartOne(input io.Reader, answer io.Writer) error {
	commands, err := commandsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read commands: %w", err)
	}

	position, depth := interpretCommandsWithoutManual(commands)

	_, err = fmt.Fprintf(answer, "%d", depth*position)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func PartTwo(input io.Reader, answer io.Writer) error {
	commands, err := commandsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read commands: %w", err)
	}

	position, depth := interpretCommandsWithManual(commands)

	_, err = fmt.Fprintf(answer, "%d", depth*position)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type command struct {
	direction string
	amount    int
}

func interpretCommandsWithoutManual(commands []command) (position, depth int) {
	position, depth = 0, 0

	for _, c := range commands {
		switch c.direction { // unknown commands are ignored
		case "forward":
			position += c.amount
		case "down":
			depth += c.amount
		case "up":
			depth -= c.amount
		}
	}

	return position, depth
}

func interpretCommandsWithManual(commands []command) (position, depth int) {
	position, depth, aim := 0, 0, 0

	for _, c := range commands {
		switch c.direction { // unknown commands are ignored
		case "forward":
			position += c.amount
			depth += aim * c.amount
		case "down":
			aim += c.amount
		case "up":
			aim -= c.amount
		}
	}

	return position, depth
}

func commandsFromReader(r io.Reader) ([]command, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading lines: %w", err)
	}

	commands := make([]command, len(lines))

	for i := range lines {
		// Each line should like "forward 123".
		parts := strings.SplitN(lines[i], " ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%q is not a valid command", lines[i])
		}

		direction := parts[0]
		amount, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("%q is not a number", parts[1])
		}

		commands[i] = command{direction, amount}
	}

	return commands, nil
}
