package service

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/in"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	params in.CreateTaskParams,
) (in.CreateTaskResult, error) {
	task := domain.NewTaskUninitialized(
		params.Title,
		params.Description,
		params.AuthorUserID,
	)

	if err := task.Validate(); err != nil {
		return in.CreateTaskResult{}, fmt.Errorf("validate task domain: %w", err)
	}

	repoParams := repository.NewSaveTaskParams(domainToRepo(task))

	repoResult, err := s.tasksRepository.SaveTask(ctx, repoParams)
	if err != nil {
		return in.CreateTaskResult{}, fmt.Errorf("save task in repository: %w", err)
	}

	return in.NewCreateTaskResult(repoTaskToIn(repoResult.Task)), nil
}
