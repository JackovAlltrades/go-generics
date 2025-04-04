package functional

// Contains checks if a slice contains a specific value.
//
// Type Parameters:
//
//	T: The type of elements in the slice, must be comparable.
//
// Parameters:
//
//	slice: The slice to search in. Can be nil or empty.
//	value: The value to search for.
//
// Returns:
//
//	true if the value is found in the slice, false otherwise.
//	For nil or empty slices, returns false.
func Contains[T comparable](slice []T, value T) bool {
	// Optimization: Check length outside the loop (Guideline #5)
	// if len(slice) == 0 { // This check is implicitly handled by the loop range
	// 	return false
	// }

	for _, item := range slice {
		if item == value {
			return true // Return early (Guideline #5)
		}
	}
	return false // Not found after checking all items
}
