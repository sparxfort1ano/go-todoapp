package domain

import (
	"errors"
	"fmt"

	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
)

// Pagination provides limit and offset parameters for listing operations.
type Pagination struct {
	Limit  *int
	Offset *int
}

// NewPagination creates a new instance of Pagination.
func NewPagination(
	limit *int,
	offset *int,
) Pagination {
	return Pagination{
		Limit:  limit,
		Offset: offset,
	}
}

const (
	Limit  = "limit"
	Offset = "offset"
)

// Validate checks whether the business rules for the limit and offset queries are met.
// It returns error if the data contradicts the rules
// such as the queries are negative.
func (p *Pagination) Validate() error {
	return errors.Join(
		validateNonNegative(Limit, p.Limit),
		validateNonNegative(Offset, p.Offset),
	)
}

func validateNonNegative(name string, val *int) error {
	if val != nil && *val < 0 {
		return fmt.Errorf(
			"%s must be non-negative: %w",
			name,
			errs.ErrInvalidArgument,
		)
	}

	return nil
}
