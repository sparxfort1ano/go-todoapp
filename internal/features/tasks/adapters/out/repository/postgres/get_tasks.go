package postgres

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	params repository.GetTasksParams,
) (repository.GetTasksResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoapp.tasks
	%s
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2;
	`
	args := []any{params.Pagination.Limit, params.Pagination.Offset}

	if params.UserID != nil {
		query = fmt.Sprintf(query, "WHERE author_user_id=$3")
		args = append(args, params.UserID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return repository.GetTasksResult{}, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel

		err := rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
		)
		if err != nil {
			return repository.GetTasksResult{}, fmt.Errorf("scan tasks: %w", err)
		}
		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return repository.GetTasksResult{}, fmt.Errorf("next rows: %w", err)
	}

	return repository.NewGetTasksResult(modelsToRepo(taskModels)), nil
}
