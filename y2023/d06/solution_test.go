package d06

import (
	"log"
	"os"
	"testing"

	"github.com/busser/adventofcode/helpers"
)

func ExamplePartOne() {
	file, err := os.Open("testdata/input.txt")
	if err != nil {
		log.Fatalf("could not open input file: %v", err)
	}
	defer file.Close()

	if err := PartOne(file, os.Stdout); err != nil {
		log.Fatalf("could not solve: %v", err)
	}
	// Output: 170000
}

func ExamplePartTwo() {
	file, err := os.Open("testdata/input.txt")
	if err != nil {
		log.Fatalf("could open input file: %v", err)
	}
	defer file.Close()

	if err := PartTwo(file, os.Stdout); err != nil {
		log.Fatalf("could not solve: %v", err)
	}
	// Output: 20537782
}

func TestISqrt(t *testing.T) {
	tests := []struct {
		n         int
		root      int
		remainder int
	}{
		{0, 0, 0},
		{1, 1, 0},
		{2, 1, 1},
		{3, 1, 2},
		{4, 2, 0},
		{5, 2, 1},
		{6, 2, 2},
		{7, 2, 3},
		{8, 2, 4},
		{9, 3, 0},
		{10, 3, 1},
		{99, 9, 18},
		{100, 10, 0},
		{101, 10, 1},
		{120, 10, 20},
		{121, 11, 0},
		{122, 11, 1},
	}

	for _, test := range tests {
		root, remainder := isqrt(test.n)
		if root != test.root || remainder != test.remainder {
			t.Errorf("isqrt(%d) = %d, %d; want %d, %d",
				test.n, root, remainder, test.root, test.remainder)
		}
	}
}

func Benchmark(b *testing.B) {
	testCases := map[string]struct {
		solution  helpers.Solution
		inputFile string
	}{
		"PartOne": {
			solution:  helpers.SolutionFunc(PartOne),
			inputFile: "testdata/input.txt",
		},

		"PartTwo": {
			solution:  helpers.SolutionFunc(PartTwo),
			inputFile: "testdata/input.txt",
		},
	}

	for name, test := range testCases {
		b.Run(name, func(b *testing.B) {
			helpers.BenchmarkSolution(b, test.solution, test.inputFile)
		})
	}
}
