package http

import (
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/request"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/response"
)

// CreateTaskRequest represents the incoming JSON body for creating a task (DTO).
type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100" example:"Домашнее задание"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000" example:"Сделать до четверга математику"`
	AuthorUserID int     `json:"author_user_id" validate:"required" example:"5"`
}

// CreateTaskResponse represents the outgoing JSON body after a task is created (JSON).
type CreateTaskResponse TaskDTOResponse

// CreateTask processes the HTTP to register a new task.
// It decodes the payload and writes the JSON response.
//
// @Summary		Создать задачу
// @Description Создать новую задачу в системе.
// @Tags 		tasks
// @Accept 		json
// @Produce 	json
// @Param 		request body 		CreateTaskRequest 	true 	"CreateUser тело запроса"
// @Success 	201 	{object} 	CreateTaskResponse 			"Успешно созданная задача"
// @Failure 	400 	{object} 	response.ErrorResponse 		"Bad request"
// @Failure 	404 	{object} 	response.ErrorResponse 		"Author not found"
// @Failure 	500 	{object} 	response.ErrorResponse 		"Internal server error"
// @Router 		/tasks 	[post]
func (h *TasksHTTPHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	var req CreateTaskRequest
	if err := request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	taskDomain := domainFromDTO(req)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)
		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateTaskRequest) domain.Task {
	return domain.NewTaskUninitialized(
		dto.Title,
		dto.Description,
		dto.AuthorUserID,
	)
}
