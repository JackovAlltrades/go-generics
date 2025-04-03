// Package functional provides generic functional programming utilities.
package functional

// Unique returns a new slice containing only the unique elements from the
// input slice, preserving the order of the first appearance of each element.
// It requires the element type T to be comparable.
//
// Type Parameters:
//
//	T: The type of elements in the input slice. Must be comparable.
//
// Parameters:
//
//	input: The slice to process. Can be nil or empty.
//
// Returns:
//
//	[]T: A new slice containing the unique elements in their original order
//	     of first appearance. Returns nil if the input slice is nil. Returns
//	     an empty slice if the input slice is empty.
//	The original input slice is never modified.
func Unique[T comparable](input []T) []T {
	if input == nil {
		return nil
	}
	if len(input) == 0 {
		// Return a non-nil empty slice for non-nil empty input
		return []T{}
	}

	// Use a map to track elements encountered so far.
	// map[T]struct{} is commonly used for sets as struct{} uses zero bytes.
	seen := make(map[T]struct{}, len(input)) // Initialize with capacity hint

	// Result slice - starting with 0 capacity.
	result := make([]T, 0)

	for _, v := range input {
		// Check if the element has been seen before.
		if _, ok := seen[v]; !ok {
			// If not seen, add it to the seen map and the result slice.
			seen[v] = struct{}{} // Mark as seen
			result = append(result, v)
		}
	}

	return result
}
