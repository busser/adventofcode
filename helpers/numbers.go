package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

// IntsFromString returns a slice of integers in str, where these numbers are
// separated by sep.
func IntsFromString(str, sep string) ([]int, error) {
	words := strings.Split(str, sep)

	ints := make([]int, len(words))

	for i, w := range words {
		n, err := strconv.Atoi(w)
		if err != nil {
			return nil, fmt.Errorf("%q is not an integer", w)
		}

		ints[i] = n
	}

	return ints, nil
}
