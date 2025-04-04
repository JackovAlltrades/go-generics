package functional_test

import (
	"fmt"
	"math/rand" // Needed for benchmark data generation
	"reflect"
	"slices"
	"sort"
	"testing"
	"time"

	"github.com/JackovAlltrades/go-generics/functional"
)

// --- Test Helper Functions ---

func sortSlice[T any](slice []T) {
	switch s := any(slice).(type) {
	case []int:
		sort.Ints(s)
	case []string:
		sort.Strings(s)
	case []float64:
		sort.Float64s(s)
	case []struct {
		ID    int
		Value string
	}:
		sort.SliceStable(s, func(i, j int) bool {
			if s[i].ID != s[j].ID {
				return s[i].ID < s[j].ID
			}
			return s[i].Value < s[j].Value
		})
	default:
		panic(fmt.Sprintf("sortSlice helper does not support type %T", slice))
	}
}

func assertSlicesEquivalent[T comparable](t *testing.T, got, want []T, msgAndArgs ...interface{}) {
	t.Helper()
	gotCopy := functional.Unique(slices.Clone(got))
	wantCopy := functional.Unique(slices.Clone(want))
	sortSlice(gotCopy)
	sortSlice(wantCopy)
	if !reflect.DeepEqual(gotCopy, wantCopy) {
		prefix := fmt.Sprintf("Slices not equivalent (ignoring order and duplicates): got=%#v, want=%#v", got, want)
		if len(msgAndArgs) > 0 {
			format, ok := msgAndArgs[0].(string)
			if !ok {
				t.Fatalf("First argument to assertSlicesEquivalent custom message was not a string format: %T", msgAndArgs[0])
			}
			args := msgAndArgs[1:]
			customMessage := fmt.Sprintf(format, args...)
			t.Errorf("%s: %s", prefix, customMessage)
		} else {
			t.Errorf("%s", prefix)
		}
	}
}

// Type used in tests needing comparable structs
type comparableStruct struct {
	ID    int
	Value string
}

// --- Test Intersection ---
func TestIntersection(t *testing.T) {
	// (Test cases remain the same)
	testCases := []struct {
		name string
		a    []any
		b    []any
		want []any
	}{
		{name: "Ints_SomeOverlap", a: []any{1, 2, 3, 4}, b: []any{3, 4, 5, 6}, want: []any{3, 4}},
		{name: "Ints_NoOverlap", a: []any{1, 2}, b: []any{3, 4}, want: []any{}},
		{name: "Ints_FullOverlap", a: []any{1, 2, 3}, b: []any{1, 2, 3}, want: []any{1, 2, 3}},
		{name: "Ints_OneEmpty", a: []any{1, 2, 3}, b: []any{}, want: []any{}},
		{name: "Ints_BothEmpty", a: []any{}, b: []any{}, want: []any{}},
		{name: "Ints_NilInputA", a: nil, b: []any{1, 2}, want: []any{}},
		{name: "Ints_NilInputB", a: []any{1, 2}, b: nil, want: []any{}},
		{name: "Ints_WithDuplicates", a: []any{1, 2, 2, 3, 4, 4}, b: []any{3, 4, 4, 5, 6}, want: []any{3, 4}},
		{name: "Strings_SomeOverlap", a: []any{"a", "b", "c"}, b: []any{"b", "c", "d"}, want: []any{"b", "c"}},
		{name: "Strings_NoOverlap", a: []any{"a", "b"}, b: []any{"c", "d"}, want: []any{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any
			var wantTyped any
			getType := func(slices ...[]any) reflect.Type {
				for _, s := range slices {
					if len(s) > 0 && s[0] != nil {
						return reflect.TypeOf(s[0])
					}
				}
				return reflect.TypeOf(int(0))
			}
			resolveType := getType(tc.a, tc.b)
			switch resolveType.Kind() {
			case reflect.Int:
				aTyped := make([]int, 0, len(tc.a))
				for _, v := range tc.a {
					if item, ok := v.(int); ok {
						aTyped = append(aTyped, item)
					}
				}
				bTyped := make([]int, 0, len(tc.b))
				for _, v := range tc.b {
					if item, ok := v.(int); ok {
						bTyped = append(bTyped, item)
					}
				}
				got = functional.Intersection(aTyped, bTyped)
				wantTypedSlice := make([]int, len(tc.want))
				for i, v := range tc.want {
					wantTypedSlice[i] = v.(int)
				}
				wantTyped = wantTypedSlice
			case reflect.String:
				aTyped := make([]string, 0, len(tc.a))
				for _, v := range tc.a {
					if item, ok := v.(string); ok {
						aTyped = append(aTyped, item)
					}
				}
				bTyped := make([]string, 0, len(tc.b))
				for _, v := range tc.b {
					if item, ok := v.(string); ok {
						bTyped = append(bTyped, item)
					}
				}
				got = functional.Intersection(aTyped, bTyped)
				wantTypedSlice := make([]string, len(tc.want))
				for i, v := range tc.want {
					wantTypedSlice[i] = v.(string)
				}
				wantTyped = wantTypedSlice
			default:
				t.Fatalf("Unhandled type in Intersection test setup: %s", resolveType.String())
			}
			switch g := got.(type) {
			case []int:
				assertSlicesEquivalent(t, g, wantTyped.([]int))
			case []string:
				assertSlicesEquivalent(t, g, wantTyped.([]string))
			default:
				t.Fatalf("Unhandled type in assertion: %T", got)
			}
		})
	}
}

// --- Test Union ---
func TestUnion(t *testing.T) {
	// (Test cases remain the same)
	testCases := []struct {
		name string
		a    []any
		b    []any
		want []any
	}{
		{name: "Ints_SomeOverlap", a: []any{1, 2, 3}, b: []any{3, 4, 5}, want: []any{1, 2, 3, 4, 5}},
		{name: "Ints_NoOverlap", a: []any{1, 2}, b: []any{3, 4}, want: []any{1, 2, 3, 4}},
		{name: "Ints_FullOverlap", a: []any{1, 2, 3}, b: []any{1, 2, 3}, want: []any{1, 2, 3}},
		{name: "Ints_OneEmpty", a: []any{1, 2, 3}, b: []any{}, want: []any{1, 2, 3}},
		{name: "Ints_BothEmpty", a: []any{}, b: []any{}, want: []any{}},
		{name: "Ints_NilInputA", a: nil, b: []any{1, 2}, want: []any{1, 2}},
		{name: "Ints_NilInputB", a: []any{1, 2}, b: nil, want: []any{1, 2}},
		{name: "Ints_WithDuplicates", a: []any{1, 2, 2}, b: []any{2, 3, 3}, want: []any{1, 2, 3}},
		{name: "Strings_SomeOverlap", a: []any{"a", "b", "c"}, b: []any{"b", "c", "d"}, want: []any{"a", "b", "c", "d"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any
			var wantTyped any
			getType := func(slices ...[]any) reflect.Type {
				for _, s := range slices {
					if len(s) > 0 && s[0] != nil {
						return reflect.TypeOf(s[0])
					}
				}
				return reflect.TypeOf(int(0))
			}
			resolveType := getType(tc.a, tc.b)
			switch resolveType.Kind() {
			case reflect.Int:
				aTyped := make([]int, 0, len(tc.a))
				for _, v := range tc.a {
					if item, ok := v.(int); ok {
						aTyped = append(aTyped, item)
					}
				}
				bTyped := make([]int, 0, len(tc.b))
				for _, v := range tc.b {
					if item, ok := v.(int); ok {
						bTyped = append(bTyped, item)
					}
				}
				got = functional.Union(aTyped, bTyped)
				wantTypedSlice := make([]int, len(tc.want))
				for i, v := range tc.want {
					wantTypedSlice[i] = v.(int)
				}
				wantTyped = wantTypedSlice
			case reflect.String:
				aTyped := make([]string, 0, len(tc.a))
				for _, v := range tc.a {
					if item, ok := v.(string); ok {
						aTyped = append(aTyped, item)
					}
				}
				bTyped := make([]string, 0, len(tc.b))
				for _, v := range tc.b {
					if item, ok := v.(string); ok {
						bTyped = append(bTyped, item)
					}
				}
				got = functional.Union(aTyped, bTyped)
				wantTypedSlice := make([]string, len(tc.want))
				for i, v := range tc.want {
					wantTypedSlice[i] = v.(string)
				}
				wantTyped = wantTypedSlice
			default:
				t.Fatalf("Unhandled type in Union test setup: %s", resolveType.String())
			}
			switch g := got.(type) {
			case []int:
				assertSlicesEquivalent(t, g, wantTyped.([]int))
			case []string:
				assertSlicesEquivalent(t, g, wantTyped.([]string))
			default:
				t.Fatalf("Unhandled type in assertion: %T", got)
			}
		})
	}
}

// --- Test Difference ---
func TestDifference(t *testing.T) {
	// (Test cases remain the same)
	testCases := []struct {
		name string
		a    []any
		b    []any
		want []any
	}{
		{name: "Ints_SomeOverlap", a: []any{1, 2, 3, 4}, b: []any{3, 4, 5, 6}, want: []any{1, 2}},
		{name: "Ints_NoOverlap", a: []any{1, 2}, b: []any{3, 4}, want: []any{1, 2}},
		{name: "Ints_FullOverlap", a: []any{1, 2, 3}, b: []any{1, 2, 3}, want: []any{}},
		{name: "Ints_OneEmpty", a: []any{1, 2, 3}, b: []any{}, want: []any{1, 2, 3}},
		{name: "Ints_A_Empty", a: []any{}, b: []any{1, 2, 3}, want: []any{}},
		{name: "Ints_BothEmpty", a: []any{}, b: []any{}, want: []any{}},
		{name: "Ints_NilInputA", a: nil, b: []any{1, 2}, want: []any{}},
		{name: "Ints_NilInputB", a: []any{1, 2}, b: nil, want: []any{1, 2}},
		{name: "Ints_WithDuplicates", a: []any{1, 2, 2, 3, 4, 4}, b: []any{3, 4, 4, 5, 6}, want: []any{1, 2}},
		{name: "Strings_SomeOverlap", a: []any{"a", "b", "c"}, b: []any{"b", "c", "d"}, want: []any{"a"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any
			var wantTyped any
			getType := func(slices ...[]any) reflect.Type {
				for _, s := range slices {
					if len(s) > 0 && s[0] != nil {
						return reflect.TypeOf(s[0])
					}
				}
				return reflect.TypeOf(int(0))
			}
			resolveType := getType(tc.a, tc.b)
			switch resolveType.Kind() {
			case reflect.Int:
				aTyped := make([]int, 0, len(tc.a))
				for _, v := range tc.a {
					if item, ok := v.(int); ok {
						aTyped = append(aTyped, item)
					}
				}
				bTyped := make([]int, 0, len(tc.b))
				for _, v := range tc.b {
					if item, ok := v.(int); ok {
						bTyped = append(bTyped, item)
					}
				}
				got = functional.Difference(aTyped, bTyped)
				wantTypedSlice := make([]int, len(tc.want))
				for i, v := range tc.want {
					wantTypedSlice[i] = v.(int)
				}
				wantTyped = wantTypedSlice
			case reflect.String:
				aTyped := make([]string, 0, len(tc.a))
				for _, v := range tc.a {
					if item, ok := v.(string); ok {
						aTyped = append(aTyped, item)
					}
				}
				bTyped := make([]string, 0, len(tc.b))
				for _, v := range tc.b {
					if item, ok := v.(string); ok {
						bTyped = append(bTyped, item)
					}
				}
				got = functional.Difference(aTyped, bTyped)
				wantTypedSlice := make([]string, len(tc.want))
				for i, v := range tc.want {
					wantTypedSlice[i] = v.(string)
				}
				wantTyped = wantTypedSlice
			default:
				t.Fatalf("Unhandled type in Difference test setup: %s", resolveType.String())
			}
			switch g := got.(type) {
			case []int:
				assertSlicesEquivalent(t, g, wantTyped.([]int))
			case []string:
				assertSlicesEquivalent(t, g, wantTyped.([]string))
			default:
				t.Fatalf("Unhandled type in assertion: %T", got)
			}
		})
	}
}

// --- Test Unique ---
func TestUnique(t *testing.T) {
	// (Test cases remain the same)
	testCases := []struct {
		name  string
		input any
		want  any
	}{
		{name: "Ints_NoDuplicates", input: []int{1, 2, 3, 4}, want: []int{1, 2, 3, 4}},
		{name: "Ints_WithDuplicates", input: []int{1, 2, 2, 3, 1, 4, 4, 4, 5, 2}, want: []int{1, 2, 3, 4, 5}},
		{name: "Ints_AllDuplicates", input: []int{5, 5, 5, 5, 5}, want: []int{5}},
		{name: "Ints_EmptyInput", input: []int{}, want: []int{}},
		{name: "Ints_NilInput", input: ([]int)(nil), want: []int{}},
		{name: "Strings_NoDuplicates", input: []string{"a", "b", "c"}, want: []string{"a", "b", "c"}},
		{name: "Strings_WithDuplicates", input: []string{"apple", "banana", "apple", "orange", "banana", "apple"}, want: []string{"apple", "banana", "orange"}},
		{name: "Strings_EmptyInput", input: []string{}, want: []string{}},
		{name: "ComparableStructs_WithDuplicates", input: []comparableStruct{{1, "A"}, {2, "B"}, {1, "A"}, {3, "C"}, {2, "B"}}, want: []comparableStruct{{1, "A"}, {2, "B"}, {3, "C"}}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any
			switch in := tc.input.(type) {
			case []int:
				got = functional.Unique(in)
			case []string:
				got = functional.Unique(in)
			case []comparableStruct:
				got = functional.Unique(in)
			case nil:
				got = functional.Unique[int](nil)
			default:
				v := reflect.ValueOf(tc.input)
				if v.Kind() == reflect.Slice && v.Len() == 0 {
					elemType := v.Type().Elem()
					switch elemType.Kind() {
					case reflect.Int:
						got = functional.Unique([]int{})
					case reflect.String:
						got = functional.Unique([]string{})
					default:
						if elemType == reflect.TypeOf(comparableStruct{}) {
							got = functional.Unique([]comparableStruct{})
						} else {
							t.Fatalf("Cannot determine type for empty slice: %s", tc.name)
						}
					}
				} else {
					t.Fatalf("Unhandled input type in test setup: %T", tc.input)
				}
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Unique() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

// --- Example Unique ---
func ExampleUnique() {
	// (Example remains the same)
	intSlice := []int{1, 2, 2, 3, 1, 4, 4, 5, 2}
	uniqueInts := functional.Unique(intSlice)
	fmt.Println("Unique ints:", uniqueInts) // Order of first appearance is preserved
	stringSlice := []string{"a", "b", "a", "c", "b", "d", "a"}
	uniqueStrings := functional.Unique(stringSlice)
	fmt.Println("Unique strings:", uniqueStrings)
	emptySlice := []int{}
	uniqueEmpty := functional.Unique(emptySlice)
	fmt.Printf("Unique empty: %#v\n", uniqueEmpty) // Use %#v for clarity
	var nilSlice []string = nil
	uniqueNil := functional.Unique(nilSlice)
	fmt.Printf("Unique nil: %#v\n", uniqueNil) // Should be []string{}
	type compStruct struct {
		ID    int
		Value string
	}
	structSlice := []compStruct{{1, "A"}, {2, "B"}, {1, "A"}}
	uniqueStructs := functional.Unique(structSlice)
	fmt.Println("Unique structs:", uniqueStructs)

	// Expected Output:
	// Unique ints: [1 2 3 4 5]
	// Unique strings: [a b c d]
	// Unique empty: []int{}
	// Unique nil: []string{}
	// Unique structs: [{1 A} {2 B}]
}

// --- Benchmarks for Unique ---
func generateIntSliceWithDuplicates(size int, uniqueRatio int) []int {
	if size == 0 {
		return []int{}
	}
	if uniqueRatio <= 0 {
		uniqueRatio = 10
	} // Default ~10% unique
	slice := make([]int, size)
	divisor := size / uniqueRatio
	if divisor == 0 {
		divisor = 1
	}
	for i := 0; i < size; i++ {
		slice[i] = (i * 7) % divisor
	}
	return slice
}

func benchmarkUniqueGeneric(size int, b *testing.B) {
	if size == 0 {
		inputSlice := []int{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = functional.Unique(inputSlice)
		}
		return
	}
	inputSlice := generateIntSliceWithDuplicates(size, 10)
	b.ResetTimer()
	var result []int
	for i := 0; i < b.N; i++ {
		result = functional.Unique(inputSlice)
	}
	_ = result
}

func benchmarkUniqueLoopMap(size int, b *testing.B) {
	if size == 0 {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = make(map[int]struct{})
			_ = make([]int, 0)
		}
		return
	}
	inputSlice := generateIntSliceWithDuplicates(size, 10)
	b.ResetTimer()
	var result []int
	for i := 0; i < b.N; i++ {
		divisor := size / 10
		if divisor == 0 {
			divisor = 1
		}
		seen := make(map[int]struct{}, divisor)
		currentResult := make([]int, 0, divisor)
		for _, item := range inputSlice {
			if _, ok := seen[item]; !ok {
				seen[item] = struct{}{}
				currentResult = append(currentResult, item)
			}
		}
		result = currentResult
	}
	_ = result
}

func benchmarkUniqueLoopNaive(size int, b *testing.B) {
	if size > 1000 {
		b.Skipf("Skipping naive O(n^2) benchmark for size %d", size)
		return
	}
	if size == 0 {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = make([]int, 0)
		}
		return
	}
	inputSlice := generateIntSliceWithDuplicates(size, 10)
	b.ResetTimer()
	var result []int
	for i := 0; i < b.N; i++ {
		currentResult := make([]int, 0)
		for _, item := range inputSlice {
			isUnique := true
			for k := 0; k < len(currentResult); k++ {
				if currentResult[k] == item {
					isUnique = false
					break
				}
			}
			if isUnique {
				currentResult = append(currentResult, item)
			}
		}
		result = currentResult
	}
	_ = result
}
func BenchmarkUnique_Generic_0(b *testing.B)       { benchmarkUniqueGeneric(0, b) }
func BenchmarkUnique_LoopMap_0(b *testing.B)       { benchmarkUniqueLoopMap(0, b) }
func BenchmarkUnique_LoopNaive_0(b *testing.B)     { benchmarkUniqueLoopNaive(0, b) }
func BenchmarkUnique_Generic_100(b *testing.B)     { benchmarkUniqueGeneric(100, b) }
func BenchmarkUnique_LoopMap_100(b *testing.B)     { benchmarkUniqueLoopMap(100, b) }
func BenchmarkUnique_LoopNaive_100(b *testing.B)   { benchmarkUniqueLoopNaive(100, b) }
func BenchmarkUnique_Generic_1000(b *testing.B)    { benchmarkUniqueGeneric(1000, b) }
func BenchmarkUnique_LoopMap_1000(b *testing.B)    { benchmarkUniqueLoopMap(1000, b) }
func BenchmarkUnique_LoopNaive_1000(b *testing.B)  { benchmarkUniqueLoopNaive(1000, b) }
func BenchmarkUnique_Generic_10000(b *testing.B)   { benchmarkUniqueGeneric(10000, b) }
func BenchmarkUnique_LoopMap_10000(b *testing.B)   { benchmarkUniqueLoopMap(10000, b) }
func BenchmarkUnique_LoopNaive_10000(b *testing.B) { benchmarkUniqueLoopNaive(10000, b) } // Skipped

// --- NEW Benchmark Helpers for Set Ops ---

// generateBenchmarkSetData creates two slices 'a' and 'b' of size 'size' with controllable overlap.
// overlapRatio: 0.0 = no overlap, 0.5 = 50% overlap, 1.0 = full overlap (a==b)
func generateBenchmarkSetData(size int, overlapRatio float64) ([]int, []int) {
	if size <= 0 {
		return []int{}, []int{}
	}
	overlapRatio = max(0.0, min(1.0, overlapRatio)) // Clamp between 0 and 1

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	a := make([]int, size)
	b := make([]int, size)
	overlapCount := int(float64(size) * overlapRatio)
	distinctCountA := size - overlapCount
	distinctCountB := size - overlapCount

	maxValue := size * 3 // Range of values to reduce accidental overlap

	// Generate overlapping elements
	overlapSet := make(map[int]struct{}, overlapCount)
	for len(overlapSet) < overlapCount {
		overlapSet[rng.Intn(maxValue)] = struct{}{}
	}
	idx := 0
	for val := range overlapSet {
		a[idx] = val
		b[idx] = val
		idx++
	}

	// Generate distinct elements for A
	distinctSetA := make(map[int]struct{}, distinctCountA)
	for len(distinctSetA) < distinctCountA {
		val := rng.Intn(maxValue)
		if _, exists := overlapSet[val]; !exists {
			distinctSetA[val] = struct{}{}
		}
	}
	for val := range distinctSetA {
		a[idx] = val
		idx++
	}

	// Generate distinct elements for B
	idx = overlapCount // Reset index for B's distinct part
	distinctSetB := make(map[int]struct{}, distinctCountB)
	for len(distinctSetB) < distinctCountB {
		val := rng.Intn(maxValue)
		if _, exists := overlapSet[val]; !exists {
			if _, existsA := distinctSetA[val]; !existsA { // Ensure not in A's distinct set either
				distinctSetB[val] = struct{}{}
			}
		}
	}
	for val := range distinctSetB {
		b[idx] = val
		idx++
	}

	// Shuffle slices to make sure overlap isn't just at the start
	rng.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	rng.Shuffle(len(b), func(i, j int) { b[i], b[j] = b[j], b[i] })

	return a, b
}

// --- Intersection Benchmarks ---

func benchmarkIntersectionGeneric(a, b []int, bench *testing.B) {
	bench.ResetTimer()
	var result []int
	for i := 0; i < bench.N; i++ {
		result = functional.Intersection(a, b)
	}
	_ = result
}

func benchmarkIntersectionLoop(a, b []int, bench *testing.B) {
	bench.ResetTimer()
	var result []int
	for i := 0; i < bench.N; i++ {
		// Manual loop implementation
		if len(a) == 0 || len(b) == 0 {
			result = []int{}
			continue
		}
		// Build map from the smaller slice for potentially better performance
		var mapSlice, iterateSlice []int
		if len(a) < len(b) {
			mapSlice = a
			iterateSlice = b
		} else {
			mapSlice = b
			iterateSlice = a
		}
		set := make(map[int]struct{}, len(mapSlice))
		for _, item := range mapSlice {
			set[item] = struct{}{}
		}
		intersection := make([]int, 0) // Capacity hard to estimate
		for _, item := range iterateSlice {
			if _, exists := set[item]; exists {
				intersection = append(intersection, item)
				// Optional: remove from set if only unique intersection needed - current impl finds all matching
				// delete(set, item) // Remove to prevent duplicates if iterateSlice has them
			}
		}
		// If true unique intersection required (vs just common elements) apply Unique
		result = functional.Unique(intersection) // Assuming result should be unique
	}
	_ = result
}

// Intersection - No Overlap
var interNoA1000, interNoB1000 = generateBenchmarkSetData(1000, 0.0)

func BenchmarkIntersection_Generic_NoOverlap_N1000(b *testing.B) {
	benchmarkIntersectionGeneric(interNoA1000, interNoB1000, b)
}

func BenchmarkIntersection_Loop_NoOverlap_N1000(b *testing.B) {
	benchmarkIntersectionLoop(interNoA1000, interNoB1000, b)
}

// Intersection - Some Overlap
var interSomeA1000, interSomeB1000 = generateBenchmarkSetData(1000, 0.5)

func BenchmarkIntersection_Generic_SomeOverlap_N1000(b *testing.B) {
	benchmarkIntersectionGeneric(interSomeA1000, interSomeB1000, b)
}

func BenchmarkIntersection_Loop_SomeOverlap_N1000(b *testing.B) {
	benchmarkIntersectionLoop(interSomeA1000, interSomeB1000, b)
}

// Intersection - Full Overlap
var interFullA1000, interFullB1000 = generateBenchmarkSetData(1000, 1.0)

func BenchmarkIntersection_Generic_FullOverlap_N1000(b *testing.B) {
	benchmarkIntersectionGeneric(interFullA1000, interFullB1000, b)
}

func BenchmarkIntersection_Loop_FullOverlap_N1000(b *testing.B) {
	benchmarkIntersectionLoop(interFullA1000, interFullB1000, b)
}

// --- Union Benchmarks ---

func benchmarkUnionGeneric(a, b []int, bench *testing.B) {
	bench.ResetTimer()
	var result []int
	for i := 0; i < bench.N; i++ {
		result = functional.Union(a, b)
	}
	_ = result
}

func benchmarkUnionLoop(a, b []int, bench *testing.B) {
	bench.ResetTimer()
	var result []int
	for i := 0; i < bench.N; i++ {
		// Manual loop implementation
		set := make(map[int]struct{}, len(a)+len(b)/2) // Estimate capacity
		currentResult := make([]int, 0, len(a)+len(b)/2)

		for _, item := range a {
			if _, ok := set[item]; !ok {
				set[item] = struct{}{}
				currentResult = append(currentResult, item)
			}
		}
		for _, item := range b {
			if _, ok := set[item]; !ok {
				set[item] = struct{}{}
				currentResult = append(currentResult, item)
			}
		}
		result = currentResult
	}
	_ = result
}

// Union - No Overlap
// Re-use data from Intersection benchmarks
func BenchmarkUnion_Generic_NoOverlap_N1000(b *testing.B) {
	benchmarkUnionGeneric(interNoA1000, interNoB1000, b)
}

func BenchmarkUnion_Loop_NoOverlap_N1000(b *testing.B) {
	benchmarkUnionLoop(interNoA1000, interNoB1000, b)
}

// Union - Some Overlap
func BenchmarkUnion_Generic_SomeOverlap_N1000(b *testing.B) {
	benchmarkUnionGeneric(interSomeA1000, interSomeB1000, b)
}

func BenchmarkUnion_Loop_SomeOverlap_N1000(b *testing.B) {
	benchmarkUnionLoop(interSomeA1000, interSomeB1000, b)
}

// Union - Full Overlap
func BenchmarkUnion_Generic_FullOverlap_N1000(b *testing.B) {
	benchmarkUnionGeneric(interFullA1000, interFullB1000, b)
}

func BenchmarkUnion_Loop_FullOverlap_N1000(b *testing.B) {
	benchmarkUnionLoop(interFullA1000, interFullB1000, b)
}

// --- Difference Benchmarks ---

func benchmarkDifferenceGeneric(a, b []int, bench *testing.B) {
	bench.ResetTimer()
	var result []int
	for i := 0; i < bench.N; i++ {
		result = functional.Difference(a, b)
	}
	_ = result
}

func benchmarkDifferenceLoop(a, b []int, bench *testing.B) {
	bench.ResetTimer()
	var result []int
	for i := 0; i < bench.N; i++ {
		// Manual loop implementation (A - B)
		setB := make(map[int]struct{}, len(b))
		for _, item := range b {
			setB[item] = struct{}{}
		}
		// Need unique elements from A first
		seenA := make(map[int]struct{}, len(a))
		difference := make([]int, 0) // Capacity hard to estimate

		for _, item := range a {
			// Only process unique items of A
			if _, exists := seenA[item]; exists {
				continue
			}
			seenA[item] = struct{}{}

			// Check if it's NOT in B
			if _, exists := setB[item]; !exists {
				difference = append(difference, item)
			}
		}
		result = difference
	}
	_ = result
}

// Difference - No Overlap (A-B = A)
func BenchmarkDifference_Generic_NoOverlap_N1000(b *testing.B) {
	benchmarkDifferenceGeneric(interNoA1000, interNoB1000, b)
}

func BenchmarkDifference_Loop_NoOverlap_N1000(b *testing.B) {
	benchmarkDifferenceLoop(interNoA1000, interNoB1000, b)
}

// Difference - Some Overlap (A-B = A_distinct)
func BenchmarkDifference_Generic_SomeOverlap_N1000(b *testing.B) {
	benchmarkDifferenceGeneric(interSomeA1000, interSomeB1000, b)
}

func BenchmarkDifference_Loop_SomeOverlap_N1000(b *testing.B) {
	benchmarkDifferenceLoop(interSomeA1000, interSomeB1000, b)
}

// Difference - Full Overlap (A-B = empty)
func BenchmarkDifference_Generic_FullOverlap_N1000(b *testing.B) {
	benchmarkDifferenceGeneric(interFullA1000, interFullB1000, b)
}

func BenchmarkDifference_Loop_FullOverlap_N1000(b *testing.B) {
	benchmarkDifferenceLoop(interFullA1000, interFullB1000, b)
}
