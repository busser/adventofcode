package d17

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 17 of Advent of Code 2024.
func PartOne(r io.Reader, w io.Writer) error {
	registerA, registerB, registerC, program, err := registersAndProgramFromReader(r)
	if err != nil {
		return err
	}

	output := runProgram(registerA, registerB, registerC, program)

	_, err = fmt.Fprintf(w, "%s", output)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 17 of Advent of Code 2024.
func PartTwo(r io.Reader, w io.Writer) error {
	_, registerB, registerC, program, err := registersAndProgramFromReader(r)
	if err != nil {
		return err
	}

	registerA := findCorrectRegisterValue(registerB, registerC, program)

	_, err = fmt.Fprintf(w, "%d", registerA)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil

	// nums := [16]int{
	// 	7, 2, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	// }
	// registerA = 0
	// for _, n := range nums {
	// 	registerA = registerA<<3 + n
	// }

	// output := runProgram(registerA, registerB, registerC, program)

	// _, err = fmt.Fprintf(w, "%s", output)
	// if err != nil {
	// 	return fmt.Errorf("could not write answer: %w", err)
	// }

	// return nil
}

type computer struct {
	registerA int
	registerB int
	registerC int

	program []int
	pointer int

	halted bool
	output []int
}

func runProgram(registerA, registerB, registerC int, program []int) string {
	c := computer{
		registerA: registerA,
		registerB: registerB,
		registerC: registerC,
		program:   program,
	}

	c.runUntilHalted()

	return c.prettyOutput()
}

func findCorrectRegisterValue(registerB, registerC int, program []int) int {
	var findWorkingValue func(depth, registerA int) (int, bool)
	findWorkingValue = func(depth, registerA int) (int, bool) {
		if depth == len(program) {
			return registerA, true
		}

		for n := 0; n < 8; n++ {
			if depth == 0 && n == 0 {
				// ensures output length is program length
				continue
			}

			newRegisterA := registerA<<3 + n

			c := computer{
				registerA: newRegisterA,
				registerB: registerB,
				registerC: registerC,
				program:   program,
			}

			c.runUntilHalted()

			desiredOutput := program[len(program)-depth-1:]

			if !slices.Equal(c.output, desiredOutput) {
				continue
			}

			if workingRegisterA, ok := findWorkingValue(depth+1, newRegisterA); ok {
				return workingRegisterA, true
			}
		}

		return 0, false
	}

	if registerA, ok := findWorkingValue(0, 0); ok {
		return registerA
	}

	panic("could not find correct register value")
}

func (c *computer) prettyOutput() string {
	nums := make([]string, len(c.output))
	for i, n := range c.output {
		nums[i] = strconv.Itoa(n)
	}
	return strings.Join(nums, ",")
}

func (c *computer) runUntilHalted() {
	for !c.halted {
		c.runSingleOp()
	}
}

func (c *computer) runSingleOp() {
	if c.pointer >= len(c.program)-1 {
		c.halted = true
		return
	}

	opcode := c.program[c.pointer]
	switch opcode {
	case 0:
		c.runADV()
	case 1:
		c.runBXL()
	case 2:
		c.runBST()
	case 3:
		c.runJNZ()
	case 4:
		c.runBXC()
	case 5:
		c.runOUT()
	case 6:
		c.runBDV()
	case 7:
		c.runCDV()
	default:
		panic(fmt.Sprintf("unknown opcode %d", opcode))
	}
}

func (c *computer) runADV() {
	c.registerA /= 1 << c.comboOperand()
	c.pointer += 2
}

func (c *computer) runBXL() {
	c.registerB ^= c.literalOperand()
	c.pointer += 2
}

func (c *computer) runBST() {
	c.registerB = c.comboOperand() % 8
	c.pointer += 2
}

func (c *computer) runJNZ() {
	if c.registerA != 0 {
		c.pointer = c.literalOperand()
	} else {
		c.pointer += 2
	}
}

func (c *computer) runBXC() {
	c.registerB ^= c.registerC
	c.pointer += 2
}

func (c *computer) runOUT() {
	c.output = append(c.output, c.comboOperand()%8)
	c.pointer += 2
}

func (c *computer) runBDV() {
	c.registerB = c.registerA / (1 << c.comboOperand())
	c.pointer += 2
}

func (c *computer) runCDV() {
	c.registerC = c.registerA / (1 << c.comboOperand())
	c.pointer += 2
}

func (c *computer) literalOperand() int {
	return c.program[c.pointer+1]
}

func (c *computer) comboOperand() int {
	operand := c.program[c.pointer+1]
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return c.registerA
	case 5:
		return c.registerB
	case 6:
		return c.registerC
	default:
		panic(fmt.Sprintf("invalid combo operand %d", operand))
	}
}

func registersAndProgramFromReader(r io.Reader) (int, int, int, []int, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("could not read input: %w", err)
	}
	if len(lines) != 5 {
		return 0, 0, 0, nil, fmt.Errorf("expected 5 lines, got %d", len(lines))
	}

	registerA := helpers.IntsFromString(lines[0])
	if len(registerA) != 1 {
		return 0, 0, 0, nil, fmt.Errorf("expected 1 number, got %d", len(registerA))
	}
	registerB := helpers.IntsFromString(lines[1])
	if len(registerB) != 1 {
		return 0, 0, 0, nil, fmt.Errorf("expected 1 number, got %d", len(registerB))
	}
	registerC := helpers.IntsFromString(lines[2])
	if len(registerC) != 1 {
		return 0, 0, 0, nil, fmt.Errorf("expected 1 number, got %d", len(registerC))
	}

	program := helpers.IntsFromString(lines[4])

	return registerA[0], registerB[0], registerC[0], program, nil
}
