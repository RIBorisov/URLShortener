package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shortener/internal/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStore struct {
	pool *pgxpool.Pool
}

func initPool(ctx context.Context, log *logger.Log, dsn string) (*pgxpool.Pool, error) {
	const (
		minConns = 1
		maxConns = 5
	)
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DatabaseDSN: %w", err)
	}
	poolCfg.ConnConfig.Tracer = &queryTracer{log: log}
	poolCfg.MinConns = minConns
	poolCfg.MaxConns = maxConns
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}
	return pool, nil
}

func New(ctx context.Context, dsn string, log *logger.Log) (*DBStore, error) {
	pool, err := initPool(ctx, log, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to init pool: %w", err)
	}

	if err = prepareDatabase(ctx, pool, log); err != nil {
		return nil, fmt.Errorf("failed to prepare database: %w", err)
	}

	return &DBStore{pool}, nil
}

func prepareDatabase(ctx context.Context, db *pgxpool.Pool, log *logger.Log) error {
	const (
		tableStmt = `CREATE TABLE IF NOT EXISTS urls (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    short VARCHAR(200) NOT NULL UNIQUE,
    long VARCHAR(200) NOT NULL
);`
		idxStmt = `CREATE UNIQUE INDEX IF NOT EXISTS idx_long_url ON urls (long);`
	)
	tx, err := db.BeginTx(ctx, pgx.TxOptions{IsoLevel: "read committed"})
	if err != nil {
		return fmt.Errorf("failed to begin the transaction: %w", err)
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				log.Err("failed to rollback the transaction: ", err)
			}
		}
	}()

	_, err = tx.Exec(ctx, tableStmt)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	_, err = tx.Exec(ctx, idxStmt)
	if err != nil {
		return fmt.Errorf("failed to set index on field long: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
	}

	return nil
}
