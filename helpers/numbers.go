package helpers

import (
	"fmt"
	"strconv"
)

// IntsFromString returns a slice of integers in str, where these numbers are
// separated by non-number runes. If a dash preceding a number is the only
// non-number rune between two numbers, it is considered a separator; otherwise,
// it is considered part of the number, which is consequently negative.
func IntsFromString(str string) []int {
	words := splitStringIntoIntStrings(str)

	ints := make([]int, len(words))

	for i, w := range words {
		n, err := strconv.Atoi(w)
		if err != nil {
			panic(fmt.Sprintf("could not parse int %q: %v", w, err))
		}

		ints[i] = n
	}

	return ints
}

func splitStringIntoIntStrings(str string) []string {
	var (
		words   []string
		wordBuf []byte
	)

	for i := range str {
		b := str[i]

		if b >= '0' && b <= '9' {
			wordBuf = append(wordBuf, b)
			continue
		}

		if b == '-' && len(wordBuf) == 0 {
			wordBuf = append(wordBuf, b)
			continue
		}

		if len(wordBuf) > 0 {
			words = append(words, string(wordBuf))
			wordBuf = wordBuf[:0] // TODO: bench vs setting to nil
		}
	}

	if len(wordBuf) > 0 {
		words = append(words, string(wordBuf))
	}

	return words
}
