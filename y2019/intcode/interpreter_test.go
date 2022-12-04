package intcode

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name string
		code []int
		want []int
	}{
		{
			name: "d02 example 1",
			code: []int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50},
			want: []int{3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50},
		},
		{
			name: "d02 example 2",
			code: []int{1, 0, 0, 0, 99},
			want: []int{2, 0, 0, 0, 99},
		},
		{
			name: "d02 example 3",
			code: []int{2, 3, 0, 3, 99},
			want: []int{2, 3, 0, 6, 99},
		},
		{
			name: "d02 example 4",
			code: []int{2, 4, 4, 5, 99, 0},
			want: []int{2, 4, 4, 5, 99, 9801},
		},
		{
			name: "d02 example 5",
			code: []int{1, 1, 1, 4, 99, 5, 6, 0, 99},
			want: []int{30, 1, 1, 4, 2, 5, 6, 0, 99},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(tt.code); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, tt.code); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
