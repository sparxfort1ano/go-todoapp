// Package postgres acts as the repository layer for the Tasks feature.
// It interacts with the PostgreSQL database to perform CRUD operations.
package postgres

import (
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/postgres"
)

// TasksRepository provides data access methods for task /entities.
// All of its Task methods apply a configured operation timeout.
type TasksRepository struct {
	pool postgres.Pool
}

// NewTasksRepository creates a new instance of TasksRepository.
func NewTasksRepository(
	pool postgres.Pool,
) *TasksRepository {
	return &TasksRepository{
		pool: pool,
	}
}
