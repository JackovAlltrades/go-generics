package functional

// Import 'sort' only if needed by other functions in this file.
// If Keys was the only user, 'sort' can be removed.
// import "sort"
// Keep cmp if other functions use it, otherwise remove.

// Keys extracts the keys from a map into a slice.
// The order of keys in the returned slice is not guaranteed.
//
// Type Parameters:
//
//	K: The type of the map keys. Must be comparable. // CORRECTED Constraint
//	V: The type of the map values (any).
//
// Parameters:
//
//	inputMap: The map from which to extract keys. Can be nil.
//
// Returns:
//
//	[]K: A slice containing the keys from the map. Returns an empty slice if
//	     the input map is nil or empty. Order is not guaranteed.
func Keys[K comparable, V any](inputMap map[K]V) []K { // CORRECTED Constraint
	// Handle nil or empty map efficiently. Check length *after* nil check for safety,
	// though len(nil map) is 0.
	if inputMap == nil {
		return []K{}
	}
	mapLen := len(inputMap)
	if mapLen == 0 {
		return []K{}
	}

	// Preallocate result slice with the correct capacity.
	keys := make([]K, 0, mapLen)

	// Iterate over the map and append keys.
	for k := range inputMap {
		keys = append(keys, k)
	}

	// --- REMOVED SORTING STEP ---
	// sort.Slice(keys, func(i, j int) bool {
	// 	return keys[i] < keys[j] // This requires K cmp.Ordered
	// })

	return keys
}

// Values extracts the values from a map into a slice.
// The order of values in the returned slice corresponds to the iteration order
// of the map, which is not guaranteed.
//
// Type Parameters:
//
//	K: The type of the map keys (must be comparable).
//	V: The type of the map values.
//
// Parameters:
//
//	inputMap: The map from which to extract values. Can be nil.
//
// Returns:
//
//	[]V: A slice containing the values from the map. Returns an empty slice if
//	     the input map is nil or empty. Order is not guaranteed.
func Values[K comparable, V any](inputMap map[K]V) []V {
	if inputMap == nil {
		return []V{}
	}
	mapLen := len(inputMap)
	if mapLen == 0 {
		return []V{}
	}

	// Preallocate result slice.
	values := make([]V, 0, mapLen)

	// Iterate over the map and append values.
	for _, v := range inputMap {
		values = append(values, v)
	}

	return values
}

// MapToSlice transforms a map into a slice by applying a function to each
// key-value pair. The order of elements in the resulting slice corresponds to the
// map's iteration order, which is not guaranteed.
//
// Type Parameters:
//
//	K: The type of the map keys (must be comparable).
//	V: The type of the map values.
//	R: The type of the elements in the resulting slice.
//
// Parameters:
//
//	inputMap: The map to process. Can be nil.
//	fn:       A function that takes a key (K) and a value (V) and returns a
//	          result (R) to be included in the output slice.
//
// Returns:
//
//	[]R: A slice containing the results of applying fn to each key-value pair.
//	     Returns an empty slice if the input map is nil or empty. Order is not guaranteed.
func MapToSlice[K comparable, V, R any](inputMap map[K]V, fn func(k K, v V) R) []R {
	if inputMap == nil {
		return []R{}
	}
	mapLen := len(inputMap)
	if mapLen == 0 {
		return []R{}
	}

	// Preallocate result slice.
	result := make([]R, 0, mapLen)

	// Iterate over the map, apply the function, and append the result.
	for k, v := range inputMap {
		result = append(result, fn(k, v))
	}

	return result
}
