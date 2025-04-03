package functional

import (
	"cmp" // Import the cmp package (Go 1.21+)
	"sort"
)

// Keys returns a slice containing all the keys from the input map.
// The order of keys in the returned slice is guaranteed to be sorted
// according to the standard Go sorting order for the key type K.
// This provides deterministic output, which is often useful.
//
// Type Parameters:
//
//	K: The type of the map keys. Must satisfy the cmp.Ordered constraint
//	   (e.g., numeric types, strings). This implies comparable.
//	V: The type of the map values (not constrained).
//
// Parameters:
//
//	inputMap: The map from which to extract keys. Can be nil or empty.
//
// Returns:
//
//	[]K: A new slice containing the keys from the map, sorted.
//	     Returns an empty slice ([]K{}) if the input map is nil or empty.
func Keys[K cmp.Ordered, V any](inputMap map[K]V) []K {
	if len(inputMap) == 0 {
		return []K{}
	}

	keys := make([]K, 0, len(inputMap))
	for k := range inputMap {
		keys = append(keys, k)
	}

	// Sort the keys using the now-safe '<' operator due to cmp.Ordered constraint.
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j] // This is now guaranteed to work for K
	})

	return keys
}

// Values returns a slice containing all the values from the input map.
// The order of values in the returned slice is arbitrary and corresponds
// to the iteration order of the map, which is not guaranteed.
//
// Type Parameters:
//
//	K: The type of the map keys. Must be comparable.
//	V: The type of the map values.
//
// Parameters:
//
//	inputMap: The map from which to extract values. Can be nil or empty.
//
// Returns:
//
//	[]V: A new slice containing the values from the map.
//	     Returns an empty slice ([]V{}) if the input map is nil or empty.
func Values[K comparable, V any](inputMap map[K]V) []V {
	if len(inputMap) == 0 {
		return []V{}
	}

	values := make([]V, 0, len(inputMap))
	for _, v := range inputMap {
		values = append(values, v)
	}

	return values
}

// MapToSlice transforms a map into a slice by applying a transformation function
// to each key-value pair. The order of elements in the resulting slice is
// arbitrary and depends on the map's iteration order.
//
// Type Parameters:
//
//	K: The type of the map keys. Must be comparable.
//	V: The type of the map values.
//	T: The type of the elements in the resulting slice.
//
// Parameters:
//
//	inputMap: The map to transform. Can be nil or empty.
//	transformFunc: A function that takes a key (K) and a value (V)
//	               and returns an element (T) for the output slice.
//
// Returns:
//
//	[]T: A new slice containing the results of applying transformFunc
//	     to each key-value pair. Returns an empty slice ([]T{}) if the
//	     input map is nil or empty.
func MapToSlice[ // Break before type params
	K comparable, // Type param 1
	V any, // Type param 2
	T any, // Type param 3
]( // Break before regular params
	inputMap map[K]V, // Param 1
	transformFunc func(key K, value V) T, // Param 2
) []T { // Closing paren and return type on new line
	if len(inputMap) == 0 {
		return []T{}
	}

	result := make([]T, 0, len(inputMap))

	for k, v := range inputMap {
		result = append(result, transformFunc(k, v))
	}

	return result
}
