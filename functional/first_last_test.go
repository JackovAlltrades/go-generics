package functional_test

import (
	"fmt"
	"reflect" // Needed for DeepEqual in tests
	"testing"

	"github.com/JackovAlltrades/go-generics/functional" // Adjust import path
)

// --- Helper Functions & Types ---
// ptr func and person struct are defined in helpers_test.go

// --- Test First ---
func TestFirst(t *testing.T) {
	p1 := person{"A", 1}
	p2 := person{"B", 2}

	testCases := []struct {
		name           string
		input          any
		wantValueCheck func(any) bool // Use any for checks, convert inside
		wantOk         bool
	}{
		{"Ints_NonEmpty", []int{10, 20, 30}, func(v any) bool { val, ok := v.(int); return ok && val == 10 }, true},
		{"Strings_NonEmpty", []string{"a", "b"}, func(v any) bool { val, ok := v.(string); return ok && val == "a" }, true},
		{"Ints_SingleElement", []int{5}, func(v any) bool { val, ok := v.(int); return ok && val == 5 }, true},
		{"Ints_Empty", []int{}, func(v any) bool { return v == nil }, false}, // Expect nil check passes when ok=false
		{"Strings_Empty", []string{}, func(v any) bool { return v == nil }, false},
		{"Ints_Nil", ([]int)(nil), func(v any) bool { return v == nil }, false},
		{"Pointers_NonEmpty", []*int{ptr(10), ptr(20)}, func(v any) bool { p, ok := v.(*int); return ok && p != nil && *p == 10 }, true},
		{"Pointers_WithNil", []*int{nil, ptr(20)}, func(v any) bool { p, ok := v.(*int); return ok && p == nil }, true}, // Check value is nil ptr
		{"Structs_NonEmpty", []person{p1, p2}, func(v any) bool { val, ok := v.(person); return ok && reflect.DeepEqual(val, p1) }, true},
		{"Structs_Empty", []person{}, func(v any) bool { return v == nil }, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotValue any // Store the actual value pointed to, or the pointer itself for []*T cases
			var gotOk bool

			// Use type switch on input to call the generic function correctly
			switch concreteInput := tc.input.(type) {
			case []int:
				ptrVal, ok := functional.First(concreteInput)
				if ok && ptrVal != nil {
					gotValue = *ptrVal
				} else {
					gotValue = nil
				}
				gotOk = ok
			case []string:
				ptrVal, ok := functional.First(concreteInput)
				if ok && ptrVal != nil {
					gotValue = *ptrVal
				} else {
					gotValue = nil
				}
				gotOk = ok
			case []*int:
				ptrVal, ok := functional.First(concreteInput) // Returns *(*int)
				// **** THE CRITICAL FIX IS HERE ****
				if ok && ptrVal != nil {
					gotValue = *ptrVal // Assign the dereferenced pointer (*int) to gotValue
				} else {
					// If ok is true but ptrVal is nil (edge case?), or ok is false
					gotValue = nil
				}
				gotOk = ok
			case []person:
				ptrVal, ok := functional.First(concreteInput)
				if ok && ptrVal != nil {
					gotValue = *ptrVal
				} else {
					gotValue = nil
				}
				gotOk = ok
			case nil: // Handle nil slice explicitly
				_, ok := functional.First[int](nil) // Use any dummy type T for nil slice
				gotValue = nil
				gotOk = ok
			default:
				v := reflect.ValueOf(tc.input)
				if v.Kind() == reflect.Slice && v.Len() == 0 {
					gotValue = nil
					gotOk = false
				} else {
					t.Fatalf("Unhandled input type in test: %T", tc.input)
				}
			}

			if gotOk != tc.wantOk {
				t.Errorf("First() ok = %v, want %v", gotOk, tc.wantOk)
			}

			// Check the value using the provided checker function
			if !tc.wantValueCheck(gotValue) {
				t.Errorf("First() value = %#v (type %T), wantValueCheck failed (ok=%v)", gotValue, gotValue, gotOk)
			}
		})
	}
}

// --- Test Last ---
func TestLast(t *testing.T) {
	p1 := person{"A", 1}
	p2 := person{"B", 2}

	testCases := []struct {
		name           string
		input          any
		wantValueCheck func(any) bool // Use any for checks, convert inside
		wantOk         bool
	}{
		{"Ints_NonEmpty", []int{10, 20, 30}, func(v any) bool { val, ok := v.(int); return ok && val == 30 }, true},
		{"Strings_NonEmpty", []string{"a", "b"}, func(v any) bool { val, ok := v.(string); return ok && val == "b" }, true},
		{"Ints_SingleElement", []int{5}, func(v any) bool { val, ok := v.(int); return ok && val == 5 }, true},
		{"Ints_Empty", []int{}, func(v any) bool { return v == nil }, false},
		{"Strings_Empty", []string{}, func(v any) bool { return v == nil }, false},
		{"Ints_Nil", ([]int)(nil), func(v any) bool { return v == nil }, false},
		{"Pointers_NonEmpty", []*int{ptr(10), ptr(20)}, func(v any) bool { p, ok := v.(*int); return ok && p != nil && *p == 20 }, true}, // Value is *int
		{"Pointers_WithNilLast", []*int{ptr(10), nil}, func(v any) bool { p, ok := v.(*int); return ok && p == nil }, true},              // Value is nil *int
		{"Structs_NonEmpty", []person{p1, p2}, func(v any) bool { val, ok := v.(person); return ok && reflect.DeepEqual(val, p2) }, true},
		{"Structs_Empty", []person{}, func(v any) bool { return v == nil }, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotValue any // Store the actual value pointed to, or the pointer itself for []*T cases
			var gotOk bool

			// Use type switch on input to call the generic function correctly
			switch concreteInput := tc.input.(type) {
			case []int:
				ptrVal, ok := functional.Last(concreteInput)
				if ok && ptrVal != nil {
					gotValue = *ptrVal
				} else {
					gotValue = nil
				}
				gotOk = ok
			case []string:
				ptrVal, ok := functional.Last(concreteInput)
				if ok && ptrVal != nil {
					gotValue = *ptrVal
				} else {
					gotValue = nil
				}
				gotOk = ok
			case []*int:
				ptrVal, ok := functional.Last(concreteInput) // Returns *(*int)
				// **** THE CRITICAL FIX IS HERE ****
				if ok && ptrVal != nil {
					gotValue = *ptrVal // Assign the dereferenced pointer (*int) to gotValue
				} else {
					gotValue = nil
				}
				gotOk = ok
			case []person:
				ptrVal, ok := functional.Last(concreteInput)
				if ok && ptrVal != nil {
					gotValue = *ptrVal
				} else {
					gotValue = nil
				}
				gotOk = ok
			case nil:
				_, ok := functional.Last[int](nil)
				gotValue = nil
				gotOk = ok
			default:
				v := reflect.ValueOf(tc.input)
				if v.Kind() == reflect.Slice && v.Len() == 0 {
					gotValue = nil
					gotOk = false
				} else {
					t.Fatalf("Unhandled input type in test: %T", tc.input)
				}
			}

			if gotOk != tc.wantOk {
				t.Errorf("Last() ok = %v, want %v", gotOk, tc.wantOk)
			}

			// Check the value using the provided checker function
			if !tc.wantValueCheck(gotValue) {
				t.Errorf("Last() value = %#v (type %T), wantValueCheck failed (ok=%v)", gotValue, gotValue, gotOk)
			}
		})
	}
}

// --- Examples ---

func ExampleFirst() { /* Example unchanged, seems correct */
	nums := []int{5, 10, 15}
	firstNumPtr, ok := functional.First(nums)
	if ok {
		fmt.Printf("First number: %d\n", *firstNumPtr)
	} else {
		fmt.Println("Slice was empty")
	}
	empty := []string{}
	firstStrPtr, ok := functional.First(empty)
	if !ok {
		fmt.Printf("First string from empty found: %v (ptr is nil: %v)\n", ok, firstStrPtr == nil)
	}
	// Output:
	// First number: 5
	// First string from empty found: false (ptr is nil: true)
}

func ExampleLast() { /* Example unchanged, seems correct */
	nums := []int{5, 10, 15}
	lastNumPtr, ok := functional.Last(nums)
	if ok {
		fmt.Printf("Last number: %d\n", *lastNumPtr)
	}
	pointers := []*string{ptr("a"), ptr("b"), nil}
	lastPtrPtr, ok := functional.Last(pointers)
	if ok {
		if lastPtrPtr != nil && *lastPtrPtr == nil {
			fmt.Println("Last element was a nil pointer")
		} else if lastPtrPtr != nil {
			fmt.Printf("Last pointer points to: %q\n", **lastPtrPtr)
		} else {
			fmt.Println("Last pointer itself was unexpectedly nil despite ok=true")
		}
	} else {
		fmt.Println("Slice was empty")
	}
	// Output:
	// Last number: 15
	// Last element was a nil pointer
}
