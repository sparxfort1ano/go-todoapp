package postgres

import (
	"context"
	"fmt"

	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
)

// DeleteUser executes the SQL query to delete the given row
// according to the user identificator.
func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE from todoapp.users
	WHERE id=$1
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id=`%d`: %w", id, errs.ErrNotFound)
	}

	return nil
}
