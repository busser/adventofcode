package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 14 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	prog, err := programFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	memory := make(map[uint64]uint64)
	prog.run(memory)

	sum := uint64(0)
	for _, v := range memory {
		sum += v
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 14 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	prog, err := programFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	memory := make(map[uint64]uint64)
	prog.runV2(memory)

	sum := uint64(0)
	for _, v := range memory {
		sum += v
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const intSize = 36

type program struct {
	blocks []block
}

type block struct {
	mask        [intSize]byte
	assignments []assignment
}

type assignment struct {
	address uint64
	value   uint64
}

func (p program) run(memory map[uint64]uint64) {
	for _, b := range p.blocks {
		b.run(memory)
	}
}

func (b block) run(memory map[uint64]uint64) {
	for _, a := range b.assignments {
		memory[a.address] = maskedValue(a.value, b.mask)
	}
}

func maskedValue(v uint64, mask [intSize]byte) uint64 {
	masked := v

	for i := 0; i < intSize; i++ {
		switch mask[i] {
		case 'X':
		case '0':
			masked = masked &^ (1 << (intSize - i - 1))
		case '1':
			masked = masked | (1 << (intSize - i - 1))
		default:
			panic(fmt.Sprintf("unknown symbol in mask: %q", mask[i]))
		}
	}

	return masked
}

func (p program) runV2(memory map[uint64]uint64) {
	for _, b := range p.blocks {
		b.runV2(memory)
	}
}

func (b block) runV2(memory map[uint64]uint64) {
	for _, a := range b.assignments {
		for _, ma := range maskedAddressList(a.address, b.mask) {
			memory[ma] = a.value
		}
	}
}

func maskedAddressList(a uint64, mask [intSize]byte) []uint64 {
	xCount := 0
	for _, b := range mask {
		if b == 'X' {
			xCount++
		}
	}

	masked := make([]uint64, 0, 1<<xCount)

	for i := uint64(0); i < 1<<xCount; i++ {
		masked = append(masked, maskedAddress(a, mask, i))
	}

	return masked
}

func maskedAddress(a uint64, mask [intSize]byte, version uint64) uint64 {
	xIndex := 0

	masked := a

	for i := range mask {
		switch mask[i] {
		case '0':
		case '1':
			masked = masked | (1 << (intSize - i - 1))
		case 'X':
			if version&(1<<xIndex) != 0 {
				masked = masked | (1 << (intSize - i - 1))
			} else {
				masked = masked &^ (1 << (intSize - i - 1))
			}
			xIndex++
		default:
			panic(fmt.Sprintf("unknown symbol in mask: %q", mask[i]))
		}
	}

	return masked
}

func programFromReader(r io.Reader) (program, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return program{}, fmt.Errorf("reading lines: %w", err)
	}

	var prog program
	var blk block

	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "mask"):
			if i > 0 {
				prog.blocks = append(prog.blocks, blk)
				blk = block{}
			}

			mask, err := maskFrom(line)
			if err != nil {
				return program{}, fmt.Errorf("reading mask on line %d: %w", i+1, err)
			}

			blk.mask = mask

		case strings.HasPrefix(line, "mem"):
			a, err := assignmentFrom(line)
			if err != nil {
				return program{}, fmt.Errorf("reading assignment on line %d: %w", i+1, err)
			}

			blk.assignments = append(blk.assignments, a)

		default:
			return program{}, errors.New("wrong format")
		}
	}

	prog.blocks = append(prog.blocks, blk)

	return prog, nil
}

func maskFrom(s string) ([intSize]byte, error) {
	const prefix = "mask = "

	if len(s) != len(prefix)+intSize {
		return [intSize]byte{}, errors.New("wrong format")
	}

	if !strings.HasPrefix(s, prefix) {
		return [intSize]byte{}, errors.New("wrong format")
	}

	var mask [intSize]byte
	copy(mask[:], s[len(prefix):])

	return mask, nil
}

func assignmentFrom(s string) (assignment, error) {
	const prefix, sep = "mem[", "] = "

	if !strings.HasPrefix(s, prefix) {
		return assignment{}, errors.New("wrong format")
	}

	parts := strings.Split(s[len(prefix):], sep)
	if len(parts) != 2 {
		return assignment{}, errors.New("wrong format")
	}

	address, err := strconv.Atoi(parts[0])
	if err != nil {
		return assignment{}, fmt.Errorf("wrong format: %w", err)
	}

	value, err := strconv.Atoi(parts[1])
	if err != nil {
		return assignment{}, fmt.Errorf("wrong format: %w", err)
	}

	return assignment{uint64(address), uint64(value)}, nil
}
