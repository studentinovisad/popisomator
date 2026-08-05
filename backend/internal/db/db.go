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
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// consumption_status is a Postgres enum; pgx needs its OID (and its array OID, for
	// ANY($1::consumption_status[]) filters) registered before it can encode/decode it.
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		enumType, err := conn.LoadType(ctx, "consumption_status")
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(enumType)

		arrayType, err := conn.LoadType(ctx, "_consumption_status")
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(arrayType)

		return nil
	}

	pool, err = pgxpool.NewWithConfig(ctx, config)
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
