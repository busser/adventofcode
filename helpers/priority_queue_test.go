package helpers

import "testing"

func TestPriorityQueue(t *testing.T) {
	testCases := []struct {
		name     string
		less     func(a, b int) bool
		input    []int
		expected []int
	}{
		{
			name:     "min_heap_ascending",
			less:     func(a, b int) bool { return a < b },
			input:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:     "min_heap_descending",
			less:     func(a, b int) bool { return a < b },
			input:    []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:     "min_heap_random_1",
			less:     func(a, b int) bool { return a < b },
			input:    []int{5, 3, 1, 2, 4, 7, 6, 9, 8, 10},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:     "min_heap_random_2",
			less:     func(a, b int) bool { return a < b },
			input:    []int{10, 8, 6, 4, 2, 1, 3, 5, 7, 9},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:     "max_heap_ascending",
			less:     func(a, b int) bool { return a > b },
			input:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name:     "max_heap_descending",
			less:     func(a, b int) bool { return a > b },
			input:    []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			expected: []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name:     "max_heap_random_1",
			less:     func(a, b int) bool { return a > b },
			input:    []int{5, 3, 1, 2, 4, 7, 6, 9, 8, 10},
			expected: []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name:     "max_heap_random_2",
			less:     func(a, b int) bool { return a > b },
			input:    []int{10, 8, 6, 4, 2, 1, 3, 5, 7, 9},
			expected: []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pq := NewPriorityQueue(tc.less)

			for _, i := range tc.input {
				pq.Push(i)
			}

			for _, e := range tc.expected {
				if actual := pq.Pop(); actual != e {
					t.Fatalf("expected %d, got %d", e, actual)
				}
			}
		})
	}
}
