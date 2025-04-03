package functional

// Flatten takes a slice of slices and returns a new slice containing all elements
// from the inner slices, concatenated in order.
//
// Type Parameters:
//
//	T: The type of elements in the inner slices (and the final result).
//
// Parameters:
//
//	input: The slice of slices ([][]T) to flatten.
//
// Returns:
//
//	[]T: A new slice containing all elements from the inner slices.
//	     Returns an empty slice ([]T{}) if the input slice is nil, empty,
//	     or contains only nil/empty inner slices.
//
// The original input slice and its inner slices are never modified.
func Flatten[T any](input [][]T) []T {
	// Handle nil or empty outer slice (Guideline #3, #5)
	if len(input) == 0 {
		return []T{} // Return empty slice
	}

	// Calculate total capacity needed for efficiency (Guideline #5)
	totalCapacity := 0
	for _, innerSlice := range input {
		totalCapacity += len(innerSlice) // len(nil slice) is 0, safe here
	}

	// Preallocate result slice with the calculated capacity
	result := make([]T, 0, totalCapacity)

	// Iterate through the outer slice
	for _, innerSlice := range input {
		// Append elements from the inner slice to the result
		// The '...' operator unpacks the innerSlice elements
		// append handles nil innerSlice correctly (appends nothing)
		result = append(result, innerSlice...)
	}

	return result
}
