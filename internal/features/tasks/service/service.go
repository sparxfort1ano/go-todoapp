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
// TODO methods commentaries
type TasksRepository interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		userID *int,
		page domain.Pagination,
	) ([]domain.Task, error)

	GetTask(
		ctx context.Context,
		id int,
	) (domain.Task, error)

	DeleteTask(
		ctx context.Context,
		id int,
	) error

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
