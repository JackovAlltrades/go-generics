package functional_test

import (
	"fmt" // Needed for Example
	"reflect"
	"strings"
	"testing"

	"github.com/JackovAlltrades/go-generics/functional" // Use variable here
)

// TestFilter runs table-driven tests for the Filter function.
func TestFilter(t *testing.T) {
	// Define helper types for tests
	type person struct {
		Name string
		Age  int
	}

	// Test integer filtering
	t.Run("Filter integers", func(t *testing.T) {
		// Create a type-specific instantiation
		filterInt := functional.Filter[int]

		input := []int{1, 8, 3, 10, 5, 6}
		expected := []int{8, 10, 6}
		result := filterInt(input, func(n int) bool { return n > 5 })
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Filter() = %v, want %v", result, expected)
		}

		// Test empty slice
		emptyInput := []int{}
		emptyResult := filterInt(emptyInput, func(n int) bool { return n > 0 })
		if len(emptyResult) != 0 {
			t.Errorf("Filter() on empty slice = %v, want empty slice", emptyResult)
		}

		// Test nil slice
		var nilInput []int
		nilResult := filterInt(nilInput, func(n int) bool { return n > 0 })
		if nilResult == nil || len(nilResult) != 0 {
			t.Errorf("Filter() on nil slice = %v, want empty non-nil slice", nilResult)
		}

		// Test all elements match
		allMatch := []int{10, 20, 30}
		allMatchResult := filterInt(allMatch, func(n int) bool { return n > 5 })
		if !reflect.DeepEqual(allMatchResult, allMatch) {
			t.Errorf("Filter() all match = %v, want %v", allMatchResult, allMatch)
		}

		// Test no elements match
		noMatch := []int{1, 2, 3}
		noMatchResult := filterInt(noMatch, func(n int) bool { return n > 5 })
		if len(noMatchResult) != 0 {
			t.Errorf("Filter() no match = %v, want empty slice", noMatchResult)
		}
	})

	// Test string filtering
	t.Run("Filter strings", func(t *testing.T) {
		// Create a type-specific instantiation
		filterString := functional.Filter[string]

		input := []string{"apple", "banana", "avocado", "apricot", "grape"}
		expected := []string{"apple", "avocado", "apricot"}
		result := filterString(input, func(s string) bool { return strings.HasPrefix(s, "a") })
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Filter() = %v, want %v", result, expected)
		}
	})

	// Test struct filtering
	t.Run("Filter structs", func(t *testing.T) {
		// Create a type-specific instantiation
		filterPerson := functional.Filter[person]

		input := []person{{"Alice", 30}, {"Bob", 45}, {"Charlie", 25}, {"David", 31}}
		expected := []person{{"Bob", 45}, {"David", 31}}
		result := filterPerson(input, func(p person) bool { return p.Age > 30 })
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Filter() = %v, want %v", result, expected)
		}
	})

	// Test pointer filtering
	t.Run("Filter pointers", func(t *testing.T) {
		// Create a type-specific instantiation
		filterStringPtr := functional.Filter[*string]

		input := []*string{strPtr("a"), nil, strPtr("b"), nil, strPtr("c")}
		expected := []*string{strPtr("a"), strPtr("b"), strPtr("c")}
		result := filterStringPtr(input, func(p *string) bool { return p != nil })

		// Need special comparison for pointers
		if len(result) != len(expected) {
			t.Errorf("Filter() length = %d, want %d", len(result), len(expected))
		}

		for i, v := range result {
			if v == nil || *v != *expected[i] {
				t.Errorf("Filter() at index %d = %v, want %v", i, v, expected[i])
			}
		}
	})
}

// Helper to create string pointers for tests
func strPtr(s string) *string {
	return &s
}

// BenchmarkFilter provides a basic benchmark for the Filter operation.
func BenchmarkFilter(b *testing.B) {
	size := 10000
	input := make([]int, size)
	for i := 0; i < size; i++ {
		input[i] = i
	}

	// Create a type-specific instantiation
	filterInt := functional.Filter[int]
	isEven := func(n int) bool { return n%2 == 0 }

	b.ResetTimer() // Start timing after setup
	for i := 0; i < b.N; i++ {
		_ = filterInt(input, isEven) // Call the function under test
	}
}

// Example usage shown as a testable example in Go documentation.
// ExampleFilter demonstrates filtering integers with the Filter function.
func ExampleFilter() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	isEven := func(n int) bool { return n%2 == 0 }

	// Create a type-specific instantiation
	filterInt := functional.Filter[int]
	evenNumbers := filterInt(numbers, isEven)
	fmt.Println(evenNumbers)
	// Output: [2 4 6 8 10]
}

// ExampleFilter_strings demonstrates filtering strings with the Filter function.
func ExampleFilter_strings() {
	words := []string{"apple", "banana", "apricot", "grape", "avocado"}
	startsWithA := func(s string) bool { return strings.HasPrefix(s, "a") }

	// Create a type-specific instantiation
	filterString := functional.Filter[string]
	aWords := filterString(words, startsWithA)
	fmt.Println(aWords)
	// Output: [apple apricot avocado]
}

// --- Benchmarks ---

// Re-use helper from map_test.go conceptually - generate a slice of ints
// (Can be copy-pasted if needed, or assume available if running ./...)
func generateIntSliceBenchFilter(size int) []int {
	slice := make([]int, size)
	for i := 0; i < size; i++ {
		slice[i] = i // Simple predictable data
	}
	return slice
}

// Predicate function for benchmarks (e.g., filter even numbers)
var isEvenPredicate = func(i int) bool {
	return i%2 == 0
}

// Benchmark for functional.Filter
func benchmarkFilterGeneric(size int, b *testing.B) {
	inputSlice := generateIntSliceBenchFilter(size)
	predicate := isEvenPredicate // Use the same predicate
	b.ResetTimer()               // Start timing after setup
	for i := 0; i < b.N; i++ {
		// Assign to a local variable to prevent compiler optimization
		_ = functional.Filter(inputSlice, predicate)
	}
}

// Benchmark for traditional for loop filter
func benchmarkFilterLoop(size int, b *testing.B) {
	inputSlice := generateIntSliceBenchFilter(size)
	predicate := isEvenPredicate // Use the same predicate
	b.ResetTimer()               // Start timing after setup
	for i := 0; i < b.N; i++ {
		// Manual loop implementation (append style)
		// Minimal preallocation is common for Filter as size is unknown
		result := make([]int, 0) // Start with zero capacity or small hint
		for _, val := range inputSlice {
			if predicate(val) {
				result = append(result, val)
			}
		}
		// Assign to prevent optimization
		_ = result
	}
}

// Benchmark for traditional for loop filter with preallocation guess
// (Less common unless you have a good estimate of filter rate)
func benchmarkFilterLoopPrealloc(size int, b *testing.B) {
	inputSlice := generateIntSliceBenchFilter(size)
	predicate := isEvenPredicate // Use the same predicate
	b.ResetTimer()               // Start timing after setup
	for i := 0; i < b.N; i++ {
		// Manual loop implementation (preallocation guess style)
		// Guessing half the size might be reasonable for evens.
		result := make([]int, 0, size/2+1) // Preallocate with guess
		for _, val := range inputSlice {
			if predicate(val) {
				result = append(result, val)
			}
		}
		// Assign to prevent optimization
		_ = result
	}
}

// --- Run Benchmarks for different sizes ---

func BenchmarkFilter_Generic_10(b *testing.B)      { benchmarkFilterGeneric(10, b) }
func BenchmarkFilter_Loop_10(b *testing.B)         { benchmarkFilterLoop(10, b) }
func BenchmarkFilter_LoopPrealloc_10(b *testing.B) { benchmarkFilterLoopPrealloc(10, b) }

func BenchmarkFilter_Generic_100(b *testing.B)      { benchmarkFilterGeneric(100, b) }
func BenchmarkFilter_Loop_100(b *testing.B)         { benchmarkFilterLoop(100, b) }
func BenchmarkFilter_LoopPrealloc_100(b *testing.B) { benchmarkFilterLoopPrealloc(100, b) }

func BenchmarkFilter_Generic_1000(b *testing.B)      { benchmarkFilterGeneric(1000, b) }
func BenchmarkFilter_Loop_1000(b *testing.B)         { benchmarkFilterLoop(1000, b) }
func BenchmarkFilter_LoopPrealloc_1000(b *testing.B) { benchmarkFilterLoopPrealloc(1000, b) }

func BenchmarkFilter_Generic_10000(b *testing.B)      { benchmarkFilterGeneric(10000, b) }
func BenchmarkFilter_Loop_10000(b *testing.B)         { benchmarkFilterLoop(10000, b) }
func BenchmarkFilter_LoopPrealloc_10000(b *testing.B) { benchmarkFilterLoopPrealloc(10000, b) }

// Consider adding larger sizes (100k, 1M) if performance at scale is critical
