package functional

// Added for ExampleFilterErr
// Added for ExampleFilterErr

// MapErr applies a function mapFunc, which can return an error, to each element
// of the input slice. If mapFunc returns an error for any element, MapErr stops
// processing immediately and returns the slice of successfully processed elements
// up to that point, along with the encountered error (fail-fast strategy).
//
// Type Parameters:
//
//	T: The type of elements in the input slice.
//	U: The type of elements in the successful output slice.
//
// Parameters:
//
//	input:   The slice to iterate over. Can be nil or empty.
//	mapFunc: The function to apply to each element. It takes an element of type T
//	         and returns a result of type U and an error.
//
// Returns:
//
//	[]U:   A new slice containing the results of successfully applying mapFunc
//	       up to the point an error occurred. If no error occurs, it contains
//	       the results for all elements. Returns an empty slice ([]U{}) if the
//	       input is nil/empty.
//	error: The first non-nil error returned by mapFunc, or nil if all elements
//	       were processed successfully.
//
// The original input slice is never modified.
func MapErr[T, U any](input []T, mapFunc func(element T) (U, error)) ([]U, error) {
	// Handle nil/empty input (Guideline #3, #5)
	if len(input) == 0 {
		return []U{}, nil
	}

	// Preallocate result slice with capacity hint (Guideline #5)
	// Length starts at 0, as we append successful results.
	result := make([]U, 0, len(input))

	// Iterate through the input slice
	for _, item := range input {
		// Use full parameter name 'mapFunc'. (Guideline #4)
		mappedValue, err := mapFunc(item)
		if err != nil {
			// Fail-fast: return results processed so far and the error
			return result, err
		}
		// Append successful result
		result = append(result, mappedValue)
	}

	// If loop completes without error
	return result, nil
}

// FilterErr filters a slice based on a predicate function that can return an error.
// It iterates through the input slice, applying the predicate function to each element.
// If the predicate returns `true` and no error, the element is included in the result.
// If the predicate returns an error for any element, FilterErr stops processing
// immediately and returns the slice of elements successfully filtered up to that
// point, along with the encountered error (fail-fast strategy).
//
// Type Parameters:
//
//	T: The type of elements in the slice.
//
// Parameters:
//
//	input:     The slice to filter. Can be nil or empty.
//	predicate: The function to apply to each element. It takes an element of type T
//	           and returns a boolean indicating inclusion and an error.
//
// Returns:
//
//	[]T:   A new slice containing the elements for which the predicate successfully
//	       returned true up to the point an error occurred. If no error occurs,
//	       it contains all elements for which the predicate returned true.
//	       Returns an empty slice ([]T{}) if the input is nil/empty.
//	error: The first non-nil error returned by the predicate, or nil if all elements
//	       were processed successfully without error.
//
// The original input slice is never modified. The order of elements is preserved.
func FilterErr[T any](input []T, predicate func(element T) (bool, error)) ([]T, error) {
	// Handle nil/empty input (Guideline #3, #5)
	if len(input) == 0 {
		return []T{}, nil
	}

	// Preallocate result slice. Capacity is hard to estimate, start with 0.
	// Alternative: make([]T, 0, len(input)/2) as a guess. Let's use 0.
	result := make([]T, 0) // Guideline #5 (minimal preallocation)

	// Iterate through the input slice
	for _, item := range input {
		// Use full parameter name 'predicate'. (Guideline #4)
		include, err := predicate(item)
		if err != nil {
			// Fail-fast: return results filtered so far and the error
			return result, err
		}
		// Append if predicate returned true (and no error)
		if include {
			result = append(result, item)
		}
	}

	// If loop completes without error
	return result, nil
}

// ReduceErr applies a reducer function to the elements of a slice, accumulating
// a result. The reducer function can return an error. If the reducer returns an
// error for any element, ReduceErr stops processing immediately and returns
// the accumulated value up to that point, along with the encountered error
// (fail-fast strategy).
//
// Type Parameters:
//
//	T: The type of elements in the input slice.
//	U: The type of the accumulator and the final result.
//
// Parameters:
//
//	input:    The slice to iterate over. Can be nil or empty.
//	initial:  The initial value of the accumulator.
//	reducer:  The function to apply to each element. It takes the current
//	          accumulator value (U) and the current element (T), and returns
//	          the next accumulator value (U) and an error.
//
// Returns:
//
//	U:     The final accumulated value. If an error occurred, this is the value
//	       accumulated *before* the error. If the input slice is nil/empty,
//	       this is the initial value.
//	error: The first non-nil error returned by the reducer, or nil if all elements
//	       were processed successfully without error.
//
// The original input slice is never modified.
func ReduceErr[T, U any](input []T, initial U, reducer func(acc U, element T) (U, error)) (U, error) {
	// Initialize accumulator with the provided initial value
	accumulator := initial // Guideline #4 (clear naming)

	// Handle nil/empty input - simply return the initial value and nil error (Guideline #3, #5)
	if len(input) == 0 {
		return accumulator, nil
	}

	// Iterate through the input slice
	for _, item := range input {
		// Apply the reducer function
		nextAccumulator, err := reducer(accumulator, item)
		if err != nil {
			// Fail-fast: return the *current* accumulator value and the error
			return accumulator, err
		}
		// Update accumulator only on success
		accumulator = nextAccumulator
	}

	// If loop completes without error, return the final accumulator and nil error
	return accumulator, nil
}
