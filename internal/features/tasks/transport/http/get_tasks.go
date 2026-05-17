package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/request"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/response"
)

// GetTasksResponse represents the outgoing JSON body after getting a tasks slice (JSON).
type GetTasksResponse []TaskDTOResponse

// GetTasks processes the HTTP request to get a list of tasks
// according to the limit, offset and user identificator parameters.
// It writes the JSON response.
//
// @Summary		Список задач
// @Description Просмотр списка задач с опциональной пагинацией и/или фильтрацией по ID автора задачи.
// @Tags 		tasks
// @Produce 	json
// @Param 		user_id query 	int false 				"Фильтрация задач по ID автора"
// @Param 		limit 	query 	int false 				"Размер страницы с задачами"
// @Param 		offset 	query 	int false 				"Смещение страницы с задачами"
// @Success 	200 	{object} GetTasksResponse 		"Успешное получение списка задач"
// @Failure 	400 	{object} response.ErrorResponse "Bad request"
// @Failure 	500 	{object} response.ErrorResponse "Internal server error"
// @Router 		/tasks 	[get]
func (h *TasksHTTPHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'userID'/'limit/'offset'")
		return
	}

	page := domain.NewPagination(limit, offset)
	taskDomain, err := h.tasksService.GetTasks(ctx, userID, page)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(taskDTOFromDomains(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		queryUserID = "user_id"
		queryLimit  = "limit"
		queryOffset = "offset"
	)

	userID, errUserID := request.GetIntQueryParam(r, queryUserID)
	if errUserID != nil {
		errUserID = fmt.Errorf(
			"get '%s' query params: %w",
			queryUserID,
			errUserID,
		)
	}

	limit, errLimit := request.GetIntQueryParam(r, queryLimit)
	if errLimit != nil {
		errLimit = fmt.Errorf(
			"get '%s' query params: %w",
			queryLimit,
			errLimit,
		)
	}

	offset, errOffset := request.GetIntQueryParam(r, queryOffset)
	if errOffset != nil {
		errOffset = fmt.Errorf(
			"get '%s' query params: %w",
			queryOffset,
			errOffset,
		)
	}

	if errs := errors.Join(
		errUserID,
		errLimit,
		errOffset,
	); errs != nil {
		return nil, nil, nil, errs
	}

	return userID, limit, offset, nil
}
