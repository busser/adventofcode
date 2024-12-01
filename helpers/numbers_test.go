package helpers

import (
	"fmt"
	"os"
	"testing"

	"github.com/magiconair/properties/assert"
)

func ExampleIntsFromString() {
	fmt.Println(IntsFromString("1 23 4 567 8 90"))
	fmt.Println(IntsFromString("1,23,4,567,8,90"))
	fmt.Println(IntsFromString("1 23,4,   567+8//90"))
	fmt.Println(IntsFromString("1234567890"))
	fmt.Println(IntsFromString("hello world"))
	fmt.Println(IntsFromString("-1 23 -4 567 -890"))
	fmt.Println(IntsFromString("-1 23-4 567-890"))
	// Output:
	// [1 23 4 567 8 90]
	// [1 23 4 567 8 90]
	// [1 23 4 567 8 90]
	// [1234567890]
	// []
	// [-1 23 -4 567 -890]
	// [-1 23 4 567 890]
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
			bytes, err := os.ReadFile(bench.dataFile)
			if err != nil {
				b.Fatalf("could not read test data file: %v", err)
			}

			str := string(bytes)

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				_ = IntsFromString(str)
			}
		})
	}
}

func TestSplitStringIntoIntStrings(t *testing.T) {
	testCases := []struct {
		str  string
		want []string
	}{
		{
			str:  "1 23 4 567 8 90",
			want: []string{"1", "23", "4", "567", "8", "90"},
		},
		{
			str:  "1,23,4,567,8,90",
			want: []string{"1", "23", "4", "567", "8", "90"},
		},
		{
			str:  "1 23,4,   567+8//90",
			want: []string{"1", "23", "4", "567", "8", "90"},
		},
		{
			str:  "1234567890",
			want: []string{"1234567890"},
		},
		{
			str:  "hello world",
			want: nil,
		},
		{
			str:  "-1 23 -4 567 -890",
			want: []string{"-1", "23", "-4", "567", "-890"},
		},
		{
			str:  "-1 23-4 567-890",
			want: []string{"-1", "23", "4", "567", "890"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.str, func(t *testing.T) {
			got := splitStringIntoIntStrings(tc.str)
			assert.Equal(t, got, tc.want)
		})
	}
}
