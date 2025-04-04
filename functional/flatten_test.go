package functional_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/JackovAlltrades/go-generics/functional" // Adjust import path
)

// --- Test Flatten ---
func TestFlatten(t *testing.T) {
	testCases := []struct {
		name  string
		input any // [][]T
		want  any // []T
	}{
		{
			name:  "FlattenInts",
			input: [][]int{{1, 2}, {3, 4, 5}, {6}},
			want:  []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:  "FlattenStrings",
			input: [][]string{{"a", "b"}, {}, {"c", "d"}, {"e"}},
			want:  []string{"a", "b", "c", "d", "e"},
		},
		{
			name:  "FlattenWithEmptyInnerSlices",
			input: [][]int{{}, {1, 2}, {}, {3}},
			want:  []int{1, 2, 3},
		},
		{
			name:  "FlattenWithNilInnerSlices",
			input: [][]int{nil, {1, 2}, nil, {3}}, // Inner slices can be nil
			want:  []int{1, 2, 3},
		},
		{
			name:  "FlattenAllEmpty",
			input: [][]int{{}, {}, {}},
			want:  []int{},
		},
		{
			name:  "FlattenAllNilInner",
			input: [][]int{nil, nil, nil},
			want:  []int{},
		},
		{
			name:  "FlattenEmptyOuter",
			input: [][]int{},
			want:  []int{},
		},
		{
			name:  "FlattenNilOuter",
			input: ([][]int)(nil),
			want:  []int{},
		},
		{
			name:  "FlattenSingleInner",
			input: [][]int{{10, 20, 30}},
			want:  []int{10, 20, 30},
		},
		{
			name:  "FlattenSingleEmptyInner",
			input: [][]int{{}},
			want:  []int{},
		},
		// Consider adding a test case with different pointer types if applicable
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any

			switch input := tc.input.(type) {
			case [][]int:
				got = functional.Flatten(input)
			case [][]string:
				got = functional.Flatten(input)
			case nil: // Handle nil outer slice
				// Assuming Flatten[T] called on nil [][]T should return empty []T
				// Let's test with int as the type T
				got = functional.Flatten[int](nil)
			default:
				// Handle empty outer slice case
				v := reflect.ValueOf(tc.input)
				if v.Kind() == reflect.Slice && v.Len() == 0 {
					elemType := v.Type().Elem().Elem() // Get type T from [][]T
					switch elemType.Kind() {
					case reflect.Int:
						got = functional.Flatten([][]int{})
					case reflect.String:
						got = functional.Flatten([][]string{})
					default:
						t.Fatalf("Unhandled type for empty outer slice: %s", elemType.String())
					}
				} else {
					t.Fatalf("Unhandled input type in test setup: %T", tc.input)
				}
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Flatten() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

// --- Flatten Examples ---
func ExampleFlatten() {
	nestedInts := [][]int{{1, 2}, {3, 4, 5}, {}, {6}}
	flattenedInts := functional.Flatten(nestedInts)
	fmt.Println("Flattened Ints:", flattenedInts)

	nestedStrings := [][]string{{"hello", "world"}, nil, {"functional", "go"}}
	flattenedStrings := functional.Flatten(nestedStrings)
	fmt.Println("Flattened Strings:", flattenedStrings)

	emptyOuter := [][]float64{}
	flattenedEmpty := functional.Flatten(emptyOuter)
	fmt.Printf("Flattened Empty Outer: %#v\n", flattenedEmpty)

	var nilOuter [][]string = nil
	flattenedNil := functional.Flatten(nilOuter)
	fmt.Printf("Flattened Nil Outer: %#v\n", flattenedNil)

	// Output:
	// Flattened Ints: [1 2 3 4 5 6]
	// Flattened Strings: [hello world functional go]
	// Flattened Empty Outer: []float64{}
	// Flattened Nil Outer: []string{}
}

// --- Benchmarks ---

// Generate nested slice data for flattening benchmarks
// Creates 'outerSize' inner slices, each of size 'innerSize'.
func generateNestedSlice[T any](outerSize, innerSize int, generator func(i, j int) T) [][]T {
	if outerSize == 0 {
		return [][]T{}
	}
	nested := make([][]T, outerSize)
	for i := 0; i < outerSize; i++ {
		if innerSize == 0 {
			nested[i] = []T{} // Empty inner slice
		} else {
			inner := make([]T, innerSize)
			for j := 0; j < innerSize; j++ {
				inner[j] = generator(i, j)
			}
			nested[i] = inner
		}
	}
	return nested
}

var intGenerator = func(i, j int) int { return i*1000 + j } // Simple unique int generator

// Benchmark runner for generic Flatten
func benchmarkFlattenGeneric(input [][]int, b *testing.B) {
	if input == nil { // Avoid issues with nil input in benchmark loop
		b.Skip("Skipping benchmark for nil input")
		return
	}
	b.ResetTimer()
	var result []int
	for i := 0; i < b.N; i++ {
		result = functional.Flatten(input)
	}
	_ = result
}

// Benchmark runner for loop-based flatten
func benchmarkFlattenLoop(input [][]int, b *testing.B) {
	if input == nil {
		b.Skip("Skipping benchmark for nil input")
		return
	}
	b.ResetTimer()
	var result []int
	for i := 0; i < b.N; i++ {
		// Manual loop implementation
		totalLen := 0
		for _, inner := range input {
			totalLen += len(inner)
		}
		currentResult := make([]int, 0, totalLen) // Preallocate with exact total length
		for _, inner := range input {
			currentResult = append(currentResult, inner...) // Use '...' to append slice contents
		}
		result = currentResult
	}
	_ = result
}

// Define benchmark scenarios
var (
	flattenData_10_10    = generateNestedSlice(10, 10, intGenerator)    // 100 total elements
	flattenData_100_10   = generateNestedSlice(100, 10, intGenerator)   // 1000 total elements
	flattenData_10_100   = generateNestedSlice(10, 100, intGenerator)   // 1000 total elements
	flattenData_1000_10  = generateNestedSlice(1000, 10, intGenerator)  // 10000 total elements
	flattenData_10_1000  = generateNestedSlice(10, 1000, intGenerator)  // 10000 total elements
	flattenData_100_100  = generateNestedSlice(100, 100, intGenerator)  // 10000 total elements
	flattenData_1000_100 = generateNestedSlice(1000, 100, intGenerator) // 100000 total elements
)

// --- Run Benchmarks ---

// Scenario: Few large inner slices
func BenchmarkFlatten_Generic_10x1000(b *testing.B) { benchmarkFlattenGeneric(flattenData_10_1000, b) }
func BenchmarkFlatten_Loop_10x1000(b *testing.B)    { benchmarkFlattenLoop(flattenData_10_1000, b) }

// Scenario: Many small inner slices
func BenchmarkFlatten_Generic_1000x10(b *testing.B) { benchmarkFlattenGeneric(flattenData_1000_10, b) }
func BenchmarkFlatten_Loop_1000x10(b *testing.B)    { benchmarkFlattenLoop(flattenData_1000_10, b) }

// Scenario: Balanced outer/inner (~10k total)
func BenchmarkFlatten_Generic_100x100(b *testing.B) { benchmarkFlattenGeneric(flattenData_100_100, b) }
func BenchmarkFlatten_Loop_100x100(b *testing.B)    { benchmarkFlattenLoop(flattenData_100_100, b) }

// Scenario: Very Many inner slices (~100k total)
func BenchmarkFlatten_Generic_1000x100(b *testing.B) {
	benchmarkFlattenGeneric(flattenData_1000_100, b)
}
func BenchmarkFlatten_Loop_1000x100(b *testing.B) { benchmarkFlattenLoop(flattenData_1000_100, b) }
