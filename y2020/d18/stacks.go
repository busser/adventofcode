package busser

type numberStack []int

func (s numberStack) len() int {
	return len(s)
}

func (s *numberStack) push(n int) {
	*s = append(*s, n)
}

func (s numberStack) peek() int {
	if len(s) == 0 {
		panic("cannot peek into empty stack")
	}

	return s[len(s)-1]
}

func (s *numberStack) pop() int {
	if len(*s) == 0 {
		panic("cannot pop from empty stack")
	}

	n := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return n
}

type operatorStack []rune

func (s operatorStack) len() int {
	return len(s)
}

func (s *operatorStack) push(op rune) {
	*s = append(*s, op)
}

func (s operatorStack) peek() rune {
	if len(s) == 0 {
		panic("cannot peek into empty stack")
	}

	return s[len(s)-1]
}

func (s *operatorStack) pop() rune {
	if len(*s) == 0 {
		panic("cannot pop from empty stack")
	}

	op := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return op
}
