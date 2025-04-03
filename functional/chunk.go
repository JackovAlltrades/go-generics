package functional

// Chunk splits a slice into smaller slices (chunks) of a specified size.
// The last chunk may contain fewer elements than the specified size if the
// total number of elements is not evenly divisible by the chunk size.
//
// Type Parameters:
//
//	T: The type of elements in the slice.
//
// Parameters:
//
//	input: The slice to split into chunks.
//	size:  The desired size of each chunk. Must be greater than 0.
//
// Returns:
//
//	[][]T: A new slice of slices, where each inner slice represents a chunk.
//	       Returns an empty slice of slices ([][]T{}) if the input slice is
//	       nil or empty, or if the specified size is less than or equal to 0.
//
// The original input slice is never modified. The returned inner slices are
// subslices of the original input slice's underlying array.
func Chunk[T any](input []T, size int) [][]T {
	inputLen := len(input)

	// Handle invalid size or empty/nil input (Guideline #3, #5)
	if size <= 0 || inputLen == 0 {
		return [][]T{} // Return empty slice of slices
	}

	// Calculate the number of chunks needed.
	// Ceiling division: (numerator + denominator - 1) / denominator
	numChunks := (inputLen + size - 1) / size

	// Preallocate the outer slice with the exact number of chunks (Guideline #5)
	result := make([][]T, 0, numChunks)

	// Iterate through the input slice, creating chunks
	for i := 0; i < inputLen; i += size {
		// Calculate the end index for the current chunk
		end := i + size
		// Ensure the end index does not exceed the slice bounds
		if end > inputLen {
			end = inputLen
		}
		// Append the subslice (chunk) to the result
		result = append(result, input[i:end])
	}

	return result
}
