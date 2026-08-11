package cached

import (
	"context"
	"errors"

	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	cache "github.com/sparxfort1ano/go-todoapp/internal/core/repository/redis"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
	"go.uber.org/zap"
)

func (r *CachedRepository) getTaskFromCache(
	ctx context.Context,
	id int,
) (repository.Task, bool) {
	log := logger.FromContext(ctx)

	key := taskKey(id)

	bytes, err := r.pool.Get(ctx, key).Bytes()
	if err != nil {
		if !errors.Is(err, cache.ErrNotFound) {
			log.Error("read a task from cache", zap.Error(err))
		}

		return repository.Task{}, false
	}

	var taskModel TaskModel
	if err := taskModel.Deserialize(bytes); err != nil {
		log.Error("deserialize cached task", zap.Error(err))
		return repository.Task{}, false
	}

	taskRepo := modelToRepo(taskModel)

	return taskRepo, true
}

func (r *CachedRepository) getTasksFromCache(
	ctx context.Context,
	page repository.Pagination,
	userID *int,
) ([]repository.Task, bool) {
	log := logger.FromContext(ctx)

	key, field := tasksListKey(userID), tasksListField(page)

	bytes, err := r.pool.HGet(ctx, key, field).Bytes()
	if err != nil {
		if !errors.Is(err, cache.ErrNotFound) {
			log.Error("hget tasks from cache", zap.Error(err))
		}

		return []repository.Task{}, false
	}

	var taskListModel TaskListModel
	if err := taskListModel.Deserialize(bytes); err != nil {
		log.Error("deserialize cached task list", zap.Error(err))
		return []repository.Task{}, false
	}

	taskListRepo := modelsToRepo(taskListModel)

	return taskListRepo, true
}

func (r *CachedRepository) cacheTask(
	ctx context.Context,
	task repository.Task,
) {
	log := logger.FromContext(ctx)

	taskModel := repoToModel(task)
	bytes, err := taskModel.Serialize()
	if err != nil {
		log.Error("serialize task", zap.Error(err))
		return
	}

	if err = r.pool.Set(
		ctx,
		taskKey(task.ID),
		bytes,
		r.pool.TTL(),
	).Err(); err != nil {
		log.Error("set task in cache", zap.Error(err))
	}
}

func (r *CachedRepository) cacheTasks(
	ctx context.Context,
	page repository.Pagination,
	userID *int,
	tasks []repository.Task,
) {
	log := logger.FromContext(ctx)

	taskListModel := repoToModels(tasks)
	bytes, err := taskListModel.Serialize()
	if err != nil {
		log.Error("serialize task list", zap.Error(err))
		return
	}

	key, field := tasksListKey(userID), tasksListField(page)
	if err := r.pool.HSet(ctx, key, field, bytes).Err(); err != nil {
		log.Error("hset tasks in cache", zap.Error(err))
		return
	}
}

// invalidateTasks removes keys associated with task lists from cache:
//   - "tasks:all"        — a list of all tasks
//   - "tasks:<userID>"   — a list of tasks for a specific user
//   - "task:<id>"        — a key for a single task, only if taskID != nil
func (r *CachedRepository) invalidateTasks(
	ctx context.Context,
	userID int,
	taskID *int,
) {
	log := logger.FromContext(ctx)

	invalidateKeys := []string{
		tasksListKey(nil),
		tasksListKey(&userID),
	}
	if taskID != nil {
		invalidateKeys = append(invalidateKeys, taskKey(*taskID))
	}

	if err := r.pool.Del(ctx, invalidateKeys...).Err(); err != nil {
		log.Error("invalidate cached tasks lists", zap.Error(err))
	}
}
