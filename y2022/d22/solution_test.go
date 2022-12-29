package d22

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
	// Output: 189140
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
	// Output: 115063
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

func TestCubeWrapping(t *testing.T) {
	tests := []struct {
		position     vector
		facing       vector
		wantPosition vector
		wantFacing   vector
	}{
		// U0 -> L1
		{vector{0, 100}, facingUp, vector{50, 50}, facingRight},
		{vector{49, 100}, facingUp, vector{50, 99}, facingRight},
		// U1 -> L3
		{vector{50, 0}, facingUp, vector{0, 150}, facingRight},
		{vector{99, 0}, facingUp, vector{0, 199}, facingRight},
		// U2 -> D0
		{vector{100, 0}, facingUp, vector{0, 199}, facingUp},
		{vector{149, 0}, facingUp, vector{49, 199}, facingUp},
		// L0 -> L2
		{vector{50, 0}, facingLeft, vector{0, 149}, facingRight},
		{vector{50, 49}, facingLeft, vector{0, 100}, facingRight},
		// L1 -> U0
		{vector{50, 50}, facingLeft, vector{0, 100}, facingDown},
		{vector{50, 99}, facingLeft, vector{49, 100}, facingDown},
		// L2 -> L0
		{vector{0, 100}, facingLeft, vector{50, 49}, facingRight},
		{vector{0, 149}, facingLeft, vector{50, 0}, facingRight},
		// L3 -> U1
		{vector{0, 150}, facingLeft, vector{50, 0}, facingDown},
		{vector{0, 199}, facingLeft, vector{99, 0}, facingDown},
		// R0 -> R2
		{vector{149, 0}, facingRight, vector{99, 149}, facingLeft},
		{vector{149, 49}, facingRight, vector{99, 100}, facingLeft},
		// R1 -> D2
		{vector{99, 50}, facingRight, vector{100, 49}, facingUp},
		{vector{99, 99}, facingRight, vector{149, 49}, facingUp},
		// R2 -> R0
		{vector{99, 100}, facingRight, vector{149, 49}, facingLeft},
		{vector{99, 149}, facingRight, vector{149, 0}, facingLeft},
		// R3 -> D1
		{vector{49, 150}, facingRight, vector{50, 149}, facingUp},
		{vector{49, 199}, facingRight, vector{99, 149}, facingUp},
		// D0 -> U2
		{vector{0, 199}, facingDown, vector{100, 0}, facingDown},
		{vector{49, 199}, facingDown, vector{149, 0}, facingDown},
		// D1 -> R3
		{vector{50, 149}, facingDown, vector{49, 150}, facingLeft},
		{vector{99, 149}, facingDown, vector{49, 199}, facingLeft},
		// D2 -> R1
		{vector{100, 49}, facingDown, vector{99, 50}, facingLeft},
		{vector{149, 49}, facingDown, vector{99, 99}, facingLeft},
	}

	b := make(board, 50*3)
	for x := range b {
		b[x] = make([]tile, 50*4)
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("pos=%v f=%v", tt.position, tt.facing), func(t *testing.T) {
			actualPosition, actualFacing := b.nextPositionAndFacingWithCubeWrapping(tt.position, tt.facing)

			if actualPosition != tt.wantPosition {
				t.Errorf("got position %#v, want %#v", actualPosition, tt.wantPosition)
			}

			if actualFacing != tt.wantFacing {
				t.Errorf("got facing %#v, want %#v", actualFacing, tt.wantFacing)
			}
		})
	}
}
