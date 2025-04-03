package functional_test

import (
	"fmt"
	"testing"

	"github.com/JackovAlltrades/go-generics/functional"
)

// Remove duplicate person struct and ptr helper function declarations
// They're already defined in other test files in the same package

func TestAny(t *testing.T) {
	testCases := []struct {
		name      string
		input     any
		predicate any
		want      bool
	}{
		{
			name:      "Any_IntEven_True",
			input:     []int{1, 3, 4, 5},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      true, // 4 is even
		},
		{
			name:      "Any_IntEven_False",
			input:     []int{1, 3, 5, 7},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      false,
		},
		{
			name:      "Any_StringNonEmpty_True",
			input:     []string{"", "a", ""},
			predicate: func(s string) bool { return s != "" },
			want:      true, // "a" is non-empty
		},
		{
			name:      "Any_StringNonEmpty_False",
			input:     []string{"", "", ""},
			predicate: func(s string) bool { return s != "" },
			want:      false,
		},
		{
			name:      "Any_StructAge_True",
			input:     []person{{Age: 10}, {Age: 20}},
			predicate: func(p person) bool { return p.Age > 15 },
			want:      true, // Age 20 > 15
		},
		{
			name:      "Any_StructAge_False",
			input:     []person{{Age: 10}, {Age: 15}},
			predicate: func(p person) bool { return p.Age > 15 },
			want:      false,
		},
		{
			name:      "Any_EmptyInput",
			input:     []int{},
			predicate: func(n int) bool { return true },
			want:      false, // No elements to satisfy
		},
		{
			name:      "Any_NilInput",
			input:     ([]float64)(nil),
			predicate: func(f float64) bool { return true },
			want:      false, // No elements to satisfy
		},
		{
			name:      "Any_Pointer_True",
			input:     []*int{nil, ptr(5), nil},
			predicate: func(p *int) bool { return p != nil && *p > 0 },
			want:      true,
		},
		{
			name:      "Any_Pointer_False",
			input:     []*int{nil, nil},
			predicate: func(p *int) bool { return p != nil },
			want:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got bool
			// Type switch for predicate to call correct instantiation
			switch p := tc.predicate.(type) {
			case func(int) bool:
				in, ok := tc.input.([]int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch")
				}
				got = functional.Any[int](in, p)
			case func(string) bool:
				in, ok := tc.input.([]string)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch")
				}
				got = functional.Any[string](in, p)
			case func(person) bool:
				in, ok := tc.input.([]person)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch")
				}
				got = functional.Any[person](in, p)
			case func(float64) bool: // For nil test case
				in, ok := tc.input.([]float64)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch")
				}
				got = functional.Any[float64](in, p)
			case func(*int) bool: // Pointer test case
				in, ok := tc.input.([]*int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch")
				}
				got = functional.Any[*int](in, p)
			default:
				t.Fatalf("Unhandled predicate type: %T", tc.predicate)
			}

			if got != tc.want {
				t.Errorf("Any() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestAll(t *testing.T) {
	testCases := []struct {
		name      string
		input     any
		predicate any
		want      bool
	}{
		{
			name:      "All_IntEven_False", // Not all are even
			input:     []int{2, 4, 5, 6},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      false, // 5 fails
		},
		{
			name:      "All_IntEven_True", // All are even
			input:     []int{2, 4, 6, 8},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      true,
		},
		{
			name:      "All_StringNonEmpty_False",
			input:     []string{"a", "", "b"}, // Contains empty string
			predicate: func(s string) bool { return s != "" },
			want:      false,
		},
		{
			name:      "All_StringNonEmpty_True",
			input:     []string{"a", "b", "c"},
			predicate: func(s string) bool { return s != "" },
			want:      true,
		},
		{
			name:      "All_StructAge_False",
			input:     []person{{Age: 20}, {Age: 15}},
			predicate: func(p person) bool { return p.Age > 18 },
			want:      false, // Age 15 fails
		},
		{
			name:      "All_StructAge_True",
			input:     []person{{Age: 20}, {Age: 19}},
			predicate: func(p person) bool { return p.Age > 18 },
			want:      true,
		},
		{
			name:      "All_EmptyInput", // Vacuously true
			input:     []int{},
			predicate: func(n int) bool { return false }, // Predicate doesn't matter
			want:      true,
		},
		{
			name:      "All_NilInput", // Vacuously true
			input:     ([]string)(nil),
			predicate: func(s string) bool { return false },
			want:      true,
		},
		{
			name:      "All_PointerNotNil_False",
			input:     []*int{ptr(1), nil, ptr(3)},
			predicate: func(p *int) bool { return p != nil },
			want:      false,
		},
		{
			name:      "All_PointerNotNil_True",
			input:     []*int{ptr(1), ptr(2)},
			predicate: func(p *int) bool { return p != nil },
			want:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got bool
			// Type switch for predicate to call correct instantiation
			switch p := tc.predicate.(type) {
			case func(int) bool:
				in, ok := tc.input.([]int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch")
				}
				got = functional.All[int](in, p)
			case func(string) bool:
				in, ok := tc.input.([]string)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch")
				}
				got = functional.All[string](in, p)
			case func(person) bool:
				in, ok := tc.input.([]person)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch")
				}
				got = functional.All[person](in, p)
			case func(*int) bool: // Pointer test case
				in, ok := tc.input.([]*int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch")
				}
				got = functional.All[*int](in, p)
			default:
				t.Fatalf("Unhandled predicate type: %T", tc.predicate)
			}

			if got != tc.want {
				t.Errorf("All() = %v, want %v", got, tc.want)
			}
		})
	}
}

// --- Go Examples ---

func ExampleAny() {
	// Example 1: Any even numbers?
	numbers := []int{1, 3, 4, 5, 7}
	hasEven := functional.Any[int](numbers, func(n int) bool { return n%2 == 0 })
	fmt.Printf("Any even? %v\n", hasEven)

	// Example 2: Any empty strings?
	words := []string{"hello", "", "world"}
	hasEmpty := functional.Any[string](words, func(s string) bool { return s == "" })
	fmt.Printf("Any empty strings? %v\n", hasEmpty)

	// Example 3: Any in empty slice?
	emptyInts := []int{}
	anyInEmpty := functional.Any[int](emptyInts, func(n int) bool { return true })
	fmt.Printf("Any in empty? %v\n", anyInEmpty)

	// Output:
	// Any even? true
	// Any empty strings? true
	// Any in empty? false
}

func ExampleAll() {
	// Example 1: All numbers positive?
	numbers1 := []int{1, 3, 4, 5, 7}
	allPositive1 := functional.All[int](numbers1, func(n int) bool { return n > 0 })
	fmt.Printf("All positive in %v? %v\n", numbers1, allPositive1)

	numbers2 := []int{1, -3, 4, 5, 7}
	allPositive2 := functional.All[int](numbers2, func(n int) bool { return n > 0 })
	fmt.Printf("All positive in %v? %v\n", numbers2, allPositive2)

	// Example 2: All strings non-empty?
	words1 := []string{"hello", "world"}
	allNonEmpty1 := functional.All[string](words1, func(s string) bool { return s != "" })
	fmt.Printf("All non-empty in %v? %v\n", words1, allNonEmpty1)

	words2 := []string{"hello", "", "world"}
	allNonEmpty2 := functional.All[string](words2, func(s string) bool { return s != "" })
	fmt.Printf("All non-empty in %v? %v\n", words2, allNonEmpty2)

	// Example 3: All in empty slice? (Vacuously true)
	emptyInts := []int{}
	allInEmpty := functional.All[int](emptyInts, func(n int) bool { return false })
	fmt.Printf("All in empty? %v\n", allInEmpty)

	// Output:
	// All positive in [1 3 4 5 7]? true
	// All positive in [1 -3 4 5 7]? false
	// All non-empty in [hello world]? true
	// All non-empty in [hello  world]? false
	// All in empty? true
}
