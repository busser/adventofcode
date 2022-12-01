package busser

import (
	"bytes"
	"fmt"
)

const (
	seaMonsterHeight = 3
	seaMonsterWidth  = 20
)

var seaMonster [seaMonsterHeight][seaMonsterWidth]pixel

func init() {
	b := []byte("                  # \n#    ##    ##    ###\n #  #  #  #  #  #   ")

	lines := bytes.Split(b, []byte("\n"))
	if len(lines) != seaMonsterHeight {
		panic("wrong sea monster format")
	}

	for row := 0; row < seaMonsterHeight; row++ {
		line := lines[row]
		if len(line) != seaMonsterWidth {
			panic("wrong sea monster format")
		}
		for col, cell := range line {
			switch cell {
			case ' ':
				seaMonster[row][col] = pixelAny
			case '#':
				seaMonster[row][col] = pixelFull
			default:
				panic(fmt.Sprintf("unknown cell value %q in sea monster image", cell))
			}
		}
	}
}

func waterRoughnessFromImage(image [][]pixel) int {
	var flaggedSeaMonsters = make([][]pixel, len(image))
	for row := range image {
		flaggedSeaMonsters[row] = make([]pixel, len(image[row]))
	}

	for row := 0; row < len(image)-seaMonsterHeight+1; row++ {
		for col := 0; col < len(image[row])-seaMonsterWidth+1; col++ {
			if seaMonsterInImageAtPosition(image, row, col) {
				addSeaMonster(flaggedSeaMonsters, row, col)
			}
		}
	}

	return fullPixelCount(image) - fullPixelCount(flaggedSeaMonsters)
}

func fullPixelCount(image [][]pixel) int {
	var count int
	for row := range image {
		for col := range image[row] {
			if image[row][col] == pixelFull {
				count++
			}
		}
	}
	return count
}

func addSeaMonster(image [][]pixel, imageRow, imageCol int) {
	for row := 0; row < seaMonsterHeight; row++ {
		for col := 0; col < seaMonsterWidth; col++ {
			if seaMonster[row][col] == pixelFull {
				image[imageRow+row][imageCol+col] = pixelFull
			}
		}
	}
}

func seaMonstersInImage(image [][]pixel) int {
	var count int

	for row := 0; row < len(image)-seaMonsterHeight+1; row++ {
		for col := 0; col < len(image[row])-seaMonsterWidth+1; col++ {
			if seaMonsterInImageAtPosition(image, row, col) {
				count++
			}
		}
	}

	return count
}

func seaMonsterInImageAtPosition(image [][]pixel, imageRow, imageCol int) bool {
	for row := 0; row < seaMonsterHeight; row++ {
		for col := 0; col < seaMonsterWidth; col++ {
			if seaMonster[row][col] == pixelFull && image[imageRow+row][imageCol+col] != pixelFull {
				return false
			}
		}
	}
	return true
}
