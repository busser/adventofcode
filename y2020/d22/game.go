package busser

const MaxDeckSize = uint8(50)

type game struct {
	player1, player2 deck
}

type deck interface {
	hash() [MaxDeckSize]uint8
	drawFromTop() uint8
	insertAtBottom(uint8)
	size() uint8
	score() uint16
	copyTopCards(uint8) deck
}

func (g game) over() bool {
	return g.player1.size() == 0 || g.player2.size() == 0
}

func (g *game) playUntilOver() deck {
	for !g.over() {
		g.playTurn()
	}

	switch {
	case g.player1.size() == 0:
		return g.player2
	case g.player2.size() == 0:
		return g.player1
	default:
		panic("game is over but neither deck is empty")
	}
}

func (g *game) playTurn() {
	card1 := g.player1.drawFromTop()
	card2 := g.player2.drawFromTop()

	if card1 > card2 {
		g.player1.insertAtBottom(card1)
		g.player1.insertAtBottom(card2)
	} else {
		g.player2.insertAtBottom(card2)
		g.player2.insertAtBottom(card1)
	}
}

type gameState struct {
	deck1, deck2 [MaxDeckSize]uint8
}

func (g *game) playRecursivelyUntilOver() deck {
	pastStates := make(map[gameState]struct{})
	for !g.over() {
		currentState := gameState{g.player1.hash(), g.player2.hash()}
		if _, seen := pastStates[currentState]; seen {
			return g.player1
		}
		pastStates[currentState] = struct{}{}

		g.playRecursiveTurn()
	}

	switch {
	case g.player1.size() == 0:
		return g.player2
	case g.player2.size() == 0:
		return g.player1
	default:
		panic("game is over but neither deck is empty")
	}
}

func (g *game) playRecursiveTurn() {
	card1 := g.player1.drawFromTop()
	card2 := g.player2.drawFromTop()

	if g.player1.size() >= card1 && g.player2.size() >= card2 {
		subGame := game{
			player1: g.player1.copyTopCards(card1),
			player2: g.player2.copyTopCards(card2),
		}
		winner := subGame.playRecursivelyUntilOver()

		switch winner {
		case subGame.player1:
			g.player1.insertAtBottom(card1)
			g.player1.insertAtBottom(card2)
		case subGame.player2:
			g.player2.insertAtBottom(card2)
			g.player2.insertAtBottom(card1)
		default:
			panic("no player won the sub-game")
		}

		return
	}

	if card1 > card2 {
		g.player1.insertAtBottom(card1)
		g.player1.insertAtBottom(card2)
	} else {
		g.player2.insertAtBottom(card2)
		g.player2.insertAtBottom(card1)
	}
}
