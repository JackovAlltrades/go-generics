package functional

// Import necessary packages. 'cmp' is removed as it's no longer needed by Union.
// 'sort' might be used by other functions or internally, keep for now if needed elsewhere.
// Consider removing 'sort' if no function in this *file* uses it directly.
// Keep ONLY if Intersection/Difference/Unique implementations happen to use it internally

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
	if len(s1) == 0 || len(s2) == 0 {
		return []T{}
	}

	s2Set := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		s2Set[v] = struct{}{}
	}

	// Using min() requires Go 1.21+. If targeting older versions, implement min manually.
	// Or simply use len(s1) as upper bound.
	addedSetCapacity := len(s1)     // Default capacity
	if len(s2) < addedSetCapacity { // Estimate based on smaller potential intersection size
		addedSetCapacity = len(s2)
	}
	addedSet := make(map[T]struct{}, addedSetCapacity)
	result := make([]T, 0)

	for _, v := range s1 {
		_, existsInS2 := s2Set[v]
		_, alreadyAdded := addedSet[v]

		if existsInS2 && !alreadyAdded {
			result = append(result, v)
			addedSet[v] = struct{}{}
		}
	}
	return result
}

// Union returns a new slice containing the unique elements present in
// *either* of the input slices (s1 or s2).
// The order of elements in the returned slice is not guaranteed.
//
// Type Parameters:
//
//	T: The type of elements in the slices. Must be comparable. // CORRECTED constraint
//
// Parameters:
//
//	s1: The first input slice.
//	s2: The second input slice.
//
// Returns:
//
//	[]T: A new slice containing the unique elements from both s1 and s2. The order
//	     is not guaranteed. Returns an empty slice ([]T{}) if both inputs are nil/empty.
func Union[T comparable](s1, s2 []T) []T { // CORRECTED constraint to comparable
	// Determine initial capacity hint
	capacityHint := len(s1) + len(s2)
	if capacityHint == 0 && (s1 != nil || s2 != nil) { // Handle case where slices have 0 len but aren't nil
		// If either is non-nil but empty, capacity can be small for the map.
		// If both nil, map shouldn't be created anyway.
		// Heuristic: just use the sum, len(nil) is 0.
	}

	unionSet := make(map[T]struct{}, capacityHint) // Use hint

	for _, v := range s1 {
		unionSet[v] = struct{}{}
	}
	for _, v := range s2 {
		unionSet[v] = struct{}{}
	}

	if len(unionSet) == 0 {
		return []T{}
	}

	result := make([]T, 0, len(unionSet)) // Preallocate result based on actual unique count
	for k := range unionSet {
		result = append(result, k)
	}

	// No sorting here! Order depends on map iteration.
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
	if len(s1) == 0 {
		return []T{}
	}

	s2Set := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		s2Set[v] = struct{}{}
	}

	// Use len(s1) as capacity hint, as max difference size is len(s1)
	addedSet := make(map[T]struct{}, len(s1))
	result := make([]T, 0)

	for _, v := range s1 {
		_, existsInS2 := s2Set[v]
		_, alreadyAdded := addedSet[v]

		if !existsInS2 && !alreadyAdded {
			result = append(result, v)
			addedSet[v] = struct{}{}
		}
	}

	return result
}

// Unique returns a new slice containing only the unique elements from the
// input slice, preserving the order of first appearance.
//
// Type Parameters:
//
//	T: The type of elements in the slice. Must be comparable.
//
// Parameters:
//
//	slice: The input slice. Can be nil or empty.
//
// Returns:
//
//	[]T: A new slice containing the unique elements from the input,
//	     in the order of their first appearance.
//	     Returns an empty slice ([]T{}) if the input is nil or empty.
func Unique[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return []T{}
	}

	// Preallocate map based on input length (upper bound)
	seen := make(map[T]struct{}, len(slice))
	// Start result slice with 0 capacity - append will handle growth.
	// Alternatively, provide a hint like make([]T, 0, len(slice)/2) - requires testing.
	result := make([]T, 0)

	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}
