package busser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 19 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	reports, err := scannerReportsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	numBeacons, _ := combineReports(reports)

	_, err = fmt.Fprintf(answer, "%d", numBeacons)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 19 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	reports, err := scannerReportsFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_, maxDistance := combineReports(reports)

	_, err = fmt.Fprintf(answer, "%d", maxDistance)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type vector [3]int

type vectorSignature struct {
	sum, max int
}

func (v vector) signature() vectorSignature {
	x, y, z := abs(v[0]), abs(v[1]), abs(v[2])
	return vectorSignature{
		sum: x + y + z,
		max: max(x, max(y, z)),
	}
}

func (v vector) plus(w vector) vector {
	return vector{
		v[0] + w[0],
		v[1] + w[1],
		v[2] + w[2],
	}
}

func (v vector) minus(w vector) vector {
	return vector{
		v[0] - w[0],
		v[1] - w[1],
		v[2] - w[2],
	}
}

func (v vector) rotate(m matrix) vector {
	var w vector
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			w[i] += m[i][j] * v[j]
		}
	}
	return w
}

type matrix [3][3]int

func (this matrix) multipliedBy(other matrix) matrix {
	var product matrix
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				product[i][j] += this[i][k] * other[k][j]
			}
		}
	}
	return product
}

var rotations = generateRotations()

func generateRotations() []matrix {
	const cos90, sin90 = 0, 1

	unity := matrix{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}

	// 90° along X axis.
	rotX90 := matrix{
		{1, 0, 0},
		{0, cos90, -sin90},
		{0, sin90, cos90},
	}

	// 90° along Y axis.
	rotY90 := matrix{
		{cos90, 0, sin90},
		{0, 1, 0},
		{-sin90, 0, cos90},
	}

	// 90° along Z axis.
	rotZ90 := matrix{
		{cos90, -sin90, 0},
		{sin90, cos90, 0},
		{0, 0, 1},
	}

	combineRotations := func(x, y, z int) matrix {
		combo := unity
		for i := 0; i < x; i++ {
			combo = combo.multipliedBy(rotX90)
		}
		for i := 0; i < y; i++ {
			combo = combo.multipliedBy(rotY90)
		}
		for i := 0; i < z; i++ {
			combo = combo.multipliedBy(rotZ90)
		}
		return combo
	}

	allRotations := make(map[matrix]struct{})
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			for z := 0; z < 4; z++ {
				r := combineRotations(x, y, z)
				allRotations[r] = struct{}{}
			}
		}
	}

	uniqueRotations := make([]matrix, 0, len(allRotations))
	for r := range allRotations {
		uniqueRotations = append(uniqueRotations, r)
	}

	return uniqueRotations
}

type scannerReport struct {
	scannerID int
	beacons   []vector
}

func combineReports(reports []scannerReport) (numBeacons, maxDistance int) {
	if len(reports) == 0 {
		return 0, 0
	}

	// List all beacon-pair signatures in each report.
	reportSignatures := make([]map[vectorSignature][2]vector, len(reports))
	for i := range reportSignatures {
		reportSignatures[i] = make(map[vectorSignature][2]vector)
	}
	for r := range reports {
		for i, b1 := range reports[r].beacons {
			for _, b2 := range reports[r].beacons[i:] {
				sig := b1.minus(b2).signature()
				reportSignatures[r][sig] = [2]vector{b1, b2}
			}
		}
	}

	// We define the first scanner to be at point (0, 0, 0).
	scannerPositions := make([]vector, len(reports))
	scannerPositions[0] = vector{0, 0, 0}
	reportProcessed := make([]bool, len(reports))
	reportProcessed[0] = true

	// We consider the beacons and the beacon-pair signatures in the first
	// report as already known.
	knownBeacons := make(map[vector]bool)
	for _, b := range reports[0].beacons {
		knownBeacons[b] = true
	}
	knownSignatures := make(map[vectorSignature][2]vector)
	for sig, pair := range reportSignatures[0] {
		knownSignatures[sig] = pair
	}

	matchNewScanner := func() (matched bool) {
		for r := range reports {
			if reportProcessed[r] {
				continue
			}

			matchingSignatures := make(map[vectorSignature][2]vector)
			for sig, pair := range reportSignatures[r] {
				if _, known := knownSignatures[sig]; known {
					matchingSignatures[sig] = pair
				}
			}

			if len(matchingSignatures) < 66 { // 66 == 12 choose 2.
				continue
			}

			for sig, newPair := range matchingSignatures {
				knownPair, ok := knownSignatures[sig]
				if !ok {
					// This new pair of beacons does not match any know pair.
					continue
				}

				var possibleRotations []matrix
				for _, rot := range rotations {
					if newPair[0].rotate(rot).minus(knownPair[0]) == newPair[1].rotate(rot).minus(knownPair[1]) {
						possibleRotations = append(possibleRotations, rot)
					}
				}
				if len(possibleRotations) == 0 {
					// The new pair having the same signature as a known pair
					// was a false positive.
					continue
				}

				for _, rot := range possibleRotations {
					newScannerPosition := newPair[0].rotate(rot).minus(knownPair[0])

					var translatedReport scannerReport
					translatedReport.beacons = make([]vector, len(reports[r].beacons))
					for i, b := range reports[r].beacons {
						translatedReport.beacons[i] = b.rotate(rot).minus(newScannerPosition)
					}

					beaconMatches := 0
					for _, b := range translatedReport.beacons {
						if knownBeacons[b] {
							beaconMatches++
						}
					}

					if beaconMatches >= 12 {
						scannerPositions[r] = newScannerPosition
						reportProcessed[r] = true
						reportSignatures[r] = make(map[vectorSignature][2]vector)
						for i, b1 := range translatedReport.beacons {
							for _, b2 := range translatedReport.beacons[i:] {
								sig := b1.minus(b2).signature()
								reportSignatures[r][sig] = [2]vector{b1, b2}
							}
						}
						for sig, pair := range reportSignatures[r] {
							knownSignatures[sig] = pair
						}
						for _, b := range translatedReport.beacons {
							knownBeacons[b] = true
						}

						return true
					}
				}
			}
		}

		return false
	}

	for matchNewScanner() {
	}

	if len(scannerPositions) != len(reports) {
		panic("failed to locate all scanners")
	}

	maxDistance = 0
	for _, s := range scannerPositions {
		for _, t := range scannerPositions {
			dist := abs(s[0]-t[0]) + abs(s[1]-t[1]) + abs(s[2]-t[2])
			if dist > maxDistance {
				maxDistance = dist
			}
		}
	}

	return len(knownBeacons), maxDistance
}

func scannerReportsFromReader(r io.Reader) ([]scannerReport, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, err
	}

	rawReports := splitSlice(lines)
	reports := make([]scannerReport, len(rawReports))

	for i := range rawReports {
		r, err := reportFromLines(rawReports[i])
		if err != nil {
			return nil, fmt.Errorf("invalid report: %w", err)
		}
		reports[i] = r
	}

	return reports, nil
}

func reportFromLines(lines []string) (scannerReport, error) {
	if len(lines) < 1 {
		return scannerReport{}, errors.New("wrong format")
	}

	headerParts := strings.SplitN(lines[0], " ", 4)
	if len(headerParts) != 4 {
		return scannerReport{}, errors.New("invalid header")
	}

	scannerID, err := strconv.Atoi(headerParts[2])
	if err != nil {
		return scannerReport{}, fmt.Errorf("invalid scanner ID %q", headerParts[2])
	}

	lines = lines[1:]

	beacons := make([]vector, len(lines))
	for i := range lines {
		b, err := vectorFromString(lines[i])
		if err != nil {
			return scannerReport{}, fmt.Errorf("invalid coordinates: %w", err)
		}
		beacons[i] = b
	}

	return scannerReport{scannerID, beacons}, nil
}

func vectorFromString(s string) (vector, error) {
	parts := strings.SplitN(s, ",", 3)
	if len(parts) != 3 {
		return vector{}, errors.New("wrong format")
	}

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return vector{}, fmt.Errorf("%q is not a number", parts[0])
	}
	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return vector{}, fmt.Errorf("%q is not a number", parts[1])
	}
	z, err := strconv.Atoi(parts[2])
	if err != nil {
		return vector{}, fmt.Errorf("%q is not a number", parts[2])
	}

	return vector{x, y, z}, nil
}

func splitSlice(s []string) [][]string {
	var parts [][]string
	var p []string
	for _, v := range s {
		if len(v) == 0 {
			parts = append(parts, p)
			p = nil
			continue
		}
		p = append(p, v)
	}
	if len(p) > 0 {
		parts = append(parts, p)
	}
	return parts
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
