package http

import (
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/request"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/response"
)

// GetTaskResponse represents the outgoing JSON body after a task is gotten (JSON).
type GetTaskResponse TaskDTOResponse

// GetTask processes the HTTP request to get a task with the given id.
// It writes the JSON response.
//
// @Summary		Получение задачи
// @Description Получение конкретной задачи по ее ID.
// @Tags 		tasks
// @Produce 	json
// @Param 		id 	path int true 					"ID удаляемого задачи"
// @Success 	200 {object} GetTaskResponse 		"Успешное получение задачи"
// @Failure 	400 {object} response.ErrorResponse "Bad request"
// @Failure 	404 {object} response.ErrorResponse "Task not found"
// @Failure 	500 {object} response.ErrorResponse "Internal server error"
// @Router 		/tasks/{id}  [get]
func (h *TasksHTTPHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	taskID, err := request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)
		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task",
		)
		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}
