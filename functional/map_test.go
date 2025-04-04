package functional_test // Use the _test package for black-box testing

import (
	"fmt"
	"reflect" // For deep comparison
	"strconv"
	"testing"

	"github.com/JackovAlltrades/go-generics/functional" // Adjust import path if needed
)

// Remove duplicate person struct declaration - already defined in find_test.go

func TestMap(t *testing.T) {
	// Define test cases using a slice of structs
	testCases := []struct {
		name    string // Name of the subtest
		input   any    // Input slice (use 'any' for type flexibility in table)
		mapFunc any    // Mapping function (use 'any')
		want    any    // Expected output slice (use 'any')
	}{
		{
			name:  "IntToString",
			input: []int{1, 2, 3, -4},
			mapFunc: func(i int) string {
				return "v" + strconv.Itoa(i)
			},
			want: []string{"v1", "v2", "v3", "v-4"},
		},
		{
			name:  "StringLength",
			input: []string{"a", "bc", "", "def"},
			mapFunc: func(s string) int {
				return len(s)
			},
			want: []int{1, 2, 0, 3},
		},
		{
			name: "StructToField",
			input: []person{
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 25},
			},
			mapFunc: func(p person) string {
				return p.Name
			},
			want: []string{"Alice", "Bob"},
		},
		{
			name:    "EmptyInput",
			input:   []int{}, // Explicitly empty, non-nil
			mapFunc: func(i int) int { return i * 2 },
			want:    []int{}, // Expect empty, non-nil
		},
		{
			name:    "NilInput",
			input:   nil,                          // Test nil input explicitly
			mapFunc: func(i int) int { return i }, // Function type needed for compiler
			want:    ([]int)(nil),                 // << FIX: Expect a typed nil slice
		},
		{
			name:    "IdentityInt",
			input:   []int{10, 20},
			mapFunc: func(i int) int { return i },
			want:    []int{10, 20},
		},
		{
			name:  "IntToInterface",
			input: []int{5, 6},
			mapFunc: func(i int) any { // Map to interface{} / any
				if i%2 == 0 {
					return float64(i)
				}
				return strconv.Itoa(i)
			},
			want: []any{"5", float64(6)}, // Note: Expecting 'any' slice
		},
		{
			name:  "PointerToStringValue",                 // Example mapping pointer elements
			input: []*int{ptr(10), ptr(20), nil, ptr(30)}, // Use ptr helper
			mapFunc: func(p *int) string {
				if p == nil {
					return "nil"
				}
				return strconv.Itoa(*p)
			},
			want: []string{"10", "20", "nil", "30"},
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- Type Assertions and Explicit Instantiation ---
			// This section uses type assertions to call the correct generic instantiation.
			var got any
			switch fn := tc.mapFunc.(type) {
			case func(int) string:
				in, ok := tc.input.([]int)
				// Special handling for nil input case before assertion
				if tc.input == nil {
					in = nil  // Ensure 'in' is nil if tc.input was nil
					ok = true // Allow nil through
				} else if !ok {
					t.Fatalf("Input type mismatch for func(int) string test case")
				}
				got = functional.Map[int, string](in, fn) // Explicit: T=int, U=string
			case func(string) int:
				in, ok := tc.input.([]string)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for func(string) int test case")
				}
				got = functional.Map[string, int](in, fn) // Explicit: T=string, U=int
			case func(person) string:
				in, ok := tc.input.([]person)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for func(person) string test case")
				}
				got = functional.Map[person, string](in, fn) // Explicit: T=person, U=string
			case func(int) int:
				in, ok := tc.input.([]int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for func(int) int test case")
				}
				got = functional.Map[int, int](in, fn) // Explicit: T=int, U=int
			case func(int) any:
				in, ok := tc.input.([]int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for func(int) any test case")
				}
				got = functional.Map[int, any](in, fn) // Explicit: T=int, U=any
			case func(*int) string:
				in, ok := tc.input.([]*int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for func(*int) string test case")
				}
				got = functional.Map[*int, string](in, fn) // Explicit: T=*int, U=string

			default:
				t.Fatalf("Unhandled mapFunc type in test setup: %T", tc.mapFunc)
			}

			// --- Comparison ---
			// Use reflect.DeepEqual for comparing slices, handles nil correctly.
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Map(%T input) = %v (%[3]T), want %v (%[4]T)", tc.input, got, got, tc.want, tc.want)
			}

			// The issue is with this check - got is not nil, it's an empty slice
			// Let's remove this check since we're already verifying the result with DeepEqual above
			// and we're now expecting empty slices for nil inputs, not nil outputs

			if tc.input != nil && got == nil {
				t.Errorf("Map(non-nil input) returned nil output")
			}
		})
	}
}

// Remove duplicate ptr helper function - already defined in find_test.go

// --- Go Example ---

func ExampleMap() {
	// Example 1: Doubling integers
	numbers := []int{1, 2, 3, 4}
	doubled := functional.Map[int, int](numbers, func(n int) int {
		return n * 2
	})
	fmt.Println("Doubled:", doubled)

	// Example 2: Converting integers to strings
	intSlice := []int{5, 6, 7}
	stringSlice := functional.Map[int, string](intSlice, func(n int) string {
		return fmt.Sprintf("Num-%d", n)
	})
	fmt.Println("Strings:", stringSlice)

	// Example 3: Getting names from structs
	people := []person{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}
	names := functional.Map[person, string](people, func(p person) string {
		return p.Name
	})
	fmt.Println("Names:", names)

	// Example 4: Nil input
	var nilInts []int = nil
	mappedNil := functional.Map[int, string](nilInts, func(n int) string { return "x" }) // Function type needed
	fmt.Println("Mapped nil:", mappedNil == nil)                                         // Check if it's actually nil

	// Output:
	// Doubled: [2 4 6 8]
	// Strings: [Num-5 Num-6 Num-7]
	// Names: [Alice Bob]
	// Mapped nil: true
}

// --- Benchmarks ---

// Helper function to generate a slice of ints for benchmarking
func generateIntSlice(size int) []int {
	slice := make([]int, size)
	for i := 0; i < size; i++ {
		slice[i] = i
	}
	return slice
}

// mapper function for benchmarks
var intToStringMapper = func(i int) string {
	return strconv.Itoa(i * 2) // Example transformation
}

// Benchmark for functional.Map
func benchmarkMapGeneric(size int, b *testing.B) {
	inputSlice := generateIntSlice(size)
	b.ResetTimer() // Start timing after setup
	for i := 0; i < b.N; i++ {
		// Assign to a local variable to prevent compiler optimization
		// from removing the function call entirely.
		_ = functional.Map(inputSlice, intToStringMapper)
	}
}

// Benchmark for traditional for loop map
func benchmarkMapLoop(size int, b *testing.B) {
	inputSlice := generateIntSlice(size)
	mapper := intToStringMapper // Use the same mapper
	b.ResetTimer()              // Start timing after setup
	for i := 0; i < b.N; i++ {
		// Manual loop implementation
		result := make([]string, len(inputSlice)) // Preallocate slice
		for j, val := range inputSlice {
			result[j] = mapper(val)
		}
		// Assign to prevent optimization
		_ = result
	}
}

// --- Run Benchmarks for different sizes ---

func BenchmarkMap_Generic_10(b *testing.B) { benchmarkMapGeneric(10, b) }
func BenchmarkMap_Loop_10(b *testing.B)    { benchmarkMapLoop(10, b) }

func BenchmarkMap_Generic_100(b *testing.B) { benchmarkMapGeneric(100, b) }
func BenchmarkMap_Loop_100(b *testing.B)    { benchmarkMapLoop(100, b) }

func BenchmarkMap_Generic_1000(b *testing.B) { benchmarkMapGeneric(1000, b) }
func BenchmarkMap_Loop_1000(b *testing.B)    { benchmarkMapLoop(1000, b) }

func BenchmarkMap_Generic_10000(b *testing.B) { benchmarkMapGeneric(10000, b) }
func BenchmarkMap_Loop_10000(b *testing.B)    { benchmarkMapLoop(10000, b) }

// Add larger sizes if needed, e.g., 100k, 1M
