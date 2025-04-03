package functional_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/JackovAlltrades/go-generics/functional"
)

func TestChunk(t *testing.T) {
	testCases := []struct {
		name  string
		input any // []T
		size  int // Chunk size
		want  any // [][]T
	}{
		{
			name:  "Ints_EvenSplit",
			input: []int{1, 2, 3, 4, 5, 6},
			size:  2,
			want:  [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name:  "Ints_UnevenSplit",
			input: []int{1, 2, 3, 4, 5, 6, 7},
			size:  3,
			want:  [][]int{{1, 2, 3}, {4, 5, 6}, {7}}, // Last chunk smaller
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
			want:  [][]int{{1, 2, 3}}, // Single chunk
		},
		{
			name:  "Size_EqualToSlice",
			input: []int{1, 2, 3},
			size:  3,
			want:  [][]int{{1, 2, 3}}, // Single chunk
		},
		{
			name:  "EmptyInput",
			input: []int{},
			size:  3,
			want:  [][]int{}, // Empty slice of slices
		},
		{
			name:  "NilInput",
			input: ([]string)(nil),
			size:  2,
			want:  [][]string{}, // Empty slice of slices
		},
		{
			name:  "Size_Zero",
			input: []int{1, 2, 3},
			size:  0,
			want:  [][]int{}, // Invalid size -> empty result
		},
		{
			name:  "Size_Negative",
			input: []int{1, 2, 3},
			size:  -1,
			want:  [][]int{}, // Invalid size -> empty result
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any

			switch in := tc.input.(type) {
			case []int:
				got = functional.Chunk[int](in, tc.size)
			case []string:
				got = functional.Chunk[string](in, tc.size)
			case nil: // Handle nil case explicitly if needed based on type T
				switch tc.want.(type) {
				case [][]int:
					got = functional.Chunk[int](nil, tc.size)
				case [][]string:
					got = functional.Chunk[string](nil, tc.size)
				default:
					t.Fatalf("Unhandled 'want' type %T for nil input", tc.want)
				}
			default:
				t.Fatalf("Unhandled input type: %T", tc.input)
			}

			// Use DeepEqual for comparison, order and content matter
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Chunk() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

// --- Go Examples ---

func ExampleChunk() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Chunk into slices of 3
	chunksOf3 := functional.Chunk(numbers, 3)
	fmt.Printf("Chunks of 3: %v\n", chunksOf3)

	// Chunk into slices of 4
	chunksOf4 := functional.Chunk(numbers, 4)
	fmt.Printf("Chunks of 4: %v\n", chunksOf4)

	// Chunk with size larger than slice
	chunksOf15 := functional.Chunk(numbers, 15)
	fmt.Printf("Chunks of 15: %v\n", chunksOf15)

	// Chunk empty slice
	empty := []string{}
	emptyChunks := functional.Chunk(empty, 2)
	fmt.Printf("Chunks of empty: %#v\n", emptyChunks)

	// Chunk with invalid size
	invalidChunks := functional.Chunk(numbers, 0)
	fmt.Printf("Chunks with size 0: %#v\n", invalidChunks)

	// Output:
	// Chunks of 3: [[1 2 3] [4 5 6] [7 8 9] [10]]
	// Chunks of 4: [[1 2 3 4] [5 6 7 8] [9 10]]
	// Chunks of 15: [[1 2 3 4 5 6 7 8 9 10]]
	// Chunks of empty: [][]string{}
	// Chunks with size 0: [][]int{}
}
