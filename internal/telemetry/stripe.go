package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// TraceStripeOperation traces a Stripe API call
func TraceStripeOperation(ctx context.Context, operation string, fn func() error) error {
	tracer := Tracer("stripe.client")
	ctx, span := tracer.Start(ctx, "stripe."+operation)
	defer span.End()

	span.SetAttributes(
		attribute.String("stripe.operation", operation),
		attribute.String("external.service", "stripe"),
	)

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	span.SetAttributes(attribute.Int64("stripe.duration_ms", duration.Milliseconds()))

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "stripe operation completed")
	return nil
}

// TraceStripeOperationWithResult traces a Stripe API call that returns a result
func TraceStripeOperationWithResult[T any](ctx context.Context, operation string, fn func() (T, error)) (T, error) {
	tracer := Tracer("stripe.client")
	ctx, span := tracer.Start(ctx, "stripe."+operation)
	defer span.End()

	span.SetAttributes(
		attribute.String("stripe.operation", operation),
		attribute.String("external.service", "stripe"),
	)

	start := time.Now()
	result, err := fn()
	duration := time.Since(start)

	span.SetAttributes(attribute.Int64("stripe.duration_ms", duration.Milliseconds()))

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		var zero T
		return zero, err
	}

	span.SetStatus(codes.Ok, "stripe operation completed")
	return result, nil
}
