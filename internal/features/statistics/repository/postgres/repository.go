// Package postgres acts as the repository layer for the Statistics feature.
// It interacts with the PostgreSQL database to perform CRUD operations.
package postgres

import "github.com/sparxfort1ano/go-todoapp/internal/core/repository/postgres"

// StatisticsRepository provides data access methods for task entities.
type StatisticsRepository struct {
	pool postgres.Pool
}

// NewStatisticsRepository creates a new instance of StatisticsRepository.
func NewStatisticsRepository(
	pool postgres.Pool,
) *StatisticsRepository {
	return &StatisticsRepository{
		pool: pool,
	}
}
