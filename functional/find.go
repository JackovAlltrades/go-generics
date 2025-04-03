// Package functional provides generic functions for common functional
// programming patterns like Map, Filter, Reduce, and other slice/map utilities.
package functional

// Find searches for an element in a slice that satisfies the predicate function.
//
// Type Parameters:
//
//	T: The type of elements in the input slice.
//
// Parameters:
//
//	input: The slice to search. Can be nil or empty.
//	predicate: The function that determines if an element matches.
//
// Returns:
//
//	A pointer to the first matching element and true if found.
//	Nil and false if no element matches or if input is nil or empty.
func Find[T any](input []T, predicate func(T) bool) (*T, bool) {
	if len(input) == 0 {
		return nil, false
	}

	for _, item := range input {
		if predicate(item) {
			return &item, true
		}
	}
	return nil, false
}
