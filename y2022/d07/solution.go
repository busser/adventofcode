package d07

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 7 of Advent of Code 2022.
func PartOne(r io.Reader, w io.Writer) error {
	root, err := filesystemFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_, err = fmt.Fprintf(w, "%d", sumOfSmallDirectories(root))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 7 of Advent of Code 2022.
func PartTwo(r io.Reader, w io.Writer) error {
	root, err := filesystemFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	spaceToFree := root.size - (70_000_000 - 30_000_000)
	if root.size < spaceToFree {
		return fmt.Errorf("cannot free %d disk space, currently use %d", spaceToFree, root.size)
	}

	dir := directoryToDelete(root, spaceToFree)

	_, err = fmt.Fprintf(w, "%d", dir.size)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type directory struct {
	name    string
	files   []file
	subdirs []directory
	size    int
}

type file struct {
	name string
	size int
}

func sumOfSmallDirectories(root directory) int {
	sum := 0

	walk(root, func(dir directory) {
		if dir.size < 100_000 {
			sum += dir.size
		}
	})

	return sum
}

func directoryToDelete(root directory, spaceToFree int) directory {
	toDelete := root

	walk(root, func(dir directory) {
		if dir.size > spaceToFree && dir.size < toDelete.size {
			toDelete = dir
		}
	})

	return toDelete
}

func walk(dir directory, f func(directory)) {
	f(dir)
	for _, sub := range dir.subdirs {
		walk(sub, f)
	}
}

func filesystemFromReader(r io.Reader) (directory, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return directory{}, err
	}

	var lineNum int

	var fill func(*directory) error
	fill = func(dir *directory) error {
		for lineNum < len(lines) {
			l := lines[lineNum]
			lineNum++

			if l == "$ cd .." {
				// Done with this directory.
				break
			}

			switch {
			case strings.HasPrefix(l, "$ cd "):
				// Moving into a subdirectory.
				name := strings.TrimPrefix(l, "$ cd ")
				for i := range dir.subdirs {
					if dir.subdirs[i].name == name {
						if err := fill(&dir.subdirs[i]); err != nil {
							return err
						}
					}
				}

			case l == "$ ls":
				// Nothing to parse, following lines are useful.

			case strings.HasPrefix(l, "dir "):
				// Found a subdirectory.
				name := strings.TrimPrefix(l, "dir ")
				dir.subdirs = append(dir.subdirs, directory{name: name})

			default:
				// Found a file.
				f, err := fileFromString(l)
				if err != nil {
					return err
				}
				dir.files = append(dir.files, f)
			}
		}

		// Compute the directory's size.
		for _, f := range dir.files {
			dir.size += f.size
		}
		for _, d := range dir.subdirs {
			dir.size += d.size
		}

		return nil
	}

	// Assuming that the first line of input is "$ cd /".
	root := directory{name: "/"}
	lineNum++

	if err := fill(&root); err != nil {
		return directory{}, err
	}

	return root, nil
}

func fileFromString(s string) (file, error) {
	parts := strings.SplitN(s, " ", 2)
	if len(parts) != 2 {
		return file{}, fmt.Errorf("wrong format")
	}

	size, err := strconv.Atoi(parts[0])
	if err != nil {
		return file{}, fmt.Errorf("%q is not a number", parts[0])
	}

	return file{name: parts[1], size: size}, nil
}
