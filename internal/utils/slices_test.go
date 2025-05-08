package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    [][]int
		expected []int
	}{
		{
			name:     "nil slice",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty slices",
			input:    [][]int{{}},
			expected: nil,
		},
		{
			name:     "single slice",
			input:    [][]int{{1, 2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "multiple slices",
			input:    [][]int{{1, 2}, {3, 4}, {5}},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "mixed empty and non-empty slices",
			input:    [][]int{{}, {1, 2}, {}, {3}, {}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "all empty slices",
			input:    [][]int{{}, {}, {}},
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := Flatten(tc.input...)
			assert.ElementsMatch(t, result, tc.expected)
		})
	}
}

func TestMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []int
		expected []string
	}{
		{
			name:     "nil slice",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: []string{},
		},
		{
			name:     "non-empty slice",
			input:    []int{1, 2, 3},
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "single element",
			input:    []int{42},
			expected: []string{"42"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := Map(tc.input, func(i int) string { return strconv.Itoa(i) }) //nolint:gocritic
			assert.ElementsMatch(t, result, tc.expected)
		})
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "nil slice",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "filter even numbers",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{2, 4},
		},
		{
			name:     "no matches",
			input:    []int{1, 3, 5, 7, 9},
			expected: []int{},
		},
		{
			name:     "all matches",
			input:    []int{2, 4, 6, 8, 10},
			expected: []int{2, 4, 6, 8, 10},
		},
		{
			name:     "single element match",
			input:    []int{1, 2, 3},
			expected: []int{2},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := Filter(tc.input, func(i int) bool { return i%2 == 0 })
			assert.ElementsMatch(t, result, tc.expected)
		})
	}
}

func TestAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []int
		pred     func(int) bool
		expected bool
	}{
		{
			name:     "nil slice",
			input:    nil,
			pred:     func(i int) bool { return true },
			expected: true,
		},
		{
			name:     "empty slice",
			input:    []int{},
			pred:     func(i int) bool { return true },
			expected: true,
		},
		{
			name:     "all even numbers",
			input:    []int{2, 4, 6, 8, 10},
			pred:     func(i int) bool { return i%2 == 0 },
			expected: true,
		},
		{
			name:     "mixed numbers",
			input:    []int{2, 3, 4, 6, 8},
			pred:     func(i int) bool { return i%2 == 0 },
			expected: false,
		},
		{
			name:     "all odd numbers",
			input:    []int{1, 3, 5, 7, 9},
			pred:     func(i int) bool { return i%2 == 1 },
			expected: true,
		},
		{
			name:     "single element true",
			input:    []int{2},
			pred:     func(i int) bool { return i%2 == 0 },
			expected: true,
		},
		{
			name:     "single element false",
			input:    []int{1},
			pred:     func(i int) bool { return i%2 == 0 },
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := All(tc.input, tc.pred)
			assert.Equal(t, result, tc.expected)
		})
	}
}
