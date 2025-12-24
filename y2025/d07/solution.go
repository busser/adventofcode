package d07

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// PartOne solves the first problem of day 7 of Advent of Code 2025.
func PartOne(r io.Reader, w io.Writer) error {
	manifold, err := tachyonManifoldFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	manifold.splitBeams()

	_, err = fmt.Fprintf(w, "%d", manifold.splitCount)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 7 of Advent of Code 2025.
func PartTwo(r io.Reader, w io.Writer) error {
	manifold, err := tachyonManifoldFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	manifold.splitBeams()

	timelineCount := manifold.getFinalTimelineCount()

	_, err = fmt.Fprintf(w, "%d", timelineCount)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const (
	start    = 'S'
	empty    = '.'
	splitter = '^'
	beam     = '|'
)

type tachyonManifold struct {
	grid              [][]byte
	beamTimelineCount [][]int
	splitCount        int
}

func (tm *tachyonManifold) String() string {
	var sb strings.Builder
	for row := range tm.grid {
		sb.Write(tm.grid[row])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (tm *tachyonManifold) get(row, col int) byte {
	if row < 0 || row >= len(tm.grid) {
		return empty
	}
	if col < 0 || col >= len(tm.grid[row]) {
		return empty
	}
	return tm.grid[row][col]
}

func (tm *tachyonManifold) set(row, col int, value byte) {
	if row < 0 || row >= len(tm.grid) {
		return
	}
	if col < 0 || col >= len(tm.grid[row]) {
		return
	}
	tm.grid[row][col] = value
}

func (tm *tachyonManifold) getBeamTimelineCount(row, col int) int {
	if row < 0 || row >= len(tm.grid) {
		return 0
	}
	if col < 0 || col >= len(tm.grid[row]) {
		return 0
	}
	return tm.beamTimelineCount[row][col]
}

func (tm *tachyonManifold) addBeamTimelines(row, col, count int) {
	if row < 0 || row >= len(tm.grid) {
		return
	}
	if col < 0 || col >= len(tm.grid[row]) {
		return
	}
	tm.beamTimelineCount[row][col] += count
}

func (tm *tachyonManifold) splitBeams() {
	for row := range tm.grid {
		tm.processRow(row, start, tm.handleStart)
		tm.processRow(row, splitter, tm.handleSplitter)
		tm.processRow(row, beam, tm.handleBeam)
	}
}

func (tm *tachyonManifold) processRow(row int, value byte, handler func(row, col int)) {
	for col := range tm.grid[row] {
		if tm.get(row, col) == value {
			handler(row, col)
		}
	}
}

func (tm *tachyonManifold) handleStart(row, col int) {
	tm.set(row+1, col, beam)
	tm.addBeamTimelines(row+1, col, 1)
}

func (tm *tachyonManifold) handleSplitter(row, col int) {
	if tm.get(row-1, col) != beam {
		return
	}

	tm.set(row, col-1, beam)
	tm.set(row, col+1, beam)

	tm.splitCount++

	beamTimelineCount := tm.getBeamTimelineCount(row-1, col)
	tm.addBeamTimelines(row, col-1, beamTimelineCount)
	tm.addBeamTimelines(row, col+1, beamTimelineCount)
}

func (tm *tachyonManifold) handleBeam(row, col int) {
	if tm.get(row+1, col) != empty {
		return
	}
	tm.set(row+1, col, beam)

	beamTimelineCount := tm.getBeamTimelineCount(row, col)
	tm.addBeamTimelines(row+1, col, beamTimelineCount)
}

func (tm *tachyonManifold) getFinalTimelineCount() int {
	lastRow := len(tm.grid) - 1

	sum := 0
	for col := range tm.grid[lastRow] {
		sum += tm.getBeamTimelineCount(lastRow, col)
	}

	return sum
}

func tachyonManifoldFromReader(r io.Reader) (*tachyonManifold, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	grid := bytes.Split(bytes.TrimSpace(input), []byte("\n"))

	if len(grid) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	startCount := 0
	for row := 0; row < len(grid); row++ {
		if len(grid[row]) != len(grid[0]) {
			return nil, fmt.Errorf("input is not a rectangle")
		}

		for col := 0; col < len(grid[row]); col++ {
			cell := grid[row][col]
			if cell != empty && cell != splitter && cell != start {
				return nil, fmt.Errorf("unknown value %q", cell)
			}
			if cell == start {
				if row > 0 {
					return nil, fmt.Errorf("start is not on first row")
				}
				startCount++
			}
		}
	}

	if startCount != 1 {
		return nil, fmt.Errorf("grid has %d starts", startCount)
	}

	beamTimelineCount := make([][]int, len(grid))
	for row := range beamTimelineCount {
		beamTimelineCount[row] = make([]int, len(grid[row]))
	}

	manifold := tachyonManifold{
		grid:              grid,
		beamTimelineCount: beamTimelineCount,
	}

	return &manifold, nil
}
