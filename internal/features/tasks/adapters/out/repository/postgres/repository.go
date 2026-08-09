// Package postgres consists of a set of adapters that adapt the Postgres repository to the outcoming port (ports/out/repository),
// interacting with models responsible for the actual network interaction with the database.
package postgres

import (
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/postgres"
)

// TasksRepository provides data access methods for task entities.
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
