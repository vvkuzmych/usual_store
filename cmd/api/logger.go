package main

import (
	"context"
	"log"
)

// ContextKey is a type for context keys
type ContextKey string

const TraceIDKey ContextKey = "TraceID"

// LogWithTrace logs a message with the trace ID from the context
func LogWithTrace(ctx context.Context, message string) {
	traceID, _ := ctx.Value(TraceIDKey).(string)
	if traceID != "" {
		log.Printf("[TraceID: %s] %s", traceID, message)
	} else {
		log.Printf("[TraceID: unknown] %s", message)
	}
}
