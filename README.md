# Go Generics Functional Utilities

[![Go Reference](https://pkg.go.dev/badge/github.com/JackovAlltrades/go-generics/functional.svg)](https://pkg.go.dev/github.com/JackovAlltrades/go-generics/functional)
[![Go Report Card](https://goreportcard.com/badge/github.com/JackovAlltrades/go-generics)](https://goreportcard.com/report/github.com/JackovAlltrades/go-generics)
<!-- [![Build Status](...) Add build status badge once CI is set up -->

A collection of common generic functional utilities implemented in Go using generics (Go 1.18+).

This project aims to provide well-documented, performant, and safe generic utilities for common slice and map operations, serving as both a practical library and a learning resource for Go generics.

## Who is this library for?

This library is useful for:
- Developers who prefer a functional programming style for data manipulation.
- Teams needing common slice/map processing utilities without writing boilerplate loops.
- Anyone looking to reduce repetitive code for tasks like mapping, filtering, reducing, or set operations.
- Developers learning or exploring the practical application and performance of Go generics.

It solves problems like:
- Processing collections declaratively.
- Transforming data structures between slices and maps.
- Implementing set logic (intersection, union, etc.).
- Simplifying code involving complex or nested loops for data handling.

## Installation

```bash
go get github.com/JackovAlltrades/go-generics/functional

Copy code
Project Structure
/functional: Generic functions for slice/map manipulation (Map, Filter, Reduce, Set Operations, etc.). This is the primary package.
/examples: Usage examples can be found as ExampleXxx functions within the functional/*_test.go files.
(Note: /ds and /concurrency mentioned in early plans are currently out of scope for this package)

Features (Current)
The functional package currently includes:

Core Functions
Map, Filter, Reduce, Any, All, Find, Contains, GroupBy
Error Handling Variants
MapErr, FilterErr, ReduceErr
Set Operations (comparable elements)
Unique, Intersection, Union, Difference
Slice Utilities
Chunk, Flatten, Reverse (in-place), ReversedCopy, First, Last
Map Utilities
Keys, Values, MapToSlice
(See the godoc reference for detailed function signatures.)

Usage Examples
// See examples directly in the godoc reference or within the _test.go files.
// The main function example from previous README versions demonstrates basic usage.
(Self-correction: Including a growing list of examples directly in the README can become unwieldy. Pointing to godoc/test files is more maintainable.)

Design Principles & Performance
Generics: Maximizes reusability and type safety using type parameters (any, comparable).
Immutability: Most functions return new collections; in-place modifications (Reverse) are explicitly named.
Nil/Empty Handling: Generally returns sensible zero values (e.g., empty, non-nil slices/maps) for nil/empty inputs. See individual function docs for specifics.
Order Guarantees:
Slice functions typically preserve relative order unless documented otherwise (Unique preserves first appearance order).
Functions operating on maps (Keys, Values, MapToSlice, GroupBy, Set Ops) do not guarantee order due to Go's map iteration behavior. Sort results explicitly if order is required.
Performance:
Benchmarks against manual Go loops show minimal overhead for most functions.
Focus is on idiomatic Go and efficient data structure use (map lookups for O(N) average set operations, slice preallocation).
Past bottlenecks (e.g., unnecessary sorting) have been identified via benchmarking and removed. See function-level godoc for specific notes.
Concurrency Safety: Functions are safe for concurrent use by multiple goroutines provided the input collection(s) are not modified concurrently by other goroutines. The library does not perform internal parallelization.
When to Choose Alternatives
While this library offers useful utilities, consider these alternatives:

Go Standard Library (slices, maps packages - Go 1.21+):

Use First: For functions like Contains, Keys, Values, SortFunc, Compact, IndexFunc, prefer the standard library versions when available and sufficient. They require no external dependencies.
See: pkg.go.dev/slices, pkg.go.dev/maps
Comprehensive Third-Party Libraries (samber/lo):

Use When: You need a wider range of utilities (e.g., FilterMap, Shuffle, Ternary operator helpers), advanced features, or built-in concurrency/async helpers (lo.Async, lo.Attempt).
See: github.com/samber/lo
Simple for Loops:

Use When: The operation is trivial (e.g., summing small int slices), maximum performance transparency is critical, or adding a library dependency feels like overkill for a single, simple task. Avoid premature optimization – the clarity gain from functional utilities is often more valuable.
Concurrency:

This library is sequential. If you need parallel execution (e.g., ParallelMap), use Go's built-in primitives (goroutines, channels, sync package, sync/errgroup) or look at libraries specifically designed for this (like samber/lo's async functions or other worker pool implementations).
Contributing
Contributions are welcome! Please feel free to submit a Pull Request or open an Issue. Running make check locally before submitting is highly recommended.

License
This project is licensed under the MIT License - see the LICENSE file for details.


---