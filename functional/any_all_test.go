package functional_test

import (
	"fmt"
	"testing" // Added testing package import

	"github.com/JackovAlltrades/go-generics/functional" // Adjust import path
)

// --- Test Any ---
func TestAny(t *testing.T) {
	testCases := []struct {
		name      string
		input     any // Use 'any' for testing different types
		predicate any // The predicate function
		want      bool
	}{
		{
			name:  "Ints_AnyEven_True",
			input: []int{1, 3, 5, 6, 7},
			predicate: func(i int) bool {
				return i%2 == 0
			},
			want: true,
		},
		{
			name:  "Ints_AnyEven_False",
			input: []int{1, 3, 5, 7, 9},
			predicate: func(i int) bool {
				return i%2 == 0
			},
			want: false,
		},
		{
			name:  "Strings_AnyLongerThan3_True",
			input: []string{"a", "bb", "cccc", "d"},
			predicate: func(s string) bool {
				return len(s) > 3
			},
			want: true,
		},
		{
			name:  "Strings_AnyLongerThan3_False",
			input: []string{"a", "bb", "ccc", "d"},
			predicate: func(s string) bool {
				return len(s) > 3
			},
			want: false,
		},
		{
			name:  "EmptySlice_ReturnsFalse",
			input: []int{},
			predicate: func(i int) bool {
				return true // Predicate doesn't matter for empty
			},
			want: false,
		},
		{
			name:  "NilSlice_ReturnsFalse",
			input: ([]string)(nil),
			predicate: func(s string) bool {
				return true // Predicate doesn't matter for nil
			},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got bool
			// Type switching to call the correct generic instantiation
			switch pred := tc.predicate.(type) {
			case func(int) bool:
				in, ok := tc.input.([]int)
				if !ok && tc.input != nil { // Allow nil input
					t.Fatalf("Input type mismatch for func(int) bool")
				}
				got = functional.Any(in, pred)
			case func(string) bool:
				in, ok := tc.input.([]string)
				if !ok && tc.input != nil { // Allow nil input
					t.Fatalf("Input type mismatch for func(string) bool")
				}
				got = functional.Any(in, pred)
			default:
				t.Fatalf("Unhandled predicate type in test setup: %T", tc.predicate)
			}

			if got != tc.want {
				t.Errorf("Any() = %v, want %v", got, tc.want)
			}
		})
	}
}

// --- Test All ---
func TestAll(t *testing.T) {
	testCases := []struct {
		name      string
		input     any
		predicate any
		want      bool
	}{
		{
			name:  "Ints_AllEven_False",
			input: []int{2, 4, 5, 6, 8},
			predicate: func(i int) bool {
				return i%2 == 0
			},
			want: false,
		},
		{
			name:  "Ints_AllEven_True",
			input: []int{2, 4, 6, 8},
			predicate: func(i int) bool {
				return i%2 == 0
			},
			want: true,
		},
		{
			name:  "Strings_AllShorterThan5_False",
			input: []string{"a", "bb", "ccccc", "d"},
			predicate: func(s string) bool {
				return len(s) < 5
			},
			want: false,
		},
		{
			name:  "Strings_AllShorterThan5_True",
			input: []string{"a", "bb", "cccc", "d"},
			predicate: func(s string) bool {
				return len(s) < 5
			},
			want: true,
		},
		{
			name:  "EmptySlice_ReturnsTrue",
			input: []int{},
			predicate: func(i int) bool {
				return false // Predicate doesn't matter, should be vacuously true
			},
			want: true, // All elements (of which there are none) satisfy the condition
		},
		{
			name:  "NilSlice_ReturnsTrue",
			input: ([]string)(nil),
			predicate: func(s string) bool {
				return false // Predicate doesn't matter, should be vacuously true
			},
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got bool
			switch pred := tc.predicate.(type) {
			case func(int) bool:
				in, ok := tc.input.([]int)
				if !ok && tc.input != nil {
					t.Fatalf("Input type mismatch for func(int) bool")
				}
				got = functional.All(in, pred)
			case func(string) bool:
				in, ok := tc.input.([]string)
				if !ok && tc.input != nil {
					t.Fatalf("Input type mismatch for func(string) bool")
				}
				got = functional.All(in, pred)
			default:
				t.Fatalf("Unhandled predicate type in test setup: %T", tc.predicate)
			}

			if got != tc.want {
				t.Errorf("All() = %v, want %v", got, tc.want)
			}
		})
	}
}

// --- Examples ---

func ExampleAny() {
	numbers := []int{1, 3, 5, 7, 8, 9}
	hasEven := functional.Any(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Println("Any even number?", hasEven)

	words := []string{"quick", "brown", "fox"}
	hasLongWord := functional.Any(words, func(s string) bool {
		return len(s) > 5
	})
	fmt.Println("Any word longer than 5 chars?", hasLongWord)

	empty := []int{}
	anyInEmpty := functional.Any(empty, func(n int) bool { return true })
	fmt.Println("Any in empty?", anyInEmpty)

	// Output:
	// Any even number? true
	// Any word longer than 5 chars? false
	// Any in empty? false
}

func ExampleAll() {
	allPositive := []int{1, 5, 10, 2}
	areAllPositive := functional.All(allPositive, func(n int) bool {
		return n > 0
	})
	fmt.Println("All positive?", areAllPositive)

	someNegative := []int{1, 5, -10, 2}
	areAllPositive2 := functional.All(someNegative, func(n int) bool {
		return n > 0
	})
	fmt.Println("All positive (with negative)?", areAllPositive2)

	allShort := []string{"cat", "dog", "bat"}
	areAllShort := functional.All(allShort, func(s string) bool {
		return len(s) < 4
	})
	fmt.Println("All shorter than 4 chars?", areAllShort)

	empty := []int{}
	allInEmpty := functional.All(empty, func(n int) bool { return false }) // Vacuously true
	fmt.Println("All in empty?", allInEmpty)

	// Output:
	// All positive? true
	// All positive (with negative)? false
	// All shorter than 4 chars? true
	// All in empty? true
}

// --- Benchmarks ---

// Helper to generate slice of ints for benchmarks
func generateIntSliceAnyAll(size int) []int {
	if size <= 0 {
		return []int{}
	}
	slice := make([]int, size)
	for i := 0; i < size; i++ {
		slice[i] = i // 0, 1, 2, ...
	}
	return slice
}

// === Any Benchmarks ===

// Predicate: checks if a number is > size/2. Met late or never.
func isLargeNumberPredicate(threshold int) func(int) bool {
	return func(n int) bool {
		return n > threshold // Condition met later in the 0..size-1 sequence
	}
}

// Predicate: checks if a number is == 5. Met early.
var isFivePredicate = func(n int) bool {
	return n == 5 // Condition met early in the 0..size-1 sequence
}

// Benchmark for functional.Any - Condition met early
func benchmarkAnyGenericTrueEarly(size int, b *testing.B) {
	if size == 0 {
		b.Skip("Skipping size 0")
		return
	} // Predicate assumes non-empty
	inputSlice := generateIntSliceAnyAll(size)
	predicate := isFivePredicate // Found at index 5 (ensure size >= 6 for this to be meaningful)
	b.ResetTimer()
	var result bool
	for i := 0; i < b.N; i++ {
		result = functional.Any(inputSlice, predicate)
	}
	_ = result
}

// Benchmark for loop Any - Condition met early
func benchmarkAnyLoopTrueEarly(size int, b *testing.B) {
	if size == 0 {
		b.Skip("Skipping size 0")
		return
	}
	inputSlice := generateIntSliceAnyAll(size)
	predicate := isFivePredicate // Found at index 5
	b.ResetTimer()
	var result bool
	for i := 0; i < b.N; i++ {
		localResult := false
		for _, v := range inputSlice {
			if predicate(v) {
				localResult = true
				break // Exit loop early
			}
		}
		result = localResult
	}
	_ = result
}

// Benchmark for functional.Any - Condition never met (Worst Case)
func benchmarkAnyGenericFalse(size int, b *testing.B) {
	inputSlice := generateIntSliceAnyAll(size) // 0..size-1
	// Predicate that will always be false for the generated data
	predicate := func(n int) bool { return n < 0 }
	b.ResetTimer()
	var result bool
	for i := 0; i < b.N; i++ {
		result = functional.Any(inputSlice, predicate)
	}
	_ = result
}

// Benchmark for loop Any - Condition never met (Worst Case)
func benchmarkAnyLoopFalse(size int, b *testing.B) {
	inputSlice := generateIntSliceAnyAll(size)
	predicate := func(n int) bool { return n < 0 }
	b.ResetTimer()
	var result bool
	for i := 0; i < b.N; i++ {
		localResult := false
		for _, v := range inputSlice {
			if predicate(v) {
				localResult = true
				break
			}
		}
		result = localResult // Will always be false here
	}
	_ = result
}

// === All Benchmarks ===

// Predicate: checks if a number is >= 0. Always true for generated data.
var isNonNegativePredicate = func(n int) bool {
	return n >= 0
}

// Predicate: checks if a number is < size/2. Fails late.
func isSmallNumberPredicate(threshold int) func(int) bool {
	return func(n int) bool {
		return n < threshold // Condition fails later in the 0..size-1 sequence
	}
}

// Benchmark for functional.All - Condition always true (Worst Case)
func benchmarkAllGenericTrue(size int, b *testing.B) {
	inputSlice := generateIntSliceAnyAll(size)
	predicate := isNonNegativePredicate // Always true for 0..size-1
	b.ResetTimer()
	var result bool
	for i := 0; i < b.N; i++ {
		result = functional.All(inputSlice, predicate)
	}
	_ = result
}

// Benchmark for loop All - Condition always true (Worst Case)
func benchmarkAllLoopTrue(size int, b *testing.B) {
	inputSlice := generateIntSliceAnyAll(size)
	predicate := isNonNegativePredicate
	b.ResetTimer()
	var result bool
	for i := 0; i < b.N; i++ {
		localResult := true // Assume true initially
		for _, v := range inputSlice {
			if !predicate(v) {
				localResult = false
				break // Exit loop early on failure
			}
		}
		result = localResult // Will always be true here
	}
	_ = result
}

// Benchmark for functional.All - Condition fails early
func benchmarkAllGenericFalseEarly(size int, b *testing.B) {
	if size < 6 {
		b.Skip("Skipping small size for FalseEarly test")
		return
	} // Ensure index 5 exists
	// Generate slightly different data so failure isn't always index 0
	inputSlice := make([]int, size)
	for i := 0; i < size; i++ {
		inputSlice[i] = i + 5
	} // 5, 6, 7...
	predicate := func(n int) bool { return n < 10 } // Fails around index 5 (when n becomes 10)
	b.ResetTimer()
	var result bool
	for i := 0; i < b.N; i++ {
		result = functional.All(inputSlice, predicate)
	}
	_ = result
}

// Benchmark for loop All - Condition fails early
func benchmarkAllLoopFalseEarly(size int, b *testing.B) {
	if size < 6 {
		b.Skip("Skipping small size for FalseEarly test")
		return
	}
	inputSlice := make([]int, size)
	for i := 0; i < size; i++ {
		inputSlice[i] = i + 5
	} // 5, 6, 7...
	predicate := func(n int) bool { return n < 10 } // Fails around index 5
	b.ResetTimer()
	var result bool
	for i := 0; i < b.N; i++ {
		localResult := true
		for _, v := range inputSlice {
			if !predicate(v) {
				localResult = false
				break
			}
		}
		result = localResult
	}
	_ = result
}

// --- Run Benchmarks for different sizes ---

// N = Number of elements in the slice
const (
	N1_AnyAll = 100 // Using different const name to avoid collision if copied elsewhere
	N2_AnyAll = 10000
)

// Any Benchmarks
func BenchmarkAny_Generic_TrueEarly_N100(b *testing.B) { benchmarkAnyGenericTrueEarly(N1_AnyAll, b) }
func BenchmarkAny_Loop_TrueEarly_N100(b *testing.B)    { benchmarkAnyLoopTrueEarly(N1_AnyAll, b) }
func BenchmarkAny_Generic_False_N100(b *testing.B)     { benchmarkAnyGenericFalse(N1_AnyAll, b) }
func BenchmarkAny_Loop_False_N100(b *testing.B)        { benchmarkAnyLoopFalse(N1_AnyAll, b) }

func BenchmarkAny_Generic_TrueEarly_N10000(b *testing.B) { benchmarkAnyGenericTrueEarly(N2_AnyAll, b) }
func BenchmarkAny_Loop_TrueEarly_N10000(b *testing.B)    { benchmarkAnyLoopTrueEarly(N2_AnyAll, b) }
func BenchmarkAny_Generic_False_N10000(b *testing.B)     { benchmarkAnyGenericFalse(N2_AnyAll, b) }
func BenchmarkAny_Loop_False_N10000(b *testing.B)        { benchmarkAnyLoopFalse(N2_AnyAll, b) }

// All Benchmarks
func BenchmarkAll_Generic_True_N100(b *testing.B)       { benchmarkAllGenericTrue(N1_AnyAll, b) }
func BenchmarkAll_Loop_True_N100(b *testing.B)          { benchmarkAllLoopTrue(N1_AnyAll, b) }
func BenchmarkAll_Generic_FalseEarly_N100(b *testing.B) { benchmarkAllGenericFalseEarly(N1_AnyAll, b) }
func BenchmarkAll_Loop_FalseEarly_N100(b *testing.B)    { benchmarkAllLoopFalseEarly(N1_AnyAll, b) }

func BenchmarkAll_Generic_True_N10000(b *testing.B) { benchmarkAllGenericTrue(N2_AnyAll, b) }
func BenchmarkAll_Loop_True_N10000(b *testing.B)    { benchmarkAllLoopTrue(N2_AnyAll, b) }
func BenchmarkAll_Generic_FalseEarly_N10000(b *testing.B) {
	benchmarkAllGenericFalseEarly(N2_AnyAll, b)
}
func BenchmarkAll_Loop_FalseEarly_N10000(b *testing.B) { benchmarkAllLoopFalseEarly(N2_AnyAll, b) }
