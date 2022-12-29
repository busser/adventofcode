package d22

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 22 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	board, path, err := boardAndPathFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	position, facing := finalPositionAndFacing(board, path, pacmanWrapping)

	pwd := password(position, facing)

	_, err = fmt.Fprintf(w, "%d", pwd)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 22 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	board, path, err := boardAndPathFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	position, facing := finalPositionAndFacing(board, path, cubeWrapping)

	pwd := password(position, facing)

	_, err = fmt.Fprintf(w, "%d", pwd)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type board [][]tile

type tile uint8

const (
	tileNone tile = iota
	tileOpen
	tileWall
)

type path []pathInstruction

type pathInstruction struct {
	kind   pathInstructionKind
	amount int
}

type pathInstructionKind uint8

const (
	pathMoveForward pathInstructionKind = iota
	pathTurnLeft
	pathTurnRight
)

type vector struct {
	x, y int
}

var (
	facingRight = vector{1, 0}
	facingDown  = vector{0, 1}
	facingLeft  = vector{-1, 0}
	facingUp    = vector{0, -1}
)

func (v vector) plus(w vector) vector {
	return vector{
		x: v.x + w.x,
		y: v.y + w.y,
	}
}

func (v vector) rotatedLeft() vector {
	return vector{
		x: v.y,
		y: -v.x,
	}
}

func (v vector) rotatedRight() vector {
	return vector{
		x: -v.y,
		y: v.x,
	}
}

type wrappingRule uint8

const (
	pacmanWrapping wrappingRule = iota
	cubeWrapping
)

func finalPositionAndFacing(b board, p path, wrap wrappingRule) (vector, vector) {
	// Start at the left-most open tile in the top row.
	var position vector
	for x := range b {
		if b[x][0] == tileOpen {
			position = vector{x, 0}
			break
		}
	}

	// Start facing right
	facing := facingRight

	// Follow the path
	for _, inst := range p {
		switch inst.kind {
		case pathMoveForward:
			for n := 0; n < inst.amount; n++ {
				nextPosition, nextFacing := b.nextPositionAndFacing(position, facing, wrap)
				if b[nextPosition.x][nextPosition.y] == tileOpen {
					position, facing = nextPosition, nextFacing
				}
			}
		case pathTurnLeft:
			facing = facing.rotatedLeft()
		case pathTurnRight:
			facing = facing.rotatedRight()
		}
	}

	return position, facing
}

func (b board) nextPositionAndFacing(position, facing vector, wrap wrappingRule) (vector, vector) {
	switch wrap {
	case pacmanWrapping:
		return b.nextPositionAndFacingWithPacmanWrapping(position, facing)
	case cubeWrapping:
		return b.nextPositionAndFacingWithCubeWrapping(position, facing)
	default:
		panic("unknown wrapping rule")
	}
}

func (b board) nextPositionAndFacingWithPacmanWrapping(position, facing vector) (vector, vector) {
	// Move one step in the direction currently facing
	position = position.plus(facing)

	// Loop around if moved out of board bounds
	position.x = (position.x + len(b)) % len(b)
	position.y = (position.y + len(b[0])) % len(b[0])

	// Repeat until we land on board again
	for b[position.x][position.y] == tileNone {
		position = position.plus(facing)
		position.x = (position.x + len(b)) % len(b)
		position.y = (position.y + len(b[0])) % len(b[0])
	}

	return position, facing
}

func (b board) nextPositionAndFacingWithCubeWrapping(position, facing vector) (vector, vector) {
	// The input always has the same layout according to Reddit, so wrapping
	// rules can be hardcoded. The example input has a different layout, which
	// we don't take into account.
	//
	//                 U1        U2
	//             ┌─────────┬─────────┐
	//             │         │         │
	//           L0│         │         │R0
	//             │         │         │
	//             ├─────────┼─────────┘
	//             │         │   D2
	//           L1│         │R1
	//       U0    │         │
	//   ┌─────────┼─────────┤
	//   │         │         │
	// L2│         │         │R2
	//   │         │         │
	//   ├─────────┼─────────┘
	//   │         │   D1
	// L3│         │R3
	//   │         │
	//   └─────────┘
	//       D0
	//
	// Pairs of sides join when the cube is formed:
	//
	// U0 <-> L1
	// U1 <-> L3
	// U2 <-> D0
	// L0 <-> L2
	// R0 <-> R2
	// R1 <-> D2
	// R3 <-> D1

	cubeSize := gcd(len(b), len(b[0]))

	// For brevity
	x, y, f, s := position.x, position.y, facing, cubeSize

	switch {
	// U0 -> L1
	case f == facingUp && y == 2*s && x >= 0 && x < s:
		return vector{x: s, y: x + s}, facingRight
	// U1 -> L3
	case f == facingUp && y == 0 && x >= s && x < 2*s:
		return vector{x: 0, y: x + 2*s}, facingRight
	// U2 -> D0
	case f == facingUp && y == 0 && x >= 2*s && x < 3*s:
		return vector{x: x - 2*s, y: 4*s - 1}, facingUp
	// L0 -> L2
	case f == facingLeft && x == s && y >= 0 && y < s:
		return vector{x: 0, y: 3*s - 1 - y}, facingRight
	// L1 -> U0
	case f == facingLeft && x == s && y >= s && y < 2*s:
		return vector{x: y - s, y: 2 * s}, facingDown
	// L2 -> L0
	case f == facingLeft && x == 0 && y >= 2*s && y < 3*s:
		return vector{x: s, y: 3*s - 1 - y}, facingRight
	// L3 -> U1
	case f == facingLeft && x == 0 && y >= 3*s && y < 4*s:
		return vector{x: y - 2*s, y: 0}, facingDown
	// R0 -> R2
	case f == facingRight && x == 3*s-1 && y >= 0 && y < s:
		return vector{x: 2*s - 1, y: 3*s - 1 - y}, facingLeft
	// R1 -> D2
	case f == facingRight && x == 2*s-1 && y >= s && y < 2*s:
		return vector{x: y + s, y: s - 1}, facingUp
	// R2 -> R0
	case f == facingRight && x == 2*s-1 && y >= 2*s && y < 3*s:
		return vector{x: 3*s - 1, y: 3*s - 1 - y}, facingLeft
	// R3 -> D1
	case f == facingRight && x == s-1 && y >= 3*s && y < 4*s:
		return vector{x: y - 2*s, y: 3*s - 1}, facingUp
	// D0 -> U2
	case f == facingDown && y == 4*s-1 && x >= 0 && x < s:
		return vector{x: x + 2*s, y: 0}, facingDown
	// D1 -> R3
	case f == facingDown && y == 3*s-1 && x >= s && x < 2*s:
		return vector{x: s - 1, y: x + 2*s}, facingLeft
	// D2 -> R1
	case f == facingDown && y == s-1 && x >= 2*s && x < 3*s:
		return vector{x: 2*s - 1, y: x - s}, facingLeft
	// Not on an edge, so no wrapping
	default:
		return position.plus(facing), facing
	}
}

func password(position, facing vector) int {
	password := 1000*(position.y+1) + 4*(position.x+1)

	switch facing {
	case facingRight:
		password += 0
	case facingDown:
		password += 1
	case facingLeft:
		password += 2
	case facingUp:
		password += 3
	}

	return password
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func boardAndPathFromReader(r io.Reader) (board, path, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, nil, err
	}

	if len(lines) < 3 {
		return nil, nil, errors.New("not enough lines")
	}

	board := boardFromLines(lines[:len(lines)-2])

	path, err := pathFromString(lines[len(lines)-1])
	if err != nil {
		return nil, nil, err
	}

	return board, path, nil
}

func boardFromLines(lines []string) board {
	maxLen := len(lines[0])
	for _, l := range lines {
		maxLen = max(maxLen, len(l))
	}

	board := make([][]tile, maxLen)
	for x := range board {
		board[x] = make([]tile, len(lines))
	}

	for y, l := range lines {
		for x, c := range l {
			switch c {
			case '.':
				board[x][y] = tileOpen
			case '#':
				board[x][y] = tileWall
			}
		}
	}

	return board
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func pathFromString(s string) (path, error) {
	// quick and dirty, could definitely be cleaner
	rawMoves := strings.FieldsFunc(s, func(r rune) bool {
		return r < '0' || r > '9'
	})
	rawTurns := strings.FieldsFunc(s, func(r rune) bool {
		return r >= '0' && r <= '9'
	})

	moves := make([]pathInstruction, len(rawMoves))
	for i, raw := range rawMoves {
		n, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("%q is not a number", raw)
		}

		moves[i] = pathInstruction{
			kind:   pathMoveForward,
			amount: n,
		}
	}

	turns := make([]pathInstruction, len(rawTurns))
	for i, raw := range rawTurns {
		switch raw {
		case "L":
			turns[i] = pathInstruction{kind: pathTurnLeft}
		case "R":
			turns[i] = pathInstruction{kind: pathTurnRight}
		default:
			return nil, fmt.Errorf("unknown turn %q", raw)
		}
	}

	diff := len(moves) - len(turns)
	if diff < 0 || diff > 1 {
		// Unclear how this could actually happen
		// but I'd rather error than panic
		return nil, fmt.Errorf("wrong format")
	}

	path := make(path, len(moves)+len(turns))
	for i := range moves {
		path[2*i] = moves[i]
	}
	for i := range turns {
		path[2*i+1] = turns[i]
	}

	return path, nil
}
