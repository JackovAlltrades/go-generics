package functional

// Find searches for an element in a slice that satisfies the predicate function.
// It returns a pointer to the *actual element within the slice's backing array*
// if found, allowing modification of the original slice element via the pointer.
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
//	*T: A pointer to the first matching element in the original slice, or nil if not found.
//	bool: true if an element was found, false if the slice is nil or empty or no element matches.
func Find[T any](input []T, predicate func(T) bool) (*T, bool) {
	// Use standard index loop to be absolutely sure we get pointer to slice element
	for i := 0; i < len(input); i++ {
		if predicate(input[i]) {
			return &input[i], true // Return address of slice element
		}
	}
	// If loop completes or slice is empty/nil, not found
	return nil, false
}
