package busser

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

// PartOne solves the first problem of day 4 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	passports, err := passportsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0
	for _, pass := range passports {
		if passportIsValid(pass) {
			count++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 4 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	passports, err := passportsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0
	for _, pass := range passports {
		if passportIsStrictlyValid(pass) {
			count++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func passportIsValid(pass map[string]string) bool {
	requiredFields := []string{
		"byr",
		"iyr",
		"eyr",
		"hgt",
		"hcl",
		"ecl",
		"pid",
	}

	for _, f := range requiredFields {
		if _, ok := pass[f]; !ok {
			return false
		}
	}

	return true
}

var fieldValidators = map[string]func(string) bool{
	"byr": byrIsValid,
	"iyr": iyrIsValid,
	"eyr": eyrIsValid,
	"hgt": hgtIsValid,
	"hcl": hclIsValid,
	"ecl": eclIsValid,
	"pid": pidIsValid,
}

func passportIsStrictlyValid(pass map[string]string) bool {
	for f, isValid := range fieldValidators {
		v, ok := pass[f]
		if !ok {
			return false
		}
		if !isValid(v) {
			return false
		}
	}

	return true
}

func byrIsValid(byr string) bool {
	if len(byr) != 4 {
		return false
	}

	num, err := strconv.Atoi(byr)
	if err != nil {
		return false
	}

	return num >= 1920 && num <= 2002
}

func iyrIsValid(iyr string) bool {
	if len(iyr) != 4 {
		return false
	}

	num, err := strconv.Atoi(iyr)
	if err != nil {
		return false
	}

	return num >= 2010 && num <= 2020
}

func eyrIsValid(eyr string) bool {
	if len(eyr) != 4 {
		return false
	}

	num, err := strconv.Atoi(eyr)
	if err != nil {
		return false
	}

	return num >= 2020 && num <= 2030
}

func hgtIsValid(hgt string) bool {
	l := len(hgt)
	if l < 3 {
		return false
	}

	num, err := strconv.Atoi(hgt[:l-2])
	if err != nil {
		return false
	}

	switch hgt[l-2:] {
	case "cm":
		return num >= 150 && num <= 193
	case "in":
		return num >= 59 && num <= 76
	default:
		return false
	}
}

func hclIsValid(hcl string) bool {
	if len(hcl) != 7 {
		return false
	}

	if hcl[0] != '#' {
		return false
	}

	for _, c := range hcl[1:] {
		if (c < 'a' || c > 'f') && (c < '0' || c > '9') {
			return false
		}
	}

	return true
}

func eclIsValid(ecl string) bool {
	validColors := [...]string{
		"amb",
		"blu",
		"brn",
		"gry",
		"grn",
		"hzl",
		"oth",
	}

	for _, v := range validColors {
		if ecl == v {
			return true
		}
	}

	return false
}

func pidIsValid(pid string) bool {
	if len(pid) != 9 {
		return false
	}

	for _, c := range pid {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

func passportsFromReader(r io.Reader) ([]map[string]string, error) {
	allInput, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	rawPassports := strings.Split(string(allInput), "\n\n")

	passports := make([]map[string]string, len(rawPassports))

	for i, rawPass := range rawPassports {
		pass := make(map[string]string)

		for _, field := range strings.Fields(rawPass) {
			splitField := strings.Split(field, ":")
			if len(splitField) != 2 {
				return nil, errors.New("invalid format")
			}

			pass[splitField[0]] = splitField[1]
		}

		passports[i] = pass
	}

	return passports, nil
}
