// Package functional provides generic functions for common functional
// programming patterns like Map, Filter, Reduce, and other slice/map utilities.
package functional

// Reduce applies a function to each element in a slice, accumulating a single result.
//
// Type Parameters:
//
//	T: The type of elements in the input slice.
//	U: The type of the accumulated result.
//
// Parameters:
//
//	input: The slice to reduce. Can be nil or empty.
//	initial: The initial value for the accumulator.
//	reducer: The function that combines the accumulator with each element.
//
// Returns:
//
//	The final accumulated value.
//	If input is nil or empty, returns the initial value.
func Reduce[T, U any](input []T, initial U, reducer func(U, T) U) U {
	if len(input) == 0 {
		return initial
	}

	accumulator := initial
	for _, item := range input {
		accumulator = reducer(accumulator, item)
	}
	return accumulator
}
