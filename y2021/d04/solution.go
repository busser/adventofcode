package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 4 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	game, err := bingoGameFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	winner, winningNumber, err := game.firstWinner()
	if err != nil {
		return fmt.Errorf("could not determine winner: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", winner.score()*winningNumber)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 4 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	game, err := bingoGameFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	winner, winningNumber, err := game.lastWinner()
	if err != nil {
		return fmt.Errorf("could not determine winner: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", winner.score()*winningNumber)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type bingoGame struct {
	boards  []*bingoBoard
	numbers []int
}

const bingoBoardSize = 5

type bingoBoard struct {
	cells [bingoBoardSize][bingoBoardSize]bingoCell
}

type bingoCell struct {
	number int
	marked bool
}

func (bg *bingoGame) firstWinner() (winner *bingoBoard, winningNumber int, err error) {
	for _, number := range bg.numbers {
		for _, board := range bg.boards {
			board.mark(number)
			if board.isComplete() {
				return board, number, nil
			}
		}
	}
	return nil, 0, errors.New("no board has won")
}

func (bg *bingoGame) lastWinner() (winner *bingoBoard, winningNumber int, err error) {
	winners := make(map[*bingoBoard]bool)
	for _, number := range bg.numbers {
		for _, board := range bg.boards {
			board.mark(number)
			if !winners[board] && board.isComplete() {
				winners[board] = true
				if len(winners) == len(bg.boards) {
					return board, number, nil
				}
			}
		}
	}
	return nil, 0, errors.New("not all boards have won")
}

func (bb *bingoBoard) mark(number int) {
	for row := range bb.cells {
		for column := range bb.cells[row] {
			if bb.cells[row][column].number == number {
				bb.cells[row][column].marked = true
			}
		}
	}
}

func (bb *bingoBoard) isComplete() bool {
	// Check if any row is complete.
	for row := 0; row < bingoBoardSize; row++ {
		rowIsComplete := true
		for column := 0; column < bingoBoardSize; column++ {
			if !bb.cells[row][column].marked {
				rowIsComplete = false
			}
		}
		if rowIsComplete {
			return true
		}
	}

	// Check if any column is complete.
	for column := 0; column < bingoBoardSize; column++ {
		columnIsComplete := true
		for row := 0; row < bingoBoardSize; row++ {
			if !bb.cells[row][column].marked {
				columnIsComplete = false
			}
		}
		if columnIsComplete {
			return true
		}
	}

	return false
}

func (bb *bingoBoard) score() int {
	s := 0
	for _, row := range bb.cells {
		for _, cell := range row {
			if !cell.marked {
				s += cell.number
			}
		}
	}
	return s
}

func bingoGameFromReader(r io.Reader) (bingoGame, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return bingoGame{}, fmt.Errorf("could not read: %w", err)
	}

	if len(lines)%(1+bingoBoardSize) != 1 {
		return bingoGame{}, fmt.Errorf("the input should have 1 line for numbers, %d lines for each board, and a blank line before each board", bingoBoardSize)
	}

	numbers := helpers.IntsFromString(lines[0])

	var boards []*bingoBoard

	for i := 1; i < len(lines); i += 1 + bingoBoardSize {
		if len(lines[i]) > 0 {
			return bingoGame{}, fmt.Errorf("each board should have an empty line before it")
		}

		var board bingoBoard

		for row := 0; row < bingoBoardSize; row++ {
			rawNumbers := strings.Fields(lines[i+1+row])
			if len(rawNumbers) != bingoBoardSize {
				return bingoGame{}, fmt.Errorf("invalid row %q", lines[i+1+row])
			}

			for column := 0; column < bingoBoardSize; column++ {
				n, err := strconv.Atoi(rawNumbers[column])
				if err != nil {
					return bingoGame{}, fmt.Errorf("%q is not a number", rawNumbers[column])
				}
				board.cells[row][column].number = n
			}
		}

		boards = append(boards, &board)
	}

	return bingoGame{boards, numbers}, nil
}
