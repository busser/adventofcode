package helpers

import "io"

// A Solution solves an Advent of Code problem.
type Solution interface {
	// Solve reads the problem input and writes the answer. If something goes
	// wrong, Solve returns an error.
	Solve(input io.Reader, answer io.Writer) error
}

// The SolutionFunc type is an adapter to allow the use of ordinary functions as
// solutions. If f is a function with the appropriate signature, SolutionFunc(f)
// is a Solution that calls f.
type SolutionFunc func(input io.Reader, answer io.Writer) error

// Solve calls f(answer, input).
func (f SolutionFunc) Solve(input io.Reader, answer io.Writer) error {
	return f(input, answer)
}
