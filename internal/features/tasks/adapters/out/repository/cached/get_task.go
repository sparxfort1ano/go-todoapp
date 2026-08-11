package cached

import (
	"context"

	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

// GetTask implements the cache-aside pattern for a single task.
func (r *CachedRepository) GetTask(
	ctx context.Context,
	params repository.GetTaskParams,
) (repository.GetTaskResult, error) {
	task, hit := r.getTaskFromCache(ctx, params.ID)
	if hit {
		return repository.NewGetTaskResult(task), nil
	}

	repoGetTaskResult, err := r.mainRepository.GetTask(ctx, params)
	if err != nil {
		return repository.GetTaskResult{}, err
	}

	r.cacheTask(ctx, repoGetTaskResult.Task)

	return repoGetTaskResult, nil
}
