package service

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/in"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	params in.GetTasksParams,
) (in.GetTasksResult, error) {
	page := domain.NewPagination(params.Pagination.Limit, params.Pagination.Offset)
	if err := page.Validate(); err != nil {
		return in.GetTasksResult{}, err
	}

	repoPage := repository.NewPagination(page.Limit, page.Offset)
	repoParams := repository.NewGetTasksParams(repoPage, params.UserID)

	repoResult, err := s.tasksRepository.GetTasks(ctx, repoParams)
	if err != nil {
		return in.GetTasksResult{}, fmt.Errorf("get tasks from repository: %w", err)
	}

	return in.NewGetTasksResult(outTasksToIn(repoResult.Tasks)), nil
}
