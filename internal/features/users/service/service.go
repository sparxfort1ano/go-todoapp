// Package service acts as the service layer for the Users feature.
// It is responsible for validating the user payload.
package service

import (
	"context"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
)

// UsersService encapsulates the core business logic for user management.
// All of its methods delegate the persistence logic to the repository layer and
// apply a configured operation timeout.
type UsersService struct {
	usersRepository UsersRepository
}

// UsersRepository defines the contract that decouples the service layer
// from the underlying repository logic.
type UsersRepository interface {
	// CreateUser executes the SQL query to insert a new user into the database.
	// It maps the resulting database row back into a domain entity.
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	// GetUsers executes the SQL query to read the given rows
	// according to the limit and offset filter.
	// It maps the resulting database row back into a domain entity.
	GetUsers(
		ctx context.Context,
		page domain.Pagination,
	) ([]domain.User, error)

	// GetUser executes the SQL query to read the given row
	// according to the user identificator.
	// It maps the resulting database row back into a domain entity.
	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	// DeleteUser executes the SQL query to delete the given row
	// according to the user identificator.
	DeleteUser(
		ctx context.Context,
		id int,
	) error

	// PatchUser executes the SQL query to patch the given row
	// according to the user identificator.
	// It uses Optimistic Concurrency Control by checking the user's Version
	// to prevent lost updates.
	// It maps the resulting database row back into a domain entity.
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
