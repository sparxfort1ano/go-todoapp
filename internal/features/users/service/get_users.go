package service

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	page domain.Pagination,
) ([]domain.User, error) {
	if err := page.Validate(); err != nil {
		return nil, err
	}

	users, err := s.usersRepository.GetUsers(ctx, page)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	return users, nil
}
