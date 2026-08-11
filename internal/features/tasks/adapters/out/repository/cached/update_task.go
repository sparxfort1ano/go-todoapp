package cached

import (
	"context"

	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

// UpdateTask updates the task in the main repository, then it updates the cache:
// 1) The updated task gets cached by its ID (write-through approach for a single task),
// 2) All of the author’s task lists get invalidated (as they are out of date).
func (r *CachedRepository) UpdateTask(
	ctx context.Context,
	params repository.UpdateTaskParams,
) (repository.UpdateTaskResult, error) {
	result, err := r.mainRepository.UpdateTask(ctx, params)
	if err != nil {
		return repository.UpdateTaskResult{}, err
	}

	r.cacheTask(ctx, result.Task)
	r.invalidateTasks(ctx, result.Task.AuthorUserID, nil)

	return result, nil
}
