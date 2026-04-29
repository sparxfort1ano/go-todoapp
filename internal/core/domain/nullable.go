package domain

// Nullable specifies a field that is
// not provided,
// provided with the null value
// or provided with any other value.
type Nullable[T any] struct {
	Value *T
	Set   bool
}
