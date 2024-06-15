package storage

import (
	"context"
	"shortener/internal/logger"

	"github.com/jackc/pgx/v5"
)

type queryTracer struct {
	log *logger.Log
}

func (t *queryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	t.log.Debug("Running query %s (%v)", data.SQL, data.Args)
	return ctx
}

func (t *queryTracer) TraceQueryEnd(_ context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	t.log.Debug("CommandTag: %s, (%v)", data.CommandTag.String(), data.CommandTag)
}
