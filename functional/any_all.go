// Package functional provides generic functions for common functional
// programming patterns like Map, Filter, Reduce, and other slice/map utilities.
package functional

// Any returns true if at least one element in the slice satisfies the predicate.
//
// Type Parameters:
//
//	T: The type of elements in the input slice.
//
// Parameters:
//
//	input: The slice to check. Can be nil or empty.
//	p: The predicate function to apply to each element.
//
// Returns:
//
//	true if any element satisfies the predicate, false otherwise.
//	For nil or empty slices, returns false.
func Any[T any](input []T, p func(T) bool) bool {
	if len(input) == 0 {
		return false
	}

	for _, v := range input {
		if p(v) {
			return true
		}
	}
	return false
}

// All returns true if all elements in the slice satisfy the predicate.
//
// Type Parameters:
//
//	T: The type of elements in the input slice.
//
// Parameters:
//
//	input: The slice to check. Can be nil or empty.
//	p: The predicate function to apply to each element.
//
// Returns:
//
//	true if all elements satisfy the predicate, false otherwise.
//	For nil or empty slices, returns true (vacuously true).
func All[T any](input []T, p func(T) bool) bool {
	if len(input) == 0 {
		return true
	}

	for _, v := range input {
		if !p(v) {
			return false
		}
	}
	return true
}
