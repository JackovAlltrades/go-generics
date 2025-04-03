// Package functional provides generic functions for common functional
// programming patterns like Map, Filter, Reduce, and other slice/map utilities.
package functional

// Map applies a transformation function to each element in a slice and returns
// a new slice with the transformed values.
//
// Type Parameters:
//
//	T: The type of elements in the input slice.
//	U: The type of elements in the output slice.
//
// Parameters:
//
//	input: The slice to transform. Can be nil or empty.
//	mapFunc: The function to apply to each element.
//
// Returns:
//
//	A new slice containing the transformed elements.
//	If input is nil, returns nil.
//	If input is empty, returns an empty slice.
func Map[T, U any](input []T, mapFunc func(T) U) []U {
	// Handle nil input slice idiomatically.
	if input == nil {
		return nil
	}

	// Preallocate the result slice with the exact required capacity and length.
	// This is a key performance optimization, avoiding repeated allocations.
	result := make([]U, len(input))

	// Iterate using index to place results directly into the preallocated slice.
	for i, v := range input {
		// Fixed: using 'mapFunc' instead of 'f'
		result[i] = mapFunc(v)
	}

	return result
}
