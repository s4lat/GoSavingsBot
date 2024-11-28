package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func NewPostgresClient(
	ctx context.Context,
	connURL string,
	maxConns int32,
	connTimeout time.Duration,
) (*pgxpool.Pool, error) {
	// Create a configuration for the connection pool
	config, err := pgxpool.ParseConfig(connURL + "")
	if err != nil {
		return nil, err
	}
	config.ConnConfig.RuntimeParams["timezone"] = "UTC"

	// Set connection pool configurations
	config.MaxConns = maxConns
	config.ConnConfig.ConnectTimeout = connTimeout

	// Initialize the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
