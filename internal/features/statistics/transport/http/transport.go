// Package http acts as the transport layer for the Statistics feature.
// It is responsible for parsing HTTP requests, formatting responses and routing.
package http

import (
	"context"
	"net/http"
	"time"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/server"
)

// StatisticsHTTPHandler handles HTTP requests related to statistics management.
type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
}

// StatisticsService defines the contract that decouples the HTTP transport layer
// from the underlying domain logic.
type StatisticsService interface {
	// GetStatistics validates the from and to parameters, calculates the statistics
	// according to the from, to and user identificator parameters
	// forming domain.Statistics as the result.
	GetStatistics(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) (domain.Statistics, error)
}

// NewStatisticsHTTPHandler creates a new instance of StatisticsHTTPHandler.
func NewStatisticsHTTPHandler(
	statisticsService StatisticsService,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: statisticsService,
	}
}

// Routes returns a list of HTTP routes to be registered in the server router.
func (h *StatisticsHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
