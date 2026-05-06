package service

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	userID *int,
	page domain.Pagination,
) ([]domain.Task, error) {
	if err := page.Validate(); err != nil {
		return nil, err
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, userID, page)
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}

	return tasks, nil
}
