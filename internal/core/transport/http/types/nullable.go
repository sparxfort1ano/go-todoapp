// Package types provides custom data structures for the HTTP transport layer.
// It handles web-specific serialization and deserialization (like custom JSON parsing).
package types

import (
	"encoding/json"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
)

// Nullable wraps the domain.Nullable type to provide custom JSON unmarshaling.
// It is primarily used for HTTP PATCH requests to distinguish
// between three states (see [domain.Nullable] for details).
type Nullable[T any] struct {
	domain.Nullable[T]
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It is only invoked by the standard encoding/json package if the field is present
// in the incoming payload. Therefore, it safely assumes the field is "Set".
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	if string(b) == "null" {
		n.Value = nil
		return nil
	}

	var value T
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

// ToDomain converts the transport-specific Nullable wrapper
// into the domain.Nullable entity for use in the service layer.
func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}
