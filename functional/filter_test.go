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
