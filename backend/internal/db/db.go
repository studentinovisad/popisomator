package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

var Queries *repository.Queries
var pool *pgxpool.Pool

func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	var err error
	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	Queries = repository.New(pool)

	return pool, nil
}

func BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	return pool.Begin(ctx)
}
