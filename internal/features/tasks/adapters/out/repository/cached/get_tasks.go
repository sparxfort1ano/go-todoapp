package cached

import (
	"context"

	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

// GetTasks implements the cache-aside pattern for a multiple tasks.
func (r *CachedRepository) GetTasks(
	ctx context.Context,
	params repository.GetTasksParams,
) (repository.GetTasksResult, error) {
	tasks, hit := r.getTasksFromCache(ctx, params.Pagination, params.UserID)
	if hit {
		return repository.NewGetTasksResult(tasks), nil
	}

	repoGetTasksResult, err := r.mainRepository.GetTasks(ctx, params)
	if err != nil {
		return repository.GetTasksResult{}, err
	}

	r.cacheTasks(ctx, params.Pagination, params.UserID, repoGetTasksResult.Tasks)

	return repository.NewGetTasksResult(repoGetTasksResult.Tasks), nil
}
