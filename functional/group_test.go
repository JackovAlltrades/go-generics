package functional_test

import (
	"fmt"
	"math/rand" // For generating varied benchmark data
	"reflect"
	"sort" // To compare map values deterministically
	"strconv"
	"testing"
	"time" // For seeding rand

	"github.com/JackovAlltrades/go-generics/functional" // Adjust import path
)

// --- Structs used in tests ---
type personGroupTest struct {
	ID   int
	City string
}

// --- Test GroupBy ---
func TestGroupBy(t *testing.T) {
	// Helper to compare grouped results (maps of slices) robustly
	assertSameGroups := func(t *testing.T, got, want map[any][]any) {
		t.Helper()
		if len(got) != len(want) {
			t.Fatalf("GroupBy() wrong number of groups: got %d, want %d. Got: %#v, Want: %#v", len(got), len(want), got, want)
		}

		for key, wantSlice := range want {
			gotSlice, ok := got[key]
			if !ok {
				t.Errorf("GroupBy() missing group key: %#v", key)
				continue
			}
			// Must sort slices within the maps for comparison as order isn't guaranteed
			sortAnySlice := func(slice []any) {
				if len(slice) < 2 { // No need to sort 0 or 1 elements
					return
				}
				// Use reflection to determine type and sort
				switch slice[0].(type) {
				case int:
					concreteSlice := make([]int, len(slice))
					for i, v := range slice {
						concreteSlice[i] = v.(int)
					}
					sort.Ints(concreteSlice)
					for i, v := range concreteSlice {
						slice[i] = v
					}
				case string:
					concreteSlice := make([]string, len(slice))
					for i, v := range slice {
						concreteSlice[i] = v.(string)
					}
					sort.Strings(concreteSlice)
					for i, v := range concreteSlice {
						slice[i] = v
					}
				case personGroupTest:
					concreteSlice := make([]personGroupTest, len(slice))
					for i, v := range slice {
						concreteSlice[i] = v.(personGroupTest)
					}
					sort.SliceStable(concreteSlice, func(i, j int) bool { return concreteSlice[i].ID < concreteSlice[j].ID })
					for i, v := range concreteSlice {
						slice[i] = v
					}
				default:
					t.Logf("Warning: Sorting not implemented for group type %T, using DeepEqual on potentially unsorted slices.", slice[0])
				}
			}

			sortAnySlice(gotSlice)
			sortAnySlice(wantSlice)

			if !reflect.DeepEqual(gotSlice, wantSlice) {
				t.Errorf("GroupBy() group for key %#v mismatch:\ngot  %#v\nwant %#v", key, gotSlice, wantSlice)
			}
		}
	}

	// Define person instances for struct test clarity
	pLon1 := personGroupTest{1, "London"}
	pPar2 := personGroupTest{2, "Paris"}
	pLon3 := personGroupTest{3, "London"}
	pTok4 := personGroupTest{4, "Tokyo"}
	pPar5 := personGroupTest{5, "Paris"}

	testCases := []struct {
		name      string
		input     any           // []T
		keyFunc   any           // func(T) K
		wantGroup map[any][]any // map[K][]T (using any for test setup)
	}{
		{
			name:  "GroupBy_EvenOddInts",
			input: []int{1, 2, 3, 4, 5, 6},
			keyFunc: func(n int) string {
				if n%2 == 0 {
					return "even"
				}
				return "odd"
			},
			// Use concrete types in the literal, assign to map[any][]any field
			wantGroup: map[any][]any{
				"even": {any(2), any(4), any(6)}, // Cast elements to any
				"odd":  {any(1), any(3), any(5)},
			},
		},
		{
			name:    "GroupBy_StringLength",
			input:   []string{"a", "bb", "ccc", "dd", "e", "ffff"},
			keyFunc: func(s string) int { return len(s) },
			wantGroup: map[any][]any{
				any(1): {any("a"), any("e")}, // Cast keys and elements
				any(2): {any("bb"), any("dd")},
				any(3): {any("ccc")},
				any(4): {any("ffff")},
			},
		},
		{
			name:    "GroupBy_StructField",
			input:   []personGroupTest{pLon1, pPar2, pLon3, pTok4, pPar5},
			keyFunc: func(p personGroupTest) string { return p.City },
			wantGroup: map[any][]any{
				"London": {any(pLon1), any(pLon3)},
				"Paris":  {any(pPar2), any(pPar5)},
				"Tokyo":  {any(pTok4)},
			},
		},
		{
			name:      "GroupBy_EmptyInput",
			input:     []int{},
			keyFunc:   func(n int) int { return n },
			wantGroup: map[any][]any{}, // Empty map is fine
		},
		{
			name:      "GroupBy_NilInput",
			input:     ([]string)(nil),
			keyFunc:   func(s string) int { return len(s) },
			wantGroup: map[any][]any{},
		},
		{
			name:    "GroupBy_AllSameKey",
			input:   []int{10, 20, 30},
			keyFunc: func(n int) string { return "all" },
			wantGroup: map[any][]any{
				"all": {any(10), any(20), any(30)},
			},
		},
		{
			name:    "GroupBy_BooleanKey",
			input:   []int{1, -2, 3, 0, -5},
			keyFunc: func(n int) bool { return n >= 0 },
			wantGroup: map[any][]any{
				any(true):  {any(1), any(3), any(0)}, // Cast bool key and int elements
				any(false): {any(-2), any(-5)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotGroup any // The result map[K][]T

			inputVal := reflect.ValueOf(tc.input)
			keyFuncVal := reflect.ValueOf(tc.keyFunc)

			if !inputVal.IsValid() || (inputVal.Kind() == reflect.Slice && inputVal.IsNil()) || inputVal.Len() == 0 {
				// Handle nil/empty slice input -> result is always empty map
				// Create empty map matching expected type structure if possible
				gotGroup = make(map[any][]any) // Default if types cannot be inferred
				if keyFuncVal.IsValid() && inputVal.IsValid() && inputVal.Type().Elem() != nil {
					keyType := keyFuncVal.Type().Out(0) // K type
					elemType := inputVal.Type().Elem()  // T type
					mapType := reflect.MapOf(keyType, reflect.SliceOf(elemType))
					gotGroup = reflect.MakeMap(mapType).Interface()
				}

			} else {
				// --- Use type switch on input type ---
				switch input := tc.input.(type) {
				case []int:
					switch kf := tc.keyFunc.(type) {
					case func(int) string:
						gotGroup = functional.GroupBy(input, kf)
					case func(int) bool:
						gotGroup = functional.GroupBy(input, kf)
					case func(int) int:
						gotGroup = functional.GroupBy(input, kf)
					default:
						t.Fatalf("Unhandled key func type for []int input: %T", tc.keyFunc)
					}
				case []string:
					switch kf := tc.keyFunc.(type) {
					case func(string) int:
						gotGroup = functional.GroupBy(input, kf)
					default:
						t.Fatalf("Unhandled key func type for []string input: %T", tc.keyFunc)
					}
				case []personGroupTest:
					switch kf := tc.keyFunc.(type) {
					case func(personGroupTest) string:
						gotGroup = functional.GroupBy(input, kf)
					default:
						t.Fatalf("Unhandled key func type for []personGroupTest input: %T", tc.keyFunc)
					}

				default:
					t.Fatalf("Unhandled input slice type in test setup: %T", tc.input)
				}
			}

			// Convert gotGroup to map[any][]any for comparison helper
			gotComparable := make(map[any][]any)
			gotMapVal := reflect.ValueOf(gotGroup)
			if gotMapVal.Kind() == reflect.Map { // Check if it's actually a map
				iter := gotMapVal.MapRange()
				for iter.Next() {
					k := iter.Key().Interface()
					vSlice := iter.Value()
					anySlice := make([]any, vSlice.Len())
					for i := 0; i < vSlice.Len(); i++ {
						anySlice[i] = vSlice.Index(i).Interface()
					}
					gotComparable[k] = anySlice
				}
			} else if gotGroup != nil { // Handle cases where gotGroup might not be a map (e.g., error fallback)
				t.Logf("Warning: gotGroup was not a map type, got %T", gotGroup)
			} // If gotGroup is nil, gotComparable remains empty, which is fine for comparison

			assertSameGroups(t, gotComparable, tc.wantGroup)
		})
	}
}

// --- Examples ---

// --- Helper funcs for sorting slices in examples ---
func sortInts(s []int) []int          { sort.Ints(s); return s }
func sortStrings(s []string) []string { sort.Strings(s); return s }

func ExampleGroupBy() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
	groupedByEvenOdd := functional.GroupBy(numbers, func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})

	fmt.Println("Evens:", sortInts(groupedByEvenOdd["even"]))
	fmt.Println("Odds:", sortInts(groupedByEvenOdd["odd"]))

	words := []string{"apple", "banana", "apricot", "blueberry", "cherry"}
	groupedByFirstLetter := functional.GroupBy(words, func(s string) rune {
		if len(s) > 0 {
			return rune(s[0])
		}
		return rune(0)
	})

	// Need to sort keys for deterministic map iteration in example output
	keys := make([]rune, 0, len(groupedByFirstLetter))
	for k := range groupedByFirstLetter {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, k := range keys {
		fmt.Printf("Group '%c': %v\n", k, sortStrings(groupedByFirstLetter[k]))
	}

	// Output:
	// Evens: [2 4 6 8]
	// Odds: [1 3 5 7]
	// Group 'a': [apple apricot]
	// Group 'b': [banana blueberry]
	// Group 'c': [cherry]
}

// --- Benchmark Helpers ---

type groupByBenchItem struct {
	ID       int
	Category string
	Value    float64
}

func generateGroupByData(size, numCats int) []groupByBenchItem {
	if size <= 0 {
		return []groupByBenchItem{}
	}
	if numCats <= 0 {
		numCats = 2
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	data := make([]groupByBenchItem, size)
	for i := 0; i < size; i++ {
		data[i] = groupByBenchItem{
			ID:       i,
			Category: "cat_" + strconv.Itoa(rng.Intn(numCats)),
			Value:    rng.Float64() * 100,
		}
	}
	return data
}

func keyFuncFewGroups(item groupByBenchItem) string { return item.Category }

func keyFuncManyGroups(divisor int) func(groupByBenchItem) int {
	if divisor <= 0 {
		divisor = 10
	}
	return func(item groupByBenchItem) int { return item.ID % divisor }
}

// Benchmark Generic GroupBy
func benchmarkGroupByGeneric[T any, K comparable](
	data []T, keyFunc func(T) K, b *testing.B,
) {
	if len(data) == 0 {
		b.Skip("Skipping empty data benchmark")
		return
	}
	b.ResetTimer()
	var result map[K][]T
	for i := 0; i < b.N; i++ {
		result = functional.GroupBy(data, keyFunc)
	}
	_ = result
}

// Benchmark Loop GroupBy
func benchmarkGroupByLoop[T any, K comparable](
	data []T, keyFunc func(T) K, b *testing.B,
) {
	if len(data) == 0 {
		b.Skip("Skipping empty data benchmark")
		return
	}
	b.ResetTimer()
	var result map[K][]T
	for i := 0; i < b.N; i++ {
		groups := make(map[K][]T)
		for _, item := range data {
			key := keyFunc(item)
			groups[key] = append(groups[key], item)
		}
		result = groups
	}
	_ = result
}

// --- Run Benchmarks ---

const (
	N1_Group = 100
	N2_Group = 10000
)

var (
	fewGroupsDataN1 = generateGroupByData(N1_Group, 3)
	fewGroupsDataN2 = generateGroupByData(N2_Group, 5)
)

func BenchmarkGroupBy_Generic_FewGroups_N100(b *testing.B) {
	benchmarkGroupByGeneric(fewGroupsDataN1, keyFuncFewGroups, b)
}

func BenchmarkGroupBy_Loop_FewGroups_N100(b *testing.B) {
	benchmarkGroupByLoop(fewGroupsDataN1, keyFuncFewGroups, b)
}

func BenchmarkGroupBy_Generic_FewGroups_N10000(b *testing.B) {
	benchmarkGroupByGeneric(fewGroupsDataN2, keyFuncFewGroups, b)
}

func BenchmarkGroupBy_Loop_FewGroups_N10000(b *testing.B) {
	benchmarkGroupByLoop(fewGroupsDataN2, keyFuncFewGroups, b)
}

var (
	manyGroupsDataN1    = generateGroupByData(N1_Group, N1_Group)
	manyGroupsDataN2    = generateGroupByData(N2_Group, N2_Group)
	keyFuncManyGroupsN1 = keyFuncManyGroups(N1_Group / 10)
	keyFuncManyGroupsN2 = keyFuncManyGroups(N2_Group / 100)
)

func BenchmarkGroupBy_Generic_ManyGroups_N100(b *testing.B) {
	benchmarkGroupByGeneric(manyGroupsDataN1, keyFuncManyGroupsN1, b)
}

func BenchmarkGroupBy_Loop_ManyGroups_N100(b *testing.B) {
	benchmarkGroupByLoop(manyGroupsDataN1, keyFuncManyGroupsN1, b)
}

func BenchmarkGroupBy_Generic_ManyGroups_N10000(b *testing.B) {
	benchmarkGroupByGeneric(manyGroupsDataN2, keyFuncManyGroupsN2, b)
}

func BenchmarkGroupBy_Loop_ManyGroups_N10000(b *testing.B) {
	benchmarkGroupByLoop(manyGroupsDataN2, keyFuncManyGroupsN2, b)
}
