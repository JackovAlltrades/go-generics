package functional_test

import (
	"fmt"
	"reflect" // For deep comparison
	"testing"

	"github.com/JackovAlltrades/go-generics/functional" // Adjust import path if needed
)

// Re-using person struct or define locally if preferred
type person struct {
	Name string
	Age  int
}

// Define the generic ptr helper function at the package level
func ptr[T any](v T) *T {
	return &v
}

func TestFind(t *testing.T) {
	// NOTE: ptr helper function is now defined outside TestFind

	testCases := []struct {
		name      string // Name of the subtest
		input     any    // Input slice
		predicate any    // Predicate function
		wantFound bool   // Expected bool result
		wantValue any    // Expected value if found (use pointer for primitives/structs)
	}{
		{
			name:      "FindEvenInt_Found",
			input:     []int{1, 3, 4, 5, 6},
			predicate: func(n int) bool { return n%2 == 0 },
			wantFound: true,
			wantValue: ptr(4), // Uses the package-level ptr helper
		},
		{
			name:      "FindEvenInt_NotFound",
			input:     []int{1, 3, 5, 7},
			predicate: func(n int) bool { return n%2 == 0 },
			wantFound: false,
			wantValue: (*int)(nil),
		},
		{
			name:      "FindString_Found",
			input:     []string{"a", "bb", "ccc"},
			predicate: func(s string) bool { return len(s) == 2 },
			wantFound: true,
			wantValue: ptr("bb"), // Uses the package-level ptr helper
		},
		{
			name:      "FindString_NotFound",
			input:     []string{"a", "bb", "ccc"},
			predicate: func(s string) bool { return len(s) > 3 },
			wantFound: false,
			wantValue: (*string)(nil),
		},
		{
			name:      "FindStruct_Found",
			input:     []person{{Name: "A", Age: 20}, {Name: "B", Age: 30}, {Name: "C", Age: 30}},
			predicate: func(p person) bool { return p.Age == 30 },
			wantFound: true,
			// Manually create pointer for struct literal comparison
			wantValue: &person{Name: "B", Age: 30},
		},
		{
			name:      "FindStruct_NotFound",
			input:     []person{{Name: "A", Age: 20}, {Name: "B", Age: 25}},
			predicate: func(p person) bool { return p.Age == 30 },
			wantFound: false,
			wantValue: (*person)(nil),
		},
		{
			name:      "EmptyInput",
			input:     []int{},
			predicate: func(n int) bool { return true },
			wantFound: false,
			wantValue: (*int)(nil),
		},
		{
			name:      "NilInput",
			input:     ([]string)(nil),
			predicate: func(s string) bool { return true },
			wantFound: false,
			wantValue: (*string)(nil),
		},
		{
			name:      "FindFirstElement",
			input:     []int{10, 20, 30},
			predicate: func(n int) bool { return n == 10 },
			wantFound: true,
			wantValue: ptr(10), // Uses the package-level ptr helper
		},
		{
			name:      "FindLastElement",
			input:     []int{10, 20, 30},
			predicate: func(n int) bool { return n == 30 },
			wantFound: true,
			wantValue: ptr(30), // Uses the package-level ptr helper
		},
		{
			name:      "FindPointerValue_Found",
			input:     []*int{ptr(1), ptr(2), nil, ptr(3)},
			predicate: func(p *int) bool { return p != nil && *p == 2 },
			wantFound: true,
			// Construct the expected pointer-to-pointer manually for clarity
			wantValue: func() **int { p := ptr(2); return &p }(),
		},
		{
			name:      "FindPointerValue_NotFound",
			input:     []*int{ptr(1), nil, ptr(3)},
			predicate: func(p *int) bool { return p != nil && *p == 2 },
			wantFound: false,
			wantValue: (**int)(nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotPtr any
			var gotFound bool

			// Type switch for explicit instantiation
			switch p := tc.predicate.(type) {
			case func(int) bool:
				in, ok := tc.input.([]int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch")
				}
				gotPtr, gotFound = functional.Find[int](in, p)
			case func(string) bool:
				in, ok := tc.input.([]string)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch")
				}
				gotPtr, gotFound = functional.Find[string](in, p)
			case func(person) bool:
				in, ok := tc.input.([]person)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch")
				}
				gotPtr, gotFound = functional.Find[person](in, p)
			case func(*int) bool: // Test case for finding pointer elements
				in, ok := tc.input.([]*int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for finding pointers")
				}
				gotPtr, gotFound = functional.Find[*int](in, p) // T = *int
			default:
				t.Fatalf("Unhandled predicate type in test setup: %T", tc.predicate)
			}

			// Compare found status
			if gotFound != tc.wantFound {
				t.Errorf("Find() found status = %v, want %v", gotFound, tc.wantFound)
			}

			// Compare the found value
			if !reflect.DeepEqual(gotPtr, tc.wantValue) {
				// Provide more debugging info on mismatch
				// Handling nil cases explicitly in the error message for clarity
				if tc.wantValue == nil || (reflect.ValueOf(tc.wantValue).Kind() == reflect.Pointer && reflect.ValueOf(tc.wantValue).IsNil()) {
					if gotPtr != nil && !(reflect.ValueOf(gotPtr).Kind() == reflect.Pointer && reflect.ValueOf(gotPtr).IsNil()) {
						t.Errorf("Find() pointer = %#v (%T), want nil pointer (%T)", gotPtr, gotPtr, tc.wantValue)
					}
				} else if gotPtr == nil || (reflect.ValueOf(gotPtr).Kind() == reflect.Pointer && reflect.ValueOf(gotPtr).IsNil()) {
					t.Errorf("Find() pointer = nil, want non-nil pointer %#v (%T)", tc.wantValue, tc.wantValue)
				} else {
					t.Errorf("Find() pointer value mismatch: got %#v (%T), want %#v (%T)", gotPtr, gotPtr, tc.wantValue, tc.wantValue)
				}
			}
		})
	}
}

// ExampleFind remains the same as before
func ExampleFind() {
	// Example 1: Find first even number
	numbers := []int{1, 3, 4, 5, 6}
	evenPtr, foundEven := functional.Find[int](numbers, func(n int) bool {
		return n%2 == 0
	})
	if foundEven {
		fmt.Printf("Found first even: %d\n", *evenPtr) // Dereference pointer
	}

	// Example 2: Find number greater than 10 (not found)
	_, foundGt10 := functional.Find[int](numbers, func(n int) bool {
		return n > 10
	})
	fmt.Printf("Found number > 10: %v\n", foundGt10)

	// Example 3: Find person by name
	people := []person{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}
	bobPtr, foundBob := functional.Find[person](people, func(p person) bool {
		return p.Name == "Bob"
	})
	if foundBob {
		fmt.Printf("Found Bob: %+v\n", *bobPtr)
	}

	// Example 4: Empty slice
	empty := []string{}
	_, foundEmpty := functional.Find[string](empty, func(s string) bool { return true })
	fmt.Printf("Found in empty: %v\n", foundEmpty)

	// Example 5: Nil slice
	var nilSlice []int = nil
	_, foundNil := functional.Find[int](nilSlice, func(n int) bool { return true })
	fmt.Printf("Found in nil: %v\n", foundNil)

	// Output:
	// Found first even: 4
	// Found number > 10: false
	// Found Bob: {Name:Bob Age:25}
	// Found in empty: false
	// Found in nil: false
}
