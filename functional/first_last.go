package functional

// First returns a pointer to the first element of a slice.
// It returns a nil pointer and false if the slice is nil or empty.
//
// Args:
//
//	slice []T: The input slice.
//
// Returns:
//
//	 *T: A pointer to the first element, or nil.
//		bool: True if an element was found, false otherwise.
func First[T any](slice []T) (*T, bool) {
	if len(slice) == 0 {
		// Explicitly return typed nil pointer for consistency if needed,
		// though 'nil' often suffices. Return false for 'ok'.
		return nil, false
	}
	// Return pointer to the first element and true for 'ok'.
	return &slice[0], true
}

// Last returns a pointer to the last element of a slice.
// It returns a nil pointer and false if the slice is nil or empty.
//
// Args:
//
//	slice []T: The input slice.
//
// Returns:
//
//	 *T: A pointer to the last element, or nil.
//		bool: True if an element was found, false otherwise.
func Last[T any](slice []T) (*T, bool) {
	if len(slice) == 0 {
		return nil, false
	}
	// Return pointer to the last element and true for 'ok'.
	return &slice[len(slice)-1], true
}
