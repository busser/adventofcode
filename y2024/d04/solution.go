package d04

import (
	"fmt"
	"io"
	"iter"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 4 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	wordSearch, err := wordSearchFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := searchForXMAS(wordSearch)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 4 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	wordSearch, err := wordSearchFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := searchForCrossMAS(wordSearch)

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type vector struct {
	row, col int
}

func (v vector) plus(w vector) vector {
	return vector{
		row: v.row + w.row,
		col: v.col + w.col,
	}
}

var (
	up    = vector{-1, 0}
	down  = vector{1, 0}
	left  = vector{0, -1}
	right = vector{0, 1}
)

func searchForCrossMAS(wordsearch [][]byte) int {
	count := 0

	for row := range wordsearch {
		for col := range wordsearch[row] {
			if wordsearch[row][col] == 'A' {
				if isCrossMAS(wordsearch, row, col) {
					count++
				}
			}
		}
	}

	return count
}

func isCrossMAS(wordSearch [][]byte, row, col int) bool {
	if row < 1 || row >= len(wordSearch)-1 || col < 1 || col >= len(wordSearch[row])-1 {
		return false
	}

	topLeft := wordSearch[row-1][col-1]
	topRight := wordSearch[row-1][col+1]
	bottomLeft := wordSearch[row+1][col-1]
	bottomRight := wordSearch[row+1][col+1]

	firstDiag := (topLeft == 'M' && bottomRight == 'S') || (topLeft == 'S' && bottomRight == 'M')
	secondDiag := (topRight == 'M' && bottomLeft == 'S') || (topRight == 'S' && bottomLeft == 'M')

	return firstDiag && secondDiag
}

func searchForXMAS(wordSearch [][]byte) int {
	return searchRowsForXMAS(wordSearch) + searchColumnsForXMAS(wordSearch) + searchDiagonalsForXMAS(wordSearch)
}

func searchRowsForXMAS(wordSearch [][]byte) int {
	count := 0

	for row := range wordSearch {
		count += searchSequenceForXMAS(iterate(wordSearch, vector{row, 0}, right))
		count += searchSequenceForXMAS(iterate(wordSearch, vector{row, len(wordSearch[row]) - 1}, left))
	}

	return count
}

func searchColumnsForXMAS(wordSearch [][]byte) int {
	count := 0

	for col := range wordSearch[0] {
		count += searchSequenceForXMAS(iterate(wordSearch, vector{0, col}, down))
		count += searchSequenceForXMAS(iterate(wordSearch, vector{len(wordSearch[0]) - 1, col}, up))
	}

	return count
}

func searchDiagonalsForXMAS(wordSearch [][]byte) int {
	count := 0

	directions := []vector{
		up.plus(left),
		up.plus(right),
		down.plus(left),
		down.plus(right),
	}

	for _, dir := range directions {
		for row := range wordSearch {
			count += searchSequenceForXMAS(iterate(wordSearch, vector{row, 0}, dir))
			count += searchSequenceForXMAS(iterate(wordSearch, vector{row, len(wordSearch[row]) - 1}, dir))
		}

		for col := range wordSearch[0] {
			// Skip starting points already covered by previous loop.
			if col == 0 || col == len(wordSearch[0])-1 {
				continue
			}

			count += searchSequenceForXMAS(iterate(wordSearch, vector{0, col}, dir))
			count += searchSequenceForXMAS(iterate(wordSearch, vector{len(wordSearch[0]) - 1, col}, dir))
		}
	}

	return count
}

func searchSequenceForXMAS(seq iter.Seq[byte]) int {
	const word = "XMAS"

	want := 0
	count := 0

	for letter := range seq {
		switch {
		case letter == word[0]:
			want = 1
		case letter == word[want]:
			want++
			if want == len(word) {
				count++
				want = 0
			}
		default:
			want = 0
		}
	}

	return count
}

func iterate(wordSearch [][]byte, start, direction vector) iter.Seq[byte] {
	return func(yield func(byte) bool) {
		pos := start
		for {
			if pos.row < 0 || pos.row >= len(wordSearch) || pos.col < 0 || pos.col >= len(wordSearch[pos.row]) {
				return
			}
			if !yield(wordSearch[pos.row][pos.col]) {
				return
			}
			pos = pos.plus(direction)
		}
	}
}

func wordSearchFromReader(r io.Reader) ([][]byte, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("empty word search")
	}

	rowLen := len(lines[0])
	for _, line := range lines {
		if len(line) != rowLen {
			return nil, fmt.Errorf("word search is not a rectangle")
		}
	}

	wordSearch := make([][]byte, len(lines))
	for row, line := range lines {
		wordSearch[row] = []byte(line)
	}

	return wordSearch, nil
}
