// Package postgres manages the configuration and connection pool lifecycle
// for the PostgreSQL database.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool defines the contract for executing SQL queries against the database.
// It abstracts the underlying pgx driver.
type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()

	OpTimeout() time.Duration

}

// ConnectionPool wraps the standard pgxpool.Pool to implement the Pool interface.
// It also stores the global timeout for database operations.
type ConnectionPool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

// NewConnectionPool establishes a connection pool to the PostgreSQL database
// using the provided configuration. It verifies the connection with a Ping.
func NewConnectionPool(
	ctx context.Context,
	cfg config,
) (*ConnectionPool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	pgxConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgxpool ping: %w", err)
	}

	return &ConnectionPool{
		Pool:      pool,
		opTimeout: cfg.Timeout,
	}, nil
}

// OpTimeout returns the configured time limit for database operations.
func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.opTimeout
}
