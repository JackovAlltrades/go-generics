package functional_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/JackovAlltrades/go-generics/functional" // Adjust import path if needed
)

// Re-using person struct from map_test for convenience, or define locally if preferred
// type person struct {
// 	Name string
// 	Age  int
// }

func TestReduce(t *testing.T) {
	testCases := []struct {
		name     string // Name of the subtest
		input    any    // Input slice (use 'any' for type flexibility)
		initial  any    // Initial value for accumulator
		reduceFn any    // Reducing function (use 'any')
		want     any    // Expected final accumulated value
	}{
		{
			name:    "SumIntegers",
			input:   []int{1, 2, 3, 4, 5},
			initial: 0, // Start sum at 0
			reduceFn: func(acc int, next int) int {
				return acc + next
			},
			want: 15, // 1+2+3+4+5
		},
		{
			name:    "SumIntegersWithInitial",
			input:   []int{1, 2, 3},
			initial: 10, // Start sum at 10
			reduceFn: func(acc int, next int) int {
				return acc + next
			},
			want: 16, // 10+1+2+3
		},
		{
			name:    "ProductFloats",
			input:   []float64{1.5, 2.0, 4.0},
			initial: 1.0,
			reduceFn: func(acc float64, next float64) float64 {
				return acc * next
			},
			want: 12.0, // 1.0 * 1.5 * 2.0 * 4.0
		},
		{
			name:    "ConcatenateStrings",
			input:   []string{"a", "b", "c"},
			initial: "",
			reduceFn: func(acc string, next string) string {
				return acc + next
			},
			want: "abc",
		},
		{
			name:    "ConcatenateStringsWithSeparator",
			input:   []string{"go", "is", "fun"},
			initial: "",
			reduceFn: func(acc string, next string) string {
				if acc == "" { // Don't prepend separator for the first element
					return next
				}
				return acc + "-" + next
			},
			want: "go-is-fun",
		},
		{
			name:    "CountEvenNumbers",
			input:   []int{1, 2, 3, 4, 5, 6},
			initial: 0, // Initial count
			reduceFn: func(count int, next int) int {
				if next%2 == 0 {
					return count + 1
				}
				return count
			},
			want: 3, // 2, 4, 6 are even
		},
		{
			name:    "EmptyInput",
			input:   []int{}, // Empty slice
			initial: 100,     // Should return initial value
			reduceFn: func(acc int, next int) int {
				t.Fatal("Reduce function should not be called for empty input") // Safety check
				return acc + next
			},
			want: 100,
		},
		{
			name:    "NilInput",
			input:   ([]int)(nil), // Typed nil slice
			initial: "default",    // Should return initial value
			reduceFn: func(acc string, next int) string {
				t.Fatal("Reduce function should not be called for nil input") // Safety check
				return acc + fmt.Sprint(next)
			},
			want: "default",
		},
		{
			name:    "StructsToTotalAge",
			input:   []person{{Name: "A", Age: 20}, {Name: "B", Age: 35}},
			initial: 0,
			reduceFn: func(totalAge int, p person) int {
				return totalAge + p.Age
			},
			want: 55,
		},
		{
			name:    "IntsToMapGroupByEvenOdd",
			input:   []int{1, 2, 3, 4, 5, 6},
			initial: map[string][]int{"even": {}, "odd": {}}, // Start with empty map structure
			reduceFn: func(groups map[string][]int, n int) map[string][]int {
				if n%2 == 0 {
					groups["even"] = append(groups["even"], n)
				} else {
					groups["odd"] = append(groups["odd"], n)
				}
				return groups // Return the modified map
			},
			// Use DeepEqual for maps
			want: map[string][]int{"even": {2, 4, 6}, "odd": {1, 3, 5}},
		},
		{
			name:    "PointersSum",
			input:   []*int{ptr(10), ptr(20), nil, ptr(30)}, // Use ptr helper from map_test
			initial: 0,
			reduceFn: func(acc int, p *int) int {
				if p != nil { // Handle nil pointers in the data
					return acc + *p
				}
				return acc
			},
			want: 60, // 10 + 20 + 30
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any

			// Type switch to call the correct instantiation of Reduce
			// This gets more complex as input (T) and accumulator (U) types vary
			switch fn := tc.reduceFn.(type) {
			case func(int, int) int: // T=int, U=int
				in, ok := tc.input.([]int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for func(int, int) int")
				}
				init, okI := tc.initial.(int)
				if !okI {
					t.Fatalf("Initial value type mismatch for func(int, int) int")
				}
				got = functional.Reduce[int, int](in, init, fn)
			case func(float64, float64) float64: // T=float64, U=float64
				in, ok := tc.input.([]float64)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for func(float64, float64) float64")
				}
				init, okI := tc.initial.(float64)
				if !okI {
					t.Fatalf("Initial value type mismatch for func(float64, float64) float64")
				}
				got = functional.Reduce[float64, float64](in, init, fn)
			case func(string, string) string: // T=string, U=string
				in, ok := tc.input.([]string)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for func(string, string) string")
				}
				init, okI := tc.initial.(string)
				if !okI {
					t.Fatalf("Initial value type mismatch for func(string, string) string")
				}
				got = functional.Reduce[string, string](in, init, fn)
			case func(int, person) int: // T=person, U=int
				in, ok := tc.input.([]person)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for func(int, person) int")
				}
				init, okI := tc.initial.(int)
				if !okI {
					t.Fatalf("Initial value type mismatch for func(int, person) int")
				}
				got = functional.Reduce[person, int](in, init, fn) // Note T=person, U=int
			case func(map[string][]int, int) map[string][]int: // T=int, U=map[string][]int
				in, ok := tc.input.([]int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for map group test")
				}
				init, okI := tc.initial.(map[string][]int)
				if !okI {
					t.Fatalf("Initial value type mismatch for map group test")
				}
				got = functional.Reduce[int, map[string][]int](in, init, fn)
			case func(string, int) string: // T=int, U=string (for nil string test)
				in, ok := tc.input.([]int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for nil string test")
				}
				init, okI := tc.initial.(string)
				if !okI {
					t.Fatalf("Initial value type mismatch for nil string test")
				}
				got = functional.Reduce[int, string](in, init, fn)
			case func(int, *int) int: // T=*int, U=int
				in, ok := tc.input.([]*int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for pointer sum test")
				}
				init, okI := tc.initial.(int)
				if !okI {
					t.Fatalf("Initial value type mismatch for pointer sum test")
				}
				got = functional.Reduce[*int, int](in, init, fn)

			default:
				t.Fatalf("Unhandled reduceFn type in test setup: %T", tc.reduceFn)
			}

			// Use DeepEqual for robust comparison, especially for slices/maps/structs
			if !reflect.DeepEqual(got, tc.want) {
				// Provide helpful output on failure
				t.Errorf("Reduce() = %#v (%[2]T), want %#v (%[3]T)", got, got, tc.want, tc.want)
			}
		})
	}
}

// --- Go Example ---

func ExampleReduce() {
	// Example 1: Summing integers
	numbers := []int{1, 2, 3, 4, 5}
	sum := functional.Reduce[int, int](numbers, 0, func(acc, next int) int {
		return acc + next
	})
	fmt.Println("Sum:", sum)

	// Example 2: Joining strings
	words := []string{"Reduce", "is", "useful"}
	sentence := functional.Reduce[string, string](words, "", func(acc, next string) string {
		if acc == "" {
			return next
		}
		return acc + " " + next
	})
	fmt.Println("Sentence:", sentence)

	// Example 3: Finding the maximum value
	scores := []int{70, 95, 88, 65, 95}
	// Start with a value guaranteed to be <= the first element,
	// or handle the empty case explicitly if needed (Reduce handles empty)
	maxScore := functional.Reduce[int, int](scores, scores[0], func(currentMax, next int) int {
		if next > currentMax {
			return next
		}
		return currentMax
	})
	// Note: This assumes scores is not empty. A safer version would use math.MinInt
	// or handle the empty case separately if Reduce didn't already return initial.
	fmt.Println("Max Score:", maxScore)

	// Example 4: Handling nil slice (returns initial value)
	var nilSlice []float64 = nil
	initialValue := 100.0
	resultFromNil := functional.Reduce[float64, float64](nilSlice, initialValue, func(acc, next float64) float64 {
		// This function won't be called
		return acc + next
	})
	fmt.Println("Result from nil:", resultFromNil)

	// Output:
	// Sum: 15
	// Sentence: Reduce is useful
	// Max Score: 95
	// Result from nil: 100
}

// --- Benchmarks ---

// Helper function to generate a slice of ints (can reuse from other files)
func generateIntSliceBenchReduce(size int) []int {
	slice := make([]int, size)
	for i := 0; i < size; i++ {
		slice[i] = i // Use simple values
	}
	return slice
}

// Reducer function for benchmarks (simple integer sum)
var intSumReducer = func(acc int, i int) int {
	return acc + i
}

// Benchmark for functional.Reduce
func benchmarkReduceGeneric(size int, b *testing.B) {
	inputSlice := generateIntSliceBenchReduce(size)
	reducer := intSumReducer
	initial := 0
	b.ResetTimer() // Start timing after setup
	var result int // Declare outside loop to avoid measuring declaration time
	for i := 0; i < b.N; i++ {
		// Assign to prevent optimization
		result = functional.Reduce(inputSlice, initial, reducer)
	}
	_ = result // Use result after loop
}

// Benchmark for traditional for loop reduce
func benchmarkReduceLoop(size int, b *testing.B) {
	inputSlice := generateIntSliceBenchReduce(size)
	reducer := intSumReducer // Use same logic for fair comparison
	initial := 0
	b.ResetTimer() // Start timing after setup
	var result int // Declare outside loop
	for i := 0; i < b.N; i++ {
		// Manual loop implementation
		accumulator := initial
		for _, val := range inputSlice {
			accumulator = reducer(accumulator, val) // Apply the 'reducer' logic
		}
		// Assign to prevent optimization
		result = accumulator
	}
	_ = result // Use result after loop
}

// --- Run Benchmarks for different sizes ---

func BenchmarkReduce_Generic_10(b *testing.B) { benchmarkReduceGeneric(10, b) }
func BenchmarkReduce_Loop_10(b *testing.B)    { benchmarkReduceLoop(10, b) }

func BenchmarkReduce_Generic_100(b *testing.B) { benchmarkReduceGeneric(100, b) }
func BenchmarkReduce_Loop_100(b *testing.B)    { benchmarkReduceLoop(100, b) }

func BenchmarkReduce_Generic_1000(b *testing.B) { benchmarkReduceGeneric(1000, b) }
func BenchmarkReduce_Loop_1000(b *testing.B)    { benchmarkReduceLoop(1000, b) }

func BenchmarkReduce_Generic_10000(b *testing.B) { benchmarkReduceGeneric(10000, b) }
func BenchmarkReduce_Loop_10000(b *testing.B)    { benchmarkReduceLoop(10000, b) }

// Consider adding larger sizes if needed
