package busser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// PartOne solves the first problem of day 6 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	groups, err := groupsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0

	for _, g := range groups {
		for question := 0; question < 26; question++ {
			var anyoneSaidYes bool

			for _, p := range g.passengers {
				if p.answers[question] {
					anyoneSaidYes = true
				}
			}

			if anyoneSaidYes {
				count++
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 6 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	groups, err := groupsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0

	for _, g := range groups {
		for question := 0; question < 26; question++ {
			var anyoneSaidYes, anyoneSaidNo bool

			for _, p := range g.passengers {
				if p.answers[question] {
					anyoneSaidYes = true
				} else {
					anyoneSaidNo = true
				}
			}

			if anyoneSaidYes && !anyoneSaidNo {
				count++
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type group struct {
	passengers []passenger
}

type passenger struct {
	answers [26]bool
}

func groupsFromReader(r io.Reader) ([]group, error) {
	allInput, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	rawGroups := strings.Split(string(bytes.TrimSpace(allInput)), "\n\n")

	groups := make([]group, len(rawGroups))
	for i, rg := range rawGroups {
		rawPassenger := strings.Split(rg, "\n")

		groups[i].passengers = make([]passenger, len(rawPassenger))
		for j, rp := range rawPassenger {

			for _, answer := range rp {
				if answer < 'a' || answer > 'z' {
					return nil, errors.New("wrong format")
				}

				groups[i].passengers[j].answers[answer-'a'] = true
			}
		}
	}

	return groups, nil
}
