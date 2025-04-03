package functional_test

import (
	"fmt"
	"reflect" // For DeepEqual
	"testing"

	"github.com/JackovAlltrades/go-generics/functional"
)

// Assume person struct and ptr helper are available from other _test files in the package
// type person struct { Name string; Age int }
// func ptr[T any](v T) *T { return &v }

func TestFirst(t *testing.T) {
	testCases := []struct {
		name      string
		input     any
		wantOk    bool
		wantValue any // Expected pointer or typed nil
	}{
		{
			name:      "First_IntSlice_NotEmpty",
			input:     []int{10, 20, 30},
			wantOk:    true,
			wantValue: ptr(10),
		},
		{
			name:      "First_StringSlice_SingleElement",
			input:     []string{"hello"},
			wantOk:    true,
			wantValue: ptr("hello"),
		},
		{
			name:      "First_StructSlice",
			input:     []person{{Name: "A", Age: 1}, {Name: "B", Age: 2}},
			wantOk:    true,
			wantValue: &person{Name: "A", Age: 1}, // Pointer to expected struct
		},
		{
			name:   "First_PointerSlice", // Slice contains pointers
			input:  []*int{ptr(5), ptr(6)},
			wantOk: true,
			// Want a pointer *to the pointer* ptr(5)
			wantValue: func() **int { p := ptr(5); return &p }(),
		},
		{
			name:      "First_EmptySlice",
			input:     []int{},
			wantOk:    false,
			wantValue: (*int)(nil), // Expect typed nil pointer
		},
		{
			name:      "First_NilSlice",
			input:     ([]string)(nil),
			wantOk:    false,
			wantValue: (*string)(nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotPtr any
			var gotOk bool

			switch in := tc.input.(type) {
			case []int:
				gotPtr, gotOk = functional.First[int](in)
			case []string:
				gotPtr, gotOk = functional.First[string](in)
			case []person:
				gotPtr, gotOk = functional.First[person](in)
			case []*int: // Slice of pointers
				gotPtr, gotOk = functional.First[*int](in) // T is *int
			case nil: // Handle nil literal if needed (though typed nil preferred)
				// Need type from wantValue to call generic function
				switch tc.wantValue.(type) {
				case (*string):
					gotPtr, gotOk = functional.First[string](nil)
				case (*int):
					gotPtr, gotOk = functional.First[int](nil)
				// Add other nil types as needed
				default:
					t.Fatalf("Unhandled nil input type for %s", tc.name)
				}
			default:
				t.Fatalf("Unhandled input type %T", tc.input)
			}

			if gotOk != tc.wantOk {
				t.Errorf("First() ok = %v, want %v", gotOk, tc.wantOk)
			}

			// Compare pointer values using DeepEqual
			if !reflect.DeepEqual(gotPtr, tc.wantValue) {
				t.Errorf("First() pointer = %#v (%T), want %#v (%T)", gotPtr, gotPtr, tc.wantValue, tc.wantValue)
			}
		})
	}
}

func TestLast(t *testing.T) {
	testCases := []struct {
		name      string
		input     any
		wantOk    bool
		wantValue any // Expected pointer or typed nil
	}{
		{
			name:      "Last_IntSlice_NotEmpty",
			input:     []int{10, 20, 30},
			wantOk:    true,
			wantValue: ptr(30),
		},
		{
			name:      "Last_StringSlice_SingleElement",
			input:     []string{"hello"},
			wantOk:    true,
			wantValue: ptr("hello"),
		},
		{
			name:      "Last_StructSlice",
			input:     []person{{Name: "A", Age: 1}, {Name: "B", Age: 2}},
			wantOk:    true,
			wantValue: &person{Name: "B", Age: 2}, // Pointer to expected struct
		},
		{
			name:   "Last_PointerSlice", // Slice contains pointers
			input:  []*int{ptr(5), ptr(6)},
			wantOk: true,
			// Want a pointer *to the pointer* ptr(6)
			wantValue: func() **int { p := ptr(6); return &p }(),
		},
		{
			name:      "Last_EmptySlice",
			input:     []int{},
			wantOk:    false,
			wantValue: (*int)(nil), // Expect typed nil pointer
		},
		{
			name:      "Last_NilSlice",
			input:     ([]string)(nil),
			wantOk:    false,
			wantValue: (*string)(nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotPtr any
			var gotOk bool

			switch in := tc.input.(type) {
			case []int:
				gotPtr, gotOk = functional.Last[int](in)
			case []string:
				gotPtr, gotOk = functional.Last[string](in)
			case []person:
				gotPtr, gotOk = functional.Last[person](in)
			case []*int: // Slice of pointers
				gotPtr, gotOk = functional.Last[*int](in) // T is *int
			case nil: // Handle nil literal if needed (though typed nil preferred)
				switch tc.wantValue.(type) {
				case (*string):
					gotPtr, gotOk = functional.Last[string](nil)
				case (*int):
					gotPtr, gotOk = functional.Last[int](nil)
				default:
					t.Fatalf("Unhandled nil input type for %s", tc.name)
				}
			default:
				t.Fatalf("Unhandled input type %T", tc.input)
			}

			if gotOk != tc.wantOk {
				t.Errorf("Last() ok = %v, want %v", gotOk, tc.wantOk)
			}

			if !reflect.DeepEqual(gotPtr, tc.wantValue) {
				t.Errorf("Last() pointer = %#v (%T), want %#v (%T)", gotPtr, gotPtr, tc.wantValue, tc.wantValue)
			}
		})
	}
}

// --- Go Examples ---

func ExampleFirst() {
	nums := []int{5, 6, 7}
	firstNumPtr, okNum := functional.First(nums)
	if okNum {
		fmt.Printf("First number: %d\n", *firstNumPtr)
	}

	strs := []string{"one"}
	firstStrPtr, okStr := functional.First(strs)
	if okStr {
		fmt.Printf("First string: %s\n", *firstStrPtr)
	}

	empty := []float64{}
	_, okEmpty := functional.First(empty)
	fmt.Printf("First from empty ok? %v\n", okEmpty)

	var nilSlice []int = nil
	_, okNil := functional.First(nilSlice)
	fmt.Printf("First from nil ok? %v\n", okNil)

	// Output:
	// First number: 5
	// First string: one
	// First from empty ok? false
	// First from nil ok? false
}

func ExampleLast() {
	nums := []int{5, 6, 7}
	lastNumPtr, okNum := functional.Last(nums)
	if okNum {
		fmt.Printf("Last number: %d\n", *lastNumPtr)
	}

	strs := []string{"one"}
	lastStrPtr, okStr := functional.Last(strs)
	if okStr {
		fmt.Printf("Last string: %s\n", *lastStrPtr)
	}

	empty := []float64{}
	_, okEmpty := functional.Last(empty)
	fmt.Printf("Last from empty ok? %v\n", okEmpty)

	var nilSlice []int = nil
	_, okNil := functional.Last(nilSlice)
	fmt.Printf("Last from nil ok? %v\n", okNil)

	// Output:
	// Last number: 7
	// Last string: one
	// Last from empty ok? false
	// Last from nil ok? false
}
