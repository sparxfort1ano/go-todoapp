package service

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/in"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	params in.GetTaskParams,
) (in.GetTaskResult, error) {
	repoParams := repository.NewGetTaskParams(params.ID)
	repoResult, err := s.tasksRepository.GetTask(ctx, repoParams)
	if err != nil {
		return in.GetTaskResult{}, fmt.Errorf("get task from repository: %w", err)
	}

	return in.NewGetTaskResult(repoTaskToIn(repoResult.Task)), nil
}
