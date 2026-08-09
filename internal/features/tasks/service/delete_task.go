package service

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/in"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

func (s *TasksService) DeleteTask(
	ctx context.Context,
	params in.DeleteTaskParams,
) (in.DeleteTaskResult, error) {
	repoParams := repository.NewDeleteTaskParams(params.ID)
	_, err := s.tasksRepository.DeleteTask(ctx, repoParams)
	if err != nil {
		return in.DeleteTaskResult{}, fmt.Errorf(
			"delete task from repository: %w",
			err,
		)
	}

	return in.NewDeleteTaskResult(), nil
}
