package domain

// Nullable specifies a field that is
//
// 1. Omitted in JSON (struct remains zero-value, Set=false).
//
// 2. Explicitly set to null (Set=true, Value=nil).
//
// 3. Provided with a value (Set=true, Value!=nil).
type Nullable[T any] struct {
	Value *T
	Set   bool
}
