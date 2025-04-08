package utils

import "slices"

// Flatten flattens a slice of slices into a single slice
func Flatten[T any](sl ...[]T) []T {
	var result []T
	for _, slice := range sl {
		result = append(result, slice...)
	}

	return result
}

// Map maps a slice of type T1 to a slice of type T2 using the given function
func Map[T1, T2 any](slice []T1, f func(T1) T2) []T2 {
	if slice == nil {
		return nil
	}

	result := make([]T2, len(slice))
	for i, e := range slice {
		result[i] = f(e)
	}

	return result
}

// Filter filters a slice of type T using the given predicate, returning a new slice with the elements that match the predicate
func Filter[T any](slice []T, f func(T) bool) []T {
	if slice == nil {
		return nil
	}
	if len(slice) == 0 {
		return slice
	}

	result := make([]T, 0)
	for _, e := range slice {
		if f(e) {
			result = append(result, e)
		}
	}

	return result
}

// All returns true if all elements match the given predicate
func All[T any](slice []T, f func(T) bool) bool {
	return slices.IndexFunc(slice, func(e T) bool { return !f(e) }) == -1
}
