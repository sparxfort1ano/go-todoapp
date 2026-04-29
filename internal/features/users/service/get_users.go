package service

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
)

// GetUsers enforces business rules (like negative values in limit or offset parameter) 
// on the user domain.
func (s *UsersService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit must be non-negative: %w",
			errs.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be non-negative: %w",
			errs.ErrInvalidArgument,
		)
	}

	users, err := s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	return users, nil
}
