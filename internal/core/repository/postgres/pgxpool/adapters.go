package pgxpool

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/postgres"
)

// pgxRows wraps the standard pgx.Rows to implement the postgres.Rows interface.
type pgxRows struct {
	pgx.Rows
}

// pgxRow wraps the standard pgx.Row to implement the postgres.Row interface.
type pgxRow struct {
	pgx.Row
}

// Scan delegates to the underlying pgx.Row's Scan method,
// translating the driver-specific pgx errors into the
// domain-agnostic errors (see errors.go for more details).
func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}

	return nil
}

// pgxCommandTag wraps the standard pgx.CommandTag to implement the postgres.CommandTag interface.
type pgxCommandTag struct {
	pgconn.CommandTag
}

func mapErrors(err error) error {
	const (
		pgxViolatesForeignKeyError = "23503"
	)

	var (
		pgErr *pgconn.PgError
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return postgres.ErrNoRows
	case errors.As(err, &pgErr):
		if pgErr.Code == pgxViolatesForeignKeyError {
			return fmt.Errorf(
				"%v: %w",
				err,
				postgres.ErrViolatesForeignKey,
			)
		}
	}

	return fmt.Errorf(
		"%v: %w",
		err,
		postgres.ErrUnknown,
	)

}
