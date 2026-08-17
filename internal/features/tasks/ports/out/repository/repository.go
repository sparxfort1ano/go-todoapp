// Package repository consists of an outcoming port for the tasks feature and its parameters.
// It acts as an intermediary through which the infrastructure communicates with the core.
// The package also provides custom data structures for the repository layer,
// which are essentially copies of domain entities (see `internal/core/domain` for further details).
package repository

import (
	"context"
)

// TasksRepository defines the contract that decouples the service layer
// from the repository logic.
type TasksRepository interface {
	SaveTask(
		ctx context.Context,
		params SaveTaskParams,
	) (SaveTaskResult, error)

	GetTasks(
		ctx context.Context,
		params GetTasksParams,
	) (GetTasksResult, error)

	GetTask(
		ctx context.Context,
		params GetTaskParams,
	) (GetTaskResult, error)

	DeleteTask(
		ctx context.Context,
		params DeleteTaskParams,
	) (DeleteTaskResult, error)

	UpdateTask(
		ctx context.Context,
		params UpdateTaskParams,
	) (UpdateTaskResult, error)
}
