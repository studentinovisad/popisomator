package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

var Queries *repository.Queries

func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
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
