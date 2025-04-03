package functional_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/JackovAlltrades/go-generics/functional"
)

// Assume person struct is available if needed for testing
// type person struct { Name string; Age int }

func TestFlatten(t *testing.T) {
	testCases := []struct {
		name  string
		input any // [][]T
		want  any // []T
	}{
		{
			name:  "FlattenInts",
			input: [][]int{{1, 2}, {3, 4, 5}, {6}},
			want:  []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:  "FlattenStrings",
			input: [][]string{{"a", "b"}, {}, {"c", "d"}, {"e"}},
			want:  []string{"a", "b", "c", "d", "e"},
		},
		{
			name:  "FlattenWithEmptyInnerSlices",
			input: [][]int{{}, {1, 2}, {}, {3}, {}},
			want:  []int{1, 2, 3},
		},
		{
			name:  "FlattenWithNilInnerSlices",
			// Remove this duplicate input field
			// input: [][]string{{"x"}, nil, {"y", "z"}, nil}, // nil needs type hint potentially
			// Keep only this input definition
			input: func() [][]string {
				s1 := []string{"x"}
				s2 := []string{"y", "z"}
				return [][]string{s1, nil, s2, nil}
			}(),
			want: []string{"x", "y", "z"},
		},
		{
			name:  "FlattenAllEmpty",
			input: [][]int{{}, {}, {}},
			want:  []int{},
		},
		{
			name:  "FlattenAllNilInner",
			input: [][]float64{nil, nil},
			want:  []float64{},
		},
		{
			name:  "FlattenEmptyOuter",
			input: [][]int{}, // Empty outer slice
			want:  []int{},
		},
		{
			name:  "FlattenNilOuter",
			input: ([][]string)(nil), // Typed nil outer slice
			want:  []string{},        // Expect empty slice (Guideline #3)
		},
		{
			name:  "FlattenSingleInner",
			input: [][]int{{10, 20}},
			want:  []int{10, 20},
		},
		{
			name:  "FlattenSingleEmptyInner",
			input: [][]int{{}},
			want:  []int{},
		},
		// Add struct test if needed, definition assumed available
		// {
		// 	name: "FlattenStructs",
		// 	input: [][]person{
		// 		{{Name: "A", Age: 1}},
		// 		{{Name: "B", Age: 2}, {Name: "C", Age: 3}},
		// 	},
		// 	want: []person{{Name: "A", Age: 1}, {Name: "B", Age: 2}, {Name: "C", Age: 3}},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any

			// Use type switch to call generic function correctly
			switch in := tc.input.(type) {
			case [][]int:
				got = functional.Flatten[int](in)
			case [][]string:
				got = functional.Flatten[string](in)
			case [][]float64: // For AllNilInner test
				got = functional.Flatten[float64](in)
			case [][]person: // If struct tests are added
				got = functional.Flatten[person](in)
			case nil: // Handle nil outer slice
				// Infer type from 'want' slice
				switch tc.want.(type) {
				case []string:
					got = functional.Flatten[string](nil)
				case []int:
					got = functional.Flatten[int](nil)
				case []float64:
					got = functional.Flatten[float64](nil)
				// Add other types as needed
				default:
					t.Fatalf("Unhandled nil input type case for want type %T", tc.want)
				}
			default:
				t.Fatalf("Unhandled input type: %T", tc.input)
			}

			// Compare using DeepEqual
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Flatten() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

// --- Go Examples ---

func ExampleFlatten() {
	intSlices := [][]int{{1, 2}, {3}, {}, {4, 5, 6}}
	flattenedInts := functional.Flatten(intSlices)
	fmt.Printf("Flattened ints: %v\n", flattenedInts)

	stringSlices := [][]string{{"a"}, {"b", "c"}, nil, {"d"}}
	flattenedStrings := functional.Flatten(stringSlices)
	fmt.Printf("Flattened strings: %v\n", flattenedStrings)

	emptyOuter := [][]float64{}
	flattenedEmpty := functional.Flatten(emptyOuter)
	fmt.Printf("Flattened empty outer: %#v\n", flattenedEmpty) // Show type

	var nilOuter [][]int = nil
	flattenedNil := functional.Flatten(nilOuter)
	fmt.Printf("Flattened nil outer: %#v\n", flattenedNil) // Show type

	// Output:
	// Flattened ints: [1 2 3 4 5 6]
	// Flattened strings: [a b c d]
	// Flattened empty outer: []float64{}
	// Flattened nil outer: []int{}
}
