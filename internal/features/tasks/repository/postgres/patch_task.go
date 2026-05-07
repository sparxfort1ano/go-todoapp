package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/postgres"
)

func (r *TasksRepository) PatchTask(
	ctx context.Context,
	id int,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.tasks
	SET 
		title=$1,
		description=$2,
		completed=$3,
		completed_at=$4,
		version=version+1
	
	WHERE id=$5 AND version=$6

	RETURNING
		id, 
		version,
		title,
		description,
		completed,
		created_at,
		completed_at,
		author_user_id;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		task.ID,
		task.Version,
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
		if errors.Is(err, postgres.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id='%d' concurrently accessed: %w",
				id,
				errs.ErrConflict,
			)
		}
		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := taskDomainFromModel(taskModel)
	return taskDomain, nil
}
