// Package http acts as the transport layer for the Tasks feature.
// It is responsible for parsing HTTP requests, formatting responses and routing.
package http

import (
	"context"
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/server"
)

// TasksHTTPHandler handles HTTP requests related to tasks management.
// All of its methods delegate the logic to the service layer.
type TasksHTTPHandler struct {
	tasksService TasksService
}

// TasksService defines the contract that decouples the HTTP transport layer
// from the underlying domain logic.
type TasksService interface {
	// CreateTask enforces business rules (like length
	// and semantically bad values) on the task domain.
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	// GetTasks enforces business rules (like negative values in a limit
	// or offset parameter) on the task domain.
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

	// PatchTask at first requests to get the given task data by the task identificator,
	// then enforces business rules on both the task patch and task domain levels
	// (see ApplyPatch for details).
	PatchTask(
		ctx context.Context,
		id int,
		patch domain.TaskPatch,
	) (domain.Task, error)
}

// NewTaskHTTPHandler creates a new instance of TasksHTTPHandler.
func NewTaskHTTPHandler(
	tasksService TasksService,
) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
	}
}

// Routes returns a list of HTTP routes to be registered in the server router.
func (h *TasksHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: h.GetTasks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/{id}",
			Handler: h.GetTask,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/tasks/{id}",
			Handler: h.DeleteTask,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/tasks/{id}",
			Handler: h.PatchTask,
		},
	}
}
