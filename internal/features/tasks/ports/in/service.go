// Package in consists of an incoming port for the tasks feature and its parameters.
// It acts as an intermediary through which the outside world communicates with the core.
// The package also provides custom data structures (see `types.go` for further details),
// which are essentially copies of domain entities (see `internal/core/domain` for further details).
package in

import (
	"context"
)

// TasksService defines the contract that decouples incoming adapters
// from the underlying domain logic.
type TasksService interface {
	// CreateTask enforces business rules (like length
	// and semantically bad values) on the task domain.
	CreateTask(
		ctx context.Context,
		params CreateTaskParams,
	) (CreateTaskResult, error)

	// GetTasks enforces business rules (like negative values in a limit
	// or offset parameter) on the task domain.
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

	// PatchTask at first requests to get the given task data by the task identificator,
	// then enforces business rules on both the task patch and task domain levels
	// (see ApplyPatch for details).
	PatchTask(
		ctx context.Context,
		params PatchTaskParams,
	) (PatchTaskResult, error)
}
