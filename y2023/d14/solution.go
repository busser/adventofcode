package d14

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 14 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	p, err := platformFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	p.tiltNorth()
	load := p.totalLoad()

	_, err = fmt.Fprintf(w, "%d", load)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 14 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	p, err := platformFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	p.spinN(1_000_000_000)
	load := p.totalLoad()

	_, err = fmt.Fprintf(w, "%d", load)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type platform [][]byte

const (
	roundRock  = 'O'
	squareRock = '#'
	empty      = '.'
)

func (p platform) String() string {
	s := ""
	for row := range p {
		s += string(p[row]) + "\n"
	}
	return s
}

func (p platform) spin() {
	p.tiltNorth()
	p.tiltWest()
	p.tiltSouth()
	p.tiltEast()
}

func (p platform) spinN(n int) {
	seen := make(map[string]int)

	for i := 0; i < n; i++ {
		s := p.String()

		if j, ok := seen[s]; ok {
			cycleLength := i - j
			remaining := (n - i) % cycleLength
			for k := 0; k < remaining; k++ {
				p.spin()
			}
			return
		}

		seen[s] = i
		p.spin()
	}
}

func (p platform) totalLoad() int {
	total := 0
	for row := range p {
		for col := range p[row] {
			if p[row][col] == roundRock {
				total += len(p) - row
			}
		}
	}
	return total
}

func (p platform) tiltNorth() {
	for col := range p[0] {
		for row := range p {
			if p[row][col] == roundRock {
				p.slideNorth(row, col)
			}
		}
	}
}

func (p platform) slideNorth(row, col int) {
	for i := row; i > 0; i-- {
		if p[i-1][col] == empty {
			p[i-1][col], p[i][col] = p[i][col], empty
		} else {
			return
		}
	}
}

func (p platform) tiltSouth() {
	for col := range p[0] {
		for row := len(p) - 1; row >= 0; row-- {
			if p[row][col] == roundRock {
				p.slideSouth(row, col)
			}
		}
	}
}

func (p platform) slideSouth(row, col int) {
	for i := row; i < len(p)-1; i++ {
		if p[i+1][col] == empty {
			p[i+1][col], p[i][col] = p[i][col], empty
		} else {
			return
		}
	}
}

func (p platform) tiltEast() {
	for row := range p {
		for col := len(p[row]) - 1; col >= 0; col-- {
			if p[row][col] == roundRock {
				p.slideEast(row, col)
			}
		}
	}
}

func (p platform) slideEast(row, col int) {
	for i := col; i < len(p[row])-1; i++ {
		if p[row][i+1] == empty {
			p[row][i+1], p[row][i] = p[row][i], empty
		} else {
			return
		}
	}
}

func (p platform) tiltWest() {
	for row := range p {
		for col := range p[row] {
			if p[row][col] == roundRock {
				p.slideWest(row, col)
			}
		}
	}
}

func (p platform) slideWest(row, col int) {
	for i := col; i > 0; i-- {
		if p[row][i-1] == empty {
			p[row][i-1], p[row][i] = p[row][i], empty
		} else {
			return
		}
	}
}

func platformFromReader(r io.Reader) (platform, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	p := make(platform, len(lines))
	for row := range lines {
		p[row] = []byte(lines[row])
		for col := range p[row] {
			switch p[row][col] {
			case roundRock, squareRock, empty:
				// do nothing
			default:
				return nil, fmt.Errorf("invalid character %c", p[row][col])
			}
		}
	}

	return p, nil
}
