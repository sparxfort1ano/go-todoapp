// Package postgres acts as the repository layer for the Users feature.
// It interacts with the PostgreSQL database to perform CRUD operations.
package postgres

import "github.com/sparxfort1ano/go-todoapp/internal/core/repository/postgres"

// UsersRepository provides data access methods for user entities.
// All of its User methods applies a configured operation timeout.
type UsersRepository struct {
	pool postgres.Pool
}

// NewUsersRepository creates a new instance of UsersRepository.
func NewUsersRepository(pool postgres.Pool) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}
