package cached

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

// DeleteTask invalidates task from the cache by its user ID
// and deletes it from the main repository.
func (r *CachedRepository) DeleteTask(
	ctx context.Context,
	params repository.DeleteTaskParams,
) (repository.DeleteTaskResult, error) {
	task, hit := r.getTaskFromCache(ctx, params.ID)
	if !hit {
		repoGetTaskResult, err := r.mainRepository.GetTask(
			ctx,
			repository.NewGetTaskParams(params.ID),
		)
		if err != nil {
			return repository.DeleteTaskResult{}, fmt.Errorf("get task info from main repo: %w", err)
		}

		task = repoGetTaskResult.Task
	}

	r.invalidateTasks(ctx, task.AuthorUserID, &task.ID)

	return r.mainRepository.DeleteTask(ctx, params)
}
