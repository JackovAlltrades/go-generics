package functional_test

import (
	"fmt"
	"testing"

	"github.com/JackovAlltrades/go-generics/functional"
)

// Note: comparablePerson struct is assumed to be defined in another _test.go file
// within this package (e.g., unique_test.go)

func TestContains(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name      string // Name of the subtest
		input     any    // Input slice (comparable type)
		value     any    // Value to search for
		wantFound bool   // Expected boolean result
	}{
		{
			name:      "ContainsInt_Found",
			input:     []int{1, 2, 3, 4, 5},
			value:     3,
			wantFound: true,
		},
		{
			name:      "ContainsInt_NotFound",
			input:     []int{1, 2, 4, 5},
			value:     3,
			wantFound: false,
		},
		{
			name:      "ContainsString_Found",
			input:     []string{"a", "b", "c"},
			value:     "b",
			wantFound: true,
		},
		{
			name:      "ContainsString_NotFound",
			input:     []string{"a", "b", "c"},
			value:     "d",
			wantFound: false,
		},
		{
			name:      "ContainsString_CaseSensitive",
			input:     []string{"a", "B", "c"},
			value:     "b",
			wantFound: false,
		},
		{
			name: "ContainsComparableStruct_Found",
			input: []comparablePerson{ // Assuming comparablePerson is available
				{ID: 1, Name: "A"},
				{ID: 2, Name: "B"},
			},
			value:     comparablePerson{ID: 2, Name: "B"},
			wantFound: true,
		},
		{
			name: "ContainsComparableStruct_NotFound",
			input: []comparablePerson{
				{ID: 1, Name: "A"},
				{ID: 2, Name: "B"},
			},
			value:     comparablePerson{ID: 3, Name: "C"},
			wantFound: false,
		},
		{
			name:      "EmptyInput",
			input:     []int{},
			value:     1,
			wantFound: false,
		},
		{
			name:      "NilInput",
			input:     ([]string)(nil),
			value:     "a",
			wantFound: false,
		},
		// ... potentially add more test cases ...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotFound bool

			switch in := tc.input.(type) {
			case []int:
				val, ok := tc.value.(int)
				if !ok {
					t.Fatalf("Value type mismatch for []int test")
				}
				gotFound = functional.Contains[int](in, val)
			case []string:
				val, ok := tc.value.(string)
				if !ok {
					t.Fatalf("Value type mismatch for []string test")
				}
				gotFound = functional.Contains[string](in, val)
			case []comparablePerson: // Assuming comparablePerson is available
				val, ok := tc.value.(comparablePerson)
				if !ok {
					t.Fatalf("Value type mismatch for []comparablePerson test")
				}
				gotFound = functional.Contains[comparablePerson](in, val)
			case nil:
				switch val := tc.value.(type) {
				case string:
					gotFound = functional.Contains[string](nil, val)
				case int:
					gotFound = functional.Contains[int](nil, val)
				// Add case for comparablePerson if needed for nil tests
				// case comparablePerson:
				// 	gotFound = functional.Contains[comparablePerson](nil, val)
				default:
					t.Fatalf("Unhandled nil input value type for %s: %T", tc.name, tc.value)
				}
			default:
				t.Fatalf("Unhandled input type in test setup: %T", tc.input)
			}

			if gotFound != tc.wantFound {
				t.Errorf("Contains(%#v, %#v) = %v, want %v", tc.input, tc.value, gotFound, tc.wantFound)
			}
		})
	}
}

// ExampleContains remains the same
func ExampleContains() {
	// Example 1: Contains integer
	numbers := []int{10, 20, 30, 40}
	has_20 := functional.Contains[int](numbers, 20)
	has_50 := functional.Contains[int](numbers, 50)
	fmt.Printf("Numbers %v contain 20? %v\n", numbers, has_20)
	fmt.Printf("Numbers %v contain 50? %v\n", numbers, has_50)

	// Example 2: Contains string
	words := []string{"apple", "banana", "cherry"}
	has_banana := functional.Contains[string](words, "banana")
	has_grape := functional.Contains[string](words, "grape")
	fmt.Printf("Words %v contain 'banana'? %v\n", words, has_banana)
	fmt.Printf("Words %v contain 'grape'? %v\n", words, has_grape)

	// Example 3: Empty slice
	empty := []int{}
	empty_has_1 := functional.Contains[int](empty, 1)
	fmt.Printf("Empty slice contains 1? %v\n", empty_has_1)

	// Example 4: Nil slice
	var nilSlice []string = nil
	nil_has_a := functional.Contains[string](nilSlice, "a")
	fmt.Printf("Nil slice contains 'a'? %v\n", nil_has_a)

	// Output:
	// Numbers [10 20 30 40] contain 20? true
	// Numbers [10 20 30 40] contain 50? false
	// Words [apple banana cherry] contain 'banana'? true
	// Words [apple banana cherry] contain 'grape'? false
	// Empty slice contains 1? false
	// Nil slice contains 'a'? false
}
