package functional_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/JackovAlltrades/go-generics/functional" // Adjust import path
)

// --- Test Chunk ---
func TestChunk(t *testing.T) {
	testCases := []struct {
		name        string
		input       any // Use any for type flexibility
		size        int
		want        any  // Expected value (use any)
		expectPanic bool // Optional: Flag if panic is expected
	}{
		{
			name:  "Ints_EvenSplit",
			input: []int{1, 2, 3, 4, 5, 6},
			size:  3,
			want:  [][]int{{1, 2, 3}, {4, 5, 6}},
		},
		{
			name:  "Ints_UnevenSplit",
			input: []int{1, 2, 3, 4, 5, 6, 7},
			size:  3,
			want:  [][]int{{1, 2, 3}, {4, 5, 6}, {7}},
		},
		{
			name:  "Strings_UnevenSplit",
			input: []string{"a", "b", "c", "d", "e"},
			size:  2,
			want:  [][]string{{"a", "b"}, {"c", "d"}, {"e"}},
		},
		{
			name:  "Size_One",
			input: []int{1, 2, 3},
			size:  1,
			want:  [][]int{{1}, {2}, {3}},
		},
		{
			name:  "Size_LargerThanSlice",
			input: []int{1, 2, 3},
			size:  5,
			want:  [][]int{{1, 2, 3}},
		},
		{
			name:  "Size_EqualToSlice",
			input: []int{1, 2, 3},
			size:  3,
			want:  [][]int{{1, 2, 3}},
		},
		{
			name:  "EmptyInput",
			input: []int{},
			size:  3,
			want:  [][]int{}, // Expect empty slice of slices
		},
		{
			name:  "NilInput",
			input: ([]int)(nil),
			size:  3,
			want:  [][]int{}, // Expect empty slice of slices
		},
		// --- Edge cases for size ---
		{
			name:        "Size_Zero",
			input:       []int{1, 2, 3},
			size:        0,
			want:        nil,  // Or whatever your implementation defines, often panic or empty
			expectPanic: true, // Assuming Chunk panics or errors on non-positive size
		},
		{
			name:        "Size_Negative",
			input:       []int{1, 2, 3},
			size:        -1,
			want:        nil,  // As above
			expectPanic: true, // Assuming Chunk panics or errors on non-positive size
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.expectPanic {
					if r == nil {
						t.Errorf("Chunk() did not panic for size %d, but expected panic", tc.size)
					}
				} else {
					if r != nil {
						t.Errorf("Chunk() panicked unexpectedly for size %d: %v", tc.size, r)
					}
				}
			}()

			var got any
			var wantTyped any

			switch input := tc.input.(type) {
			case []int:
				got = functional.Chunk(input, tc.size)
				// Ensure want is also typed for comparison
				if w, ok := tc.want.([][]int); ok {
					wantTyped = w
				} else if tc.want == nil && tc.expectPanic {
					wantTyped = nil // Explicitly nil for panic case comparison if needed
				} else if tc.want == nil && !tc.expectPanic { // Check if want should be empty [][]int
					if w, ok := got.([][]int); ok && len(w) == 0 {
						wantTyped = got // If got is [][]int{}, want is effectively [][]int{}
					} else {
						t.Fatalf("Want type mismatch for []int case: expected [][]int, got %T", tc.want)
					}
				} else {
					t.Fatalf("Want type mismatch for []int case: expected [][]int, got %T", tc.want)
				}
			case []string:
				got = functional.Chunk(input, tc.size)
				if w, ok := tc.want.([][]string); ok {
					wantTyped = w
				} else if tc.want == nil && tc.expectPanic {
					wantTyped = nil
				} else if tc.want == nil && !tc.expectPanic {
					if w, ok := got.([][]string); ok && len(w) == 0 {
						wantTyped = got
					} else {
						t.Fatalf("Want type mismatch for []string case: expected [][]string, got %T", tc.want)
					}
				} else {
					t.Fatalf("Want type mismatch for []string case: expected [][]string, got %T", tc.want)
				}
			case nil: // Handle nil input specifically
				// What type should Chunk[T] return for nil input? Assuming empty slice of slices of T.
				// Need to infer T somehow if possible, or test for a specific T.
				// Let's test with int as a default example for nil input
				got = functional.Chunk[int](nil, tc.size) // Explicit type T=int
				if w, ok := tc.want.([][]int); ok {       // Check tc.want is correct type
					wantTyped = w
				} else {
					t.Fatalf("Want type mismatch for nil input int case: expected [][]int, got %T", tc.want)
				}
			default:
				v := reflect.ValueOf(tc.input)
				if v.Kind() == reflect.Slice && v.Len() == 0 {
					// Test with int as a default example for empty slice
					got = functional.Chunk[int]([]int{}, tc.size)
					if w, ok := tc.want.([][]int); ok {
						wantTyped = w
					} else {
						t.Fatalf("Want type mismatch for empty []int case: expected [][]int, got %T", tc.want)
					}
				} else {
					t.Fatalf("Unhandled input type in test setup: %T", tc.input)
				}
			}

			// Only compare if no panic was expected
			if !tc.expectPanic {
				if !reflect.DeepEqual(got, wantTyped) {
					t.Errorf("Chunk() = %#v, want %#v", got, wantTyped)
				}
			}
		})
	}
}

// --- Chunk Examples ---
func ExampleChunk() {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	chunksOf3 := functional.Chunk(ints, 3)
	fmt.Printf("Ints chunked by 3: %v\n", chunksOf3)

	strings := []string{"a", "b", "c", "d", "e"}
	chunksOf2 := functional.Chunk(strings, 2)
	fmt.Printf("Strings chunked by 2: %v\n", chunksOf2)

	singleChunk := functional.Chunk(ints, 15) // Size larger than slice
	fmt.Printf("Chunk larger than slice: %v\n", singleChunk)

	empty := []int{}
	chunkEmpty := functional.Chunk(empty, 5)
	fmt.Printf("Chunk empty slice: %#v\n", chunkEmpty) // Use %#v for clarity

	// Output:
	// Ints chunked by 3: [[1 2 3] [4 5 6] [7 8 9] [10]]
	// Strings chunked by 2: [[a b] [c d] [e]]
	// Chunk larger than slice: [[1 2 3 4 5 6 7 8 9 10]]
	// Chunk empty slice: [][]int{}
}

// --- Benchmarks ---

// Create a helper to generate benchmark data
func generateIntSliceForChunk(size int) []int {
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = i
	}
	return data
}

// Generic Chunk benchmark runner
func benchmarkChunkGeneric(input []int, size int, b *testing.B) {
	if size <= 0 {
		b.Skip("Skipping chunk benchmark for non-positive size")
		return
	}
	b.ResetTimer()
	var result [][]int
	for i := 0; i < b.N; i++ {
		result = functional.Chunk(input, size)
	}
	_ = result
}

// Loop Chunk benchmark runner
func benchmarkChunkLoop(input []int, size int, b *testing.B) {
	if size <= 0 {
		b.Skip("Skipping chunk benchmark for non-positive size")
		return
	}
	b.ResetTimer()
	var result [][]int
	for i := 0; i < b.N; i++ {
		// Manual loop implementation
		var chunks [][]int
		sliceLen := len(input)
		if sliceLen == 0 {
			chunks = [][]int{} // Explicitly empty
			// continue // Skip rest of loop body for this iteration (though the outer check covers it)
		} else {
			// Estimate number of chunks: Ceil(len/size)
			numChunks := (sliceLen + size - 1) / size
			chunks = make([][]int, 0, numChunks) // Preallocate outer slice

			for i := 0; i < sliceLen; i += size {
				end := i + size
				if end > sliceLen {
					end = sliceLen
				}
				// Appending slice segment - this involves copying the segment
				chunks = append(chunks, input[i:end])
			}
		}
		result = chunks
	}
	_ = result
}

// Define benchmark scenarios
var (
	chunkDataN1000  = generateIntSliceForChunk(1000)
	chunkDataN10000 = generateIntSliceForChunk(10000)
)

// Small chunk size
const smallChunkSize = 5

func BenchmarkChunk_Generic_N1000_Size5(b *testing.B) {
	benchmarkChunkGeneric(chunkDataN1000, smallChunkSize, b)
}

func BenchmarkChunk_Loop_N1000_Size5(b *testing.B) {
	benchmarkChunkLoop(chunkDataN1000, smallChunkSize, b)
}

func BenchmarkChunk_Generic_N10000_Size5(b *testing.B) {
	benchmarkChunkGeneric(chunkDataN10000, smallChunkSize, b)
}

func BenchmarkChunk_Loop_N10000_Size5(b *testing.B) {
	benchmarkChunkLoop(chunkDataN10000, smallChunkSize, b)
}

// Medium chunk size
const mediumChunkSize = 50

func BenchmarkChunk_Generic_N1000_Size50(b *testing.B) {
	benchmarkChunkGeneric(chunkDataN1000, mediumChunkSize, b)
}

func BenchmarkChunk_Loop_N1000_Size50(b *testing.B) {
	benchmarkChunkLoop(chunkDataN1000, mediumChunkSize, b)
}

func BenchmarkChunk_Generic_N10000_Size50(b *testing.B) {
	benchmarkChunkGeneric(chunkDataN10000, mediumChunkSize, b)
}

func BenchmarkChunk_Loop_N10000_Size50(b *testing.B) {
	benchmarkChunkLoop(chunkDataN10000, mediumChunkSize, b)
}

// Large chunk size
const largeChunkSize = 200

func BenchmarkChunk_Generic_N1000_Size200(b *testing.B) {
	benchmarkChunkGeneric(chunkDataN1000, largeChunkSize, b)
}

func BenchmarkChunk_Loop_N1000_Size200(b *testing.B) {
	benchmarkChunkLoop(chunkDataN1000, largeChunkSize, b)
}

func BenchmarkChunk_Generic_N10000_Size200(b *testing.B) {
	benchmarkChunkGeneric(chunkDataN10000, largeChunkSize, b)
}

func BenchmarkChunk_Loop_N10000_Size200(b *testing.B) {
	benchmarkChunkLoop(chunkDataN10000, largeChunkSize, b)
}
