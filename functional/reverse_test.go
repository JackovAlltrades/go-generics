package functional_test

import (
	"fmt"
	"reflect" // For DeepEqual
	"testing"

	"github.com/JackovAlltrades/go-generics/functional"
)

// Assume person struct is available if needed for testing
// type person struct { Name string; Age int }

func TestReverse(t *testing.T) {
	testCases := []struct {
		name  string
		input any // []T
		want  any // []T
	}{
		{
			name:  "ReverseInts",
			input: []int{1, 2, 3, 4, 5},
			want:  []int{5, 4, 3, 2, 1},
		},
		{
			name:  "ReverseStrings",
			input: []string{"a", "b", "c"},
			want:  []string{"c", "b", "a"},
		},
		{
			name:  "ReverseSingleElement",
			input: []int{42},
			want:  []int{42},
		},
		{
			name:  "ReverseEmpty",
			input: []string{},
			want:  []string{}, // Expect empty slice
		},
		{
			name:  "ReverseNil",
			input: ([]int)(nil),
			want:  []int{}, // Expect empty slice (Guideline #3)
		},
		// Add struct test if needed
		// {
		// 	name: "ReverseStructs",
		// 	input: []person{{Name: "A", Age: 1}, {Name: "B", Age: 2}},
		// 	want:  []person{{Name: "B", Age: 2}, {Name: "A", Age: 1}},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any

			// Use type switch to call generic function correctly
			switch in := tc.input.(type) {
			case []int:
				got = functional.Reverse[int](in)
			case []string:
				got = functional.Reverse[string](in)
			case []person: // If struct tests are added
				got = functional.Reverse[person](in)
			case nil: // Handle nil outer slice
				// Infer type from 'want' slice
				switch tc.want.(type) {
				case []string:
					got = functional.Reverse[string](nil)
				case []int:
					got = functional.Reverse[int](nil)
				// Add other types as needed
				default:
					t.Fatalf("Unhandled nil input type case for want type %T", tc.want)
				}
			default:
				t.Fatalf("Unhandled input type: %T", tc.input)
			}

			// Compare using DeepEqual
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Reverse() = %#v, want %#v", got, tc.want)
			}

			// Also verify the original input was not modified (for non-nil cases)
			// This requires comparing the ORIGINAL input slice (before the Reverse call)
			// with its state after the call. The table test setup makes this tricky.
			// A simpler check outside the table loop, or separate tests for mutation,
			// might be better if this guarantee is critical to enforce via tests.
			// For now, we rely on the implementation clearly showing no modification.
			// Example explicit check (less ideal within table test):
			// if tc.input != nil {
			// 	originalInputCopy := make([]T, len(tc.input.([]T))) // Requires type assertion again
			// 	copy(originalInputCopy, tc.input.([]T))
			// 	// ... call functional.Reverse ...
			// 	if !reflect.DeepEqual(originalInputCopy, tc.input.([]T)) {
			// 		 t.Errorf("Input slice was modified!")
			// 	}
			// }
		})
	}
}

// --- Go Examples ---

func ExampleReverse() {
	nums := []int{10, 20, 30, 40}
	reversedNums := functional.Reverse(nums)
	fmt.Printf("Original: %v, Reversed: %v\n", nums, reversedNums)

	strs := []string{"one", "two", "three"}
	reversedStrs := functional.Reverse(strs)
	fmt.Printf("Original: %v, Reversed: %v\n", strs, reversedStrs)

	empty := []int{}
	reversedEmpty := functional.Reverse(empty)
	fmt.Printf("Original: %v, Reversed: %#v\n", empty, reversedEmpty) // Show type

	var nilSlice []string = nil
	reversedNil := functional.Reverse(nilSlice)
	fmt.Printf("Original: %v, Reversed: %#v\n", nilSlice, reversedNil) // Show type

	// Note: Verify original slices are unchanged
	fmt.Printf("Original nums after reverse: %v\n", nums)
	fmt.Printf("Original strs after reverse: %v\n", strs)

	// Output:
	// Original: [10 20 30 40], Reversed: [40 30 20 10]
	// Original: [one two three], Reversed: [three two one]
	// Original: [], Reversed: []int{}
	// Original: [], Reversed: []string{}
	// Original nums after reverse: [10 20 30 40]
	// Original strs after reverse: [one two three]
}
