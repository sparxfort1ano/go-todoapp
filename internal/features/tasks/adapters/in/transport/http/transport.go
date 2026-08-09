// Package http consists of a set of adapters that adapt the HTTP transport to the incoming port (ports/in),
// interacting with DTOs responsible for receiving, decoding and validating HTTP requests and sending HTTP responses.
package http

import (
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/server"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/in"
)

// TasksHTTPHandler handles HTTP requests related to tasks management.
// All of its methods delegate the logic to the service layer.
type TasksHTTPHandler struct {
	tasksService in.TasksService
}

// NewTaskHTTPHandler creates a new instance of TasksHTTPHandler.
func NewTaskHTTPHandler(
	tasksService in.TasksService,
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
