package functional_test

import (
	"fmt"
	"reflect" // Needed for DeepEqual
	"testing"

	"github.com/JackovAlltrades/go-generics/functional" // Adjust import path
)

// --- Test Find ---
func TestFind(t *testing.T) {
	// Define a local struct for this test
	type person struct {
		Name string
		Age  int
	}
	p1 := person{"Alice", 30}
	p2 := person{"Bob", 25}
	p3 := person{"Charlie", 35}

	testCases := []struct {
		name      string
		input     any // Use 'any' for testing different types
		predicate any // The predicate function
		// wantValue *any // Expected pointer to the value (difficult with 'any')
		wantValueCheck func(any) bool // Function to check the dereferenced value
		wantOk         bool           // Expected boolean result
	}{
		// Int Tests
		{
			name:  "Ints_FindFirstEven",
			input: []int{1, 3, 5, 6, 7, 8},
			predicate: func(i int) bool {
				return i%2 == 0
			},
			wantValueCheck: func(v any) bool { return v.(int) == 6 },
			wantOk:         true,
		},
		{
			name:  "Ints_FindNonExistent",
			input: []int{1, 3, 5, 7, 9},
			predicate: func(i int) bool {
				return i%2 == 0
			},
			wantValueCheck: func(v any) bool { return false }, // Not checked if not found
			wantOk:         false,
		},
		{
			name:           "Ints_EmptySlice",
			input:          []int{},
			predicate:      func(i int) bool { return true },
			wantValueCheck: func(v any) bool { return false },
			wantOk:         false,
		},
		{
			name:           "Ints_NilSlice",
			input:          ([]int)(nil),
			predicate:      func(i int) bool { return true },
			wantValueCheck: func(v any) bool { return false },
			wantOk:         false,
		},
		// String Tests
		{
			name:  "Strings_FindFirstLongerThan3",
			input: []string{"a", "bee", "cat", "door", "xyz"},
			predicate: func(s string) bool {
				return len(s) > 3
			},
			wantValueCheck: func(v any) bool { return v.(string) == "door" },
			wantOk:         true,
		},
		{
			name:  "Strings_FindNonExistent",
			input: []string{"a", "b", "c"},
			predicate: func(s string) bool {
				return len(s) > 1
			},
			wantValueCheck: func(v any) bool { return false },
			wantOk:         false,
		},
		// Struct Tests
		{
			name:  "Structs_FindPersonByName",
			input: []person{p1, p2, p3}, // Use defined structs
			predicate: func(p person) bool {
				return p.Name == "Bob"
			},
			// Check if the дереференced value equals p2
			wantValueCheck: func(v any) bool { return reflect.DeepEqual(v, p2) },
			wantOk:         true,
		},
		{
			name: "Structs_FindPersonNonExistent",
			input: []person{
				{"Alice", 30}, {"Bob", 25},
			},
			predicate: func(p person) bool {
				return p.Age > 40
			},
			wantValueCheck: func(v any) bool { return false },
			wantOk:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotPtr any // Changed to gotPtr to reflect pointer return
			var gotOk bool

			// Type switching to call the correct generic instantiation
			switch pred := tc.predicate.(type) {
			case func(int) bool:
				in, ok := tc.input.([]int)
				if !ok && tc.input != nil {
					t.Fatalf("Input type mismatch for func(int) bool")
				}
				// Find returns (*T, bool)
				gotPtr, gotOk = functional.Find(in, pred)
			case func(string) bool:
				in, ok := tc.input.([]string)
				if !ok && tc.input != nil {
					t.Fatalf("Input type mismatch for func(string) bool")
				}
				gotPtr, gotOk = functional.Find(in, pred)
			case func(person) bool:
				in, ok := tc.input.([]person)
				if !ok && tc.input != nil {
					t.Fatalf("Input type mismatch for func(person) bool")
				}
				gotPtr, gotOk = functional.Find(in, pred)
			default:
				t.Fatalf("Unhandled predicate type in test setup: %T", tc.predicate)
			}

			// Check boolean result
			if gotOk != tc.wantOk {
				t.Errorf("Find() ok = %v, want %v", gotOk, tc.wantOk)
			}

			// Check pointer and value
			if gotOk {
				// Check if pointer is non-nil when ok is true
				if gotPtr == nil {
					t.Errorf("Find() pointer is nil, but ok is true")
				} else {
					// Check the value the pointer points to
					// Need to get the value via reflection since gotPtr is 'any'
					val := reflect.ValueOf(gotPtr).Elem().Interface()
					if !tc.wantValueCheck(val) {
						t.Errorf("Find() value = %#v, but wantValueCheck failed", val)
					}

					// Specific check for struct case to ensure pointer is to original slice data
					if persons, ok := tc.input.([]person); ok {
						found := false
						for i := range persons {
							// Check if the returned pointer points to the memory address of any element in the original slice
							if gotPtr == &persons[i] {
								// Additionally verify the content matches (already done by wantValueCheck)
								if reflect.DeepEqual(persons[i], reflect.ValueOf(gotPtr).Elem().Interface()) {
									found = true
									break
								}
							}
						}
						if !found {
							t.Errorf("Find() returned pointer %p does not point to an element within the original slice", gotPtr)
						}
					}
				}
			} else {
				// Check if pointer is nil when ok is false
				// Note: Checking gotPtr == nil is tricky when gotPtr is 'any'. Need reflection.
				if !reflect.ValueOf(gotPtr).IsNil() {
					t.Errorf("Find() pointer = %p (%#v), want nil when ok is false", gotPtr, reflect.ValueOf(gotPtr).Elem().Interface())
				}
			}
		})
	}
}

// NOTE: TestFindPtr has been removed as functional.Find returns the pointer.

// --- Examples ---

func ExampleFind() {
	// Example 1: Find integer
	numbers := []int{10, 25, 31, 45, 50}
	foundNumPtr, ok := functional.Find(numbers, func(n int) bool {
		return n > 30 // Find first number > 30
	})
	if ok {
		fmt.Printf("Found number > 30: %d\n", *foundNumPtr) // Dereference pointer
	} else {
		fmt.Println("No number > 30 found")
	}

	// Example 2: Find struct and modify
	users := []struct {
		ID   int
		Name string
	}{
		{1, "Alice"},
		{2, "Bob"},
		{3, "Charlie"},
	}

	userPtr, ok := functional.Find(users, func(u struct {
		ID   int
		Name string
	},
	) bool {
		return u.ID == 2
	}) // Find user with ID 2

	if ok {
		fmt.Printf("Found user: ID=%d, Name=%s\n", userPtr.ID, userPtr.Name)
		// Modify the found user through the pointer
		userPtr.Name = "Robert"
		fmt.Printf("Original slice modified: User at index 1 is now %#v\n", users[1]) // Bob is at index 1
	} else {
		fmt.Println("User with ID 2 not found")
	}

	// Example 3: Not found
	strPtr, ok := functional.Find(users, func(u struct {
		ID   int
		Name string
	},
	) bool {
		return u.ID == 4
	}) // Find user with ID 4

	if !ok {
		fmt.Println("User with ID 4 not found")
		fmt.Printf("Pointer returned when not found: %v\n", strPtr) // Should be nil
	}

	// Output:
	// Found number > 30: 31
	// Found user: ID=2, Name=Bob
	// Original slice modified: User at index 1 is now struct { ID int; Name string }{ID:2, Name:"Robert"}
	// User with ID 4 not found
	// Pointer returned when not found: <nil>
}

// NOTE: ExampleFindPtr has been removed.

// --- Benchmarks ---

// Helper: Generate slice of simple structs for benchmarks
type benchStruct struct {
	ID    int
	Value string
}

func generateBenchStructSlice(size int) []benchStruct {
	if size <= 0 {
		return []benchStruct{}
	}
	slice := make([]benchStruct, size)
	for i := 0; i < size; i++ {
		slice[i] = benchStruct{ID: i, Value: "v" + fmt.Sprint(i)}
	}
	return slice
}

// Predicate: find struct with ID near the start
var findBenchStructEarlyPred = func(item benchStruct) bool {
	return item.ID == 5 // Found early (index 5)
}

// Predicate: find struct with ID near the end
func findBenchStructLatePred(size int) func(benchStruct) bool {
	targetID := size - 2 // Found late (ensure size >= 2)
	if targetID < 0 {
		targetID = 0
	}
	return func(item benchStruct) bool {
		return item.ID == targetID
	}
}

// Predicate: find struct that doesn't exist
var findBenchStructNotFoundPred = func(item benchStruct) bool {
	return item.ID == -1 // Not found
}

// Benchmark functional.Find - Found Early
func benchmarkFindGenericEarly(size int, b *testing.B) {
	if size < 6 {
		b.Skip("Skipping small size for Early test")
		return
	}
	inputSlice := generateBenchStructSlice(size)
	predicate := findBenchStructEarlyPred
	b.ResetTimer()
	var ptr *benchStruct // Changed variable name
	var ok bool
	for i := 0; i < b.N; i++ {
		// Receives (*T, bool)
		ptr, ok = functional.Find(inputSlice, predicate)
	}
	_ = ptr // Use result
	_ = ok
}

// Benchmark loop Find - Found Early (Adjusted to return *T, bool)
func benchmarkFindLoopEarly(size int, b *testing.B) {
	if size < 6 {
		b.Skip("Skipping small size for Early test")
		return
	}
	inputSlice := generateBenchStructSlice(size)
	predicate := findBenchStructEarlyPred
	b.ResetTimer()
	var ptr *benchStruct // Changed variable name
	var ok bool
	for i := 0; i < b.N; i++ {
		localPtr := (*benchStruct)(nil) // Explicitly *benchStruct nil
		localOk := false
		for idx := range inputSlice { // Iterate by index to get pointer to element
			if predicate(inputSlice[idx]) {
				localPtr = &inputSlice[idx] // Get pointer to item in slice
				localOk = true
				break // Found
			}
		}
		ptr = localPtr
		ok = localOk
	}
	_ = ptr // Use result
	_ = ok
}

// Benchmark functional.Find - Found Late
func benchmarkFindGenericLate(size int, b *testing.B) {
	if size < 2 {
		b.Skip("Skipping small size for Late test")
		return
	}
	inputSlice := generateBenchStructSlice(size)
	predicate := findBenchStructLatePred(size)
	b.ResetTimer()
	var ptr *benchStruct
	var ok bool
	for i := 0; i < b.N; i++ {
		ptr, ok = functional.Find(inputSlice, predicate)
	}
	_ = ptr
	_ = ok
}

// Benchmark loop Find - Found Late (Adjusted to return *T, bool)
func benchmarkFindLoopLate(size int, b *testing.B) {
	if size < 2 {
		b.Skip("Skipping small size for Late test")
		return
	}
	inputSlice := generateBenchStructSlice(size)
	predicate := findBenchStructLatePred(size)
	b.ResetTimer()
	var ptr *benchStruct
	var ok bool
	for i := 0; i < b.N; i++ {
		localPtr := (*benchStruct)(nil)
		localOk := false
		for idx := range inputSlice {
			if predicate(inputSlice[idx]) {
				localPtr = &inputSlice[idx]
				localOk = true
				break // Found
			}
		}
		ptr = localPtr
		ok = localOk
	}
	_ = ptr
	_ = ok
}

// Benchmark functional.Find - Not Found (Worst Case)
func benchmarkFindGenericNotFound(size int, b *testing.B) {
	inputSlice := generateBenchStructSlice(size)
	predicate := findBenchStructNotFoundPred
	b.ResetTimer()
	var ptr *benchStruct
	var ok bool
	for i := 0; i < b.N; i++ {
		ptr, ok = functional.Find(inputSlice, predicate)
	}
	_ = ptr
	_ = ok
}

// Benchmark loop Find - Not Found (Worst Case) (Adjusted to return *T, bool)
func benchmarkFindLoopNotFound(size int, b *testing.B) {
	inputSlice := generateBenchStructSlice(size)
	predicate := findBenchStructNotFoundPred
	b.ResetTimer()
	var ptr *benchStruct
	var ok bool
	for i := 0; i < b.N; i++ {
		localPtr := (*benchStruct)(nil)
		localOk := false
		for idx := range inputSlice {
			if predicate(inputSlice[idx]) {
				localPtr = &inputSlice[idx]
				localOk = true
				break // Found (shouldn't happen)
			}
		}
		ptr = localPtr
		ok = localOk // Will be false
	}
	_ = ptr
	_ = ok
}

// --- Run Benchmarks ---
const (
	N1_Find = 100
	N2_Find = 10000
)

// Found Early
func BenchmarkFind_Generic_Early_N100(b *testing.B)   { benchmarkFindGenericEarly(N1_Find, b) }
func BenchmarkFind_Loop_Early_N100(b *testing.B)      { benchmarkFindLoopEarly(N1_Find, b) }
func BenchmarkFind_Generic_Early_N10000(b *testing.B) { benchmarkFindGenericEarly(N2_Find, b) }
func BenchmarkFind_Loop_Early_N10000(b *testing.B)    { benchmarkFindLoopEarly(N2_Find, b) }

// Found Late
func BenchmarkFind_Generic_Late_N100(b *testing.B)   { benchmarkFindGenericLate(N1_Find, b) }
func BenchmarkFind_Loop_Late_N100(b *testing.B)      { benchmarkFindLoopLate(N1_Find, b) }
func BenchmarkFind_Generic_Late_N10000(b *testing.B) { benchmarkFindGenericLate(N2_Find, b) }
func BenchmarkFind_Loop_Late_N10000(b *testing.B)    { benchmarkFindLoopLate(N2_Find, b) }

// Not Found
func BenchmarkFind_Generic_NotFound_N100(b *testing.B)   { benchmarkFindGenericNotFound(N1_Find, b) }
func BenchmarkFind_Loop_NotFound_N100(b *testing.B)      { benchmarkFindLoopNotFound(N1_Find, b) }
func BenchmarkFind_Generic_NotFound_N10000(b *testing.B) { benchmarkFindGenericNotFound(N2_Find, b) }
func BenchmarkFind_Loop_NotFound_N10000(b *testing.B)    { benchmarkFindLoopNotFound(N2_Find, b) }
