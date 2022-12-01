package busser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

// PartOne solves the first problem of day 20 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	tiles, err := tilesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	corners, _, _, err := classifyTiles(tiles)
	if err != nil {
		return fmt.Errorf("could not classify tiles: %w", err)
	}

	product := 1
	for _, t := range corners {
		product *= t.id
	}

	_, err = fmt.Fprintf(answer, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 20 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	tiles, err := tilesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	image, err := assembleTiles(tiles)
	if err != nil {
		return fmt.Errorf("could not assemble tiles: %w", err)
	}

	roughness := waterRoughnessFromImage(image)

	_, err = fmt.Fprintf(answer, "%d", roughness)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func classifyTiles(tiles []tile) (corners, edges, others []*tile, err error) {
	cardinalityBySide := make(map[[tileSize]pixel]int)
	for _, t := range tiles {
		for i := 0; i < 8; i++ {
			cardinalityBySide[t.topSide()]++
			t.reorient()
		}
	}

	for i, t := range tiles {
		var uniqueSides int

		for i := 0; i < 4; i++ {
			if cardinalityBySide[t.topSide()] == 1 {
				uniqueSides++
			}
			t.reorient()
		}

		switch uniqueSides {
		case 2:
			corners = append(corners, &tiles[i])
		case 1:
			edges = append(edges, &tiles[i])
		case 0:
			others = append(others, &tiles[i])
		default:
			return nil, nil, nil, fmt.Errorf("tile #%d has %d unique sides", i, uniqueSides)
		}
	}

	return corners, edges, others, nil
}

func assembleTiles(tiles []tile) (image [][]pixel, err error) {
	corners, edges, others, err := classifyTiles(tiles)
	if err != nil {
		return nil, fmt.Errorf("classifying tiles: %w", err)
	}

	if len(corners) != 4 {
		return nil, fmt.Errorf("expected 4 corners, found %d", len(corners))
	}

	puzzleSize := (len(edges) / 4) + 2
	if puzzleSize*puzzleSize != len(tiles) {
		return nil, errors.New("number of edge tiles indicates final image is not a square")
	}

	tileArrangement := make([][]*tile, puzzleSize)
	for i := range tileArrangement {
		tileArrangement[i] = make([]*tile, puzzleSize)
	}

	alreadyPlaced := make(map[*tile]bool)

	var placeTile func(int, int) bool
	placeTile = func(row, col int) bool {
		var candidates []*tile

		switch {
		case (row == 0 || row == puzzleSize-1) && (col == 0 || col == puzzleSize-1):
			candidates = corners
		case (row == 0 || row == puzzleSize-1) || (col == 0 || col == puzzleSize-1):
			candidates = edges
		default:
			candidates = others
		}

		for _, t := range candidates {
			if alreadyPlaced[t] {
				continue
			}

			for j := 0; j < 8; j++ {
				t.reorient()

				if row > 0 && t.topSide() != tileArrangement[row-1][col].bottomSide() {
					// tile does not match tile placed above
					continue
				}
				if col > 0 && t.leftSide() != tileArrangement[row][col-1].rightSide() {
					// tile does not match tile placed to the left
					continue
				}

				tileArrangement[row][col] = t
				alreadyPlaced[t] = true

				var nextRow, nextCol int
				switch {
				case row == puzzleSize-1 && col == puzzleSize-1:
					return true
				case col == puzzleSize-1:
					nextRow, nextCol = row+1, 0
				default:
					nextRow, nextCol = row, col+1
				}

				if placeTile(nextRow, nextCol) {
					return true
				}

				tileArrangement[row][col] = nil
				alreadyPlaced[t] = false
			}
		}

		return false
	}

	if ok := placeTile(0, 0); !ok {
		return nil, errors.New("failed to arrange tiles correctly")
	}

	image = make([][]pixel, puzzleSize*(tileSize-2))
	for i := range image {
		image[i] = make([]pixel, puzzleSize*(tileSize-2))
	}

	for row := range tileArrangement {
		for col, t := range tileArrangement[row] {
			for x := 0; x < tileSize-2; x++ {
				for y := 0; y < tileSize-2; y++ {
					image[row*(tileSize-2)+x][col*(tileSize-2)+y] = t.content[x+1][y+1]
				}
			}
		}
	}

	return image, nil
}

func tilesFromReader(r io.Reader) ([]tile, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	rawTiles := bytes.Split(b, []byte("\n\n"))
	rawTiles = rawTiles[:len(rawTiles)-1] // last value is empty

	tiles := make([]tile, len(rawTiles))
	for i := range rawTiles {
		if err := tiles[i].fromBytes(rawTiles[i]); err != nil {
			return nil, fmt.Errorf("reading tile #%d: %w", i, err)
		}
	}

	return tiles, nil
}
