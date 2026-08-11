package cached

import (
	"context"

	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

// SaveTask saves the task to the main repository, then it updates the cache:
// 1) The task gets cached by its ID (write-through approach for a single task),
// 2) All of the author’s task lists get invalidated (as they are out of date).
func (r *CachedRepository) SaveTask(
	ctx context.Context,
	params repository.SaveTaskParams,
) (repository.SaveTaskResult, error) {
	result, err := r.mainRepository.SaveTask(ctx, params)
	if err != nil {
		return repository.SaveTaskResult{}, err
	}

	r.cacheTask(ctx, result.Task)
	r.invalidateTasks(ctx, result.Task.AuthorUserID, nil)

	return result, nil
}
