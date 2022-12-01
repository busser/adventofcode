package busser

import (
	"fmt"
	"io"
	"sort"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 10 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := codeFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	syntaxErrorScore := 0

	for _, l := range lines {
		illegalChar, missingChars := syntaxErrors(l)
		if len(missingChars) > 0 {
			continue
		}

		switch illegalChar {
		case ')':
			syntaxErrorScore += 3
		case ']':
			syntaxErrorScore += 57
		case '}':
			syntaxErrorScore += 1197
		case '>':
			syntaxErrorScore += 25137
		}
	}

	_, err = fmt.Fprintf(answer, "%d", syntaxErrorScore)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 10 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := codeFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var scores []int
	for _, l := range lines {
		_, missingChars := syntaxErrors(l)
		if len(missingChars) == 0 {
			continue
		}
		scores = append(scores, completionStringScore(missingChars))
	}
	if len(scores)%2 == 0 {
		return fmt.Errorf("expected odd number of incomplete lines, found %d", len(scores))
	}

	sort.Ints(scores)
	middleScore := scores[len(scores)/2]

	_, err = fmt.Fprintf(answer, "%d", middleScore)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type runeStack []rune

func (s *runeStack) push(r rune) {
	*s = append(*s, r)
}

func (s *runeStack) pop() rune {
	if len(*s) == 0 {
		panic("cannot pop from empty stack")
	}
	r := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return r
}

func (s *runeStack) peek() rune {
	if len(*s) == 0 {
		panic("cannot peek into empty stack")
	}
	return (*s)[len(*s)-1]
}

func syntaxErrors(line string) (illegal rune, missing []rune) {
	var stack runeStack

	for _, char := range line {
		switch char {
		case '(', '[', '{', '<':
			stack.push(char)
		case ')':
			if len(stack) == 0 || stack.pop() != '(' {
				return char, nil
			}
		case ']':
			if len(stack) == 0 || stack.pop() != '[' {
				return char, nil
			}
		case '}':
			if len(stack) == 0 || stack.pop() != '{' {
				return char, nil
			}
		case '>':
			if len(stack) == 0 || stack.pop() != '<' {
				return char, nil
			}
		}
	}

	for len(stack) > 0 {
		char := stack.pop()
		switch char {
		case '(':
			missing = append(missing, ')')
		case '[':
			missing = append(missing, ']')
		case '{':
			missing = append(missing, '}')
		case '<':
			missing = append(missing, '>')
		}
	}

	return 0, missing
}

func completionStringScore(str []rune) int {
	score := 0
	for _, char := range str {
		score *= 5
		switch char {
		case ')':
			score += 1
		case ']':
			score += 2
		case '}':
			score += 3
		case '>':
			score += 4
		}
	}
	return score
}

func codeFromReader(r io.Reader) ([]string, error) {
	return helpers.LinesFromReader(r)
}
