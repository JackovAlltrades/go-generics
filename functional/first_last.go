// Package functional provides generic functional programming utilities.
package functional

// First returns a pointer to the first element in a slice and a boolean indicating success.
//
// Type Parameters:
//
//	T: The type of elements in the input slice.
//
// Parameters:
//
//	slice: The slice to get the first element from. Can be nil or empty.
//
// Returns:
//
//	*T: A pointer to the first element if the slice is not empty.
//	bool: true if an element was found, false if the slice is nil or empty.
func First[T any](slice []T) (*T, bool) {
	if len(slice) == 0 {
		return nil, false
	}
	return &slice[0], true
}

// Last returns a pointer to the last element in a slice and a boolean indicating success.
//
// Type Parameters:
//
//	T: The type of elements in the input slice.
//
// Parameters:
//
//	slice: The slice to get the last element from. Can be nil or empty.
//
// Returns:
//
//	*T: A pointer to the last element if the slice is not empty.
//	bool: true if an element was found, false if the slice is nil or empty.
func Last[T any](slice []T) (*T, bool) {
	if len(slice) == 0 {
		return nil, false
	}
	return &slice[len(slice)-1], true
}
