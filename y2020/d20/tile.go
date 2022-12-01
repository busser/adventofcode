package busser

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

const tileSize = 10

type pixel uint8

const (
	pixelEmpty pixel = iota
	pixelFull
	pixelAny
)

type tile struct {
	id                 int
	content            [tileSize][tileSize]pixel
	reorientationCount int
}

func (t *tile) rotate() {
	for row := 0; row < tileSize/2; row++ {
		for col := row; col < tileSize-row-1; col++ {
			tmp := t.content[row][col]

			// top <- left
			t.content[row][col] = t.content[tileSize-1-col][row]

			// left <- bottom
			t.content[tileSize-1-col][row] = t.content[tileSize-1-row][tileSize-1-col]

			// bottom <- right
			t.content[tileSize-1-row][tileSize-1-col] = t.content[col][tileSize-1-row]

			// right <- top
			t.content[col][tileSize-1-row] = tmp
		}
	}
}

func (t *tile) flip() {
	for row := 0; row < tileSize; row++ {
		left, right := 0, tileSize-1
		for left < right {
			t.content[row][left], t.content[row][right] = t.content[row][right], t.content[row][left]
			left++
			right--
		}
	}
}

func (t *tile) reorient() {
	t.reorientationCount = (t.reorientationCount + 1) % 4

	t.rotate()
	if t.reorientationCount == 0 {
		t.flip()
	}
}

// topSide returns a copy of t's top side from left to right.
func (t *tile) topSide() [tileSize]pixel {
	var side [tileSize]pixel
	for i := range side {
		side[i] = t.content[0][i]
	}
	return side
}

// bottomSide returns a copy of t's bottom side from left to right.
func (t *tile) bottomSide() [tileSize]pixel {
	var side [tileSize]pixel
	for i := range side {
		side[i] = t.content[tileSize-1][i]
	}
	return side
}

// leftSide returns a copy of t's left side from top to bottom.
func (t *tile) leftSide() [tileSize]pixel {
	var side [tileSize]pixel
	for i := range side {
		side[i] = t.content[i][0]
	}
	return side
}

// rightSide returns a copy of t's right side from top to bottom.
func (t *tile) rightSide() [tileSize]pixel {
	var side [tileSize]pixel
	for i := range side {
		side[i] = t.content[i][tileSize-1]
	}
	return side
}

func (t *tile) fromBytes(b []byte) error {
	lines := bytes.Split(b, []byte("\n"))
	if len(lines) != 1+tileSize {
		return errors.New("wrong format")
	}

	line := lines[0]

	prefix, suffix := []byte("Tile "), []byte(":")
	if !bytes.HasPrefix(line, prefix) && !bytes.HasSuffix(line, suffix) {
		return errors.New("wrong format")
	}

	rawID := string(line[len(prefix) : len(line)-len(suffix)])
	id, err := strconv.Atoi(rawID)
	if err != nil {
		return fmt.Errorf("wrong format: %w", err)
	}
	t.id = id

	for row := 0; row < tileSize; row++ {
		line := lines[row+1]
		if len(line) != tileSize {
			return errors.New("wrong format")
		}
		for col, cell := range line {
			switch cell {
			case '.':
				t.content[row][col] = pixelEmpty
			case '#':
				t.content[row][col] = pixelFull
			default:
				return fmt.Errorf("unknown cell value %q", cell)
			}
		}
	}

	return nil
}
