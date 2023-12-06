package d05

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
	// Output: 165788812
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
	// Output: 1928058
}

func TestMergeIntervals(t *testing.T) {
	tests := []struct {
		name      string
		intervals []interval
		want      []interval
	}{
		{
			name:      "empty",
			intervals: []interval{},
			want:      []interval{},
		},
		{
			name: "single",
			intervals: []interval{
				{0, 1},
			},
			want: []interval{
				{0, 1},
			},
		},
		{
			name: "two",
			intervals: []interval{
				{0, 1},
				{2, 3},
			},
			want: []interval{
				{0, 1},
				{2, 3},
			},
		},
		{
			name: "two overlapping",
			intervals: []interval{
				{0, 2},
				{1, 3},
			},
			want: []interval{
				{0, 3},
			},
		},
		{
			name: "three overlapping",
			intervals: []interval{
				{0, 2},
				{1, 3},
				{2, 4},
			},
			want: []interval{
				{0, 4},
			},
		},
		{
			name: "some overlapping",
			intervals: []interval{
				{0, 2},
				{1, 3},
				{2, 4},
				{5, 6},
				{7, 8},
				{8, 9},
			},
			want: []interval{
				{0, 4},
				{5, 6},
				{7, 9},
			},
		},
		{
			name: "overlapping subset",
			intervals: []interval{
				{0, 5},
				{1, 3},
				{6, 8},
			},
			want: []interval{
				{0, 5},
				{6, 8},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := mergeIntervals(test.intervals)
			if !slicesEqual(got, test.want) {
				t.Errorf("mergeIntervals(%v) = %v, want %v", test.intervals, got, test.want)
			}
		})
	}
}

func slicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if b[i] != v {
			return false
		}
	}

	return true
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
