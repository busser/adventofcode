package busser

import (
	"fmt"
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
	// Output: 208
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
	// Output: 167
}

func TestFieldValidators(t *testing.T) {
	testCases := []struct {
		field string
		value string
		valid bool
	}{
		{"byr", "2002", true},
		{"byr", "2003", false},

		{"hgt", "60in", true},
		{"hgt", "190cm", true},
		{"hgt", "190in", false},
		{"hgt", "190", false},

		{"hcl", "#123abc", true},
		{"hcl", "#123abz", false},
		{"hcl", "123abc", false},

		{"ecl", "brn", true},
		{"ecl", "wat", false},

		{"pid", "000000001", true},
		{"pid", "0123456789", false},
	}

	for _, test := range testCases {
		t.Run(fmt.Sprintf("%s:%s", test.field, test.value), func(t *testing.T) {
			actual := fieldValidators[test.field](test.value)
			if actual != test.valid {
				t.Errorf("expected: %t, got: %t", test.valid, actual)
			}
		})
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
