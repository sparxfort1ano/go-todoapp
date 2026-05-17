package http

import (
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/request"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/response"
)

// DeleteTask processes the HTTP request to delete all the data
// about the task with the given identificator.
//
// @Summary		Удаление задачи
// @Description Удалить существующую в системе задачу по ее ID.
// @Tags 		tasks
// @Param 		id 			path int true 					"ID удаляемой задачи"
// @Success 	204 										"Успешное удаление задачи"
// @Failure 	400 		{object} response.ErrorResponse "Bad request"
// @Failure 	404 		{object} response.ErrorResponse "Task not found"
// @Failure 	500 		{object} response.ErrorResponse "Internal server error"
// @Router 		/tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
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

	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)
		return
	}

	responseHandler.NoContentResponse()
}
