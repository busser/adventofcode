package d03

import (
	"fmt"
	"io"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 3 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	schematic, err := schematicFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	schematic.determinePartNumbers()

	_, err = fmt.Fprintf(w, "%d", schematic.sumPartNumbers())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 3 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	schematic, err := schematicFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	schematic.determinePartNumbers()
	schematic.determineGearParts()
	schematic.determineGearRatios()

	_, err = fmt.Fprintf(w, "%d", schematic.sumGearRatios())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type schematic struct {
	elements    []*schematicElement
	lookupTable [][]*schematicElement
}

type schematicElement struct {
	row, col int
	length   int
	value    string

	// For part numbers.
	isNumber     bool
	numberValue  int
	isPartNumber bool

	// For gears.
	isSymbol  bool
	isGear    bool
	gearParts []*schematicElement
	gearRatio int
}

func (s *schematic) sumGearRatios() int {
	sum := 0

	for _, element := range s.elements {
		if element.isGear {
			sum += element.gearRatio
		}
	}

	return sum
}

func (s *schematic) determineGearRatios() {
	for _, element := range s.elements {
		if element.isGear {
			element.gearRatio = element.gearParts[0].numberValue * element.gearParts[1].numberValue
		}
	}
}

func (s *schematic) determineGearParts() {
	for _, element := range s.elements {
		s.findGearParts(element)
	}
}

func (s *schematic) findGearParts(element *schematicElement) {
	if !element.isSymbol {
		return
	}

	if element.value != "*" {
		return
	}

	var gearParts []*schematicElement
	for col := element.col - 1; col <= element.col+1; col++ {
		for row := element.row - 1; row <= element.row+1; row++ {
			neighbor := s.lookup(row, col)
			if neighbor != nil && neighbor.isPartNumber {
				gearParts = append(gearParts, neighbor)
			}
		}
	}

	gearParts = unique(gearParts)
	if len(gearParts) != 2 {
		return
	}

	element.isGear = true
	element.gearParts = gearParts
}

func unique(elems []*schematicElement) []*schematicElement {
	seen := make(map[*schematicElement]bool)
	var result []*schematicElement

	for _, elem := range elems {
		if !seen[elem] {
			seen[elem] = true
			result = append(result, elem)
		}
	}

	return result
}

func (s *schematic) sumPartNumbers() int {
	sum := 0

	for _, element := range s.elements {
		if element.isPartNumber {
			sum += element.numberValue
		}
	}

	return sum
}

func (s *schematic) determinePartNumbers() {
	for _, element := range s.elements {
		if s.fitsPartNumberCriteria(element) {
			element.isPartNumber = true
		}
	}
}

func (s *schematic) fitsPartNumberCriteria(element *schematicElement) bool {
	if !element.isNumber {
		return false
	}

	// Check left of number.
	if s.isSymbolAtPosition(element.row, element.col-1) {
		return true
	}

	// Check right of number.
	if s.isSymbolAtPosition(element.row, element.col+element.length) {
		return true
	}

	// Check above and below number.
	for col := element.col - 1; col < element.col+element.length+1; col++ {
		if s.isSymbolAtPosition(element.row-1, col) {
			return true
		}
		if s.isSymbolAtPosition(element.row+1, col) {
			return true
		}
	}

	return false
}

func (s *schematic) isSymbolAtPosition(row, col int) bool {
	elem := s.lookup(row, col)
	return elem != nil && elem.isSymbol
}

func (s *schematic) lookup(row, col int) *schematicElement {
	if row < 0 || row >= len(s.lookupTable) || col < 0 || col >= len(s.lookupTable[row]) {
		return nil
	}
	return s.lookupTable[row][col]
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isEmpty(b byte) bool {
	return b == '.'
}

func isSymbol(b byte) bool {
	return !isDigit(b) && !isEmpty(b)
}

func schematicFromReader(r io.Reader) (*schematic, error) {
	rawSchematic, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read schematic: %w", err)
	}

	var elements []*schematicElement
	lookupTable := make([][]*schematicElement, len(rawSchematic))
	for row := range rawSchematic {
		lookupTable[row] = make([]*schematicElement, len(rawSchematic[row]))
	}

	insert := func(e *schematicElement) {
		elements = append(elements, e)
		for col := e.col; col < e.col+e.length; col++ {
			lookupTable[e.row][col] = e
		}
	}

	for row := range rawSchematic {
		var numberWIP *schematicElement

		for col := range rawSchematic[row] {
			cellValue := rawSchematic[row][col]

			if numberWIP == nil && isDigit(cellValue) {
				numberWIP = &schematicElement{
					row:         row,
					col:         col,
					length:      1,
					value:       string(cellValue),
					isNumber:    true,
					numberValue: int(cellValue - '0'),
				}
				continue
			}

			if numberWIP != nil && isDigit(cellValue) {
				numberWIP.length++
				numberWIP.value += string(cellValue)
				numberWIP.numberValue = numberWIP.numberValue*10 + int(cellValue-'0')
				continue
			}

			if numberWIP != nil {
				insert(numberWIP)
				numberWIP = nil
			}

			if isSymbol(cellValue) {
				insert(&schematicElement{
					row:      row,
					col:      col,
					length:   1,
					value:    string(cellValue),
					isSymbol: true,
				})
			}
		}

		if numberWIP != nil {
			insert(numberWIP)
			numberWIP = nil
		}
	}

	return &schematic{
		elements:    elements,
		lookupTable: lookupTable,
	}, nil
}
