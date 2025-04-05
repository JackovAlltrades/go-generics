package functional_test

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	// "fmt"

	"github.com/JackovAlltrades/go-generics/functional" // Adjust import path if needed
)

// --- Test Setup ---
// Uses errTestSentinel, errRateTest defined below
var (
	errTestSentinel = errors.New("test error condition met") // Use a sentinel error
	errRateTest     = 3                                      // Trigger error on elements divisible by 3 (for predictability)
)

// ptr func and person struct are defined in helpers_test.go

// --- Test MapErr ---

func TestMapErr(t *testing.T) {
	// Define testCases specific to MapErr
	testCases := []struct {
		name                string
		input               any
		mapper              any // func(T) (U, error)
		want                any // Expected result on success OR partial result on error
		wantErr             error
		checkPartialOnError bool // Flag to indicate we should verify the partial result
	}{
		{
			name:  "NoError_IntToString",
			input: []int{1, 2, 4}, // None divisible by 3
			mapper: func(n int) (string, error) {
				return strconv.Itoa(n), nil
			},
			want:    []string{"1", "2", "4"},
			wantErr: nil,
		},
		{
			name:  "Error_Middle_IntToString",
			input: []int{1, 2, 3, 4, 5}, // 3 will cause error
			mapper: func(n int) (string, error) {
				if n%errRateTest == 0 {
					return "", errTestSentinel
				}
				return strconv.Itoa(n), nil
			},
			want:                []string{"1", "2"}, // Partial result before error on '3'
			wantErr:             errTestSentinel,
			checkPartialOnError: true,
		},
		{
			name:  "Error_FirstElement_IntToString",
			input: []int{3, 4, 5}, // 3 will cause error immediately
			mapper: func(n int) (string, error) {
				if n%errRateTest == 0 {
					return "", errTestSentinel
				}
				return strconv.Itoa(n), nil
			},
			want:                []string{}, // Empty partial result
			wantErr:             errTestSentinel,
			checkPartialOnError: true,
		},
		{
			name:  "Error_LastElement_IntToString",
			input: []int{1, 2, 4, 3}, // 3 will cause error at the end
			mapper: func(n int) (string, error) {
				if n%errRateTest == 0 {
					return "", errTestSentinel
				}
				return strconv.Itoa(n), nil
			},
			want:                []string{"1", "2", "4"}, // Partial result
			wantErr:             errTestSentinel,
			checkPartialOnError: true,
		},
		{
			name:    "NilInput",
			input:   ([]int)(nil),
			mapper:  func(n int) (string, error) { return strconv.Itoa(n), nil },
			want:    []string{}, // Expect empty, non-nil slice for nil/empty input
			wantErr: nil,
		},
		{
			name:    "EmptyInput",
			input:   []int{},
			mapper:  func(n int) (string, error) { return strconv.Itoa(n), nil },
			want:    []string{}, // Expect empty, non-nil slice for nil/empty input
			wantErr: nil,
		},
		{
			name:  "StructInput_ToString_NoError",
			input: []person{{"A", 20}, {"B", 40}}, // Ages not divisible by 3
			mapper: func(p person) (string, error) {
				return p.Name + ":" + strconv.Itoa(p.Age), nil
			},
			want:    []string{"A:20", "B:40"},
			wantErr: nil,
		},
		{
			name:  "StructInput_ToString_WithError",
			input: []person{{"A", 20}, {"B", 30}, {"C", 40}}, // B's age 30 is divisible by 3
			mapper: func(p person) (string, error) {
				if p.Age%errRateTest == 0 {
					return "", errTestSentinel
				}
				return p.Name + ":" + strconv.Itoa(p.Age), nil
			},
			want:                []string{"A:20"}, // Partial result
			wantErr:             errTestSentinel,
			checkPartialOnError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any
			var err error

			// Type switch to call MapErr with correct types
			switch mapper := tc.mapper.(type) {
			case func(int) (string, error):
				in, ok := tc.input.([]int)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for func(int)(string, error)")
				}
				got, err = functional.MapErr[int, string](in, mapper)

			case func(person) (string, error):
				in, ok := tc.input.([]person)
				if tc.input == nil {
					in = nil
					ok = true
				} else if !ok {
					t.Fatalf("Input type mismatch for func(person)(string, error)")
				}
				got, err = functional.MapErr[person, string](in, mapper)

			default:
				t.Fatalf("Unhandled mapper type in MapErr test setup: %T", tc.mapper)
			}

			// Check if the error matches the expected error (using errors.Is for sentinel checks)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("MapErr() error mismatch: got %v, want %v", err, tc.wantErr)
			}
			// Check the result value
			if tc.wantErr == nil {
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("MapErr() success value mismatch: got %#v, want %#v", got, tc.want)
				}
			} else if tc.checkPartialOnError {
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("MapErr() partial result on error mismatch: got %#v, want %#v", got, tc.want)
				}
			}
		})
	}
}

// --- Test FilterErr ---

func TestFilterErr(t *testing.T) {
	// Define testCases specific to FilterErr
	testCases := []struct {
		name                string
		input               []int
		predicate           func(int) (bool, error)
		want                []int
		wantErr             error
		checkPartialOnError bool
	}{
		{
			name:  "NoError_FilterEvens",
			input: []int{1, 2, 4, 5}, // Input that doesn't trigger the sentinel
			predicate: func(n int) (bool, error) {
				if n%errRateTest == 0 {
					return false, errTestSentinel
				}
				return n%2 == 0, nil // Keep evens
			},
			want:    []int{2, 4}, // Only evens from the input
			wantErr: nil,
		},
		{
			name:  "Error_Middle_FilterEvens",
			input: []int{1, 2, 3, 4, 5, 6}, // 3 causes error
			predicate: func(n int) (bool, error) {
				if n%errRateTest == 0 {
					return false, errTestSentinel
				}
				return n%2 == 0, nil
			},
			want:                []int{2}, // Partial result before error
			wantErr:             errTestSentinel,
			checkPartialOnError: true,
		},
		{
			name:  "Error_Middle_FilteringOutElementWithError",
			input: []int{2, 4, 3, 6, 8}, // 3 causes error
			predicate: func(n int) (bool, error) {
				if n%errRateTest == 0 {
					return false, errTestSentinel
				}
				return n%2 == 0, nil
			},
			want:                []int{2, 4}, // Partial result before error
			wantErr:             errTestSentinel,
			checkPartialOnError: true,
		},
		{
			name:  "Error_FirstElement",
			input: []int{3, 2, 4}, // 3 causes error immediately
			predicate: func(n int) (bool, error) {
				if n%errRateTest == 0 {
					return false, errTestSentinel
				}
				return n%2 == 0, nil
			},
			want:                []int{}, // Empty partial
			wantErr:             errTestSentinel,
			checkPartialOnError: true,
		},
		{
			name:  "Error_LastElement",
			input: []int{2, 4, 5, 3}, // 3 causes error at the end
			predicate: func(n int) (bool, error) {
				if n%errRateTest == 0 {
					return false, errTestSentinel
				}
				return n%2 == 0, nil
			},
			want:                []int{2, 4}, // Partial result
			wantErr:             errTestSentinel,
			checkPartialOnError: true,
		},
		{name: "NilInput", input: nil, predicate: func(n int) (bool, error) { return n%2 == 0, nil }, want: []int{}, wantErr: nil},
		{name: "EmptyInput", input: []int{}, predicate: func(n int) (bool, error) { return n%2 == 0, nil }, want: []int{}, wantErr: nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := functional.FilterErr(tc.input, tc.predicate)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("FilterErr() error mismatch: got %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr == nil {
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("FilterErr() success value mismatch: got %#v, want %#v", got, tc.want)
				}
			} else if tc.checkPartialOnError {
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("FilterErr() partial result on error mismatch: got %#v, want %#v", got, tc.want)
				}
			}
		})
	}
}

// --- Test ReduceErr ---

func TestReduceErr(t *testing.T) {
	testCases := []struct {
		name    string
		input   any
		initial any
		reducer any
		want    any
		wantErr error
	}{
		{
			name: "NoError_Summing", input: []int{1, 2, 4, 5}, initial: 0, want: 12, wantErr: nil,
			reducer: func(acc, next int) (int, error) {
				if next%errRateTest == 0 {
					return acc, errTestSentinel
				}
				return acc + next, nil
			},
		},
		{
			name: "Error_Middle_Summing", input: []int{1, 2, 3, 4, 5}, initial: 0, want: 3, wantErr: errTestSentinel,
			reducer: func(acc, next int) (int, error) {
				if next%errRateTest == 0 {
					return acc, errTestSentinel
				}
				return acc + next, nil
			},
		},
		{
			name: "Error_FirstElement_Summing", input: []int{3, 1, 2}, initial: 10, want: 10, wantErr: errTestSentinel,
			reducer: func(acc, next int) (int, error) {
				if next%errRateTest == 0 {
					return acc, errTestSentinel
				}
				return acc + next, nil
			},
		},
		{
			name: "Error_LastElement_Summing", input: []int{1, 2, 4, 3}, initial: 0, want: 7, wantErr: errTestSentinel,
			reducer: func(acc, next int) (int, error) {
				if next%errRateTest == 0 {
					return acc, errTestSentinel
				}
				return acc + next, nil
			},
		},
		{
			name: "NoError_StringConcat", input: []string{"a", "b", "d"}, initial: "", want: "abd", wantErr: nil,
			reducer: func(acc, next string) (string, error) {
				if next == "c" {
					return acc, errTestSentinel
				}
				return acc + next, nil
			},
		},
		{
			name: "Error_StringConcat", input: []string{"a", "b", "c", "d"}, initial: "start:", want: "start:ab", wantErr: errTestSentinel,
			reducer: func(acc, next string) (string, error) {
				if next == "c" {
					return acc, errTestSentinel
				}
				return acc + next, nil
			},
		},
		{
			name: "NilInput", input: ([]int)(nil), initial: 100, want: 100, wantErr: nil,
			reducer: func(acc, next int) (int, error) { return acc + next, nil }, // Reducer type doesn't matter for nil input
		},
		{
			name: "EmptyInput", input: []string{}, initial: "empty", want: "empty", wantErr: nil, // Input type matches initial/reducer
			reducer: func(acc, next string) (string, error) { return acc + next, nil },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any
			var err error

			switch reducer := tc.reducer.(type) {
			case func(int, int) (int, error):
				inputSlice, okI := tc.input.([]int)
				if tc.input == nil {
					okI = true
					inputSlice = nil
				}
				if !okI {
					t.Fatalf("Input type mismatch for int reducer")
				}
				initialValue, okInit := tc.initial.(int)
				if !okInit {
					t.Fatalf("Initial type mismatch for int reducer")
				}
				got, err = functional.ReduceErr[int, int](inputSlice, initialValue, reducer)

			case func(string, string) (string, error):
				inputSlice, okI := tc.input.([]string)
				if tc.input == nil {
					okI = true
					inputSlice = nil
				}
				if !okI {
					t.Fatalf("Input type mismatch for string reducer")
				}
				initialValue, okInit := tc.initial.(string)
				if !okInit {
					t.Fatalf("Initial type mismatch for string reducer")
				}
				got, err = functional.ReduceErr[string, string](inputSlice, initialValue, reducer)
			default:
				if tc.input == nil { // Handle nil specifically
					got = tc.initial
					err = nil
				} else {
					t.Fatalf("Unhandled reducer type or mismatched input/initial types: Reducer=%T, Input=%T, Initial=%T", tc.reducer, tc.input, tc.initial)
				}
			}

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ReduceErr() error mismatch: got %v, want %v", err, tc.wantErr)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ReduceErr() value mismatch: got %#v (type %T), want %#v (type %T) (wantErr=%v, actualErr=%v)", got, got, tc.want, tc.want, tc.wantErr, err)
			}
		})
	}
}

// --- Benchmarks ---
var benchInputInts = make([]int, 1000)

func init() {
	for i := range benchInputInts {
		benchInputInts[i] = i
	}
}
func mapperNoErrorBench(n int) (string, error) { return strconv.Itoa(n), nil }
func mapperWithErrorBench(n int) (string, error) {
	if n != 0 && n%(errRateTest*10) == 0 {
		return "", errTestSentinel
	}
	return strconv.Itoa(n), nil
}
func predicateNoErrorBench(n int) (bool, error) { return n%2 == 0, nil }
func predicateWithErrorBench(n int) (bool, error) {
	if n != 0 && n%(errRateTest*10) == 0 {
		return false, errTestSentinel
	}
	return n%2 == 0, nil
}
func reducerNoErrorBench(acc, current int) (int, error) { return acc + current, nil }
func reducerWithErrorBench(acc, current int) (int, error) {
	if current != 0 && current%(errRateTest*10) == 0 {
		return acc, errTestSentinel
	}
	return acc + current, nil
}

func BenchmarkMapErr_NoError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	b.ResetTimer()
	var r []string
	var e error
	for i := 0; i < b.N; i++ {
		r, e = functional.MapErr(data, mapperNoErrorBench)
	}
	_, _ = r, e
}

func BenchmarkMapErr_WithError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	b.ResetTimer()
	var r []string
	var e error
	for i := 0; i < b.N; i++ {
		r, e = functional.MapErr(data, mapperWithErrorBench)
	}
	_, _ = r, e
}

func BenchmarkFilterErr_NoError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	b.ResetTimer()
	var r []int
	var e error
	for i := 0; i < b.N; i++ {
		r, e = functional.FilterErr(data, predicateNoErrorBench)
	}
	_, _ = r, e
}

func BenchmarkFilterErr_WithError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	b.ResetTimer()
	var r []int
	var e error
	for i := 0; i < b.N; i++ {
		r, e = functional.FilterErr(data, predicateWithErrorBench)
	}
	_, _ = r, e
}

func BenchmarkReduceErr_NoError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	initial := 0
	b.ResetTimer()
	var r int
	var e error
	for i := 0; i < b.N; i++ {
		r, e = functional.ReduceErr(data, initial, reducerNoErrorBench)
	}
	_, _ = r, e
}

func BenchmarkReduceErr_WithError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	initial := 0
	b.ResetTimer()
	var r int
	var e error
	for i := 0; i < b.N; i++ {
		r, e = functional.ReduceErr(data, initial, reducerWithErrorBench)
	}
	_, _ = r, e
}

func BenchmarkMapLoop_WithError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	b.ResetTimer()
	var res []string
	var err error
	for i := 0; i < b.N; i++ {
		r := make([]string, 0, len(data))
		var lE error
		for _, item := range data {
			v, mE := mapperWithErrorBench(item)
			if mE != nil {
				lE = mE
				break
			}
			r = append(r, v)
		}
		res = r
		err = lE
	}
	_, _ = res, err
}

func BenchmarkFilterLoop_WithError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	b.ResetTimer()
	var res []int
	var err error
	for i := 0; i < b.N; i++ {
		r := make([]int, 0)
		var lE error
		for _, item := range data {
			inc, pE := predicateWithErrorBench(item)
			if pE != nil {
				lE = pE
				break
			}
			if inc {
				r = append(r, item)
			}
		}
		res = r
		err = lE
	}
	_, _ = res, err
}

// CORRECTED BenchmarkReduceLoop_WithError_N1000
func BenchmarkReduceLoop_WithError_N1000(b *testing.B) {
	data := benchInputInts[:1000]
	initial := 0
	b.ResetTimer()
	var res int
	var err error
	for i := 0; i < b.N; i++ {
		acc := initial
		var lE error // loop error
		for _, item := range data {
			nA, rE := reducerWithErrorBench(acc, item) // next Accumulator, reduce Error
			if rE != nil {
				lE = rE                     // Store the error
				res = acc                   // Store the accumulator state *before* the error
				goto endLoopWithErrorReduce // Exit inner loop
			}
			acc = nA // Update accumulator only on success
		}
		// If loop finished OR goto was taken
		res = acc // Store final/last accumulator
		err = lE  // Assign whatever lE holds (nil if no error, sentinel if error)
	endLoopWithErrorReduce: // Label for goto
	}
	// Use the results to prevent compiler optimization removing the loops
	_, _ = res, err
}
