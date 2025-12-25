package telemetry

import (
	"context"
	"database/sql"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracedDB wraps a sql.DB with tracing capabilities
type TracedDB struct {
	*sql.DB
	tracer trace.Tracer
}

// NewTracedDB creates a new traced database connection
func NewTracedDB(db *sql.DB, serviceName string) *TracedDB {
	return &TracedDB{
		DB:     db,
		tracer: Tracer(serviceName + ".database"),
	}
}

// QueryContext executes a query with tracing
func (tdb *TracedDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	ctx, span := tdb.tracer.Start(ctx, "database.query")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.statement", query),
	)

	start := time.Now()
	rows, err := tdb.DB.QueryContext(ctx, query, args...)
	duration := time.Since(start)

	span.SetAttributes(attribute.Int64("db.duration_ms", duration.Milliseconds()))

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "query executed successfully")
	return rows, nil
}

// QueryRowContext executes a query that returns a single row with tracing
func (tdb *TracedDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	ctx, span := tdb.tracer.Start(ctx, "database.query_row")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.statement", query),
	)

	start := time.Now()
	row := tdb.DB.QueryRowContext(ctx, query, args...)
	duration := time.Since(start)

	span.SetAttributes(attribute.Int64("db.duration_ms", duration.Milliseconds()))
	span.SetStatus(codes.Ok, "query row executed")

	return row
}

// ExecContext executes a command with tracing
func (tdb *TracedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ctx, span := tdb.tracer.Start(ctx, "database.exec")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.statement", query),
	)

	start := time.Now()
	result, err := tdb.DB.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	span.SetAttributes(attribute.Int64("db.duration_ms", duration.Milliseconds()))

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "exec completed successfully")
	return result, nil
}
