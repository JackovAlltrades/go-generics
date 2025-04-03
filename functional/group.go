package functional

// GroupBy takes a slice and a classifier function, returning a map where keys
// are the results of applying the classifier function to each element, and
// values are slices containing the elements that produced that key.
//
// Type Parameters:
//
//	T: The type of elements in the input slice.
//	K: The type returned by the classifier function. Must be comparable
//	   to be used as map keys.
//
// Parameters:
//
//	input:      The slice to group.
//	classifier: A function that takes an element of type T and returns a key
//	            of type K.
//
// Returns:
//
//	map[K][]T: A new map where keys are the classification results (K) and
//	           values are slices ([]T) of elements that classified to that key.
//	           Returns an empty, non-nil map if the input slice is nil or empty.
//
// The order of elements within the value slices in the resulting map corresponds
// to their original order in the input slice. The iteration order over the map
// itself is not guaranteed.
func GroupBy[T any, K comparable](input []T, classifier func(element T) K) map[K][]T {
	// Create the result map. Always return non-nil map (Guideline #3)
	// Use capacity hint (Guideline #5) - difficult to predict accurately,
	// len(input)/N might be reasonable, but 0 is safe. Let's start with 0.
	// Alternative: result := make(map[K][]T, len(input)/some_avg_group_size)
	result := make(map[K][]T) // Start empty, grows as needed

	// Handle nil/empty input slice: loop won't run, empty map is returned. (Guideline #5)
	// if len(input) == 0 { return result } // This check is redundant

	// Iterate through the input slice
	for _, item := range input {
		// Use full parameter name 'classifier'. (Guideline #4)
		key := classifier(item)

		// Append the item to the slice associated with its key.
		// If the key doesn't exist yet, a new slice is implicitly created.
		result[key] = append(result[key], item)
	}

	return result
}
