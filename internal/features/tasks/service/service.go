// Package service acts as the service layer for the Tasks feature.
// It is responsible for validating the task payload.
package service

import (
	"context"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
)

// TasksService encapsulates the core business logic for task management.
// All of its methods delegate the persistence logic to the repository layer and
// apply a configured operation timeout.
type TasksService struct {
	tasksRepository TasksRepository
}

// TasksRepository defines the contract that decouples the service layer
// from the underlying repository logic.
type TasksRepository interface {
	// CreateTask executes the SQL query to insert a new task into the database.
	// It maps the resulting database row back into a domain entity.
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	// GetTasks executes the SQL query to read the given rows
	// according to the limit, offset and user identificator filter.
	// It maps the resulting database row back into a domain entity.
	GetTasks(
		ctx context.Context,
		userID *int,
		page domain.Pagination,
	) ([]domain.Task, error)

	// GetTask executes the SQL query to read the given row
	// according to the task identificator.
	// It maps the resulting database row back into a domain entity.
	GetTask(
		ctx context.Context,
		id int,
	) (domain.Task, error)

	// DeleteTask executes the SQL query to delete the given row
	// according to the task identificator.
	DeleteTask(
		ctx context.Context,
		id int,
	) error

	// PatchTask executes the SQL query to patch the given row
	// according to the task identificator.
	// It uses Optimistic Concurrency Control by checking the task's Version
	// to prevent lost updates.
	// It maps the resulting database row back into a domain entity.
	PatchTask(
		ctx context.Context,
		id int,
		task domain.Task,
	) (domain.Task, error)
}

// NewTaskService creates a new instance of TaskService.
func NewTaskService(
	tasksRepository TasksRepository,
) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}
