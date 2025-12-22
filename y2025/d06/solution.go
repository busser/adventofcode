package d06

import (
	"bytes"
	"fmt"
	"io"
)

// PartOne solves the first problem of day 6 of Advent of Code 2025.
func PartOne(r io.Reader, w io.Writer) error {
	worksheet, err := worksheetFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	problems := readProblemsAsHuman(worksheet)

	sum := 0
	for _, p := range problems {
		sum += p.answer()
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 6 of Advent of Code 2025.
func PartTwo(r io.Reader, w io.Writer) error {
	worksheet, err := worksheetFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	problems := readProblemsAsCephalopod(worksheet)

	sum := 0
	for _, p := range problems {
		sum += p.answer()
	}

	_, err = fmt.Fprintf(w, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const (
	add      = '+'
	multiply = '*'
)

type problem struct {
	numbers   []int
	operation byte
}

func (p problem) answer() int {
	switch p.operation {
	case add:
		return sumOf(p.numbers)
	case multiply:
		return productOf(p.numbers)
	default:
		panic(fmt.Sprintf("unknown operation %q", p.operation))
	}
}

func (p problem) String() string {
	return fmt.Sprintf("(%v %c)", p.numbers, p.operation)
}

func sumOf(numbers []int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return sum
}

func productOf(numbers []int) int {
	product := 1
	for _, n := range numbers {
		product *= n
	}
	return product
}

func worksheetFromReader(r io.Reader) ([][]byte, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	worksheet := bytes.Split(input, []byte("\n"))
	for len(worksheet[len(worksheet)-1]) == 0 {
		worksheet = worksheet[:len(worksheet)-1]
	}

	if len(worksheet) < 2 {
		return nil, fmt.Errorf("too few lines")
	}
	for row := range worksheet {
		if len(worksheet[row]) != len(worksheet[0]) {
			return nil, fmt.Errorf("input in not a rectangle")
		}
	}

	return worksheet, nil
}

func readProblemsAsHuman(worksheet [][]byte) []problem {
	rawProblems := splitWorksheetIntoParts(worksheet)

	problems := make([]problem, len(rawProblems))
	for i, rawProblem := range rawProblems {
		problems[i] = readProblemAsHuman(rawProblem)
	}

	return problems
}

func readProblemAsHuman(rawProblem [][]byte) problem {
	numbers := make([]int, len(rawProblem)-1)
	for row := 0; row < len(rawProblem)-1; row++ {
		for col := 0; col < len(rawProblem[row]); col++ {
			cell := rawProblem[row][col]
			if cell == ' ' {
				continue
			}
			numbers[row] *= 10
			numbers[row] += int(cell - '0')
		}
	}

	operation := rawProblem[len(rawProblem)-1][0]

	return problem{numbers, operation}
}

func readProblemsAsCephalopod(worksheet [][]byte) []problem {
	rawProblems := splitWorksheetIntoParts(worksheet)

	problems := make([]problem, len(rawProblems))
	for i, rawProblem := range rawProblems {
		problems[i] = readProblemAsCephalopod(rawProblem)
	}

	return problems
}

func readProblemAsCephalopod(rawProblem [][]byte) problem {
	numbers := make([]int, len(rawProblem[0]))
	for row := 0; row < len(rawProblem)-1; row++ {
		for col := 0; col < len(rawProblem[row]); col++ {
			cell := rawProblem[row][col]
			if cell == ' ' {
				continue
			}
			numbers[col] *= 10
			numbers[col] += int(cell - '0')
		}
	}

	operation := rawProblem[len(rawProblem)-1][0]

	return problem{numbers, operation}
}

func splitWorksheetIntoParts(worksheet [][]byte) [][][]byte {
	var parts [][][]byte

	numRows := len(worksheet)
	numCols := len(worksheet[0])

	part := make([][]byte, numRows)
	for col := range numCols {
		columnIsEmpty := true
		for row := range worksheet {
			if worksheet[row][col] != ' ' {
				columnIsEmpty = false
				break
			}
		}

		if columnIsEmpty {
			parts = append(parts, part)
			part = make([][]byte, numRows)
			continue
		}

		for row := range numRows {
			part[row] = append(part[row], worksheet[row][col])
		}
	}

	parts = append(parts, part)

	return parts
}
