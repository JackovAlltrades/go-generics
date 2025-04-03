package functional_test

import (
	"fmt"
	"reflect" // For DeepEqual
	"testing"

	"github.com/JackovAlltrades/go-generics/functional"
)

// Define a comparable struct for testing
type comparablePerson struct {
	ID   int
	Name string // Assuming Name doesn't affect comparability here if ID is unique
}

// Note: We don't need the 'ptr' helper here as Unique deals with values.

func TestUnique(t *testing.T) {
	testCases := []struct {
		name  string
		input any // Use []T where T is comparable
		want  any // Expected []T
	}{
		{
			name:  "UniqueInts_WithDuplicates",
			input: []int{1, 2, 2, 3, 1, 4, 5, 4},
			want:  []int{1, 2, 3, 4, 5}, // Order of first appearance preserved
		},
		{
			name:  "UniqueInts_NoDuplicates",
			input: []int{1, 2, 3, 4, 5},
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			name:  "UniqueInts_AllDuplicates",
			input: []int{7, 7, 7, 7},
			want:  []int{7},
		},
		{
			name:  "UniqueStrings_WithDuplicates",
			input: []string{"a", "b", "a", "c", "b", "d"},
			want:  []string{"a", "b", "c", "d"},
		},
		{
			name:  "UniqueStrings_NoDuplicates",
			input: []string{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name: "UniqueComparableStructs",
			input: []comparablePerson{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
				{ID: 1, Name: "Alice"}, // Duplicate based on ID and Name (struct is comparable)
				{ID: 3, Name: "Charlie"},
				{ID: 2, Name: "Bob"},
			},
			want: []comparablePerson{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
				{ID: 3, Name: "Charlie"},
			},
		},
		{
			name:  "EmptyInput",
			input: []int{},
			want:  []int{}, // Expect empty, non-nil slice
		},
		{
			name:  "NilInput",
			input: ([]string)(nil), // Typed nil
			want:  ([]string)(nil), // Expect nil
		},
		{
			name:  "SingleElement",
			input: []int{100},
			want:  []int{100},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any

			// --- Type Assertions and Explicit Instantiation ---
			// Requires comparable types
			switch in := tc.input.(type) {
			case []int:
				// Explicit type parameter instantiation not strictly necessary
				// for the call if compiler can infer, but good for clarity/consistency
				got = functional.Unique[int](in)
			case []string:
				got = functional.Unique[string](in)
			case []comparablePerson:
				got = functional.Unique[comparablePerson](in)
			case nil: // Handle the nil input case explicitly for type inference
				// Need to know the *intended* type if input is nil to call Unique
				// We infer from `want` type for the test setup if possible
				switch tc.want.(type) {
				case []string:
					got = functional.Unique[string](nil)
				case []int: // Add other types used in nil tests if necessary
					got = functional.Unique[int](nil)
				default:
					t.Fatalf("Unhandled nil input type for %s", tc.name)
				}
			default:
				t.Fatalf("Unhandled input type in test setup: %T", tc.input)
			}

			// --- Comparison ---
			// Use reflect.DeepEqual for comparing slices, handles nil/empty correctly.
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Unique() = %#v (%[2]T), want %#v (%[3]T)", got, got, tc.want, tc.want)
			}
		})
	}
}

// --- Go Example ---

func ExampleUnique() {
	// Example 1: Unique integers
	numbers := []int{1, 2, 2, 3, 1, 4, 4, 5}
	uniqueNumbers := functional.Unique[int](numbers)
	fmt.Println("Unique numbers:", uniqueNumbers)

	// Example 2: Unique strings
	words := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	uniqueWords := functional.Unique[string](words)
	fmt.Println("Unique words:", uniqueWords)

	// Example 3: Empty slice
	emptySlice := []int{}
	uniqueEmpty := functional.Unique[int](emptySlice)
	fmt.Printf("Unique empty: %#v (is nil: %v)\n", uniqueEmpty, uniqueEmpty == nil)

	// Example 4: Nil slice
	var nilSlice []string = nil
	uniqueNil := functional.Unique[string](nilSlice)
	fmt.Printf("Unique nil: %#v (is nil: %v)\n", uniqueNil, uniqueNil == nil)

	// Output:
	// Unique numbers: [1 2 3 4 5]
	// Unique words: [apple banana orange grape]
	// Unique empty: []int{} (is nil: false)
	// Unique nil: []string(nil) (is nil: true)
}
