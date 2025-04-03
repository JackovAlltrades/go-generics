package functional_test

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/JackovAlltrades/go-generics/functional"
)

// Assume person struct is available (e.g., from find_test.go)
// type person struct { Name string; Age int }

func TestGroupBy(t *testing.T) {
	testCases := []struct {
		name       string // Test case name
		input      any    // The input slice ( []T )
		classifier any    // The classifier function ( func(T) K )
		want       any    // The expected map ( map[K][]T )
		checkOrder bool   // Whether to check element order within map value slices
	}{
		{
			name:  "GroupBy_EvenOddInts",
			input: []int{1, 2, 3, 4, 5, 6},
			classifier: func(n int) string {
				if n%2 == 0 {
					return "even"
				}
				return "odd"
			},
			want: map[string][]int{
				"odd":  {1, 3, 5},
				"even": {2, 4, 6},
			},
			checkOrder: true, // Order within {1,3,5} and {2,4,6} matters
		},
		{
			name:  "GroupBy_StringLength",
			input: []string{"a", "bb", "ccc", "d", "ee", "fff"},
			classifier: func(s string) int {
				return len(s)
			},
			want: map[int][]string{
				1: {"a", "d"},
				2: {"bb", "ee"},
				3: {"ccc", "fff"},
			},
			checkOrder: true, // Order within value slices matters
		},
		{
			name: "GroupBy_StructField", // Assuming person struct is available
			input: []person{
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 20},
				{Name: "Charlie", Age: 30},
				{Name: "David", Age: 20},
			},
			classifier: func(p person) int {
				return p.Age
			},
			want: map[int][]person{
				30: {{Name: "Alice", Age: 30}, {Name: "Charlie", Age: 30}},
				20: {{Name: "Bob", Age: 20}, {Name: "David", Age: 20}},
			},
			checkOrder: true,
		},
		{
			name:       "GroupBy_EmptyInput",
			input:      []int{},
			classifier: func(n int) string { return "key" },
			want:       map[string][]int{}, // Expect empty map
			checkOrder: false,
		},
		{
			name:       "GroupBy_NilInput",
			input:      ([]string)(nil),
			classifier: func(s string) int { return 0 },
			want:       map[int][]string{}, // Expect empty map
			checkOrder: false,
		},
		{
			name:  "GroupBy_AllSameKey",
			input: []int{1, 1, 1, 1},
			classifier: func(n int) string {
				return "ones"
			},
			want: map[string][]int{
				"ones": {1, 1, 1, 1},
			},
			checkOrder: true,
		},
		{
			name:  "GroupBy_BooleanKey",
			input: []int{-1, 0, 1, 2, -5},
			classifier: func(n int) bool {
				return n > 0
			},
			want: map[bool][]int{
				false: {-1, 0, -5}, // Includes 0
				true:  {1, 2},
			},
			checkOrder: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any // Result will be map[K][]T

			// --- Call GroupBy using type switch ---
			// This is complex because we need T and K from input/classifier
			switch fn := tc.classifier.(type) {
			case func(int) string: // T=int, K=string
				in, ok := tc.input.([]int)
				if !ok && tc.input != nil {
					t.Fatalf("Input type mismatch")
				}
				got = functional.GroupBy[int, string](in, fn)
			case func(string) int: // T=string, K=int
				in, ok := tc.input.([]string)
				if !ok && tc.input != nil {
					t.Fatalf("Input type mismatch")
				}
				got = functional.GroupBy[string, int](in, fn)
			case func(person) int: // T=person, K=int
				in, ok := tc.input.([]person)
				if !ok && tc.input != nil {
					t.Fatalf("Input type mismatch")
				}
				got = functional.GroupBy[person, int](in, fn)
			case func(int) bool: // T=int, K=bool
				in, ok := tc.input.([]int)
				if !ok && tc.input != nil {
					t.Fatalf("Input type mismatch")
				}
				got = functional.GroupBy[int, bool](in, fn)

			default:
				// Handle nil/empty cases based on signature if needed,
				// but often easier just to check the result map.
				// If input is nil/empty, GroupBy returns empty map anyway.
				// Add specific cases if complex types are used for nil/empty tests.
				if tc.input == nil || reflect.ValueOf(tc.input).Len() == 0 {
					// Need to call GroupBy to get typed empty map for comparison
					switch tc.want.(type) {
					case map[string][]int:
						// Infer T from input type if possible, K from func, T from input
						// Need concrete types to call generic func
						// Let's rely on the fact that the code path for nil/empty
						// works regardless of the classifier signature here.
						// We just need *a* call to get the correctly typed empty map.
						emptyInput := ([]int)(nil)                          // Example T
						emptyClassifier := func(n int) string { return "" } // Example K
						got = functional.GroupBy[int, string](emptyInput, emptyClassifier)

					case map[int][]string:
						emptyInput := ([]string)(nil)
						emptyClassifier := func(s string) int { return 0 }
						got = functional.GroupBy[string, int](emptyInput, emptyClassifier)

					case map[int][]person:
						emptyInput := ([]person)(nil)
						emptyClassifier := func(p person) int { return 0 }
						got = functional.GroupBy[person, int](emptyInput, emptyClassifier)

					case map[bool][]int:
						emptyInput := ([]int)(nil)
						emptyClassifier := func(n int) bool { return false }
						got = functional.GroupBy[int, bool](emptyInput, emptyClassifier)

					default:
						t.Fatalf("Cannot determine types for empty/nil case with want type %T", tc.want)
					}
				} else {
					t.Fatalf("Unhandled classifier function type: %T", tc.classifier)
				}
			}

			// --- Compare Maps ---
			// Need to compare keys and the slices within. Order of elements in slices matters if checkOrder=true.
			if !mapsDeepEqual(t, got, tc.want, tc.checkOrder) {
				t.Errorf("GroupBy() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

// Helper function for comparing map[K][]T - more robust than reflect.DeepEqual for this structure.
func mapsDeepEqual(t *testing.T, gotMap, wantMap any, checkOrder bool) bool {
	t.Helper()
	gotVal := reflect.ValueOf(gotMap)
	wantVal := reflect.ValueOf(wantMap)

	if gotVal.Kind() != reflect.Map || wantVal.Kind() != reflect.Map {
		t.Errorf("Inputs are not maps: got %T, want %T", gotMap, wantMap)
		return false
	}
	if gotVal.IsNil() != wantVal.IsNil() {
		t.Errorf("Nil mismatch: got nil=%v, want nil=%v", gotVal.IsNil(), wantVal.IsNil())
		return false
	}
	if gotVal.Len() != wantVal.Len() {
		t.Errorf("Map lengths differ: got %d, want %d", gotVal.Len(), wantVal.Len())
		return false
	}
	if gotVal.Len() == 0 {
		return true // Both empty/nil
	}

	// Iterate over the keys of the 'want' map
	for _, key := range wantVal.MapKeys() {
		gotSliceVal := gotVal.MapIndex(key)
		wantSliceVal := wantVal.MapIndex(key)

		if !gotSliceVal.IsValid() {
			t.Errorf("Key %v missing in 'got' map", key.Interface())
			return false
		}
		if !wantSliceVal.IsValid() {
			// Should not happen if lengths match, but check anyway
			t.Errorf("Key %v somehow missing in 'want' map", key.Interface())
			return false
		}

		// Compare the slices
		if checkOrder {
			// If order matters, use DeepEqual directly on the slices
			if !reflect.DeepEqual(gotSliceVal.Interface(), wantSliceVal.Interface()) {
				t.Errorf("Slice mismatch for key %v (order matters):\nGot:  %#v\nWant: %#v",
					key.Interface(), gotSliceVal.Interface(), wantSliceVal.Interface())
				return false
			}
		} else {
			// If order doesn't matter, need frequency check (Adapt assertSliceContentsEqual?)
			// For GroupBy, order *should* typically matter as it reflects input order.
			// Let's assume checkOrder=true for most GroupBy tests. If not, need order-independent comparison.
			t.Logf("Warning: Order-independent slice comparison not fully implemented in mapsDeepEqual helper.")
			// Fallback to DeepEqual even if checkOrder is false for now
			if !reflect.DeepEqual(gotSliceVal.Interface(), wantSliceVal.Interface()) {
				t.Errorf("Slice mismatch for key %v (fallback DeepEqual):\nGot:  %#v\nWant: %#v",
					key.Interface(), gotSliceVal.Interface(), wantSliceVal.Interface())
				return false
			}
		}
	}
	return true // All keys and slices matched
}

// --- Go Examples ---

func ExampleGroupBy() {
	// Example 1: Group numbers by even/odd
	numbers := []int{1, 2, 3, 4, 5, 6, 7}
	groupedByEvenOdd := functional.GroupBy(numbers, func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})
	// Print map content (order of keys is not guaranteed)
	// To make output deterministic for testing, we can print keys sorted
	keys := make([]string, 0, len(groupedByEvenOdd))
	for k := range groupedByEvenOdd {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Requires Go 1.8+
	for _, k := range keys {
		fmt.Printf("%s: %v\n", k, groupedByEvenOdd[k])
	}
	fmt.Println("---")

	// Example 2: Group strings by first letter
	words := []string{"apple", "ant", "banana", "bat", "cat"}
	groupedByFirstLetter := functional.GroupBy(words, func(s string) string {
		if len(s) == 0 {
			return ""
		}
		return strings.ToLower(string(s[0]))
	})
	// Print map content sorted by key
	keys2 := make([]string, 0, len(groupedByFirstLetter))
	for k := range groupedByFirstLetter {
		keys2 = append(keys2, k)
	}
	sort.Strings(keys2)
	for _, k := range keys2 {
		fmt.Printf("%s: %v\n", k, groupedByFirstLetter[k])
	}
	fmt.Println("---")

	// Example 3: Empty input
	empty := []int{}
	groupedEmpty := functional.GroupBy(empty, func(n int) int { return n })
	fmt.Printf("Grouped empty: %v (len %d)\n", groupedEmpty, len(groupedEmpty))
	fmt.Println("---")

	// Example 4: Nil input
	var nilSlice []string = nil
	groupedNil := functional.GroupBy(nilSlice, func(s string) int { return len(s) })
	fmt.Printf("Grouped nil: %v (len %d)\n", groupedNil, len(groupedNil))

	// Output:
	// even: [2 4 6]
	// odd: [1 3 5 7]
	// ---
	// a: [apple ant]
	// b: [banana bat]
	// c: [cat]
	// ---
	// Grouped empty: map[] (len 0)
	// ---
	// Grouped nil: map[] (len 0)
}
