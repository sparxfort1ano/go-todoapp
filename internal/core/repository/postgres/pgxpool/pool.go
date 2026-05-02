// Package pgxpool provides a concrete implentation of the postgres.Pool interface
// using the github.com/jackc/pgx/v5 driver.
// It acts as an adapter, translating abstract database operations into pgx-specific calls,
// manages configuration for a PostgreSQL database.
package pgxpool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/postgres"
)

// Pool wraps the standard pgxpool.Pool to implement the postgres.Pool interface.
// It also stores the global timeout for database operations.
type Pool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

// NewPool establishes a connection pool to the PostgreSQL database
// using the provided configuration. It verifies the connection with a Ping.
func NewPool(
	ctx context.Context,
	cfg config,
) (*Pool, error) {
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

	return &Pool{
		Pool:      pool,
		opTimeout: cfg.Timeout,
	}, nil
}

func (p *Pool) Query(ctx context.Context, sql string, args ...any) (postgres.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxRows{rows}, nil
}

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) postgres.Row {
	row := p.Pool.QueryRow(ctx, sql, args...)

	return pgxRow{row}
}

func (p *Pool) Exec(ctx context.Context, sql string, args ...any) (postgres.CommandTag, error) {
	cmdTag, err := p.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxCommandTag{cmdTag}, nil
}

func (p *Pool) OpTimeout() time.Duration {
	return p.opTimeout
}
