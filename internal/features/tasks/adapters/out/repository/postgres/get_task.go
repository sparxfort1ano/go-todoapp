package postgres

import (
	"context"
	"errors"
	"fmt"

	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/postgres"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

func (r *TasksRepository) GetTask(
	ctx context.Context,
	params repository.GetTaskParams,
) (repository.GetTaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoapp.tasks
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, params.ID)

	var taskModel TaskModel
	if err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	); err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return repository.GetTaskResult{}, fmt.Errorf(
				"task with id='%d': %w",
				params.ID,
				errs.ErrNotFound,
			)
		}
		return repository.GetTaskResult{}, fmt.Errorf("scan error: %w", err)
	}

	return repository.NewGetTaskResult(modelToRepo(taskModel)), nil
}
