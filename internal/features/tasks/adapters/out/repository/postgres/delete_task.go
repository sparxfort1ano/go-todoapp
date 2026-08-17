package postgres

import (
	"context"
	"fmt"

	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

// DeleteTask executes the SQL query to delete the given row
// according to the task identificator.
func (r *TasksRepository) DeleteTask(
	ctx context.Context,
	params repository.DeleteTaskParams,
) (repository.DeleteTaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.tasks
	WHERE id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, params.ID)
	if err != nil {
		return repository.DeleteTaskResult{}, fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return repository.DeleteTaskResult{}, fmt.Errorf(
			"task with id=`%d`: %w",
			params.ID,
			errs.ErrNotFound,
		)
	}

	return repository.NewDeleteTaskResult(), nil
}
