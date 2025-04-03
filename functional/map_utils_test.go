package functional_test

import (
	// Import cmp (Go 1.21+)
	"fmt"
	"reflect"
	"sort" // Needed for comparing Values in examples and Keys testing
	"testing"

	"github.com/JackovAlltrades/go-generics/functional"
)

// Assume comparablePerson struct is available from another _test file in this package
// type comparablePerson struct { ID int; Name string }
// Assume person struct is available if testing maps with struct values
// type person struct { Name string; Age int }

// Helper to compare slice contents using frequency maps (handles duplicates, order-independent)
// Requires T to be comparable.
func assertSliceContentsEqual[T comparable](t *testing.T, got, want []T) {
	t.Helper() // Mark as test helper
	if len(got) != len(want) {
		t.Errorf("Slice lengths differ: got %d, want %d. Got=%#v, Want=%#v", len(got), len(want), got, want)
		return
	}
	if len(got) == 0 {
		return // Both empty
	}

	gotFreq := make(map[T]int)
	for _, item := range got {
		gotFreq[item]++
	}

	wantFreq := make(map[T]int)
	for _, item := range want {
		wantFreq[item]++
	}

	if !reflect.DeepEqual(gotFreq, wantFreq) {
		// Provide more detail on frequency mismatch
		t.Errorf("Slice contents differ (frequency map comparison):\nGot:  %#v (freq: %v)\nWant: %#v (freq: %v)", got, gotFreq, want, wantFreq)
	}
}

func TestKeys(t *testing.T) {
	testCases := []struct {
		name     string
		inputMap any // map[K]V where K is cmp.Ordered
		wantKeys any // []K, sorted
	}{
		{
			name:     "Keys_IntKeys",
			inputMap: map[int]string{3: "c", 1: "a", 2: "b"},
			wantKeys: []int{1, 2, 3}, // Expect sorted keys
		},
		{
			name:     "Keys_StringKeys",
			inputMap: map[string]int{"banana": 2, "apple": 1, "cherry": 3},
			wantKeys: []string{"apple", "banana", "cherry"}, // Expect sorted keys
		},
		{
			name:     "Keys_EmptyMap",
			inputMap: map[int]bool{},
			wantKeys: []int{},
		},
		{
			name:     "Keys_NilMap",
			inputMap: (map[string]int)(nil),
			wantKeys: []string{}, // Expect empty slice
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotKeys any

			switch m := tc.inputMap.(type) {
			case map[int]string:
				gotKeys = functional.Keys[int, string](m)
			case map[string]int:
				gotKeys = functional.Keys[string, int](m)
			case map[int]bool: // Empty map case
				gotKeys = functional.Keys[int, bool](m)
			case nil:
				switch tc.wantKeys.(type) { // Infer key type from want
				case []string:
					gotKeys = functional.Keys[string, int](nil) // Value type doesn't matter for Keys(nil)
				case []int:
					gotKeys = functional.Keys[int, bool](nil)
				default:
					t.Fatalf("Unhandled nil map key type for %s", tc.name)
				}

			default:
				t.Fatalf("Unhandled map type: %T", tc.inputMap)
			}

			// Keys should be sorted, so DeepEqual works directly
			if !reflect.DeepEqual(gotKeys, tc.wantKeys) {
				t.Errorf("Keys() = %#v, want %#v (sorted)", gotKeys, tc.wantKeys)
			}
		})
	}
}

func TestValues(t *testing.T) {
	testCases := []struct {
		name       string
		inputMap   any // map[K]V where K is comparable
		wantValues any // []V - order doesn't matter, only content
	}{
		{
			name:       "Values_IntValues",
			inputMap:   map[string]int{"a": 10, "b": 20, "c": 10},
			wantValues: []int{10, 20, 10}, // Exact values, order may vary
		},
		{
			name:       "Values_StringValues",
			inputMap:   map[int]string{1: "hello", 2: "world", 3: "hello"},
			wantValues: []string{"hello", "world", "hello"},
		},
		// Example with comparable struct values (if comparablePerson is defined and comparable)
		{
			name: "Values_ComparableStructValues",
			inputMap: map[int]comparablePerson{
				1: {ID: 1, Name: "A"},
				2: {ID: 2, Name: "B"},
				3: {ID: 1, Name: "A"}, // Duplicate value
			},
			wantValues: []comparablePerson{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}, {ID: 1, Name: "A"}},
		},
		{
			name:       "Values_EmptyMap",
			inputMap:   map[int]string{},
			wantValues: []string{},
		},
		{
			name:       "Values_NilMap",
			inputMap:   (map[string]int)(nil),
			wantValues: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotValues any

			switch m := tc.inputMap.(type) {
			case map[string]int:
				gotValues = functional.Values[string, int](m)
			case map[int]string:
				gotValues = functional.Values[int, string](m)
			case map[int]comparablePerson: // Struct value case (Must be comparable)
				gotValues = functional.Values[int, comparablePerson](m)
			case nil:
				switch tc.wantValues.(type) { // Infer value type from want
				case []int:
					gotValues = functional.Values[string, int](nil) // Key type doesn't matter for Values(nil)
				case []string:
					gotValues = functional.Values[int, string](nil)
				case []comparablePerson:
					gotValues = functional.Values[int, comparablePerson](nil)
				default:
					t.Fatalf("Unhandled nil map value type for %s", tc.name)
				}
			default:
				t.Fatalf("Unhandled map type: %T", tc.inputMap)
			}

			// --- Use frequency map comparison for Values ---
			switch want := tc.wantValues.(type) {
			case []int:
				got, ok := gotValues.([]int)
				if !ok {
					t.Fatalf("Got values type mismatch, expected []int")
				}
				assertSliceContentsEqual[int](t, got, want) // Use frequency helper
			case []string:
				got, ok := gotValues.([]string)
				if !ok {
					t.Fatalf("Got values type mismatch, expected []string")
				}
				assertSliceContentsEqual[string](t, got, want) // Use frequency helper
			case []comparablePerson: // comparablePerson must be comparable for this helper
				got, ok := gotValues.([]comparablePerson)
				if !ok {
					t.Fatalf("Got values type mismatch, expected []comparablePerson")
				}
				assertSliceContentsEqual[comparablePerson](t, got, want) // Use frequency helper
			default:
				// If result type V is not comparable, this helper won't work.
				// Need a different comparison strategy (e.g., sorting if ordered, or custom logic).
				t.Logf("Warning: Cannot use frequency map comparison for result type %T (requires comparable). Using DeepEqual (may fail on order).", tc.wantValues)
				if !reflect.DeepEqual(gotValues, tc.wantValues) {
					t.Errorf("Values() = %#v, want %#v (order might differ, frequency check skipped)", gotValues, tc.wantValues)
				}
			}
		})
	}
}

func TestMapToSlice(t *testing.T) {
	testCases := []struct {
		name          string
		inputMap      any // map[K]V
		transformFunc any // func(K, V) T
		wantSlice     any // []T - order doesn't matter, use helper
	}{
		{
			name:     "MapIntStringToString",
			inputMap: map[int]string{1: "a", 2: "bb", 3: "ccc"},
			transformFunc: func(k int, v string) string {
				return fmt.Sprintf("%d:%s", k, v)
			},
			wantSlice: []string{"1:a", "2:bb", "3:ccc"},
		},
		{
			name:     "MapStringIntToInt",
			inputMap: map[string]int{"apple": 5, "banana": 6},
			transformFunc: func(k string, v int) int {
				return len(k) + v
			},
			wantSlice: []int{10, 12},
		},
		// Assume 'person' struct exists and 'comparablePerson' exists if needed
		{
			name: "MapIntPersonToString",
			inputMap: map[int]person{
				10: {Name: "X", Age: 1},
				20: {Name: "Y", Age: 2},
			},
			transformFunc: func(k int, v person) string {
				return fmt.Sprintf("ID%d->%s", k, v.Name)
			},
			wantSlice: []string{"ID10->X", "ID20->Y"},
		},
		{
			name:     "EmptyMap",
			inputMap: map[int]string{},
			transformFunc: func(k int, v string) string {
				t.Fatal("Transform func should not be called on empty map")
				return ""
			},
			wantSlice: []string{},
		},
		{
			name:     "NilMap",
			inputMap: (map[string]int)(nil),
			transformFunc: func(k string, v int) int {
				t.Fatal("Transform func should not be called on nil map")
				return 0
			},
			wantSlice: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotSlice any

			// Use type switch on the transform function to call MapToSlice correctly
			switch fn := tc.transformFunc.(type) {
			case func(int, string) string: // Handles MapIntStringToString AND EmptyMap
				m, ok := tc.inputMap.(map[int]string)
				if tc.inputMap != nil && !ok {
					t.Fatalf("Input map type mismatch for func(int, string) string")
				}
				gotSlice = functional.MapToSlice[int, string, string](m, fn)
			case func(string, int) int: // Handles MapStringIntToInt AND NilMap
				m, ok := tc.inputMap.(map[string]int)
				if tc.inputMap != nil && !ok {
					t.Fatalf("Input map type mismatch for func(string, int) int")
				}
				gotSlice = functional.MapToSlice[string, int, int](m, fn)
			case func(int, person) string: // Struct value case
				m, ok := tc.inputMap.(map[int]person)
				if tc.inputMap != nil && !ok {
					t.Fatalf("Input map type mismatch for func(int, person) string")
				}
				gotSlice = functional.MapToSlice[int, person, string](m, fn)

			default:
				t.Fatalf("Unhandled transformFunc type: %T", tc.transformFunc)
			}

			// --- Compare Slice Contents (Order Independent) ---
			switch want := tc.wantSlice.(type) {
			case []string: // Result slice is []string (comparable)
				got, ok := gotSlice.([]string)
				if !ok {
					t.Fatalf("Got slice type mismatch, expected []string")
				}
				assertSliceContentsEqual[string](t, got, want)
			case []int: // Result slice is []int (comparable)
				got, ok := gotSlice.([]int)
				if !ok {
					t.Fatalf("Got slice type mismatch, expected []int")
				}
				assertSliceContentsEqual[int](t, got, want)
			// Add cases for other comparable result types (T) if needed
			default:
				// If result type T is not comparable or helper not applicable
				t.Logf("Warning: Cannot use order-independent comparison for result type %T. Using DeepEqual (may fail on order).", tc.wantSlice)
				if !reflect.DeepEqual(gotSlice, tc.wantSlice) {
					t.Errorf("MapToSlice() = %#v, want %#v (order might differ)", gotSlice, tc.wantSlice)
				}
			}
		})
	}
}

// --- Go Examples ---

func ExampleKeys() {
	intKeyMap := map[int]string{30: "z", 10: "x", 20: "y"}
	keysInt := functional.Keys(intKeyMap)
	fmt.Printf("Int keys (sorted): %v\n", keysInt)

	strKeyMap := map[string]bool{"kiwi": true, "apple": false, "banana": true}
	keysStr := functional.Keys(strKeyMap)
	fmt.Printf("String keys (sorted): %v\n", keysStr)

	emptyMap := map[float64]int{}
	keysEmpty := functional.Keys(emptyMap)
	fmt.Printf("Empty map keys: %#v\n", keysEmpty)

	var nilMap map[int]int = nil
	keysNil := functional.Keys(nilMap)
	fmt.Printf("Nil map keys: %#v\n", keysNil)

	// Output:
	// Int keys (sorted): [10 20 30]
	// String keys (sorted): [apple banana kiwi]
	// Empty map keys: []float64{}
	// Nil map keys: []int{}
}

func ExampleValues() {
	strValMap := map[int]string{1: "one", 2: "two", 3: "one"}
	valuesStr := functional.Values(strValMap)
	sort.Strings(valuesStr) // Sort for predictable example output
	fmt.Printf("String values (sorted): %v\n", valuesStr)

	intValMap := map[string]int{"a": 100, "b": 50, "c": 100}
	valuesInt := functional.Values(intValMap)
	sort.Ints(valuesInt) // Sort for predictable example output
	fmt.Printf("Int values (sorted): %v\n", valuesInt)

	emptyMap := map[string]bool{}
	valuesEmpty := functional.Values(emptyMap)
	fmt.Printf("Empty map values: %#v\n", valuesEmpty)

	var nilMap map[int]float64 = nil
	valuesNil := functional.Values(nilMap)
	fmt.Printf("Nil map values: %#v\n", valuesNil)

	// Output:
	// String values (sorted): [one one two]
	// Int values (sorted): [50 100 100]
	// Empty map values: []bool{}
	// Nil map values: []float64{}
}

func ExampleMapToSlice() {
	ageMap := map[string]int{"Alice": 30, "Bob": 25}
	infoSlice := functional.MapToSlice(ageMap, func(name string, age int) string {
		return fmt.Sprintf("%s is %d", name, age)
	})
	sort.Strings(infoSlice) // Sort for predictable example output
	fmt.Printf("Info slice (sorted): %v\n", infoSlice)

	numMap := map[int]int{2: 4, 3: 9, 4: 16}
	sumSlice := functional.MapToSlice(numMap, func(k, v int) int {
		return k + v
	})
	sort.Ints(sumSlice) // Sort for predictable example output
	fmt.Printf("Sum slice (sorted): %v\n", sumSlice)

	emptyMap := map[string]int{}
	emptyResult := functional.MapToSlice(emptyMap, func(k string, v int) string { return "" })
	fmt.Printf("Empty map result: %#v\n", emptyResult)

	var nilMap map[int]string = nil
	nilResult := functional.MapToSlice(nilMap, func(k int, v string) bool { return false })
	fmt.Printf("Nil map result: %#v\n", nilResult)

	// Output:
	// Info slice (sorted): [Alice is 30 Bob is 25]
	// Sum slice (sorted): [6 12 20]
	// Empty map result: []string{}
	// Nil map result: []bool{}
}
