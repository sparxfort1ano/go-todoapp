package postgres

import (
	"context"
	"errors"
	"fmt"

	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/postgres"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

// SaveTask executes the SQL query to insert a new task into the database.
// It maps the resulting database row back into a repository entity.
func (r *TasksRepository) SaveTask(
	ctx context.Context,
	params repository.SaveTaskParams,
) (repository.SaveTaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.tasks (title, description, completed, created_at, completed_at, author_user_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id;
	`

	task := params.Task
	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)

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
		if errors.Is(err, postgres.ErrViolatesForeignKey) {
			return repository.SaveTaskResult{}, fmt.Errorf(
				"%v: user with id='%d': %w",
				err,
				task.AuthorUserID,
				errs.ErrNotFound,
			)
		}
		return repository.SaveTaskResult{}, fmt.Errorf("scan error: %w", err)
	}

	return repository.NewSaveTaskResult(modelToRepo(taskModel)), nil
}
