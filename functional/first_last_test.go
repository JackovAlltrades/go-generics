package functional_test

import (
	"fmt"
	"reflect" // Needed for DeepEqual in tests
	"testing"

	"github.com/JackovAlltrades/go-generics/functional" // Adjust import path
)

// --- Helper Functions ---

// Helper function to create a pointer to a value.
func ptr[T any](v T) *T {
	return &v
}

// person struct used in tests.
type person struct {
	Name string
	Age  int
}

// --- Test Find ---
// NOTE: This is defined in find_test.go, but needed here for comparison if kept separate
// It's better practice to define shared test types once. Assuming it's defined elsewhere.

// --- Test First ---
func TestFirst(t *testing.T) {
	p1 := person{"A", 1}
	p2 := person{"B", 2}

	testCases := []struct {
		name           string
		input          any
		wantValueCheck func(any) bool
		wantOk         bool
	}{
		{"Ints_NonEmpty", []int{10, 20, 30}, func(v any) bool { return v.(int) == 10 }, true},
		{"Strings_NonEmpty", []string{"a", "b"}, func(v any) bool { return v.(string) == "a" }, true},
		{"Ints_SingleElement", []int{5}, func(v any) bool { return v.(int) == 5 }, true},
		{"Ints_Empty", []int{}, func(v any) bool { return false }, false},
		{"Strings_Empty", []string{}, func(v any) bool { return false }, false},
		{"Ints_Nil", ([]int)(nil), func(v any) bool { return false }, false},
		{"Pointers_NonEmpty", []*int{ptr(10), ptr(20)}, func(v any) bool { p, ok := v.(*int); return ok && p != nil && *p == 10 }, true},
		{"Pointers_WithNil", []*int{nil, ptr(20)}, func(v any) bool { p, ok := v.(*int); return ok && p == nil }, true},
		{"Structs_NonEmpty", []person{p1, p2}, func(v any) bool { return reflect.DeepEqual(v, p1) }, true},
		{"Structs_Empty", []person{}, func(v any) bool { return false }, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotPtr any
			var gotOk bool

			v := reflect.ValueOf(tc.input)
			if !v.IsValid() || (v.Kind() == reflect.Slice && v.IsNil()) {
				_, gotOk = functional.First[int](nil)
				gotPtr = (any)(nil)
			} else {
				switch concreteInput := tc.input.(type) {
				case []int:
					ptrVal, ok := functional.First(concreteInput)
					gotPtr = ptrVal
					gotOk = ok
				case []string:
					ptrVal, ok := functional.First(concreteInput)
					gotPtr = ptrVal
					gotOk = ok
				case []*int:
					ptrVal, ok := functional.First(concreteInput)
					gotPtr = ptrVal
					gotOk = ok
				case []person:
					ptrVal, ok := functional.First(concreteInput)
					gotPtr = ptrVal
					gotOk = ok
				default:
					if v.Kind() == reflect.Slice && v.Len() == 0 {
						switch v.Type().Elem().Kind() {
						case reflect.Int:
							_, gotOk = functional.First([]int{})
							gotPtr = (*int)(nil)
						case reflect.String:
							_, gotOk = functional.First([]string{})
							gotPtr = (*string)(nil)
						case reflect.Pointer:
							_, gotOk = functional.First([]*int{})
							gotPtr = nil
						case reflect.Struct:
							_, gotOk = functional.First([]person{})
							gotPtr = (*person)(nil)
						default:
							t.Fatalf("Unhandled empty slice element type: %s", v.Type().Elem().Kind())
						}
					} else {
						t.Fatalf("Unhandled input type in test: %T", tc.input)
					}
				}
			}

			if gotOk != tc.wantOk {
				t.Errorf("First() ok = %v, want %v", gotOk, tc.wantOk)
			}

			if gotOk {
				isGotPtrNil := (gotPtr == nil)
				if !isGotPtrNil {
					// Handle typed nils potentially held by interface
					vPtr := reflect.ValueOf(gotPtr)
					if vPtr.Kind() == reflect.Pointer && vPtr.IsNil() {
						isGotPtrNil = true
					}
				}

				if isGotPtrNil {
					if !tc.wantValueCheck(nil) {
						t.Errorf("First() pointer is nil, but ok is true and test did not expect nil value")
					}
				} else {
					// Pointer is not nil, check the value
					val := reflect.ValueOf(gotPtr).Elem().Interface()
					if !tc.wantValueCheck(val) {
						t.Errorf("First() value = %#v, but wantValueCheck failed", val)
					}
				}
			} else {
				// Check if gotPtr is actually nil when ok is false
				if gotPtr != nil {
					vPtr := reflect.ValueOf(gotPtr)
					if !(vPtr.Kind() == reflect.Pointer && vPtr.IsNil()) { // Allow typed nil pointers
						t.Errorf("First() pointer = %v, want nil when ok is false", gotPtr)
					}
				}
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
		wantValueCheck func(any) bool
		wantOk         bool
	}{
		{"Ints_NonEmpty", []int{10, 20, 30}, func(v any) bool { return v.(int) == 30 }, true},
		{"Strings_NonEmpty", []string{"a", "b"}, func(v any) bool { return v.(string) == "b" }, true},
		{"Ints_SingleElement", []int{5}, func(v any) bool { return v.(int) == 5 }, true},
		{"Ints_Empty", []int{}, func(v any) bool { return false }, false},
		{"Strings_Empty", []string{}, func(v any) bool { return false }, false},
		{"Ints_Nil", ([]int)(nil), func(v any) bool { return false }, false},
		{"Pointers_NonEmpty", []*int{ptr(10), ptr(20)}, func(v any) bool { p, ok := v.(*int); return ok && p != nil && *p == 20 }, true},
		{"Pointers_WithNilLast", []*int{ptr(10), nil}, func(v any) bool { p, ok := v.(*int); return ok && p == nil }, true},
		{"Structs_NonEmpty", []person{p1, p2}, func(v any) bool { return reflect.DeepEqual(v, p2) }, true},
		{"Structs_Empty", []person{}, func(v any) bool { return false }, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotPtr any
			var gotOk bool

			v := reflect.ValueOf(tc.input)
			if !v.IsValid() || (v.Kind() == reflect.Slice && v.IsNil()) {
				_, gotOk = functional.Last[int](nil)
				gotPtr = (any)(nil)
			} else {
				switch concreteInput := tc.input.(type) {
				case []int:
					ptrVal, ok := functional.Last(concreteInput)
					gotPtr = ptrVal
					gotOk = ok
				case []string:
					ptrVal, ok := functional.Last(concreteInput)
					gotPtr = ptrVal
					gotOk = ok
				case []*int:
					ptrVal, ok := functional.Last(concreteInput)
					gotPtr = ptrVal
					gotOk = ok
				case []person:
					ptrVal, ok := functional.Last(concreteInput)
					gotPtr = ptrVal
					gotOk = ok
				default:
					if v.Kind() == reflect.Slice && v.Len() == 0 {
						switch v.Type().Elem().Kind() {
						case reflect.Int:
							_, gotOk = functional.Last([]int{})
							gotPtr = (*int)(nil)
						case reflect.String:
							_, gotOk = functional.Last([]string{})
							gotPtr = (*string)(nil)
						case reflect.Pointer:
							_, gotOk = functional.Last([]*int{})
							gotPtr = nil
						case reflect.Struct:
							_, gotOk = functional.Last([]person{})
							gotPtr = (*person)(nil)
						default:
							t.Fatalf("Unhandled empty slice element type: %s", v.Type().Elem().Kind())
						}
					} else {
						t.Fatalf("Unhandled input type in test: %T", tc.input)
					}
				}
			}

			if gotOk != tc.wantOk {
				t.Errorf("Last() ok = %v, want %v", gotOk, tc.wantOk)
			}

			if gotOk {
				isGotPtrNil := (gotPtr == nil)
				if !isGotPtrNil {
					vPtr := reflect.ValueOf(gotPtr)
					if vPtr.Kind() == reflect.Pointer && vPtr.IsNil() {
						isGotPtrNil = true
					}
				}

				if isGotPtrNil {
					if !tc.wantValueCheck(nil) {
						t.Errorf("Last() pointer is nil, but ok is true and test did not expect nil value")
					}
				} else {
					val := reflect.ValueOf(gotPtr).Elem().Interface()
					if !tc.wantValueCheck(val) {
						t.Errorf("Last() value = %#v, but wantValueCheck failed", val)
					}
				}
			} else {
				if gotPtr != nil {
					vPtr := reflect.ValueOf(gotPtr)
					if !(vPtr.Kind() == reflect.Pointer && vPtr.IsNil()) { // Allow typed nil pointers
						t.Errorf("Last() pointer = %v, want nil when ok is false", gotPtr)
					}
				}
			}
		})
	}
}

// --- Examples ---

func ExampleFirst() {
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
		fmt.Printf("First string from empty found: %v (ptr: %v)\n", ok, firstStrPtr)
	}
	// Corrected Output:
	// First number: 5
	// First string from empty found: false (ptr: <nil>)
}

func ExampleLast() {
	nums := []int{5, 10, 15}
	lastNumPtr, ok := functional.Last(nums)
	if ok {
		fmt.Printf("Last number: %d\n", *lastNumPtr)
	}

	pointers := []*string{ptr("a"), ptr("b"), nil}
	lastPtr, ok := functional.Last(pointers) // Last element is nil pointer
	if ok {
		// CORRECTED Logic: Check if the returned pointer itself is nil
		if lastPtr == nil {
			fmt.Println("Last pointer was nil")
		} else {
			// Should not happen for this specific example, but good practice
			fmt.Printf("Last pointer points to: %v\n", *lastPtr)
		}
	} else {
		// Should not happen for this specific example
		fmt.Println("Last element not found (slice was empty?)")
	}
	// Corrected Output:
	// Last number: 15
	// Last pointer was nil
}
