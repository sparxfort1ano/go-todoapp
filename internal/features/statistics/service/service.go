// Package service acts as the service layer for the Statistics feature.
// It is responsible for validating the task payload.
package service

import (
	"context"
	"time"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
)

// StatisticsService encapsulates the core business logic for task management.
type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	// GetTasks executes the SQL query to read the given rows
	// according to the from, to (both define the time range for creating tasks)
	// and user identificator filter. It applies configured operation
	// timeout and maps the resulting database row back into a domain entity.
	GetTasks(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

// NewStatisticsService creates a new instance of StatisticsService.
func NewStatisticsService(
	statisticsRepository StatisticsRepository,
) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
