package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func ExampleIntsFromString() {
	str := "1 23 4 567 8 90"
	sep := " "

	ints, err := IntsFromString(str, sep)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ints)
	// Output: [1 23 4 567 8 90]
}

func BenchmarkIntsFromString(b *testing.B) {
	benchmarks := []struct {
		name     string
		dataFile string
	}{
		{
			name:     "10-small-integers",
			dataFile: "testdata/numbers/10-small-integers.txt",
		},
		{
			name:     "500-small-integers",
			dataFile: "testdata/numbers/500-small-integers.txt",
		},
		{
			name:     "25000-small-integers",
			dataFile: "testdata/numbers/25000-small-integers.txt",
		},
		{
			name:     "10-large-integers",
			dataFile: "testdata/numbers/10-large-integers.txt",
		},
		{
			name:     "500-large-integers",
			dataFile: "testdata/numbers/500-large-integers.txt",
		},
		{
			name:     "25000-large-integers",
			dataFile: "testdata/numbers/25000-large-integers.txt",
		},
	}

	for _, bench := range benchmarks {
		b.Run(bench.name, func(b *testing.B) {
			bytes, err := ioutil.ReadFile(bench.dataFile)
			if err != nil {
				b.Fatalf("could not read test data file: %v", err)
			}

			str := string(bytes)

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				_, _ = IntsFromString(str, " ")
			}
		})
	}
}
