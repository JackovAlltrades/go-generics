package functional

import (
	"cmp" // Import cmp (Go 1.21+)
	"sort"
)

// Intersection returns a new slice containing elements that are present
// in *both* input slices (s1 and s2). The order of elements in the result
// corresponds to their first appearance in the first slice (s1).
// Duplicates within each input slice are effectively ignored; the result contains
// each unique common element only once.
//
// Type Parameters:
//
//	T: The type of elements in the slices. Must be comparable.
//
// Parameters:
//
//	s1: The first input slice.
//	s2: The second input slice.
//
// Returns:
//
//	[]T: A new slice containing the unique elements present in both s1 and s2,
//	     in the order of their first appearance in s1.
//	     Returns an empty slice ([]T{}) if either input is nil/empty, or if
//	     there is no intersection.
func Intersection[T comparable](s1, s2 []T) []T {
	// Handle nil/empty inputs efficiently (Guideline #3, #5)
	if len(s1) == 0 || len(s2) == 0 {
		return []T{}
	}

	// Create a set from the second slice for efficient lookup
	// Use capacity hint (Guideline #5)
	s2Set := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		s2Set[v] = struct{}{}
	}

	// Create a set to track elements added to the result, ensuring uniqueness
	// and preserving s1 order. Use min() from Go 1.21+.
	addedSet := make(map[T]struct{}, min(len(s1), len(s2)))

	// Result slice. Capacity is hard to predict, start with 0.
	result := make([]T, 0)

	// Iterate through the first slice (s1)
	for _, v := range s1 {
		// Check if element exists in s2 AND hasn't been added to result yet
		_, existsInS2 := s2Set[v]
		_, alreadyAdded := addedSet[v]

		if existsInS2 && !alreadyAdded {
			result = append(result, v)
			addedSet[v] = struct{}{} // Mark as added
		}
	}

	return result
}

// Union returns a new slice containing the unique elements present in
// *either* of the input slices (s1 or s2).
// The elements in the returned slice are sorted according to the standard Go
// sorting order for the element type T.
//
// Type Parameters:
//
//	T: The type of elements in the slices. Must satisfy the cmp.Ordered
//	   constraint (e.g., numeric types, strings) for sorting. This
//	   implies comparable.
//
// Parameters:
//
//	s1: The first input slice.
//	s2: The second input slice.
//
// Returns:
//
//	[]T: A new slice containing the unique elements from both s1 and s2, sorted.
//	     Returns an empty slice ([]T{}) if both inputs are nil/empty.
func Union[T cmp.Ordered](s1, s2 []T) []T {
	// Determine initial capacity hint (sum of lengths is upper bound) (Guideline #5)
	capacityHint := len(s1) + len(s2) // len(nil) is 0, safe

	// Use a map to store unique elements encountered.
	unionSet := make(map[T]struct{}, capacityHint) // Use hint

	// Add elements from the first slice
	for _, v := range s1 {
		unionSet[v] = struct{}{}
	}

	// Add elements from the second slice (duplicates are ignored by map)
	for _, v := range s2 {
		unionSet[v] = struct{}{}
	}

	// Handle case where both inputs were empty/nil resulting in empty set
	if len(unionSet) == 0 {
		return []T{} // Return empty slice (Guideline #3)
	}

	// Extract unique elements into a slice
	result := make([]T, 0, len(unionSet))
	for k := range unionSet {
		result = append(result, k)
	}

	// Sort the result slice for deterministic output (requires cmp.Ordered)
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j] // Safe due to cmp.Ordered constraint
	})

	return result
}

// Difference returns a new slice containing the unique elements that are present
// in the first slice (s1) but *not* present in the second slice (s2).
// The order of elements in the result corresponds to their first appearance in s1.
//
// Type Parameters:
//
//	T: The type of elements in the slices. Must be comparable.
//
// Parameters:
//
//	s1: The slice from which to subtract elements.
//	s2: The slice containing elements to subtract.
//
// Returns:
//
//	[]T: A new slice containing the unique elements from s1 that are not in s2,
//	     in the order of their first appearance in s1.
//	     Returns an empty slice ([]T{}) if s1 is nil/empty. Returns a slice
//	     containing the unique elements of s1 (ordered by first appearance) if s2 is nil/empty.
func Difference[T comparable](s1, s2 []T) []T {
	// If s1 is empty/nil, the difference is always empty
	if len(s1) == 0 {
		return []T{} // Guideline #3
	}

	// Create a set from the second slice for efficient lookup
	// Handle nil s2 gracefully - the set will be empty.
	s2Set := make(map[T]struct{}, len(s2)) // Capacity hint (Guideline #5)
	for _, v := range s2 {
		s2Set[v] = struct{}{}
	}

	// Create a set to track elements added to the result (handles duplicates in s1)
	addedSet := make(map[T]struct{}, len(s1)) // Capacity hint (Guideline #5)

	// Result slice.
	result := make([]T, 0) // Start with 0 capacity

	// Iterate through the first slice (s1)
	for _, v := range s1 {
		// Check if element exists in s2 AND hasn't been added to result yet
		_, existsInS2 := s2Set[v]
		_, alreadyAdded := addedSet[v]

		// Only add if it's NOT in s2 AND NOT already added from s1
		if !existsInS2 && !alreadyAdded {
			result = append(result, v)
			addedSet[v] = struct{}{} // Mark as added
		}
	}

	return result
}
