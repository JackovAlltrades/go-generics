# Go Generics Reference Module

[![Go Reference](https://pkg.go.dev/badge/github.com/JackovAlltrades/go-generics/functional.svg)](https://pkg.go.dev/github.com/JackovAlltrades/go-generics/functional)
[![Go Report Card](https://goreportcard.com/badge/github.com/JackovAlltrades/go-generics)](https://goreportcard.com/report/github.com/JackovAlltrades/go-generics)
<!-- Add build status badge once CI is set up -->

A collection of common generic functional utilities, data structures, and concurrency patterns implemented in Go using generics (Go 1.18+).

This project aims to provide well-documented, performant, and safe examples for learning and using Go generics effectively.

## Installation

```bash
go get github.com/JackovAlltrades/go-generics/functional
```

## Project Structure
*(Note: /ds and /concurrency are planned future additions)*
- `/functional`: Generic functions (Map, Filter, Reduce, etc.)
- `/ds`: Generic data structures (Stack, Queue, Set, etc.) - *Future*
- `/concurrency`: Concurrency patterns using generics (ParallelMap, WorkerPool, etc.) - *Future*
- `/examples`: Usage examples found within `_test.go` files.

## Features (v0.1.0)

The `functional` package currently includes:

### Slice Utilities
- **Core:** `Map`, `Filter`, `Reduce`, `Find`
- **Checks:** `Any`, `All`, `Contains`
- **Manipulation:** `Unique`, `Reverse`, `ReverseInPlace`, `Chunk`
- **Transformation:** `Flatten`
- **Accessors:** `First`, `Last`
- **Set Operations:** `Intersection`, `Union`, `Difference` (Require `comparable` elements)
- **Grouping:** `GroupBy`

### Map Utilities
- **Extraction:** `Keys` (sorted `cmp.Ordered` keys), `Values` (order arbitrary)
- **Transformation:** `MapToSlice`

## Usage Examples

```go
package main

import (
	"fmt"
	"github.com/JackovAlltrades/go-generics/functional"
    // "strings" // Only needed if using string example from README below
)

func main() {
	// --- Map ---
	nums := []int{1, 2, 3, 4}
	doubled := functional.Map(nums, func(n int) int { return n * 2 })
	fmt.Printf("Doubled: %v\n", doubled) // Output: Doubled: [2 4 6 8]

	// --- Filter ---
	evens := functional.Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Printf("Evens: %v\n", evens) // Output: Evens: [2 4]

	// --- Reduce ---
	sum := functional.Reduce(nums, 0, func(acc int, n int) int { return acc + n })
	fmt.Printf("Sum: %d\n", sum) // Output: Sum: 10

	// --- Unique ---
	duplicates := []string{"a", "b", "a", "c", "b", "a"}
	unique := functional.Unique(duplicates)
	fmt.Printf("Unique: %v\n", unique) // Output: Unique: [a b c]

	// --- GroupBy ---
	words := []string{"apple", "ant", "banana", "bat"}
	grouped := functional.GroupBy(words, func(s string) string { return string(s[0]) })
	fmt.Printf("Grouped: %v\n", grouped) // Output: Grouped: map[a:[apple ant] b:[banana bat]] (map order varies)

	// --- Intersection ---
	s1 := []int{1, 2, 3, 4}
	s2 := []int{3, 4, 5, 6}
	intersect := functional.Intersection(s1, s2)
	fmt.Printf("Intersection: %v\n", intersect) // Output: Intersection: [3 4]
}
```

For detailed usage of all functions, please see the [Go Package Documentation](https://pkg.go.dev/github.com/JackovAlltrades/go-generics/functional).

## Design Principles
- **Generics:** Functions are implemented using type parameters for maximum reusability and type safety.
- **Immutability:** Most functions operate on input data without modification, returning new slices or maps. Functions performing in-place modifications (e.g., `ReverseInPlace`) are explicitly named.
- **Nil/Empty Handling:** Functions generally return sensible zero values (often empty, non-nil slices/maps) for nil or empty inputs. Behavior is documented for each function.
- **Constraints:** Appropriate type constraints (`comparable`, `cmp.Ordered`) are used where necessary.

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
This project is licensed under the MIT License - see the `LICENSE` file for details.
