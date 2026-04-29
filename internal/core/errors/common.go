// Package errors defines application-wide sentinel errors.
package errors

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid agrument")
	ErrConflict        = errors.New("conflict")
)
