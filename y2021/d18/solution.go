package busser

import (
	"errors"
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 18 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	numbers, err := snailfishNumbersFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	num := numbers[0]
	for _, n := range numbers[1:] {
		num = addNumbers(num, n)
	}

	_, err = fmt.Fprintf(answer, "%d", magnitude(num))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 18 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	numbers, err := snailfishNumbersFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	max := 0

	for i := range numbers {
		for j := range numbers {
			if i == j {
				continue
			}
			a := deepCopy(numbers[i])
			b := deepCopy(numbers[j])

			if mag := magnitude(addNumbers(a, b)); mag > max {
				max = mag
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", max)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type snailfishNumber struct {
	value int

	parent     *snailfishNumber
	leftChild  *snailfishNumber
	rightChild *snailfishNumber
}

func (num *snailfishNumber) String() string {
	if num.leftChild == nil {
		return fmt.Sprintf("%d", num.value)
	}
	return fmt.Sprintf("[%s,%s]", num.leftChild, num.rightChild)
}

func deepCopy(number *snailfishNumber) *snailfishNumber {
	newNumber := new(snailfishNumber)
	newNumber.value = number.value

	if number.leftChild != nil {
		newNumber.leftChild = deepCopy(number.leftChild)
		newNumber.leftChild.parent = newNumber
	}
	if number.rightChild != nil {
		newNumber.rightChild = deepCopy(number.rightChild)
		newNumber.rightChild.parent = newNumber
	}

	return newNumber
}

func addNumbers(a, b *snailfishNumber) *snailfishNumber {
	parent := snailfishNumber{
		leftChild:  a,
		rightChild: b,
	}
	a.parent = &parent
	b.parent = &parent
	reduce(&parent)
	return &parent
}

func magnitude(node *snailfishNumber) int {
	if node.leftChild == nil {
		return node.value
	}
	return 3*magnitude(node.leftChild) + 2*magnitude(node.rightChild)
}

func reduce(root *snailfishNumber) {
	for reduceOnce(root) {
	}
}

func reduceOnce(root *snailfishNumber) (reduced bool) {
	shouldExplode := firstNumberThatShouldExplode(root)
	if shouldExplode != nil {
		explode(shouldExplode)
		return true
	}

	shouldSplit := firstNumberThatShouldSplit(root)
	if shouldSplit != nil {
		split(shouldSplit)
		return true
	}

	return false
}

func explode(node *snailfishNumber) {
	toTheLeft := firstNumberToTheLeftOf(node.leftChild)
	toTheRight := firstNumberToTheRightOf(node.rightChild)

	if toTheLeft != nil {
		toTheLeft.value += node.leftChild.value
	}
	if toTheRight != nil {
		toTheRight.value += node.rightChild.value
	}

	node.leftChild, node.rightChild = nil, nil
	node.value = 0
}

func firstNumberThatShouldExplode(root *snailfishNumber) *snailfishNumber {
	var walk func(*snailfishNumber, int) *snailfishNumber
	walk = func(node *snailfishNumber, nesting int) (shouldExplode *snailfishNumber) {
		if node.leftChild == nil && node.rightChild == nil {
			return nil
		}
		if nesting == 4 {
			return node
		}
		if shouldExplode = walk(node.leftChild, nesting+1); shouldExplode != nil {
			return shouldExplode
		}
		if shouldExplode = walk(node.rightChild, nesting+1); shouldExplode != nil {
			return shouldExplode
		}
		return nil
	}

	return walk(root, 0)
}

func split(node *snailfishNumber) {
	left, right := new(snailfishNumber), new(snailfishNumber)

	left.value = node.value / 2
	right.value = node.value - left.value

	node.value = 0
	node.leftChild = left
	node.rightChild = right
	left.parent = node
	right.parent = node
}

func firstNumberThatShouldSplit(root *snailfishNumber) *snailfishNumber {
	var walk func(*snailfishNumber) *snailfishNumber
	walk = func(node *snailfishNumber) (shouldSplit *snailfishNumber) {
		if node.leftChild == nil && node.rightChild == nil {
			if node.value >= 10 {
				return node
			}
			return nil
		}
		if shouldSplit = walk(node.leftChild); shouldSplit != nil {
			return shouldSplit
		}
		if shouldSplit = walk(node.rightChild); shouldSplit != nil {
			return shouldSplit
		}
		return nil
	}

	return walk(root)
}

func leftMostChildOf(node *snailfishNumber) *snailfishNumber {
	for node.leftChild != nil {
		node = node.leftChild
	}
	return node
}

func rightMostChildOf(node *snailfishNumber) *snailfishNumber {
	for node.rightChild != nil {
		node = node.rightChild
	}
	return node
}

func firstNumberToTheLeftOf(node *snailfishNumber) *snailfishNumber {
	if node.leftChild != nil {
		// In this case, node is a pair and not a regular number.
		return rightMostChildOf(node.leftChild)
	}

	p := node.parent
	for p != nil && node == p.leftChild {
		node, p = p, p.parent
	}
	if p == nil {
		return nil
	}
	return rightMostChildOf(p.leftChild)
}

func firstNumberToTheRightOf(node *snailfishNumber) *snailfishNumber {
	if node.rightChild != nil {
		// In this case, node is a pair and not a regular number.
		return leftMostChildOf(node.rightChild)
	}

	p := node.parent
	for p != nil && node == p.rightChild {
		node, p = p, p.parent
	}
	if p == nil {
		return nil
	}
	return leftMostChildOf(p.rightChild)
}

func snailfishNumbersFromReader(r io.Reader) ([]*snailfishNumber, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	numbers := make([]*snailfishNumber, len(lines))
	for i := range lines {
		numbers[i], err = snailfishNumberFromString(lines[i])
		if err != nil {
			return nil, fmt.Errorf("invalid number on line %d: %w", i+1, err)
		}
	}

	return numbers, nil
}

func snailfishNumberFromString(s string) (*snailfishNumber, error) {
	var number snailfishNumber
	d := snailfishNumberDecoder{data: []rune(s)}

	d.parse(&number)
	if d.err != nil {
		return nil, d.err
	}
	if d.index != len(s) {
		return nil, errors.New("invalid syntax")
	}

	return &number, nil
}
