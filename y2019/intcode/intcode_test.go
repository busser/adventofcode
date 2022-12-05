package intcode

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		code        []int
		inputs      []int
		wantState   []int
		wantOutputs []int
	}{
		{
			name:      "d02 example 1",
			code:      []int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50},
			wantState: []int{3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50},
		},
		{
			name:      "d02 example 2",
			code:      []int{1, 0, 0, 0, 99},
			wantState: []int{2, 0, 0, 0, 99},
		},
		{
			name:      "d02 example 3",
			code:      []int{2, 3, 0, 3, 99},
			wantState: []int{2, 3, 0, 6, 99},
		},
		{
			name:      "d02 example 4",
			code:      []int{2, 4, 4, 5, 99, 0},
			wantState: []int{2, 4, 4, 5, 99, 9801},
		},
		{
			name:      "d02 example 5",
			code:      []int{1, 1, 1, 4, 99, 5, 6, 0, 99},
			wantState: []int{30, 1, 1, 4, 2, 5, 6, 0, 99},
		},
		{
			name:        "d05 example 1",
			code:        []int{3, 0, 4, 0, 99},
			inputs:      []int{123},
			wantState:   []int{123, 0, 4, 0, 99},
			wantOutputs: []int{123},
		},
		{
			name:      "d05 example 2",
			code:      []int{1002, 4, 3, 4, 33},
			wantState: []int{1002, 4, 3, 4, 99},
		},
		{
			name:        "d05 example 3 true",
			code:        []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			inputs:      []int{8},
			wantOutputs: []int{1},
		},
		{
			name:        "d05 example 3 false",
			code:        []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			inputs:      []int{123},
			wantOutputs: []int{0},
		},
		{
			name:        "d05 example 4 true",
			code:        []int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8},
			inputs:      []int{5},
			wantOutputs: []int{1},
		},
		{
			name:        "d05 example 4 false",
			code:        []int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8},
			inputs:      []int{123},
			wantOutputs: []int{0},
		},
		{
			name:        "d05 example 5 true",
			code:        []int{3, 3, 1108, -1, 8, 3, 4, 3, 99},
			inputs:      []int{8},
			wantOutputs: []int{1},
		},
		{
			name:        "d05 example 5 false",
			code:        []int{3, 3, 1108, -1, 8, 3, 4, 3, 99},
			inputs:      []int{123},
			wantOutputs: []int{0},
		},
		{
			name:        "d05 example 6 true",
			code:        []int{3, 3, 1107, -1, 8, 3, 4, 3, 99},
			inputs:      []int{5},
			wantOutputs: []int{1},
		},
		{
			name:        "d05 example 6 false",
			code:        []int{3, 3, 1107, -1, 8, 3, 4, 3, 99},
			inputs:      []int{123},
			wantOutputs: []int{0},
		},
		{
			name:        "d05 example 7 true",
			code:        []int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9},
			inputs:      []int{123},
			wantOutputs: []int{1},
		},
		{
			name:        "d05 example 7 false",
			code:        []int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9},
			inputs:      []int{0},
			wantOutputs: []int{0},
		},
		{
			name:        "d05 example 8 true",
			code:        []int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1},
			inputs:      []int{123},
			wantOutputs: []int{1},
		},
		{
			name:        "d05 example 8 false",
			code:        []int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1},
			inputs:      []int{0},
			wantOutputs: []int{0},
		},
		{
			name: "d05 example 9 less",
			code: []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
				1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
				999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99},
			inputs:      []int{5},
			wantOutputs: []int{999},
		},
		{
			name: "d05 example 9 equal",
			code: []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
				1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
				999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99},
			inputs:      []int{8},
			wantOutputs: []int{1000},
		},
		{
			name: "d05 example 9 greater",
			code: []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
				1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
				999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99},
			inputs:      []int{123},
			wantOutputs: []int{1001},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualOutputs, err := Run(tt.code, tt.inputs)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.wantState != nil {
				if diff := cmp.Diff(tt.wantState, tt.code); diff != "" {
					t.Errorf("state mismatch (-want +got):\n%s", diff)
				}
			}

			if diff := cmp.Diff(tt.wantOutputs, actualOutputs); diff != "" {
				t.Errorf("outputs mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
