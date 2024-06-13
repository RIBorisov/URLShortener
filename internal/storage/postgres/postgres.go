package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shortener/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBPool struct {
	*pgxpool.Pool
}

func initPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	const (
		minConns = 1
		maxConns = 5
	)
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DatabaseDSN: %w", err)
	}
	poolCfg.ConnConfig.Tracer = &queryTracer{}
	poolCfg.MinConns = minConns
	poolCfg.MaxConns = maxConns
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}
	return pool, nil
}

func New(ctx context.Context, dsn string) (*DBPool, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection with database: %w", err)
	}

	if err = prepareDatabase(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to prepare database: %w", err)
	}
	pool, err := initPool(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to init pool: %w", err)
	}
	return &DBPool{pool}, nil
}

func prepareDatabase(ctx context.Context, db *sql.DB) error {
	const (
		tableStmt = `CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short TEXT NOT NULL UNIQUE,
    long TEXT NOT NULL
);`
		idxStmt = `CREATE UNIQUE INDEX IF NOT EXISTS idx_long_url ON urls (long);`
	)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin the transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				logger.Err("failed to rollback the transaction", err)
			}
		}
	}()

	_, err = tx.ExecContext(ctx, tableStmt)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	_, err = tx.ExecContext(ctx, idxStmt)
	if err != nil {
		return fmt.Errorf("failed to set index on field long: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
	}

	return nil
}
