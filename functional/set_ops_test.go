package functional_test

import (
	// Import cmp (Go 1.21+)
	"fmt"
	"reflect" // Import sort
	"testing"

	"github.com/JackovAlltrades/go-generics/functional"
)

// Note: comparablePerson struct is assumed to be defined in another _test.go file
// within this package (e.g., unique_test.go).
// Note: person struct is assumed to be defined in another _test.go file
// within this package (e.g., find_test.go).

// --- Intersection Tests ---

func TestIntersection(t *testing.T) {
	testCases := []struct {
		name string
		s1   any // []T comparable
		s2   any // []T comparable
		want any // []T
	}{
		{
			name: "Ints_SomeIntersection",
			s1:   []int{1, 2, 3, 4, 5},
			s2:   []int{4, 5, 6, 7, 8},
			want: []int{4, 5}, // Order from s1
		},
		{
			name: "Ints_FullIntersection",
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "Ints_NoIntersection",
			s1:   []int{1, 2, 3},
			s2:   []int{4, 5, 6},
			want: []int{},
		},
		{
			name: "Ints_WithDuplicates_Input",
			s1:   []int{1, 2, 2, 3, 1, 4, 4},
			s2:   []int{3, 4, 4, 5, 3},
			want: []int{3, 4}, // Unique common elements, order from s1's first appearance
		},
		{
			name: "Strings_SomeIntersection",
			s1:   []string{"a", "b", "c", "d"},
			s2:   []string{"c", "d", "e", "f"},
			want: []string{"c", "d"},
		},
		{
			name: "Strings_OrderCheck", // Ensure order comes from s1
			s1:   []string{"d", "c", "b", "a"},
			s2:   []string{"a", "b", "c", "d"},
			want: []string{"d", "c", "b", "a"}, // Reversed order because s1 is reversed
		},
		// Assuming comparablePerson is available and comparable
		{
			name: "ComparableStructs_Intersection",
			s1: []comparablePerson{
				{ID: 1, Name: "A"},
				{ID: 2, Name: "B"},
				{ID: 3, Name: "C"},
			},
			s2: []comparablePerson{
				{ID: 3, Name: "C"}, // Match
				{ID: 4, Name: "D"},
				{ID: 1, Name: "A"}, // Match
			},
			want: []comparablePerson{ // Order from s1
				{ID: 1, Name: "A"},
				{ID: 3, Name: "C"},
			},
		},
		{
			name: "S1_Empty_Intersection",
			s1:   []int{},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			name: "S2_Empty_Intersection",
			s1:   []int{1, 2, 3},
			s2:   []int{},
			want: []int{},
		},
		{
			name: "Both_Empty_Intersection",
			s1:   []int{},
			s2:   []int{},
			want: []int{},
		},
		{
			name: "S1_Nil_Intersection",
			s1:   ([]string)(nil),
			s2:   []string{"a", "b"},
			want: []string{},
		},
		{
			name: "S2_Nil_Intersection",
			s1:   []string{"a", "b"},
			s2:   ([]string)(nil),
			want: []string{},
		},
		{
			name: "Both_Nil_Intersection",
			s1:   ([]int)(nil),
			s2:   ([]int)(nil),
			want: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any

			switch s1 := tc.s1.(type) {
			case []int:
				s2, ok := tc.s2.([]int)
				if !ok && tc.s2 != nil {
					t.Fatalf("s2 type mismatch for []int test")
				}
				got = functional.Intersection[int](s1, s2)
			case []string:
				s2, ok := tc.s2.([]string)
				if !ok && tc.s2 != nil {
					t.Fatalf("s2 type mismatch for []string test")
				}
				got = functional.Intersection[string](s1, s2)
			case []comparablePerson:
				s2, ok := tc.s2.([]comparablePerson)
				if !ok && tc.s2 != nil {
					t.Fatalf("s2 type mismatch for []comparablePerson test")
				}
				got = functional.Intersection[comparablePerson](s1, s2)
			case nil:
				switch tc.want.(type) {
				case []int:
					s2, ok := tc.s2.([]int)
					if !ok && tc.s2 != nil {
						t.Fatalf("s2 type mismatch for nil []int test")
					}
					got = functional.Intersection[int](nil, s2)
				case []string:
					s2, ok := tc.s2.([]string)
					if !ok && tc.s2 != nil {
						t.Fatalf("s2 type mismatch for nil []string test")
					}
					got = functional.Intersection[string](nil, s2)
				case []comparablePerson:
					s2, ok := tc.s2.([]comparablePerson)
					if !ok && tc.s2 != nil {
						t.Fatalf("s2 type mismatch for nil []comparablePerson test")
					}
					got = functional.Intersection[comparablePerson](nil, s2)
				default:
					t.Fatalf("Unhandled type for nil s1 case: want type %T", tc.want)
				}
			default:
				t.Fatalf("Unhandled s1 type in test setup: %T", tc.s1)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Intersection() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func ExampleIntersection() {
	slice1 := []int{1, 2, 3, 4, 4, 5}
	slice2 := []int{4, 5, 6, 7, 5, 4}
	intersection1 := functional.Intersection(slice1, slice2)
	fmt.Printf("Intersection of %v and %v: %v\n", slice1, slice2, intersection1)

	slice3 := []string{"a", "b", "c"}
	slice4 := []string{"d", "e", "f"}
	intersection2 := functional.Intersection(slice3, slice4)
	fmt.Printf("Intersection of %v and %v: %v\n", slice3, slice4, intersection2)

	slice5 := []string{"x", "y", "z", "x"}
	var slice6 []string = nil
	intersection3 := functional.Intersection(slice5, slice6)
	fmt.Printf("Intersection of %v and %v: %#v\n", slice5, slice6, intersection3) // Show type

	// Output:
	// Intersection of [1 2 3 4 4 5] and [4 5 6 7 5 4]: [4 5]
	// Intersection of [a b c] and [d e f]: []
	// Intersection of [x y z x] and []: []string{}
}

// --- Union Tests ---

func TestUnion(t *testing.T) {
	testCases := []struct {
		name string
		s1   any // []T where T is cmp.Ordered
		s2   any // []T where T is cmp.Ordered
		want any // []T, sorted
	}{
		{
			name: "Ints_SomeOverlap_Union",
			s1:   []int{1, 2, 3, 4},
			s2:   []int{3, 4, 5, 6},
			want: []int{1, 2, 3, 4, 5, 6}, // Sorted unique elements
		},
		{
			name: "Ints_S2_Contains_S1_Union",
			s1:   []int{1, 2},
			s2:   []int{0, 1, 2, 3},
			want: []int{0, 1, 2, 3},
		},
		{
			name: "Ints_NoOverlap_Union",
			s1:   []int{1, 2, 3},
			s2:   []int{4, 5, 6},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "Ints_WithDuplicates_Input_Union",
			s1:   []int{1, 2, 2, 1},
			s2:   []int{3, 1, 3},
			want: []int{1, 2, 3}, // Unique elements, sorted
		},
		{
			name: "Strings_SomeOverlap_Union",
			s1:   []string{"a", "b", "c"},
			s2:   []string{"c", "d", "e", "a"},
			want: []string{"a", "b", "c", "d", "e"}, // Sorted
		},
		{
			name: "S1_Empty_Union",
			s1:   []int{},
			s2:   []int{1, 3, 2},
			want: []int{1, 2, 3}, // Sorted s2
		},
		{
			name: "S2_Empty_Union",
			s1:   []int{3, 1, 2},
			s2:   []int{},
			want: []int{1, 2, 3}, // Sorted s1
		},
		{
			name: "Both_Empty_Union",
			s1:   []int{},
			s2:   []int{},
			want: []int{},
		},
		{
			name: "S1_Nil_Union",
			s1:   ([]string)(nil),
			s2:   []string{"c", "a", "b"},
			want: []string{"a", "b", "c"}, // Sorted s2
		},
		{
			name: "S2_Nil_Union",
			s1:   []string{"c", "a", "b"},
			s2:   ([]string)(nil),
			want: []string{"a", "b", "c"}, // Sorted s1
		},
		{
			name: "Both_Nil_Union",
			s1:   ([]int)(nil),
			s2:   ([]int)(nil),
			want: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any

			switch s1 := tc.s1.(type) {
			case []int:
				s2, ok := tc.s2.([]int)
				if !ok && tc.s2 != nil {
					t.Fatalf("s2 type mismatch for []int test")
				}
				got = functional.Union[int](s1, s2)
			case []string:
				s2, ok := tc.s2.([]string)
				if !ok && tc.s2 != nil {
					t.Fatalf("s2 type mismatch for []string test")
				}
				got = functional.Union[string](s1, s2)
			// Cannot test comparablePerson unless it satisfies cmp.Ordered
			// case []comparablePerson: ...
			case nil:
				switch s2 := tc.s2.(type) {
				case []int:
					got = functional.Union[int](nil, s2)
				case []string:
					got = functional.Union[string](nil, s2)
				case nil:
					switch tc.want.(type) {
					case []int:
						got = functional.Union[int](nil, nil)
					case []string:
						got = functional.Union[string](nil, nil)
					// Add float64 if needed for nil example test
					case []float64:
						got = functional.Union[float64](nil, nil)
					default:
						t.Fatalf("Cannot infer type for both nil case, want type %T", tc.want)
					}
				default:
					t.Fatalf("Unhandled s2 type for nil s1 case: %T", tc.s2)
				}
			default:
				t.Fatalf("Unhandled s1 type in test setup (must be slice of cmp.Ordered): %T", tc.s1)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Union() = %#v, want %#v (sorted)", got, tc.want)
			}
		})
	}
}

func ExampleUnion() {
	slice1 := []int{1, 2, 2, 3, 4}
	slice2 := []int{3, 4, 4, 5, 6}
	union1 := functional.Union(slice1, slice2)
	fmt.Printf("Union of %v and %v: %v\n", slice1, slice2, union1)

	slice3 := []string{"apple", "banana"}
	slice4 := []string{"cherry", "apple", "date"}
	union2 := functional.Union(slice3, slice4)
	fmt.Printf("Union of %v and %v: %v\n", slice3, slice4, union2)

	var slice5 []int = nil
	slice6 := []int{10, 5}
	union3 := functional.Union(slice5, slice6)
	fmt.Printf("Union of %v and %v: %v\n", slice5, slice6, union3)

	slice7 := []string{"z"}
	var slice8 []string = nil
	union4 := functional.Union(slice7, slice8)
	fmt.Printf("Union of %v and %v: %v\n", slice7, slice8, union4)

	var slice9 []float64 = nil
	var slice10 []float64 = nil
	union5 := functional.Union(slice9, slice10)
	fmt.Printf("Union of %v and %v: %#v\n", slice9, slice10, union5) // Show type

	// Output:
	// Union of [1 2 2 3 4] and [3 4 4 5 6]: [1 2 3 4 5 6]
	// Union of [apple banana] and [cherry apple date]: [apple banana cherry date]
	// Union of [] and [10 5]: [5 10]
	// Union of [z] and []: [z]
	// Union of [] and []: []float64{}
}

// --- Difference Tests ---

func TestDifference(t *testing.T) {
	testCases := []struct {
		name string
		s1   any // []T comparable
		s2   any // []T comparable
		want any // []T
	}{
		{
			name: "Ints_SomeDifference",
			s1:   []int{1, 2, 3, 4, 5},
			s2:   []int{4, 5, 6, 7},
			want: []int{1, 2, 3}, // Elements in s1 but not s2, order from s1
		},
		{
			name: "Ints_NoDifference_S1_Subset_S2",
			s1:   []int{4, 5},
			s2:   []int{1, 2, 3, 4, 5, 6},
			want: []int{},
		},
		{
			name: "Ints_FullDifference_NoOverlap",
			s1:   []int{1, 2, 3},
			s2:   []int{4, 5, 6},
			want: []int{1, 2, 3}, // All of s1
		},
		{
			name: "Ints_WithDuplicates_Input_Difference",
			s1:   []int{1, 2, 2, 3, 4, 1, 4}, // Unique in s1: 1, 2, 3, 4
			s2:   []int{3, 4, 4, 5},          // Elements to remove: 3, 4
			want: []int{1, 2},                // In s1, not s2: 1, 2. Order from s1.
		},
		{
			name: "Strings_SomeDifference",
			s1:   []string{"a", "b", "c", "d"},
			s2:   []string{"c", "e", "a"},
			want: []string{"b", "d"}, // Order from s1
		},
		// Assuming comparablePerson is available and comparable
		{
			name: "ComparableStructs_Difference",
			s1: []comparablePerson{
				{ID: 1, Name: "A"}, {ID: 2, Name: "B"}, {ID: 3, Name: "C"}, {ID: 1, Name: "A"},
			},
			s2: []comparablePerson{
				{ID: 3, Name: "C"}, {ID: 4, Name: "D"},
			},
			want: []comparablePerson{ // Order from s1, unique
				{ID: 1, Name: "A"}, {ID: 2, Name: "B"},
			},
		},
		{
			name: "S1_Empty_Difference",
			s1:   []int{},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			name: "S2_Empty_Difference", // Should return unique elements of s1
			s1:   []int{1, 2, 1, 3, 2},
			s2:   []int{},
			want: []int{1, 2, 3}, // Unique s1 elements, order preserved
		},
		{
			name: "Both_Empty_Difference",
			s1:   []int{},
			s2:   []int{},
			want: []int{},
		},
		{
			name: "S1_Nil_Difference",
			s1:   ([]string)(nil),
			s2:   []string{"a", "b"},
			want: []string{},
		},
		{
			name: "S2_Nil_Difference", // Should return unique elements of s1
			s1:   []string{"a", "b", "a", "c"},
			s2:   ([]string)(nil),
			want: []string{"a", "b", "c"}, // Unique s1 elements, order preserved
		},
		{
			name: "Both_Nil_Difference",
			s1:   ([]int)(nil),
			s2:   ([]int)(nil),
			want: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any

			switch s1 := tc.s1.(type) {
			case []int:
				s2, ok := tc.s2.([]int)
				if !ok && tc.s2 != nil {
					t.Fatalf("s2 type mismatch for []int test")
				}
				got = functional.Difference[int](s1, s2)
			case []string:
				s2, ok := tc.s2.([]string)
				if !ok && tc.s2 != nil {
					t.Fatalf("s2 type mismatch for []string test")
				}
				got = functional.Difference[string](s1, s2)
			case []comparablePerson:
				s2, ok := tc.s2.([]comparablePerson)
				if !ok && tc.s2 != nil {
					t.Fatalf("s2 type mismatch for []comparablePerson test")
				}
				got = functional.Difference[comparablePerson](s1, s2)
			case nil:
				switch tc.want.(type) {
				case []int:
					s2, ok := tc.s2.([]int)
					if !ok && tc.s2 != nil {
						t.Fatalf("s2 type mismatch for nil []int test")
					}
					got = functional.Difference[int](nil, s2)
				case []string:
					s2, ok := tc.s2.([]string)
					if !ok && tc.s2 != nil {
						t.Fatalf("s2 type mismatch for nil []string test")
					}
					got = functional.Difference[string](nil, s2)
				case []comparablePerson:
					s2, ok := tc.s2.([]comparablePerson)
					if !ok && tc.s2 != nil {
						t.Fatalf("s2 type mismatch for nil []comparablePerson test")
					}
					got = functional.Difference[comparablePerson](nil, s2)
				default:
					t.Fatalf("Unhandled type for nil s1 case: want type %T", tc.want)
				}
			default:
				t.Fatalf("Unhandled s1 type in test setup: %T", tc.s1)
			}

			// Use DeepEqual because Difference preserves order from s1
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Difference() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func ExampleDifference() {
	slice1 := []int{1, 2, 3, 4, 4, 5}
	slice2 := []int{4, 5, 6, 7, 5}
	diff1 := functional.Difference(slice1, slice2)
	fmt.Printf("Difference of %v \\ %v: %v\n", slice1, slice2, diff1)

	slice3 := []string{"a", "b", "c"}
	slice4 := []string{"c", "d", "e", "a"}
	diff2 := functional.Difference(slice3, slice4)
	fmt.Printf("Difference of %v \\ %v: %v\n", slice3, slice4, diff2)

	slice5 := []string{"x", "y", "x", "z"}
	var slice6 []string = nil // Subtracting nothing
	diff3 := functional.Difference(slice5, slice6)
	fmt.Printf("Difference of %v \\ %v: %v\n", slice5, slice6, diff3)

	var slice7 []int = nil
	slice8 := []int{1, 2}
	diff4 := functional.Difference(slice7, slice8)
	fmt.Printf("Difference of %v \\ %v: %#v\n", slice7, slice8, diff4) // Show type

	// Output:
	// Difference of [1 2 3 4 4 5] \ [4 5 6 7 5]: [1 2 3]
	// Difference of [a b c] \ [c d e a]: [b]
	// Difference of [x y x z] \ []: [x y z]
	// Difference of [] \ [1 2]: []int{}
}
