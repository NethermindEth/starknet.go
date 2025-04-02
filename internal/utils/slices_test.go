package utils

import (
	"fmt"
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := Map(tc.input, func(i int) string { return fmt.Sprintf("%d", i) })
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := Filter(tc.input, func(i int) bool { return i%2 == 0 })
			assert.ElementsMatch(t, result, tc.expected)
		})
	}
}
