// functional/chunk.go (or wherever your Chunk function is defined)

package functional

// Chunk divides a slice into smaller slices of a specified size.
// The last chunk may have fewer elements if the input slice's length
// is not evenly divisible by the size.
// Panics if size is not positive.
//
// Type Parameters:
//
//	T: The type of elements in the slice.
//
// Parameters:
//
//	slice: The input slice. Can be nil or empty.
//	size:  The desired size of each chunk. Must be positive.
//
// Returns:
//
//	[][]T: A new slice containing slices (chunks) of the original data.
//	       Returns an empty slice of slices ([][]T{}) if the input is nil/empty.
//
// The original input slice is never modified. The returned inner slices are
// subslices of the original input slice's underlying array.
func Chunk[T any](slice []T, size int) [][]T {
	// Panic if size is not positive
	if size <= 0 {
		panic("functional.Chunk: size must be positive")
	}

	inputLen := len(slice)

	// Handle nil or empty input slice after the size check
	if inputLen == 0 {
		return [][]T{}
	}

	// Calculate the number of chunks needed using ceiling division
	numChunks := (inputLen + size - 1) / size

	// Preallocate the outer slice with the exact number of chunks (Optimization)
	result := make([][]T, 0, numChunks)

	// Loop through the input slice, creating chunks
	for i := 0; i < inputLen; i += size {
		end := i + size
		// Ensure 'end' does not exceed the slice bounds for the last chunk
		if end > inputLen {
			end = inputLen
		}
		// Append the sub-slice (chunk) to the result
		result = append(result, slice[i:end])
	}

	return result
}
