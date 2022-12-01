package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 21 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	player1, player2, err := startingPositionsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	losingScore, diceRolls := playDeterministicGame(player1, player2)

	_, err = fmt.Fprintf(answer, "%d", losingScore*diceRolls)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 21 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	player1, player2, err := startingPositionsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	p1Wins, p2Wins := playQuantumGame(player1, player2)

	max := p1Wins
	if p2Wins > max {
		max = p2Wins
	}

	_, err = fmt.Fprintf(answer, "%d", max)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func playDeterministicGame(p1Pos, p2Pos int) (losingScore, diceRolls int) {
	p1Score, p2Score, diceRolls := 0, 0, 0

	rollDie := func() int {
		total := 0
		for i := 0; i < 3; i++ {
			total += diceRolls%100 + 1
			diceRolls++
		}
		return total
	}

	for {
		// Player 1 plays.
		p1Pos = (p1Pos+rollDie()-1)%10 + 1
		p1Score += p1Pos
		if p1Score >= 1000 {
			return p2Score, diceRolls
		}

		// Player 2 plays.
		p2Pos = (p2Pos+rollDie()-1)%10 + 1
		p2Score += p2Pos
		if p2Score >= 1000 {
			return p1Score, diceRolls
		}
	}
}

func playQuantumGame(p1Pos, p2Pos int) (p1Wins, p2Wins int) {
	type gameState struct {
		p1Pos, p2Pos     int
		p1Score, p2Score int
		p1ToPlay         bool
		diceRolls        int
		diceTotal        int
	}

	memory := make(map[gameState][2]int)

	var helper func(gameState) (int, int)
	helper = func(state gameState) (p1Wins, p2Wins int) {
		if wins, ok := memory[state]; ok {
			return wins[0], wins[1]
		}

		// Symetrical games have symetrical outcomes.
		mirrorState := state
		mirrorState.p1Pos = state.p2Pos
		mirrorState.p2Pos = state.p1Pos
		mirrorState.p1Score = state.p2Score
		mirrorState.p2Score = state.p1Score
		mirrorState.p1ToPlay = !state.p1ToPlay
		if wins, ok := memory[mirrorState]; ok {
			memory[state] = [2]int{wins[1], wins[0]}
			return wins[1], wins[0]
		}

		// If the game is over, tally the winner.
		if state.p1Score >= 21 {
			memory[state] = [2]int{1, 0}
			return 1, 0
		}
		if state.p2Score >= 21 {
			memory[state] = [2]int{0, 1}
			return 0, 1
		}

		originalState := state

		// If playerX has rolled the dice three times, move him forward and
		// switch to playerY.
		if state.diceRolls == 3 {
			if state.p1ToPlay {
				state.p1Pos = (state.p1Pos+state.diceTotal-1)%10 + 1
				state.p1Score += state.p1Pos
			} else {
				state.p2Pos = (state.p2Pos+state.diceTotal-1)%10 + 1
				state.p2Score += state.p2Pos
			}
			state.diceRolls = 0
			state.diceTotal = 0
			state.p1ToPlay = !state.p1ToPlay

			p1Wins, p2Wins = helper(state)

			memory[originalState] = [2]int{p1Wins, p2Wins}
			return helper(state)
		}

		// If a player's turn is in progress, roll the die.
		state.diceRolls++
		for i := 0; i < 3; i++ {
			state.diceTotal++
			p1, p2 := helper(state)
			p1Wins, p2Wins = p1Wins+p1, p2Wins+p2
		}

		memory[originalState] = [2]int{p1Wins, p2Wins}
		return p1Wins, p2Wins
	}

	state := gameState{
		p1Pos: p1Pos, p2Pos: p2Pos,
		p1ToPlay: true,
	}
	return helper(state)
}

func startingPositionsFromReader(r io.Reader) (player1, player2 int, err error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return 0, 0, err
	}

	if len(lines) != 2 {
		return 0, 0, fmt.Errorf("expected 2 lines, got %d", len(lines))
	}

	prefix := "Player X starting position: "
	if len(lines[0]) < len(prefix) || len(lines[1]) < len(prefix) {
		return 0, 0, errors.New("wrong format")
	}

	rawP1, rawP2 := lines[0][len(prefix):], lines[1][len(prefix):]

	player1, err = strconv.Atoi(rawP1)
	if err != nil {
		return 0, 0, fmt.Errorf("%q is not a number", rawP1)
	}
	player2, err = strconv.Atoi(rawP2)
	if err != nil {
		return 0, 0, fmt.Errorf("%q is not a number", rawP2)
	}

	return player1, player2, nil
}
