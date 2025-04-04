package functional_test

import (
	"fmt"
	"reflect"
	"slices" // For slices.Clone (Go 1.21+)
	"testing"

	"github.com/JackovAlltrades/go-generics/functional"
)

// --- Test Reverse ---
func TestReverse(t *testing.T) {
	testCases := []struct {
		name  string
		input any // Slice type T
		want  any // Expected slice of type T after reversing
	}{
		{
			name:  "ReverseInts",
			input: []int{1, 2, 3, 4, 5},
			want:  []int{5, 4, 3, 2, 1},
		},
		{
			name:  "ReverseStrings",
			input: []string{"a", "b", "c", "d"},
			want:  []string{"d", "c", "b", "a"},
		},
		{
			name:  "ReverseSingleElement",
			input: []int{42},
			want:  []int{42},
		},
		{
			name:  "ReverseEmpty",
			input: []int{},
			want:  []int{},
		},
		{
			name:  "ReverseNil",
			input: ([]int)(nil),
			want:  ([]int)(nil), // Reversing nil should probably result in nil
		},
		// A test case with an even number of elements
		{
			name:  "ReverseEvenLength",
			input: []float64{1.1, 2.2, 3.3, 4.4},
			want:  []float64{4.4, 3.3, 2.2, 1.1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Crucial: Reverse modifies in-place. We need a copy for the test.
			var inputCopy any

			// Make a copy based on the type
			switch v := tc.input.(type) {
			case []int:
				inputCopy = slices.Clone(v) // Use slices.Clone (Go 1.21+)
			case []string:
				inputCopy = slices.Clone(v)
			case []float64:
				inputCopy = slices.Clone(v)
			case nil:
				inputCopy = nil // Copy of nil is nil
			default:
				// Handle empty slice case which might be typed differently
				rv := reflect.ValueOf(tc.input)
				if rv.Kind() == reflect.Slice && rv.Len() == 0 {
					// Create an empty slice of the correct type
					inputCopy = reflect.MakeSlice(rv.Type(), 0, 0).Interface()
				} else {
					t.Fatalf("Unhandled type for copying in test setup: %T", tc.input)
				}
			}

			// Call Reverse on the copy
			switch c := inputCopy.(type) {
			case []int:
				functional.Reverse(c)
			case []string:
				functional.Reverse(c)
			case []float64:
				functional.Reverse(c)
			case nil:
				functional.Reverse[any](nil) // Call on nil (with explicit type if needed by impl)
			default:
				// Check again if it was an empty slice
				rv := reflect.ValueOf(inputCopy)
				if !(rv.Kind() == reflect.Slice && rv.Len() == 0) {
					// If it's not nil and not an empty slice we explicitly handled, error
					t.Fatalf("Unhandled type for calling Reverse in test setup: %T", inputCopy)
				}
				// If it was an empty slice, calling Reverse on it is fine, no specific call needed.
			}

			// Now compare the modified copy (inputCopy) with the expected want value
			if !reflect.DeepEqual(inputCopy, tc.want) {
				// Special handling for comparing nil vs empty slice potentially

				// Consider nil == empty slice for reversal? Usually no. nil reversed is nil. Empty reversed is empty.
				// DeepEqual handles nil vs nil and empty vs empty correctly.
				// So a direct DeepEqual should be sufficient.

				t.Errorf("Reverse() resulted in %#v, want %#v", inputCopy, tc.want)
			}
		})
	}
}

// --- Reverse Examples ---
func ExampleReverse() {
	nums := []int{10, 20, 30, 40, 50}
	fmt.Println("Original nums:", nums)
	functional.Reverse(nums) // Modifies in-place
	fmt.Println("Reversed nums:", nums)

	letters := []string{"x", "y", "z"}
	fmt.Println("Original letters:", letters)
	functional.Reverse(letters)
	fmt.Println("Reversed letters:", letters)

	empty := []float32{}
	fmt.Println("Original empty:", empty)
	functional.Reverse(empty) // Should do nothing
	fmt.Println("Reversed empty:", empty)

	var nilSlice []int = nil
	fmt.Printf("Original nil: %#v\n", nilSlice)
	functional.Reverse(nilSlice) // Should do nothing (no panic)
	fmt.Printf("Reversed nil: %#v\n", nilSlice)

	// Output:
	// Original nums: [10 20 30 40 50]
	// Reversed nums: [50 40 30 20 10]
	// Original letters: [x y z]
	// Reversed letters: [z y x]
	// Original empty: []
	// Reversed empty: []
	// Original nil: []int(nil)
	// Reversed nil: []int(nil)
}

// --- Benchmarks ---

// Helper to generate slice data
func generateSliceForReverse(size int) []int {
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = i
	}
	return data
}

// NOTE: Benchmarking in-place operations requires care.
// We must ensure each benchmark iteration operates on fresh or reset data,
// otherwise, we're just benchmarking reversing an already reversed slice repeatedly.
// The easiest way is often to CLONE the data inside the benchmark loop,
// although this adds the overhead of cloning to the measurement.
// A more complex alternative is to reverse it back manually after each call,
// but that adds its own overhead. Cloning is simpler to implement correctly here.

// Benchmark runner for generic Reverse (with cloning inside)
func benchmarkReverseGeneric(b *testing.B, data []int) {
	if data == nil {
		b.Skip("Skipping nil data")
		return
	}
	b.ResetTimer() // Start timing before the loop

	for i := 0; i < b.N; i++ {
		// Clone ensures we reverse original data each time
		dataCopy := slices.Clone(data)
		functional.Reverse(dataCopy) // Operate on the copy
	}
}

// Benchmark runner for loop-based Reverse (with cloning inside)
func benchmarkReverseLoop(b *testing.B, data []int) {
	if data == nil {
		b.Skip("Skipping nil data")
		return
	}
	b.ResetTimer() // Start timing before the loop

	for i := 0; i < b.N; i++ {
		// Clone ensures we reverse original data each time
		dataCopy := slices.Clone(data)
		// Manual loop implementation (in-place)
		l := len(dataCopy)
		for i := 0; i < l/2; i++ {
			dataCopy[i], dataCopy[l-1-i] = dataCopy[l-1-i], dataCopy[i]
		}
	}
}

// Define benchmark scenarios
var (
	reverseDataN10    = generateSliceForReverse(10)
	reverseDataN100   = generateSliceForReverse(100)
	reverseDataN1000  = generateSliceForReverse(1000)
	reverseDataN10000 = generateSliceForReverse(10000)
)

// --- Run Benchmarks ---

func BenchmarkReverse_Generic_N10(b *testing.B) { benchmarkReverseGeneric(b, reverseDataN10) }
func BenchmarkReverse_Loop_N10(b *testing.B)    { benchmarkReverseLoop(b, reverseDataN10) }

func BenchmarkReverse_Generic_N100(b *testing.B) { benchmarkReverseGeneric(b, reverseDataN100) }
func BenchmarkReverse_Loop_N100(b *testing.B)    { benchmarkReverseLoop(b, reverseDataN100) }

func BenchmarkReverse_Generic_N1000(b *testing.B) { benchmarkReverseGeneric(b, reverseDataN1000) }
func BenchmarkReverse_Loop_N1000(b *testing.B)    { benchmarkReverseLoop(b, reverseDataN1000) }

func BenchmarkReverse_Generic_N10000(b *testing.B) { benchmarkReverseGeneric(b, reverseDataN10000) }
func BenchmarkReverse_Loop_N10000(b *testing.B)    { benchmarkReverseLoop(b, reverseDataN10000) }
