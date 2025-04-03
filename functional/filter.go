// Package functional provides generic functional programming utilities.
package functional

// Filter returns a new slice containing only the elements from the input slice
// that satisfy the predicate function.
//
// Type Parameters:
//
//	T: The type of elements in the input slice.
//
// Parameters:
//
//	input: The slice to filter. Can be nil or empty.
//	predicate: The function that determines if an element should be included.
//
// Returns:
//
//	A new slice containing only the elements that satisfy the predicate.
//	If input is nil or empty, returns an empty non-nil slice.
func Filter[T any](input []T, predicate func(T) bool) []T {
	// Simplified check - len(nil) is already 0 in Go
	if len(input) == 0 {
		return []T{}
	}

	// Create a result slice with initial capacity of 0
	// We don't know in advance how many elements will pass the filter
	result := make([]T, 0)

	for _, v := range input {
		if predicate(v) {
			result = append(result, v)
		}
	}

	return result
}
