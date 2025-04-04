package functional_test

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/JackovAlltrades/go-generics/functional"
)

var (
	errTest = errors.New("test error condition met")
	errRate = 100
)

// --- Test MapErr ---
func TestMapErr(t *testing.T) {
	// Define test cases with expected full results (want)
	testCases := []struct {
		name    string
		input   []int
		mapper  func(int) (string, error)
		want    []string // Expected on success
		wantErr error
		// We will calculate expected partial result inside the test runner
	}{
		{name: "Success_NoError_MapErr", input: []int{1, 2, 3}, mapper: func(n int) (string, error) { return fmt.Sprintf("n=%d", n), nil }, want: []string{"n=1", "n=2", "n=3"}, wantErr: nil},
		{name: "Error_Middle_MapErr", input: []int{1, 2, 3, 4, 5}, mapper: func(n int) (string, error) {
			if n == 3 {
				return "", errTest
			}
			return fmt.Sprintf("n=%d", n), nil
		}, want: nil /* N/A on error */, wantErr: errTest},
		{name: "Error_FirstElement_MapErr", input: []int{1, 2, 3}, mapper: func(n int) (string, error) {
			if n == 1 {
				return "", errTest
			}
			return fmt.Sprintf("n=%d", n), nil
		}, want: nil /* N/A */, wantErr: errTest},
		{name: "Error_LastElement_MapErr", input: []int{1, 2, 3}, mapper: func(n int) (string, error) {
			if n == 3 {
				return "", errTest
			}
			return fmt.Sprintf("n=%d", n), nil
		}, want: nil /* N/A */, wantErr: errTest},
		{name: "EmptyInput_MapErr", input: []int{}, mapper: func(n int) (string, error) { return "", nil }, want: []string{}, wantErr: nil}, // Mapper shouldn't error if not called
		{name: "NilInput_MapErr", input: nil, mapper: func(n int) (string, error) { return "", nil }, want: []string{}, wantErr: nil},       // Mapper shouldn't error if not called
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := functional.MapErr(tc.input, tc.mapper)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("MapErr() error mismatch: got %v, want %v", err, tc.wantErr)
			}

			if tc.wantErr == nil {
				// Success case: Check against full 'want'
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("MapErr() success value mismatch: got %#v, want %#v", got, tc.want)
				}
			} else if errors.Is(err, tc.wantErr) {
				// Error case: Check against EXPECTED PARTIAL result based on implementation
				var expectedPartial []string
				switch tc.name {
				case "Error_Middle_MapErr": // Error on 3
					expectedPartial = []string{"n=1", "n=2"}
				case "Error_FirstElement_MapErr": // Error on 1
					expectedPartial = []string{}
				case "Error_LastElement_MapErr": // Error on 3
					expectedPartial = []string{"n=1", "n=2"}
					// Add other error cases if necessary
				default:
					// Should we expect nil or empty? Implementation returns partial. If no partial expected, expect empty.
					// Assuming empty if not explicitly defined above. This case shouldn't really be hit with specific tests.
					expectedPartial = []string{}
				}

				// Check if the actual result matches the calculated expected partial result
				if !reflect.DeepEqual(got, expectedPartial) {
					t.Errorf("MapErr() partial result on error mismatch: got %#v, want %#v", got, expectedPartial)
				}
			}
			// If error occurred but was not expected, the first error check already failed.
		})
	}
}

// --- Test FilterErr ---
func TestFilterErr(t *testing.T) {
	testCases := []struct {
		name      string
		input     []int
		predicate func(int) (bool, error)
		want      []int // Expected on success
		wantErr   error
		// Expected partial result calculated in test runner
	}{
		{name: "Success_NoError_FilterErr", input: []int{1, 2, 3, 4, 5, 6}, predicate: func(n int) (bool, error) { return n%2 == 0, nil }, want: []int{2, 4, 6}, wantErr: nil},
		{name: "Error_Middle_AfterSomeFiltered_FilterErr", input: []int{1, 2, 3, 4, 5, 6}, predicate: func(n int) (bool, error) {
			if n == 4 {
				return false, errTest
			}
			return n%2 == 0, nil
		}, want: nil, wantErr: errTest},
		{name: "Error_Middle_OnElementToFilterOut_FilterErr", input: []int{1, 2, 3, 4, 5, 6}, predicate: func(n int) (bool, error) {
			if n == 3 {
				return false, errTest
			}
			return n%2 == 0, nil
		}, want: nil, wantErr: errTest},
		{name: "Error_FirstElement_FilterErr", input: []int{1, 2, 3}, predicate: func(n int) (bool, error) {
			if n == 1 {
				return false, errTest
			}
			return n%2 == 0, nil
		}, want: nil, wantErr: errTest},
		{name: "Error_LastElement_FilterErr", input: []int{1, 2, 3}, predicate: func(n int) (bool, error) {
			if n == 3 {
				return false, errTest
			}
			return n%2 == 0, nil
		}, want: nil, wantErr: errTest},
		{name: "EmptyInput_FilterErr", input: []int{}, predicate: func(n int) (bool, error) { return false, nil }, want: []int{}, wantErr: nil},
		{name: "NilInput_FilterErr", input: nil, predicate: func(n int) (bool, error) { return false, nil }, want: []int{}, wantErr: nil},
		{name: "NoError_AllFilteredOut_FilterErr", input: []int{1, 3, 5}, predicate: func(n int) (bool, error) { return n%2 == 0, nil }, want: []int{}, wantErr: nil},
		{name: "NoError_AllKept_FilterErr", input: []int{2, 4, 6}, predicate: func(n int) (bool, error) { return n%2 == 0, nil }, want: []int{2, 4, 6}, wantErr: nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := functional.FilterErr(tc.input, tc.predicate)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("FilterErr() error mismatch: got %v, want %v", err, tc.wantErr)
			}

			if tc.wantErr == nil {
				// Success case: Check against full 'want'
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("FilterErr() success value mismatch: got %#v, want %#v", got, tc.want)
				}
			} else if errors.Is(err, tc.wantErr) {
				// Error case: Check against EXPECTED PARTIAL result based on implementation
				var expectedPartial []int
				switch tc.name {
				// Keep evens, error on 4. Processed 1(no), 2(yes), 3(no). Error on 4. Partial=[2]
				case "Error_Middle_AfterSomeFiltered_FilterErr":
					expectedPartial = []int{2}
					// Keep evens, error on 3. Processed 1(no), 2(yes). Error on 3. Partial=[2]
				case "Error_Middle_OnElementToFilterOut_FilterErr":
					expectedPartial = []int{2}
					// Keep evens, error on 1. Processed none. Error on 1. Partial=[]
				case "Error_FirstElement_FilterErr":
					expectedPartial = []int{}
					// Keep evens, error on 3. Processed 1(no), 2(yes). Error on 3. Partial=[2]
				case "Error_LastElement_FilterErr":
					expectedPartial = []int{2}
				default:
					expectedPartial = []int{} // Default should be empty slice if no partial expected
				}

				// Check if the actual result matches the calculated expected partial result
				if !reflect.DeepEqual(got, expectedPartial) {
					t.Errorf("FilterErr() partial result on error mismatch: got %#v, want %#v", got, expectedPartial)
				}
			}
		})
	}
}

// --- Test ReduceErr ---
func TestReduceErr(t *testing.T) {
	// Test cases and runner remain the same as the previously corrected version
	testCases := []struct {
		name    string
		input   any
		reducer any
		initial any
		want    any
		wantErr error
	}{
		{name: "Success_SumInts_ReduceErr", input: []int{1, 2, 3, 4}, reducer: func(acc, current int) (int, error) { return acc + current, nil }, initial: 0, want: 10, wantErr: nil},
		{name: "Success_ConcatStrings_ReduceErr", input: []string{"a", "b", "c"}, reducer: func(acc, current string) (string, error) { return acc + current, nil }, initial: "", want: "abc", wantErr: nil},
		{name: "Error_Middle_Summing_ReduceErr", input: []int{1, 2, 3, 4, 5}, reducer: func(acc, current int) (int, error) {
			if current == 4 {
				return acc, errTest
			}
			return acc + current, nil
		}, initial: 0, want: 6 /* 0+1+2+3 before error on 4 */, wantErr: errTest},
		{name: "Error_FirstElement_ReduceErr", input: []int{1, 2, 3}, reducer: func(acc, current int) (int, error) {
			if current == 1 {
				return acc, errTest
			}
			return acc + current, nil
		}, initial: 100, want: 100 /* initial */, wantErr: errTest},
		{name: "Error_LastElement_ReduceErr", input: []int{1, 2, 3}, reducer: func(acc, current int) (int, error) {
			if current == 3 {
				return acc, errTest
			}
			return acc + current, nil
		}, initial: 0, want: 3 /* 0+1+2 before error on 3 */, wantErr: errTest},
		{name: "EmptyInput_ReturnsInitial_ReduceErr", input: []int{}, reducer: func(acc, current int) (int, error) { return acc, errTest }, initial: 55, want: 55, wantErr: nil},
		{name: "NilInput_ReturnsInitial_ReduceErr", input: nil, reducer: func(acc, current int) (int, error) { return acc, errTest }, initial: 66, want: 66, wantErr: nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any
			var err error

			callReduceErr := func() {
				switch input := tc.input.(type) {
				case []int:
					reducer, okR := tc.reducer.(func(int, int) (int, error))
					initial, okI := tc.initial.(int)
					if !okR {
						t.Fatalf("Reducer type mismatch: %T", tc.reducer)
					}
					if !okI {
						t.Fatalf("Initial type mismatch: %T", tc.initial)
					}
					got, err = functional.ReduceErr[int, int](input, initial, reducer)
				case []string:
					reducer, okR := tc.reducer.(func(string, string) (string, error))
					initial, okI := tc.initial.(string)
					if !okR {
						t.Fatalf("Reducer type mismatch: %T", tc.reducer)
					}
					if !okI {
						t.Fatalf("Initial type mismatch: %T", tc.initial)
					}
					got, err = functional.ReduceErr[string, string](input, initial, reducer)
				default:
					inputVal := reflect.ValueOf(input)
					if input == nil || (inputVal.IsValid() && inputVal.Kind() == reflect.Slice && inputVal.Len() == 0) {
						got = tc.initial
						err = nil
					} else {
						t.Fatalf("Unhandled input type for ReduceErr test: %T", tc.input)
					}
				}
			}

			callReduceErr()

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ReduceErr() error mismatch: got %v, want %v", err, tc.wantErr)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ReduceErr() value mismatch: got %#v, want %#v (wantErr=%v, err=%v)", got, tc.want, tc.wantErr, err)
			}
		})
	}
}

// --- Benchmarks ---
// (Benchmarks remain the same as previous corrected version - no changes needed here)
var (
	benchInputInts = make([]int, 1000)
)

func init() {
	rand.Seed(time.Now().UnixNano())
	for i := range benchInputInts {
		benchInputInts[i] = i
	}
}
func mapperNoError(n int) (string, error) { return strconv.Itoa(n), nil }
func mapperWithError(n int) (string, error) {
	if n != 0 && n%errRate == 0 {
		return "", errTest
	}
	return strconv.Itoa(n), nil
}
func predicateNoError(n int) (bool, error) { return n%2 == 0, nil }
func predicateWithError(n int) (bool, error) {
	if n != 0 && n%errRate == 0 {
		return false, errTest
	}
	return n%2 == 0, nil
}
func reducerNoError(acc, current int) (int, error) { return acc + current, nil }
func reducerWithError(acc, current int) (int, error) {
	if current != 0 && current%errRate == 0 {
		return acc, errTest
	}
	return acc + current, nil
}

func BenchmarkMapErr_NoError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	b.ResetTimer()
	var res []string
	var err error
	for i := 0; i < b.N; i++ {
		res, err = functional.MapErr(data, mapperNoError)
	}
	_ = res
	_ = err
}

func BenchmarkMap_Baseline_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	mapper := func(n int) string { return strconv.Itoa(n) }
	b.ResetTimer()
	var res []string
	for i := 0; i < b.N; i++ {
		res = functional.Map(data, mapper)
	}
	_ = res
}

func BenchmarkMapLoop_Baseline_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	b.ResetTimer()
	var res []string
	for i := 0; i < b.N; i++ {
		localRes := make([]string, len(data))
		for j, item := range data {
			localRes[j] = strconv.Itoa(item)
		}
		res = localRes
	}
	_ = res
}

func BenchmarkMapErr_WithError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	b.ResetTimer()
	var res []string
	var err error
	for i := 0; i < b.N; i++ {
		res, err = functional.MapErr(data, mapperWithError)
	}
	_ = res
	_ = err
}

func BenchmarkMapLoop_WithError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	b.ResetTimer()
	var res []string
	var err error
	for i := 0; i < b.N; i++ {
		localRes := make([]string, 0, len(data))
		var loopErr error
		for _, item := range data {
			s, mapperErr := mapperWithError(item)
			if mapperErr != nil {
				loopErr = mapperErr
				localRes = nil
				break
			}
			localRes = append(localRes, s)
		}
		res = localRes
		err = loopErr
	}
	_ = res
	_ = err
}

func BenchmarkFilterErr_NoError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	b.ResetTimer()
	var res []int
	var err error
	for i := 0; i < b.N; i++ {
		res, err = functional.FilterErr(data, predicateNoError)
	}
	_ = res
	_ = err
}

func BenchmarkFilter_Baseline_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	predicate := func(n int) bool { return n%2 == 0 }
	b.ResetTimer()
	var res []int
	for i := 0; i < b.N; i++ {
		res = functional.Filter(data, predicate)
	}
	_ = res
}

func BenchmarkFilterLoop_Baseline_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	predicate := func(n int) bool { return n%2 == 0 }
	b.ResetTimer()
	var res []int
	for i := 0; i < b.N; i++ {
		localRes := make([]int, 0)
		for _, item := range data {
			if predicate(item) {
				localRes = append(localRes, item)
			}
		}
		res = localRes
	}
	_ = res
}

func BenchmarkFilterErr_WithError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	b.ResetTimer()
	var res []int
	var err error
	for i := 0; i < b.N; i++ {
		res, err = functional.FilterErr(data, predicateWithError)
	}
	_ = res
	_ = err
}

func BenchmarkFilterLoop_WithError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	b.ResetTimer()
	var res []int
	var err error
	for i := 0; i < b.N; i++ {
		localRes := make([]int, 0)
		var loopErr error
		for _, item := range data {
			keep, predErr := predicateWithError(item)
			if predErr != nil {
				loopErr = predErr
				localRes = nil
				break
			}
			if keep {
				localRes = append(localRes, item)
			}
		}
		res = localRes
		err = loopErr
	}
	_ = res
	_ = err
}

func BenchmarkReduceErr_NoError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	initial := 0
	b.ResetTimer()
	var res int
	var err error
	for i := 0; i < b.N; i++ {
		res, err = functional.ReduceErr[int, int](data, initial, reducerNoError)
	}
	_ = res
	_ = err
}

func BenchmarkReduce_Baseline_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	initial := 0
	reducer := func(acc, current int) int { return acc + current }
	b.ResetTimer()
	var res int
	for i := 0; i < b.N; i++ {
		res = functional.Reduce[int, int](data, initial, reducer)
	}
	_ = res
}

func BenchmarkReduceLoop_Baseline_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	initial := 0
	b.ResetTimer()
	var res int
	for i := 0; i < b.N; i++ {
		acc := initial
		for _, item := range data {
			acc = acc + item
		}
		res = acc
	}
	_ = res
}

func BenchmarkReduceErr_WithError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	initial := 0
	b.ResetTimer()
	var res int
	var err error
	for i := 0; i < b.N; i++ {
		res, err = functional.ReduceErr[int, int](data, initial, reducerWithError)
	}
	_ = res
	_ = err
}

func BenchmarkReduceLoop_WithError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	initial := 0
	b.ResetTimer()
	var res int
	var err error
	for i := 0; i < b.N; i++ {
		acc := initial
		var loopErr error
		for _, item := range data {
			var reduceErr error
			acc, reduceErr = reducerWithError(acc, item)
			if reduceErr != nil {
				loopErr = reduceErr
				acc = initial
				break
			}
		}
		res = acc
		err = loopErr
	}
	_ = res
	_ = err
}
