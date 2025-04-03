package functional

// Reverse returns a new slice with the elements of the input slice in reverse order.
//
// Type Parameters:
//
//	T: The type of elements in the slice.
//
// Parameters:
//
//	input: The slice to reverse.
//
// Returns:
//
//	[]T: A new slice containing the elements of the input in reverse order.
//	     Returns an empty slice ([]T{}) if the input slice is nil or empty.
//
// The original input slice is never modified.
func Reverse[T any](input []T) []T {
	length := len(input)
	// Handle nil or empty input slice (Guideline #3, #5)
	if length == 0 {
		return []T{} // Return empty slice
	}

	// Preallocate result slice with the correct length and capacity (Guideline #5)
	result := make([]T, length)

	// Populate the result slice in reverse order
	for i, v := range input {
		// Place element v (from input index i) into the corresponding reverse position
		// in the result slice (index length - 1 - i).
		result[length-1-i] = v
	}

	return result
}

// ReverseInPlace reverses the elements of the input slice directly (in-place).
//
// Type Parameters:
//
//	T: The type of elements in the slice.
//
// Parameters:
//
//	input: The slice to reverse in-place. If the slice is nil or has
//	       fewer than 2 elements, the function does nothing.
//
// Returns:
//
//	None. The input slice itself is modified.
func ReverseInPlace[T any](input []T) {
	length := len(input)
	// No action needed for nil, empty, or single-element slices (Guideline #5)
	if length < 2 {
		return
	}

	// Swap elements from the start and end, moving inwards.
	// Iterate only up to the middle of the slice.
	for i := 0; i < length/2; i++ {
		// Calculate the index of the corresponding element from the end.
		j := length - 1 - i
		// Perform the swap using parallel assignment.
		input[i], input[j] = input[j], input[i]
	}
	// No return value needed as the slice is modified directly.
}
