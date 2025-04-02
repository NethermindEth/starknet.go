package utils

import (
	"fmt"
	"reflect"
	"testing"
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
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := Flatten(tc.input...)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Flatten(%v) = %v, want %v", tc.input, result, tc.expected)
			}
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
			name:     "empty slice",
			input:    []int{},
			expected: []string{},
		},
		{
			name:     "non-empty slice",
			input:    []int{1, 2, 3},
			expected: []string{"1", "2", "3"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := Map(tc.input, func(i int) string { return fmt.Sprintf("%d", i) })
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Map(%v) = %v, want %v", tc.input, result, tc.expected)
			}
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
			name:     "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "filter even numbers",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{2, 4},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := Filter(tc.input, func(i int) bool { return i%2 == 0 })
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Filter(%v) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}
