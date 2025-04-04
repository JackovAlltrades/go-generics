package functional_test

import (
	"fmt"
	"testing" // Make sure testing is imported

	"github.com/JackovAlltrades/go-generics/functional"
	// Add if needed by benchmarks later
)

// comparablePerson is a simple struct used in tests, must be comparable.
// Ensure it's defined here so it's accessible within this test file.
// If used across multiple _test.go files, it should ideally be in a
// shared test helper file or defined consistently.
type comparablePerson struct {
	ID   int
	Name string // Use Name field as shown in the map_utils_test example
}

func TestContains(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name      string // Name of the subtest
		input     any    // Input slice (comparable type)
		value     any    // Value to search for
		wantFound bool   // Expected boolean result
	}{
		{
			name:      "ContainsInt_Found",
			input:     []int{1, 2, 3, 4, 5},
			value:     3,
			wantFound: true,
		},
		{
			name:      "ContainsInt_NotFound",
			input:     []int{1, 2, 4, 5},
			value:     3,
			wantFound: false,
		},
		{
			name:      "ContainsString_Found",
			input:     []string{"a", "b", "c"},
			value:     "b",
			wantFound: true,
		},
		{
			name:      "ContainsString_NotFound",
			input:     []string{"a", "b", "c"},
			value:     "d",
			wantFound: false,
		},
		{
			name:      "ContainsString_CaseSensitive",
			input:     []string{"a", "B", "c"},
			value:     "b",
			wantFound: false,
		},
		{
			name: "ContainsComparableStruct_Found",
			input: []comparablePerson{ // Assuming comparablePerson is available
				{ID: 1, Name: "A"},
				{ID: 2, Name: "B"},
			},
			value:     comparablePerson{ID: 2, Name: "B"},
			wantFound: true,
		},
		{
			name: "ContainsComparableStruct_NotFound",
			input: []comparablePerson{
				{ID: 1, Name: "A"},
				{ID: 2, Name: "B"},
			},
			value:     comparablePerson{ID: 3, Name: "C"},
			wantFound: false,
		},
		{
			name:      "EmptyInput",
			input:     []int{},
			value:     1,
			wantFound: false,
		},
		{
			name:      "NilInput",
			input:     ([]string)(nil),
			value:     "a",
			wantFound: false,
		},
		// ... potentially add more test cases ...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotFound bool

			switch in := tc.input.(type) {
			case []int:
				val, ok := tc.value.(int)
				if !ok {
					t.Fatalf("Value type mismatch for []int test")
				}
				gotFound = functional.Contains[int](in, val)
			case []string:
				val, ok := tc.value.(string)
				if !ok {
					t.Fatalf("Value type mismatch for []string test")
				}
				gotFound = functional.Contains[string](in, val)
			case []comparablePerson: // Assuming comparablePerson is available
				val, ok := tc.value.(comparablePerson)
				if !ok {
					t.Fatalf("Value type mismatch for []comparablePerson test")
				}
				gotFound = functional.Contains[comparablePerson](in, val)
			case nil:
				switch val := tc.value.(type) {
				case string:
					gotFound = functional.Contains[string](nil, val)
				case int:
					gotFound = functional.Contains[int](nil, val)
				// Add case for comparablePerson if needed for nil tests
				// case comparablePerson:
				// 	gotFound = functional.Contains[comparablePerson](nil, val)
				default:
					t.Fatalf("Unhandled nil input value type for %s: %T", tc.name, tc.value)
				}
			default:
				t.Fatalf("Unhandled input type in test setup: %T", tc.input)
			}

			if gotFound != tc.wantFound {
				t.Errorf("Contains(%#v, %#v) = %v, want %v", tc.input, tc.value, gotFound, tc.wantFound)
			}
		})
	}
}

// ExampleContains remains the same
func ExampleContains() {
	// Example 1: Contains integer
	numbers := []int{10, 20, 30, 40}
	has_20 := functional.Contains[int](numbers, 20)
	has_50 := functional.Contains[int](numbers, 50)
	fmt.Printf("Numbers %v contain 20? %v\n", numbers, has_20)
	fmt.Printf("Numbers %v contain 50? %v\n", numbers, has_50)

	// Example 2: Contains string
	words := []string{"apple", "banana", "cherry"}
	has_banana := functional.Contains[string](words, "banana")
	has_grape := functional.Contains[string](words, "grape")
	fmt.Printf("Words %v contain 'banana'? %v\n", words, has_banana)
	fmt.Printf("Words %v contain 'grape'? %v\n", words, has_grape)

	// Example 3: Empty slice
	empty := []int{}
	empty_has_1 := functional.Contains[int](empty, 1)
	fmt.Printf("Empty slice contains 1? %v\n", empty_has_1)

	// Example 4: Nil slice
	var nilSlice []string = nil
	nil_has_a := functional.Contains[string](nilSlice, "a")
	fmt.Printf("Nil slice contains 'a'? %v\n", nil_has_a)

	// Output:
	// Numbers [10 20 30 40] contain 20? true
	// Numbers [10 20 30 40] contain 50? false
	// Words [apple banana cherry] contain 'banana'? true
	// Words [apple banana cherry] contain 'grape'? false
	// Empty slice contains 1? false
	// Nil slice contains 'a'? false
}

// --- Benchmarks ---

// Helper function to generate a slice of ints
func generateIntSliceBenchContains(size int) []int {
	slice := make([]int, size)
	for i := 0; i < size; i++ {
		slice[i] = i * 2 // e.g., 0, 2, 4...
	}
	return slice
}

// Benchmark for functional.Contains - Target Found Early
func benchmarkContainsGenericFoundEarly(size int, b *testing.B) {
	inputSlice := generateIntSliceBenchContains(size)
	target := 10 // Typically found early in the generated slice (at index 5)
	b.ResetTimer()
	var found bool
	for i := 0; i < b.N; i++ {
		found = functional.Contains(inputSlice, target)
	}
	_ = found
}

// Benchmark for loop - Target Found Early
func benchmarkContainsLoopFoundEarly(size int, b *testing.B) {
	inputSlice := generateIntSliceBenchContains(size)
	target := 10 // Found early
	b.ResetTimer()
	var found bool
	for i := 0; i < b.N; i++ {
		localFound := false
		for _, val := range inputSlice {
			if val == target {
				localFound = true
				break // Exit loop early once found
			}
		}
		found = localFound
	}
	_ = found
}

// Benchmark for functional.Contains - Target Found Late
func benchmarkContainsGenericFoundLate(size int, b *testing.B) {
	inputSlice := generateIntSliceBenchContains(size)
	target := (size - 1) * 2 // The last element
	b.ResetTimer()
	var found bool
	for i := 0; i < b.N; i++ {
		found = functional.Contains(inputSlice, target)
	}
	_ = found
}

// Benchmark for loop - Target Found Late
func benchmarkContainsLoopFoundLate(size int, b *testing.B) {
	inputSlice := generateIntSliceBenchContains(size)
	target := (size - 1) * 2 // The last element
	b.ResetTimer()
	var found bool
	for i := 0; i < b.N; i++ {
		localFound := false
		for _, val := range inputSlice {
			if val == target {
				localFound = true
				break
			}
		}
		found = localFound
	}
	_ = found
}

// Benchmark for functional.Contains - Target Not Found
func benchmarkContainsGenericNotFound(size int, b *testing.B) {
	inputSlice := generateIntSliceBenchContains(size)
	target := 999999 // Value guaranteed not to be in the slice
	b.ResetTimer()
	var found bool
	for i := 0; i < b.N; i++ {
		found = functional.Contains(inputSlice, target)
	}
	_ = found
}

// Benchmark for loop - Target Not Found
func benchmarkContainsLoopNotFound(size int, b *testing.B) {
	inputSlice := generateIntSliceBenchContains(size)
	target := 999999 // Not found
	b.ResetTimer()
	var found bool
	for i := 0; i < b.N; i++ {
		localFound := false
		for _, val := range inputSlice {
			if val == target {
				localFound = true // Should not happen
				break
			}
		}
		found = localFound
	}
	_ = found
}

// --- Run Benchmarks for different sizes ---
// Suffix indicates scenario: FE=Found Early, FL=Found Late, NF=NotFound

func BenchmarkContains_Generic_FE_100(b *testing.B) { benchmarkContainsGenericFoundEarly(100, b) }
func BenchmarkContains_Loop_FE_100(b *testing.B)    { benchmarkContainsLoopFoundEarly(100, b) }
func BenchmarkContains_Generic_FL_100(b *testing.B) { benchmarkContainsGenericFoundLate(100, b) }
func BenchmarkContains_Loop_FL_100(b *testing.B)    { benchmarkContainsLoopFoundLate(100, b) }
func BenchmarkContains_Generic_NF_100(b *testing.B) { benchmarkContainsGenericNotFound(100, b) }
func BenchmarkContains_Loop_NF_100(b *testing.B)    { benchmarkContainsLoopNotFound(100, b) }

func BenchmarkContains_Generic_FE_10000(b *testing.B) { benchmarkContainsGenericFoundEarly(10000, b) }
func BenchmarkContains_Loop_FE_10000(b *testing.B)    { benchmarkContainsLoopFoundEarly(10000, b) }
func BenchmarkContains_Generic_FL_10000(b *testing.B) { benchmarkContainsGenericFoundLate(10000, b) }
func BenchmarkContains_Loop_FL_10000(b *testing.B)    { benchmarkContainsLoopFoundLate(10000, b) }
func BenchmarkContains_Generic_NF_10000(b *testing.B) { benchmarkContainsGenericNotFound(10000, b) }
func BenchmarkContains_Loop_NF_10000(b *testing.B)    { benchmarkContainsLoopNotFound(10000, b) }
