// Package postgres defines the core contract for PostgreSQL database operations
// used by the repository layer, allowing the underlying implentation (e.g. pgx) to be easily swapped or mocked.
package postgres

import (
	"context"
	"time"
)

// Pool defines the driver-agnostic contract for executing SQL queries against a PostgreSQL database.
// It abstracts the underlying connection management and query execution.
type Pool interface {
	// Query executes a query that returns multiple rows.
	// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
	// The caller must ensure that the returned Rows are closed after processing.
	Query(ctx context.Context, sql string, args ...any) (Rows, error)

	// QueryRow executes a query that is expected to return at most one row.
	// Errors are deferred until Row's Scan method is called.
	// If the query selects no rows, calling Scan on the returned Row will yield ErrNoRows.
	// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
	QueryRow(ctx context.Context, sql string, args ...any) Row

	// Exec acquires a query without returning any rows.
	// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
	Exec(ctx context.Context, sql string, args ...any) (CommandTag, error)

	// Close gracefully closes all connections in the pool.
	// It blocks until all connections are safely terminated.
	Close()

	// OpTimeout returns the configured time limit for database operations.
	OpTimeout() time.Duration
}

// Rows is the result set returned from the Query method.
type Rows interface {
	// Close closes the rows, making the connection ready for use again.
	// It is safe to call Close multiple times.
	Close()

	// Err returns any error that occured during query execution or row iteration.
	// It should be checked after the iteration loop finishes.
	Err() error

	// Next prepares the next row for reading.
	// It returns true if a row is available or false if no more rows are
	// available or a fatal error has occurred.
	// It automatically closes rows upon returning false.
	Next() bool

	// Scan reads the values from the current row into the provided destination pointers.
	// It must only be called after successful call to Next.
	Scan(dest ...any) error
}

// Row is the result returned from the QueryRow method.
type Row interface {
	// Scan read the values from the single row into the provided destination pointers.
	// If no rows were found, it returns ErrNoRows.
	Scan(dest ...any) error
}

// CommandTag is the result from the Exec method.
// It provides information about the executed statement,
// such as the number of rows affected.
type CommandTag interface {
	// RowsAffected returns the number of rows inserted, updated or deleted.
	RowsAffected() int64
}
