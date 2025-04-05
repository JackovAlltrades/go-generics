package functional_test

// --- Shared Test Helper Functions & Types ---

// Helper function to create a pointer to a value.
// Defined once here for use by all tests in this package.
func ptr[T any](v T) *T {
	return &v
}

// person struct used in multiple tests.
// Defined once here for use by all tests in this package.
type person struct {
	Name string
	Age  int
}

// You could also put shared error variables here if used across multiple files
// import "errors"
// var errSharedTest = errors.New("shared test error")
