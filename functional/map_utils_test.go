package functional_test // Note the _test suffix for the package name

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/JackovAlltrades/go-generics/functional" // Adjust import path if needed
)

// --- Structs Used in Tests ---
type personValue struct {
	Name string
	Age  int
}
type personToMap struct {
	Name string
	Age  int
}

// --- Helper functions ---
func sortSliceAny(slice any) bool {
	switch s := slice.(type) {
	case []int:
		sort.Ints(s)
		return true
	case []string:
		sort.Strings(s)
		return true
	case []float64:
		sort.Float64s(s)
		return true
	case []personValue:
		// Example sorting for personValue (e.g., by Age)
		sort.Slice(s, func(i, j int) bool { return s[i].Age < s[j].Age })
		return true
	default:
		return false
	}
}

func compareUnorderedSlices[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("Unordered slice length mismatch: len(got)=%d, len(want)=%d (got=%#v, want=%#v)", len(got), len(want), got, want)
		return
	}
	if len(got) == 0 {
		return
	} // Both empty, they are equal

	gotMap := make(map[T]int, len(got))
	wantMap := make(map[T]int, len(want))
	for _, v := range got {
		gotMap[v]++
	}
	for _, v := range want {
		wantMap[v]++
	}

	if !reflect.DeepEqual(gotMap, wantMap) {
		t.Errorf("Unordered slice element mismatch: got elements=%v, want elements=%v (counts differ or elements differ)", got, want)
	}
}

// --- Test Keys ---
func TestKeys(t *testing.T) {
	t.Run("Keys_IntKeys", func(t *testing.T) {
		input := map[int]string{1: "a", 2: "b", 3: "c"}
		want := []int{1, 2, 3}
		got := functional.Keys(input) // Calling the function from the functional package
		sortSliceAny(got)
		sortSliceAny(want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Keys() = %#v, want %#v", got, want)
		}
	})

	t.Run("Keys_StringKeys", func(t *testing.T) {
		input := map[string]int{"one": 1, "two": 2, "three": 3}
		want := []string{"one", "two", "three"}
		got := functional.Keys(input) // Calling the function from the functional package
		sortSliceAny(got)
		sortSliceAny(want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Keys() = %#v, want %#v", got, want)
		}
	})

	t.Run("Keys_EmptyMap", func(t *testing.T) {
		input := map[int]int{}
		want := []int{}
		got := functional.Keys(input) // Calling the function from the functional package
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Keys() = %#v, want %#v", got, want)
		}
	})

	t.Run("Keys_NilMap", func(t *testing.T) {
		var input map[string]bool = nil
		want := []string{}
		got := functional.Keys(input) // Calling the function from the functional package
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Keys() = %#v, want %#v", got, want)
		}
	})
}

// --- Test Values ---
func TestValues(t *testing.T) {
	t.Run("Values_IntValues", func(t *testing.T) {
		input := map[string]int{"a": 1, "b": 2, "c": 3}
		want := []int{1, 2, 3}
		got := functional.Values(input) // Calling the function from the functional package
		sortSliceAny(got)
		sortSliceAny(want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Values() = %#v, want %#v", got, want)
		}
	})

	t.Run("Values_StringValues", func(t *testing.T) {
		input := map[int]string{1: "one", 2: "two", 3: "three"}
		want := []string{"one", "two", "three"}
		got := functional.Values(input) // Calling the function from the functional package
		sortSliceAny(got)
		sortSliceAny(want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Values() = %#v, want %#v", got, want)
		}
	})

	t.Run("Values_ComparableStructValues", func(t *testing.T) {
		input := map[int]personValue{1: {"A", 20}, 2: {"B", 30}, 3: {"C", 25}}
		want := []personValue{{"A", 20}, {"B", 30}, {"C", 25}}
		got := functional.Values(input) // Calling the function from the functional package
		compareUnorderedSlices(t, got, want)
	})

	t.Run("Values_EmptyMap", func(t *testing.T) {
		input := map[string]float64{}
		want := []float64{}
		got := functional.Values(input) // Calling the function from the functional package
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Values() = %#v, want %#v", got, want)
		}
	})

	t.Run("Values_NilMap", func(t *testing.T) {
		var input map[int]string = nil
		want := []string{}
		got := functional.Values(input) // Calling the function from the functional package
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Values() = %#v, want %#v", got, want)
		}
	})
}

// --- Test MapToSlice ---
func TestMapToSlice(t *testing.T) {
	mapperIntStr := func(k int, v string) string { return fmt.Sprintf("%d:%s", k, v) }
	mapperStrInt := func(k string, v int) int { return v * len(k) }
	mapperIntPerson := func(k int, v personToMap) string { return fmt.Sprintf("%s (%d)", v.Name, v.Age+k) }

	t.Run("MapIntStringToString", func(t *testing.T) {
		input := map[int]string{1: "a", 2: "bb", 3: "ccc"}
		want := []string{"1:a", "2:bb", "3:ccc"}
		got := functional.MapToSlice(input, mapperIntStr) // Calling the function from the functional package
		sortSliceAny(got)
		sortSliceAny(want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("MapToSlice() = %#v, want %#v", got, want)
		}
	})

	t.Run("MapStringIntToInt", func(t *testing.T) {
		input := map[string]int{"one": 1, "two": 2, "three": 3}
		want := []int{3, 6, 15}
		got := functional.MapToSlice(input, mapperStrInt) // Calling the function from the functional package
		sortSliceAny(got)
		sortSliceAny(want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("MapToSlice() = %#v, want %#v", got, want)
		}
	})

	t.Run("MapIntPersonToString", func(t *testing.T) {
		input := map[int]personToMap{10: {"A", 20}, 20: {"B", 30}}
		want := []string{"A (30)", "B (50)"}
		got := functional.MapToSlice(input, mapperIntPerson) // Calling the function from the functional package
		sortSliceAny(got)
		sortSliceAny(want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("MapToSlice() = %#v, want %#v", got, want)
		}
	})

	t.Run("EmptyMap", func(t *testing.T) {
		input := map[int]string{}
		want := []string{}
		got := functional.MapToSlice(input, mapperIntStr) // Calling the function from the functional package
		if !reflect.DeepEqual(got, want) {
			t.Errorf("MapToSlice() = %#v, want %#v", got, want)
		}
	})

	t.Run("NilMap", func(t *testing.T) {
		var input map[int]string = nil
		want := []string{}
		got := functional.MapToSlice(input, mapperIntStr) // Calling the function from the functional package
		if !reflect.DeepEqual(got, want) {
			t.Errorf("MapToSlice() = %#v, want %#v", got, want)
		}
	})
}

// --- Map Utils Examples ---
func ExampleKeys() {
	m := map[string]int{"apple": 1, "banana": 2, "cherry": 3}
	keys := functional.Keys(m) // Calling the function from the functional package
	sort.Strings(keys)
	fmt.Println("Keys:", keys)
	emptyMap := map[int]bool{}
	emptyKeys := functional.Keys(emptyMap) // Calling the function from the functional package
	fmt.Println("Empty Map Keys:", emptyKeys)
	var nilMap map[float64]string = nil
	nilKeys := functional.Keys(nilMap) // Calling the function from the functional package
	fmt.Println("Nil Map Keys:", nilKeys)
	// Output:
	// Keys: [apple banana cherry]
	// Empty Map Keys: []
	// Nil Map Keys: []
}

func ExampleValues() {
	m := map[string]int{"one": 10, "two": 20, "three": 30}
	values := functional.Values(m) // Calling the function from the functional package
	sort.Ints(values)
	fmt.Println("Values:", values)
	emptyMap := map[int]string{}
	emptyValues := functional.Values(emptyMap) // Calling the function from the functional package
	fmt.Println("Empty Map Values:", emptyValues)
	var nilMap map[bool]int = nil
	nilValues := functional.Values(nilMap) // Calling the function from the functional package
	fmt.Println("Nil Map Values:", nilValues)
	// Output:
	// Values: [10 20 30]
	// Empty Map Values: []
	// Nil Map Values: []
}

func ExampleMapToSlice() {
	scores := map[string]int{"Alice": 85, "Bob": 92, "Charlie": 78}
	descriptions := functional.MapToSlice(scores, func(name string, score int) string { // Calling the function from the functional package
		return fmt.Sprintf("%s scored %d", name, score)
	})
	sort.Strings(descriptions)
	fmt.Println("Descriptions:", descriptions)
	coords := map[int]string{1: "x", 2: "y", 3: "z"}
	byteSlices := functional.MapToSlice(coords, func(k int, v string) []byte { // Calling the function from the functional package
		return []byte(fmt.Sprintf("%d-%s", k, v))
	})
	sort.Slice(byteSlices, func(i, j int) bool { return string(byteSlices[i]) < string(byteSlices[j]) })
	byteStrings := make([]string, len(byteSlices))
	for i, b := range byteSlices {
		byteStrings[i] = string(b)
	}
	fmt.Println("Byte Slices:", byteStrings)
	emptyMap := map[float64]bool{}
	emptyResult := functional.MapToSlice(emptyMap, func(k float64, v bool) string { return "never" }) // Calling the function from the functional package
	fmt.Println("Empty Map Result:", emptyResult)
	// Output:
	// Descriptions: [Alice scored 85 Bob scored 92 Charlie scored 78]
	// Byte Slices: [1-x 2-y 3-z]
	// Empty Map Result: []
}

// --- Benchmarks ---

// Helper to generate maps
func generateMap[K comparable, V any](size int, keyGen func(int) K, valGen func(int) V) map[K]V {
	m := make(map[K]V, size)
	for i := 0; i < size; i++ {
		m[keyGen(i)] = valGen(i)
	}
	return m
}

var (
	intKeyGen    = func(i int) int { return i }
	stringKeyGen = func(i int) string { return fmt.Sprintf("key_%d", i) }
	intValGen    = func(i int) int { return i*10 + 1 }
	stringValGen = func(i int) string { return fmt.Sprintf("val:%d", i) }
)

// --- Benchmark Keys ---
// Specific versions for each map type - These are helpers called by Benchmark* functions
func benchmarkKeysGeneric_IntStr(m map[int]string, b *testing.B) {
	if m == nil {
		b.Skip("nil map")
		return
	}
	b.ResetTimer()
	var result []int
	for i := 0; i < b.N; i++ {
		result = functional.Keys(m) // Calling the function from the functional package
	}
	_ = result // Keep result used
}

func benchmarkKeysLoop_IntStr(m map[int]string, b *testing.B) {
	if m == nil {
		b.Skip("nil map")
		return
	}
	b.ResetTimer()
	var result []int
	for i := 0; i < b.N; i++ {
		keys := make([]int, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		result = keys
	}
	_ = result // Keep result used
}

func benchmarkKeysGeneric_StrInt(m map[string]int, b *testing.B) {
	if m == nil {
		b.Skip("nil map")
		return
	}
	b.ResetTimer()
	var result []string
	for i := 0; i < b.N; i++ {
		result = functional.Keys(m) // Calling the function from the functional package
	}
	_ = result // Keep result used
}

func benchmarkKeysLoop_StrInt(m map[string]int, b *testing.B) {
	if m == nil {
		b.Skip("nil map")
		return
	}
	b.ResetTimer()
	var result []string
	for i := 0; i < b.N; i++ {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		result = keys
	}
	_ = result // Keep result used
}

var (
	keysMapIntStr_N1000 = generateMap(1000, intKeyGen, stringValGen)
	keysMapStrInt_N1000 = generateMap(1000, stringKeyGen, intValGen)
)

// Top-level Benchmark functions (these are discovered by `go test`)
func BenchmarkKeys_Generic_IntStr_N1000(b *testing.B) {
	benchmarkKeysGeneric_IntStr(keysMapIntStr_N1000, b)
}
func BenchmarkKeys_Loop_IntStr_N1000(b *testing.B) { benchmarkKeysLoop_IntStr(keysMapIntStr_N1000, b) }
func BenchmarkKeys_Generic_StrInt_N1000(b *testing.B) {
	benchmarkKeysGeneric_StrInt(keysMapStrInt_N1000, b)
}
func BenchmarkKeys_Loop_StrInt_N1000(b *testing.B) { benchmarkKeysLoop_StrInt(keysMapStrInt_N1000, b) }

// --- Benchmark Values ---
// Specific versions for each map type - These are helpers
func benchmarkValuesGeneric_IntStr(m map[int]string, b *testing.B) {
	if m == nil {
		b.Skip("nil map")
		return
	}
	b.ResetTimer()
	var result []string
	for i := 0; i < b.N; i++ {
		result = functional.Values(m) // Calling the function from the functional package
	}
	_ = result
}

func benchmarkValuesLoop_IntStr(m map[int]string, b *testing.B) {
	if m == nil {
		b.Skip("nil map")
		return
	}
	b.ResetTimer()
	var result []string
	for i := 0; i < b.N; i++ {
		values := make([]string, 0, len(m))
		for _, v := range m {
			values = append(values, v)
		}
		result = values
	}
	_ = result
}

func benchmarkValuesGeneric_StrInt(m map[string]int, b *testing.B) {
	if m == nil {
		b.Skip("nil map")
		return
	}
	b.ResetTimer()
	var result []int
	for i := 0; i < b.N; i++ {
		result = functional.Values(m) // Calling the function from the functional package
	}
	_ = result
}

func benchmarkValuesLoop_StrInt(m map[string]int, b *testing.B) {
	if m == nil {
		b.Skip("nil map")
		return
	}
	b.ResetTimer()
	var result []int
	for i := 0; i < b.N; i++ {
		values := make([]int, 0, len(m))
		for _, v := range m {
			values = append(values, v)
		}
		result = values
	}
	_ = result
}

// Top-level Benchmark functions
func BenchmarkValues_Generic_IntStr_N1000(b *testing.B) {
	benchmarkValuesGeneric_IntStr(keysMapIntStr_N1000, b)
}

func BenchmarkValues_Loop_IntStr_N1000(b *testing.B) {
	benchmarkValuesLoop_IntStr(keysMapIntStr_N1000, b)
}

func BenchmarkValues_Generic_StrInt_N1000(b *testing.B) {
	benchmarkValuesGeneric_StrInt(keysMapStrInt_N1000, b)
}

func BenchmarkValues_Loop_StrInt_N1000(b *testing.B) {
	benchmarkValuesLoop_StrInt(keysMapStrInt_N1000, b)
}

// --- Benchmark MapToSlice ---
// Specific versions for each map/mapper type - These are helpers
func benchmarkMapToSliceGeneric_IntStr_ToStr(m map[int]string, fn func(int, string) string, b *testing.B) {
	if m == nil {
		b.Skip("nil map")
		return
	}
	b.ResetTimer()
	var result []string
	for i := 0; i < b.N; i++ {
		result = functional.MapToSlice(m, fn) // Calling the function from the functional package
	}
	_ = result
}

func benchmarkMapToSliceLoop_IntStr_ToStr(m map[int]string, fn func(int, string) string, b *testing.B) {
	if m == nil {
		b.Skip("nil map")
		return
	}
	b.ResetTimer()
	var result []string
	for i := 0; i < b.N; i++ {
		mapped := make([]string, 0, len(m))
		for k, v := range m {
			mapped = append(mapped, fn(k, v))
		}
		result = mapped
	}
	_ = result
}

func benchmarkMapToSliceGeneric_StrInt_ToInt(m map[string]int, fn func(string, int) int, b *testing.B) {
	if m == nil {
		b.Skip("nil map")
		return
	}
	b.ResetTimer()
	var result []int
	for i := 0; i < b.N; i++ {
		result = functional.MapToSlice(m, fn) // Calling the function from the functional package
	}
	_ = result
}

func benchmarkMapToSliceLoop_StrInt_ToInt(m map[string]int, fn func(string, int) int, b *testing.B) {
	if m == nil {
		b.Skip("nil map")
		return
	}
	b.ResetTimer()
	var result []int
	for i := 0; i < b.N; i++ {
		mapped := make([]int, 0, len(m))
		for k, v := range m {
			mapped = append(mapped, fn(k, v))
		}
		result = mapped
	}
	_ = result
}

// Define mappers for benchmarks
var (
	mapperIntStrToStr = func(k int, v string) string { return fmt.Sprintf("%d-%s", k, v) }
	mapperStrIntToInt = func(k string, v int) int { return v + len(k) }
)

// Top-level Benchmark functions
func BenchmarkMapToSlice_Generic_IntStr_ToStr_N1000(b *testing.B) {
	benchmarkMapToSliceGeneric_IntStr_ToStr(keysMapIntStr_N1000, mapperIntStrToStr, b)
}

func BenchmarkMapToSlice_Loop_IntStr_ToStr_N1000(b *testing.B) {
	benchmarkMapToSliceLoop_IntStr_ToStr(keysMapIntStr_N1000, mapperIntStrToStr, b)
}

func BenchmarkMapToSlice_Generic_StrInt_ToInt_N1000(b *testing.B) {
	benchmarkMapToSliceGeneric_StrInt_ToInt(keysMapStrInt_N1000, mapperStrIntToInt, b)
}

func BenchmarkMapToSlice_Loop_StrInt_ToInt_N1000(b *testing.B) {
	benchmarkMapToSliceLoop_StrInt_ToInt(keysMapStrInt_N1000, mapperStrIntToInt, b)
}
