package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/request"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/response"
)

// GetStatisticsResponse represents the outgoing JSON body after getting the statistics (JSON).
type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletionTime *string  `json:"task_average_completion_time"`
}

// GetStatistics processes the HTTP request to get task statistics
// according to the from and to parameters and user identificator.
// It delegates the logic to the service layer and writes the JSON response.
func (h *StatisticsHTTPHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'userID'/'from'/'to'",
		)
		return
	}

	statisticsDomain, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)
		return
	}

	response := toDTOfromDomain(statisticsDomain)
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		queryUserID = "user_id"
		queryFrom   = "from"
		queryTo     = "to"
	)

	userID, errUserID := request.GetIntQueryParam(r, queryUserID)
	if errUserID != nil {
		errUserID = fmt.Errorf(
			"get '%s' query params: %w",
			queryUserID,
			errUserID,
		)
	}

	from, errFrom := request.GetDateQueryParam(r, queryFrom)
	if errFrom != nil {
		errFrom = fmt.Errorf(
			"get '%s' query params: %w",
			queryFrom,
			errFrom,
		)
	}

	to, errTo := request.GetDateQueryParam(r, queryTo)
	if errTo != nil {
		errTo = fmt.Errorf(
			"get '%s' query params: %w",
			queryTo,
			errTo,
		)
	}

	if errs := errors.Join(
		errUserID,
		errFrom,
		errTo,
	); errs != nil {
		return nil, nil, nil, errs
	}

	return userID, from, to, nil
}

func toDTOfromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}
	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}
