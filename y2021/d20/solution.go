package busser

import (
	"errors"
	"fmt"
	"io"
	"math"
	"runtime"
	"sync"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 20 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	algorithm, image, err := algorithmAndImageFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	ie := newImageEnhancer(algorithm, image, 2)
	for i := 0; i < 2; i++ {
		ie.enhance()
	}

	_, err = fmt.Fprintf(answer, "%d", ie.pixelCount())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 20 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	algorithm, image, err := algorithmAndImageFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	ie := newImageEnhancer(algorithm, image, 50)
	for i := 0; i < 50; i++ {
		ie.enhance()
	}

	_, err = fmt.Fprintf(answer, "%d", ie.pixelCount())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type imageEnhancer struct {
	algorithm                   []bool
	currentImage, nextImage     [][]bool
	currentOutside, nextOutside bool
}

func newImageEnhancer(algorithm []bool, image [][]bool, margin int) *imageEnhancer {
	ie := new(imageEnhancer)

	ie.algorithm = algorithm

	ie.currentImage = newImage(len(image) + 2*margin)
	ie.nextImage = newImage(len(image) + 2*margin)

	for y := range image {
		for x := range image[y] {
			ie.currentImage[y+margin][x+margin] = image[y][x]
		}
	}

	return ie
}

func (ie *imageEnhancer) enhance() {
	// Rows can be enhanced in parallel for faster results.
	rows := make(chan int, len(ie.currentImage))
	for y := range ie.currentImage {
		rows <- y
	}
	close(rows)

	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			for y := range rows {
				for x := range ie.currentImage[y] {
					ie.nextImage[y][x] = ie.nextPixelState(x, y)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	switch ie.currentOutside {
	case true:
		ie.nextOutside = ie.algorithm[1<<9-1]
	case false:
		ie.nextOutside = ie.algorithm[0]
	}

	ie.currentImage, ie.nextImage = ie.nextImage, ie.currentImage
	ie.currentOutside, ie.nextOutside = ie.nextOutside, ie.currentOutside
}

func (ie *imageEnhancer) nextPixelState(x, y int) bool {
	var key int
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			key *= 2
			if ie.pixelIsOn(x+dx, y+dy) {
				key++
			}
		}
	}
	return ie.algorithm[key]
}

func (ie *imageEnhancer) pixelIsOn(x, y int) bool {
	if y < 0 || y >= len(ie.currentImage) || x < 0 || x >= len(ie.currentImage[y]) {
		return ie.currentOutside
	}
	return ie.currentImage[y][x]
}

func (ie *imageEnhancer) pixelCount() int {
	if ie.currentOutside {
		return math.MaxInt
	}

	count := 0
	for y := range ie.currentImage {
		for x := range ie.currentImage[y] {
			if ie.currentImage[y][x] {
				count++
			}
		}
	}

	return count
}

func newImage(size int) [][]bool {
	img := make([][]bool, size)
	for i := range img {
		img[i] = make([]bool, size)
	}
	return img
}

func algorithmAndImageFromReader(r io.Reader) (algorithm []bool, image [][]bool, err error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, nil, err
	}

	if len(lines) < 3 {
		return nil, nil, errors.New("too few lines")
	}

	algorithm = make([]bool, len(lines[0]))
	if len(algorithm) != 1<<9 {
		return nil, nil, fmt.Errorf("invalid algorithm: expected %d symbols, got %d", 1<<9, len(algorithm))
	}
	for i, c := range lines[0] {
		switch c {
		case '.':
			algorithm[i] = false
		case '#':
			algorithm[i] = true
		default:
			return nil, nil, fmt.Errorf("unknown symbol %q", c)
		}
	}

	image = make([][]bool, len(lines)-2)
	for y, l := range lines[2:] {
		image[y] = make([]bool, len(l))
		for x, c := range l {
			switch c {
			case '.':
				image[y][x] = false
			case '#':
				image[y][x] = true
			default:
				return nil, nil, fmt.Errorf("unknown symbol %q", c)
			}
		}
	}

	return algorithm, image, nil
}
