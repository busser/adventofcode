package helpers

import (
	"bytes"
	"io/ioutil"
	"testing"
)

// TestSolution tests whether s, when provided with input, provides the expected
// answer.
func TestSolution(t *testing.T, s Solution, inputFile, answerFile string) {
	t.Helper()

	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Fatalf("could not read input file: %v", err)
	}

	answer, err := ioutil.ReadFile(answerFile)
	if err != nil {
		t.Fatalf("could not read answer file: %v", err)
	}
	if len(answer) == 0 {
		t.Fatalf("\n👉 Write the answer to %s\n", answerFile)
	}

	r := bytes.NewReader(input)
	w := &bytes.Buffer{}

	if err := s.Solve(r, w); err != nil {
		t.Fatalf("error running solution: %v", err)
	}

	actual := w.Bytes()
	if !bytes.Equal(answer, actual) {
		t.Fatalf("did not get expected answer:\n\texpected: %q\n\tgot: %q", answer, actual)
	}
}

// BenchmarkSolution runs a benchmark of s with the provided input.
func BenchmarkSolution(b *testing.B, s Solution, inputFile string) {
	b.Helper()

	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		b.Fatalf("could not read input file: %v", err)
	}

	r := bytes.NewReader(input)
	w := ioutil.Discard

	for n := 0; n < b.N; n++ {
		r.Reset(input)
		_ = s.Solve(r, w)
	}
}
