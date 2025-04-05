package functional

// Intersection returns a new slice containing elements present in both s1 and s2.
// It requires the element type T to be comparable. The result contains unique elements.
// The order of elements in the result is not guaranteed.
//
// Args:
//
//	s1 ([]T): The first input slice.
//	s2 ([]T): The second input slice.
//
// Returns:
//
//	[]T: A slice containing the common unique elements. Returns an empty slice if no common elements or if inputs are nil/empty.
func Intersection[T comparable](s1, s2 []T) []T {
	if len(s1) == 0 || len(s2) == 0 {
		return []T{}
	}

	// Build map from the smaller slice for potentially better performance
	var mapSlice, iterateSlice []T
	if len(s1) < len(s2) {
		mapSlice = s1
		iterateSlice = s2
	} else {
		mapSlice = s2
		iterateSlice = s1
	}

	set := make(map[T]struct{}, len(mapSlice))
	for _, item := range mapSlice {
		set[item] = struct{}{}
	}

	intersectionMap := make(map[T]struct{}) // To store unique intersection results
	for _, item := range iterateSlice {
		if _, exists := set[item]; exists {
			intersectionMap[item] = struct{}{}
		}
	}

	if len(intersectionMap) == 0 {
		return []T{}
	}

	result := make([]T, 0, len(intersectionMap))
	for k := range intersectionMap {
		result = append(result, k)
	}
	return result
}

// Union returns a new slice containing unique elements from both s1 and s2.
// It requires the element type T to be comparable.
// The order of elements in the result is not guaranteed.
//
// Args:
//
//	s1 ([]T): The first input slice. Can be nil or empty.
//	s2 ([]T): The second input slice. Can be nil or empty.
//
// Returns:
//
//	[]T: A new slice containing the unique elements from both s1 and s2.
//	     Returns an empty slice ([]T{}) if both inputs are nil/empty.
func Union[T comparable](s1, s2 []T) []T {
	capacityHint := len(s1) + len(s2) // Over-estimation is okay for map capacity
	unionSet := make(map[T]struct{}, capacityHint)

	for _, v := range s1 {
		unionSet[v] = struct{}{}
	}
	for _, v := range s2 {
		unionSet[v] = struct{}{}
	}

	if len(unionSet) == 0 {
		return []T{}
	}

	result := make([]T, 0, len(unionSet))
	for k := range unionSet {
		result = append(result, k)
	}
	return result
}

// Difference returns a new slice containing unique elements present in s1 but not in s2 (s1 - s2).
// It requires the element type T to be comparable.
// The order of elements in the result is not guaranteed.
//
// Args:
//
//	s1 ([]T): The slice to subtract from.
//	s2 ([]T): The slice containing elements to remove.
//
// Returns:
//
//	[]T: A slice containing unique elements from s1 that are not in s2.
//	     Returns an empty slice if s1 is nil/empty or if all elements of s1 are also in s2.
func Difference[T comparable](s1, s2 []T) []T {
	if len(s1) == 0 {
		return []T{}
	}

	setB := make(map[T]struct{}, len(s2))
	for _, item := range s2 {
		setB[item] = struct{}{}
	}

	// Use a map to collect unique results from s1 that are not in setB
	resultSet := make(map[T]struct{})
	for _, item := range s1 {
		if _, existsInB := setB[item]; !existsInB {
			resultSet[item] = struct{}{} // Add to result if not in B
		}
	}

	if len(resultSet) == 0 {
		return []T{}
	}

	result := make([]T, 0, len(resultSet))
	for k := range resultSet {
		result = append(result, k)
	}
	return result
}

// Unique returns a new slice containing only the unique elements from the input slice,
// preserving the order of first appearance.
// It requires the element type T to be comparable.
//
// Args:
//
//	slice []T: The input slice, which may contain duplicates.
//
// Returns:
//
//	[]T: A new slice with duplicates removed. Returns an empty slice if the input is nil or empty.
func Unique[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return []T{}
	}

	seen := make(map[T]struct{}, len(slice)) // Optimized capacity guess
	result := make([]T, 0, len(slice))       // Capacity can be up to original length
	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
