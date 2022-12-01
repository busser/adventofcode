package scaffolding

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var (
	//go:embed templates/solution.go.tmpl
	solutionTemplate string
	//go:embed templates/solution_test.go.tmpl
	solutionTestTemplate string
)

// A Generator creates a directory with all contents required to kickstart
// a solution to a puzzle in the Advent of Code calendar.
type Generator struct {
	// The day and year to generate code for.
	day, year int

	// The directory where all Advent of Code solutions are stored.
	workdir string

	// Session cookie for adventofcode.com.
	cookie string

	// Whether to overwrite existing files.
	overwrite bool

	// Path to scaffolded directory.
	packageDir string
}

// NewGenerator builds a generator for the given date and author. If overwrite
// is true, the generator will overwrite existing files.
func NewGenerator(day, year int, author, workdir, cookie string, overwrite bool) (*Generator, error) {
	gen := &Generator{
		day:       day,
		year:      year,
		workdir:   workdir,
		cookie:    cookie,
		overwrite: overwrite,
	}

	if err := gen.Initialize(); err != nil {
		return nil, fmt.Errorf("failed initialization: %w", err)
	}

	return gen, nil
}

// Initialize validates gen's parameters and pre-computes useful values.
func (gen *Generator) Initialize() error {
	if gen.day <= 0 || gen.day > 25 {
		return fmt.Errorf("invalid day: %d", gen.day)
	}
	if gen.year <= 0 {
		return fmt.Errorf("invalid year: %d", gen.year)
	}
	if gen.workdir == "" {
		return errors.New("working directory unknown")
	}

	gen.setPackageDir()

	return nil
}

// Run uses gen to build scaffolding for an Advent of Code solution.
func (gen *Generator) Run() error {
	if err := gen.CreatePackage(); err != nil {
		return fmt.Errorf("creating package: %w", err)
	}
	if err := gen.WriteCode(); err != nil {
		return fmt.Errorf("writing code: %w", err)
	}
	if err := gen.DownloadInput(); err != nil {
		return fmt.Errorf("downloading input: %w", err)
	}
	return nil
}

// CreatePackage creates a directory/package to put scaffolding into.
func (gen *Generator) CreatePackage() error {
	err := os.MkdirAll(gen.packageDir, 0755)
	if err != nil {
		return fmt.Errorf("creating directory %q: %w", gen.packageDir, err)
	}

	fmt.Printf("🏗️  Scaffolding package: %s\n", gen.packageDir)
	return nil
}

// WriteCode builds Go scaffolding for implementing, testing, and benchmarking
// solutions to Advent of Code problems.
func (gen *Generator) WriteCode() error {
	if err := gen.renderTemplateIntoFile(solutionTemplate, "solution.go"); err != nil {
		return fmt.Errorf("creating %q: %w", "solution.go", err)
	}
	if err := gen.renderTemplateIntoFile(solutionTestTemplate, "solution_test.go"); err != nil {
		return fmt.Errorf("creating %q: %w", "solution_test.go", err)
	}
	return nil
}

// DownloadInput fetches the Advent of Code's daily input and writes it to a
// testdata directory.
func (gen *Generator) DownloadInput() error {
	path := filepath.Join(gen.packageDir, "testdata", "input.txt")
	if fileExists(path) && !gen.overwrite {
		fmt.Println("  👉 Skipping input download; file already exists.")
		return nil
	}

	if gen.cookie == "" {
		fmt.Println("  👉 Skipping input download; no session cookie provided.")
		return nil
	}

	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", gen.year, gen.day)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("preparing GET request to %q: %w", url, err)
	}

	req.AddCookie(&http.Cookie{Name: "session", Value: gen.cookie})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending GET request to %q: %w", url, err)
	}
	defer resp.Body.Close()

	input, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response from %q: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("adventofcode.com responded with %d: %s", resp.StatusCode, input)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory %q: %w", filepath.Dir(path), err)
	}

	err = os.WriteFile(path, input, 0644)
	if err != nil {
		return fmt.Errorf("writing input to file %q: %w", path, err)
	}

	fmt.Printf("  👉 Downloaded input.\n")

	return nil
}

func (gen *Generator) setPackageDir() {
	gen.packageDir = filepath.Join(
		gen.workdir,
		fmt.Sprintf("y%04d", gen.year),
		fmt.Sprintf("d%02d", gen.day),
	)
}

func (gen *Generator) renderTemplateIntoFile(templateText, filename string) error {
	path := filepath.Join(gen.packageDir, filename)

	if fileExists(path) && !gen.overwrite {
		fmt.Printf("  👉 Skipping existing file %s.\n", filename)
		return nil
	}

	tmpl, err := template.New("aoc").Parse(templateText)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating or opening file %q: %w", path, err)
	}
	defer f.Close()

	data := struct {
		Day, Year   int
		PackageName string
	}{
		Day:  gen.day,
		Year: gen.year,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	fmt.Printf("  👉 Scaffolded %s.\n", filename)

	return nil
}

// fileExists checks whether filename exists.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		fmt.Println(filename)
		panic(err)
	}
	return !info.IsDir()
}
