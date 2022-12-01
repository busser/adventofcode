package busser

import (
	"errors"
	"fmt"
)

type snailfishNumberDecoder struct {
	data  []rune
	index int
	err   error
}

func (decoder *snailfishNumberDecoder) read() rune {
	if decoder.err != nil {
		return 0
	}

	if decoder.index >= len(decoder.data) {
		decoder.err = errors.New("no more data to read")
		return 0
	}

	r := decoder.data[decoder.index]
	decoder.index++

	return r
}

func (decoder *snailfishNumberDecoder) mustRead(r rune) {
	actual := decoder.read()
	if r != actual {
		decoder.err = fmt.Errorf("expected to read %q, got %q", r, actual)
	}
}

func (decoder *snailfishNumberDecoder) parse(number *snailfishNumber) {
	if decoder.err != nil {
		return
	}

	r := decoder.read()
	switch {
	case r == '[':
		number.leftChild = new(snailfishNumber)
		number.leftChild.parent = number
		decoder.parse(number.leftChild)

		decoder.mustRead(',')

		number.rightChild = new(snailfishNumber)
		number.rightChild.parent = number
		decoder.parse(number.rightChild)

		decoder.mustRead(']')
	case r >= '0' && r <= '9':
		number.value = int(r - '0')
	default:
		decoder.err = fmt.Errorf("unexpected character %q", r)
	}
}
