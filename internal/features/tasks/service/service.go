// Package service acts as the service layer for the Tasks feature.
// It is responsible for validating the task payload.
package service

import "github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"

// TasksService encapsulates the core business logic for task management.
// All of its methods delegate the persistence logic to the repository layer and
// apply a configured operation timeout.
type TasksService struct {
	tasksRepository repository.TasksRepository
}

// NewTaskService creates a new instance of TaskService.
func NewTaskService(
	tasksRepository repository.TasksRepository,
) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}
