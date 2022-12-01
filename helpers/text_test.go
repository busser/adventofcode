package helpers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func ExampleLinesFromReader() {
	f, err := os.Open("testdata/text/go-proverbs.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines, err := LinesFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of lines: %d\n", len(lines))
	fmt.Printf("The first one is: %q\n", lines[0])
	// Output:
	// Number of lines: 19
	// The first one is: "Don't communicate by sharing memory, share memory by communicating."
}

func BenchmarkLinesFromReader(b *testing.B) {
	benchmarks := []struct {
		name     string
		dataFile string
	}{
		{
			name:     "go-proverbs",
			dataFile: "testdata/text/go-proverbs.txt",
		},
	}

	for _, bench := range benchmarks {
		b.Run(bench.name, func(b *testing.B) {
			data, err := ioutil.ReadFile(bench.dataFile)
			if err != nil {
				b.Fatalf("could not read test data file: %v", err)
			}

			r := bytes.NewReader(data)

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				b.StopTimer()
				r.Reset(data)
				b.StartTimer()

				_, _ = LinesFromReader(r)
			}
		})
	}
}
