package service

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/in"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

func (s *TasksService) PatchTask(
	ctx context.Context,
	params in.PatchTaskParams,
) (in.PatchTaskResult, error) {
	repoGetTaskResult, err := s.tasksRepository.GetTask(ctx, repository.NewGetTaskParams(params.ID))
	if err != nil {
		return in.PatchTaskResult{}, fmt.Errorf("get task from repository: %w", err)
	}

	task, patch := repoToDomain(repoGetTaskResult.Task), domain.NewTaskPatch(
		domain.Nullable[string](params.Patch.Title),
		domain.Nullable[string](params.Patch.Description),
		domain.Nullable[bool](params.Patch.Completed),
	)
	if err := task.ApplyPatch(patch); err != nil {
		return in.PatchTaskResult{}, fmt.Errorf("apply task patch: %w", err)
	}

	repoUpdateTaskParams := repository.NewUpdateTaskParams(domainToRepo(task))
	repoUpdateTaskResult, err := s.tasksRepository.UpdateTask(ctx, repoUpdateTaskParams)
	if err != nil {
		return in.PatchTaskResult{}, fmt.Errorf("update task from repository: %w", err)
	}

	return in.NewPatchTaskResult(repoTaskToIn(repoUpdateTaskResult.Task)), nil
}
