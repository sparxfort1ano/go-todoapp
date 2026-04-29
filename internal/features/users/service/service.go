// Package service acts as the service layer for the Users feature.
// It is responsible for validating the user payload.
package service

import (
	"context"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
)

// UsersService encapsulates the core business logic for user management.
// All of its User methods delegate the persistence logic to the repository layer.
type UsersService struct {
	usersRepository UsersRepository
}

// UsersRepository defines the contract that decouples the service layer
// from the underlying repository logic.
type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		user domain.User,
	) (domain.User, error)
}

// NewUsersService creates a new instance of UsersService.
func NewUsersService(
	usersRepository UsersRepository,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
